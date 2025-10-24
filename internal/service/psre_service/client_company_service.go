package psreservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ClientCompanyService interface {
	RegisterCompany(token string, body interface{}) ([]byte, int, error)
}
type clientCompanyService struct {
}

func NewClientCompanyService() ClientCompanyService {
	return &clientCompanyService{}
}
func (s *clientCompanyService) RegisterCompany(token string, body interface{}) ([]byte, int, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to marshal request body: %w", err)
	}
	psreUrl := os.Getenv("PSRE_BACKEND_URL")
	if psreUrl == "" {
		psreUrl = "http://10.100.20.14:2000" // fallback default
	}

	url := fmt.Sprintf("%s/client/company/create", psreUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Kalau bukan sukses, kirim balik ke handler biar bisa dibungkus JSON rapi
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return respBody, resp.StatusCode, fmt.Errorf("psre_error")
	}

	return respBody, resp.StatusCode, nil
}
