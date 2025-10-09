package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	tokenCache     string
	tokenExpiresAt time.Time
	mu             sync.Mutex
)

type loginResponse struct {
	Code int `json:"code,omitempty"`
	Data struct {
		AccessToken string `json:"accessToken,omitempty"`
		Token       string `json:"token,omitempty"`
	} `json:"data,omitempty"`
	Token string `json:"token,omitempty"`
}

func DefaultPassword() string{

	password := os.Getenv("PSRE_DEFAULT_PASSWORD")
	if password == "" {
		password = "DefaultP@ssw0rd!" // fallback default
	}
	return password
}

func ExpireDate() time.Time {

	expDate := time.Now()
	if days, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_DAYS")); days > 0 {
		expDate = expDate.AddDate(0, 0, days)
	} else if months, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_MONTHS")); months > 0 {
		expDate = expDate.AddDate(0, months, 0)
	} else if years, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_YEARS")); years > 0 {
		expDate = expDate.AddDate(years, 0, 0)
	} else {
		// default 1 tahun
		expDate = expDate.AddDate(1, 0, 0)
	}
	return expDate
}

// psreLogin → khusus untuk ambil token dari PSRE
func psreLogin(username, password string) (string, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	raw, err := DoPsreRequest("POST", "/backend/login", payload,nil)
	if err != nil {
		return "", err
	}

	var res loginResponse
	if err := json.Unmarshal(raw, &res); err != nil {
		return "", fmt.Errorf("gagal decode response login: %v | body=%s", err, string(raw))
	}

	switch {
	case res.Data.AccessToken != "":
		return res.Data.AccessToken, nil
	case res.Data.Token != "":
		return res.Data.Token, nil
	case res.Token != "":
		return res.Token, nil
	default:
		return "", fmt.Errorf("access token tidak ditemukan. response=%s", string(raw))
	}
}

func GetAdministratorToken() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if tokenCache != "" && time.Now().Before(tokenExpiresAt) {
		return tokenCache, nil
	}

	username := os.Getenv("PSRE_ADMIN_USERNAME")
	password := os.Getenv("PSRE_ADMIN_PASSWORD")
	if username == "" || password == "" {
		return "", errors.New("PSRE_ADMIN_USERNAME / PSRE_ADMIN_PASSWORD env belum di set")
	}

	token, err := psreLogin(username, password)
	if err != nil {
		return "", err
	}

	tokenCache = token
	tokenExpiresAt = time.Now().Add(55 * time.Minute)

	return tokenCache, nil
}

// DoPsreRequest adalah core utilitas untuk kirim request ke PSRE API
func DoPsreRequest(method, path string, payload any, headers map[string]string) ([]byte, error) {
	// marshal payload jadi JSON
	var bodyBytes []byte
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %v", err)
		}
		bodyBytes = b
	}

url := os.Getenv("PSRE_BACKEND_URL") + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// set default headers
	req.Header.Set("Content-Type", "application/json")

	// tambahan header custom (misal Authorization)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return resBody, fmt.Errorf("PSRE error %d: %s", resp.StatusCode, string(resBody))
	}

	return resBody, nil
}
