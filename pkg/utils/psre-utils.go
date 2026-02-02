package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ====== Token Cache ======
var (
	tokenCache     string
	tokenExpiresAt time.Time
	mu             sync.Mutex
)

// ====== PSRE Log Repository ======
var psreLogRepo interface {
	Create(data interface{}) error
}

// ====== Structs ======
type loginResponse struct {
	Code int `json:"code,omitempty"`
	Data struct {
		AccessToken string `json:"accessToken,omitempty"`
		Token       string `json:"token,omitempty"`
	} `json:"data,omitempty"`
	Token string `json:"token,omitempty"`
}

// ====== Utility Functions ======
func DefaultPassword() string {
	if pass := os.Getenv("PSRE_DEFAULT_PASSWORD"); pass != "" {
		return pass
	}
	return "DefaultP@ssw0rd!"
}

func ExpireDate() time.Time {
	exp := time.Now()
	if days, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_DAYS")); days > 0 {
		return exp.AddDate(0, 0, days)
	}
	if months, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_MONTHS")); months > 0 {
		return exp.AddDate(0, months, 0)
	}
	if years, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_YEARS")); years > 0 {
		return exp.AddDate(years, 0, 0)
	}
	return exp.AddDate(1, 0, 0) // default 1 tahun
}

func ExtractExternalID(authData any) (string, error) {
	m1, ok := authData.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid authData")
	}

	level1, ok := m1["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("missing data")
	}

	id, ok := level1["id"].(string)
	if !ok {
		return "", errors.New("missing id field")
	}

	return id, nil
}

// SetPsreLogRepository sets the repository for logging PSRE requests
func SetPsreLogRepository(repo interface {
	Create(data interface{}) error
}) {
	psreLogRepo = repo
}

// ====== Core HTTP Utility ======
// LogCallback is a function type for logging PSRE requests and responses
type LogCallback func(description, jsonHeader, jsonContent, fullURL string)

func PsreRequest(method, path string, payload any, token string, queryParams map[string]string) ([]byte, int, error) {
	return PsreRequestWithLogging(method, path, payload, token, queryParams, nil)
}

// PsreRequestWithLogging sends a request to PSRE and logs it using the provided callback
func PsreRequestWithLogging(method, path string, payload any, token string, queryParams map[string]string, logCallback LogCallback) ([]byte, int, error) {
	// Base URL
	baseURL := os.Getenv("PSRE_BACKEND_URL")
	if baseURL == "" {
		baseURL = "http://10.100.20.14:2000"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Buat URL lengkap
	reqURL, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("invalid PSRE URL: %w", err)
	}

	// Tambahkan query params jika ada
	if len(queryParams) > 0 {
		q := reqURL.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		reqURL.RawQuery = q.Encode()
	}

	// Siapkan payload
	var body io.Reader
	var payloadBytes []byte
	if payload != nil {
		var err error
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(payloadBytes)
	}

	// Buat request
	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	// Extract headers for logging
	headerJSON, _ := json.Marshal(req.Header)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Log failed request to database
		if psreLogRepo != nil {
			logDescription := "PSRE Request Failed"
			logFullURL := reqURL.String()
			headerJSONStr := string(headerJSON)

			logEntry := map[string]interface{}{
				"description":  logDescription,
				"json_header":  &headerJSONStr,
				"json_content": string(payloadBytes),
				"full_url":     &logFullURL,
			}

			_ = psreLogRepo.Create(logEntry)
		}

		// Log failed request via callback
		if logCallback != nil {
			logCallback("PSRE Request Failed", string(headerJSON), string(payloadBytes), reqURL.String())
		}
		return nil, http.StatusBadGateway, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Auto-log the request and response to database
	if psreLogRepo != nil {
		respLog := map[string]interface{}{
			"status_code": resp.StatusCode,
			"body":        string(respBody),
		}
		respLogJSON, _ := json.Marshal(respLog)

		// Create log entry
		logDescription := "PSRE Request"
		logFullURL := reqURL.String()
		headerJSONStr := string(headerJSON)

		logEntry := map[string]interface{}{
			"description":  logDescription,
			"json_header":  &headerJSONStr,
			"json_content": string(respLogJSON),
			"full_url":     &logFullURL,
		}

		// Ignore logging errors - don't let them affect the main request
		_ = psreLogRepo.Create(logEntry)
	}

	// Log the request and response (if callback provided)
	if logCallback != nil {
		respLog := map[string]interface{}{
			"status_code": resp.StatusCode,
			"body":        string(respBody),
		}
		respLogJSON, _ := json.Marshal(respLog)
		logCallback("PSRE Request", string(headerJSON), string(respLogJSON), reqURL.String())
	}

	// Jika error dari PSRE
	if resp.StatusCode >= 400 {
		return respBody, resp.StatusCode, fmt.Errorf("PSRE error %d: %s", resp.StatusCode, string(respBody))
	}
	// fmt.Println(string(respBody))

	return respBody, resp.StatusCode, nil
}

// ====== Token Management ======
func psreLogin(username, password string) (string, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	respBody, _, err := PsreRequest("POST", "/backend/login", payload, "", nil)
	if err != nil {
		return "", err
	}

	var res loginResponse
	if err := json.Unmarshal(respBody, &res); err != nil {
		return "", fmt.Errorf("gagal decode response login: %v | body=%s", err, string(respBody))
	}

	switch {
	case res.Data.AccessToken != "":
		return res.Data.AccessToken, nil
	case res.Data.Token != "":
		return res.Data.Token, nil
	case res.Token != "":
		return res.Token, nil
	default:
		return "", fmt.Errorf("access token tidak ditemukan. response=%s", string(respBody))
	}
}

func GetAdministratorToken() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	// Jika token masih valid, return cached
	if tokenCache != "" && time.Now().Before(tokenExpiresAt) {
		return tokenCache, nil
	}

	username := os.Getenv("PSRE_ADMIN_USERNAME")
	password := os.Getenv("PSRE_ADMIN_PASSWORD")
	if username == "" || password == "" {
		return "", errors.New("PSRE_ADMIN_USERNAME / PSRE_ADMIN_PASSWORD belum di-set")
	}

	token, err := psreLogin(username, password)
	if err != nil {
		return "", err
	}

	tokenCache = token
	tokenExpiresAt = time.Now().Add(55 * time.Minute)
	return tokenCache, nil
}
