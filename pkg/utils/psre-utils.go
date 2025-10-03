package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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

// doPsreRequest → core request ke PSRE backend (generic)
func doPsreRequest(method, path string, payload any) ([]byte, error) {
	psreURL := os.Getenv("PSRE_BACKEND_URL")
	if psreURL == "" {
		psreURL = "http://10.100.20.14:2000"
	}
	fullURL := psreURL + path

	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("gagal encode payload: %w", err)
		}
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// PSRE balikin 200 atau 201 untuk sukses
	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request ke PSRE gagal. status=%d, body=%s", resp.StatusCode, string(raw))
	}

	return io.ReadAll(resp.Body)
}

// psreLogin → khusus untuk ambil token dari PSRE
func psreLogin(username, password string) (string, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	raw, err := doPsreRequest("POST", "/backend/login", payload)
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

func GetPsreToken() (string, error) {
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
