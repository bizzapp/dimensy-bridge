package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/pkg/utils"
	"fmt"
	"net/url"
)

type BackendService interface {
	LoginBackend(req dto.PsreBackendLoginRequest) ([]byte, int, error)
	CreateClient(token string, req *dto.PsreBackendCreateClientRequest) ([]byte, int, error)
	ListClient(token string, filter, page, limit string) ([]byte, int, error)
	UpdateClient(id, token string, req *dto.PsreBackendCreateClientRequest) ([]byte, int, error)
	UpdateClientStatus(id, token string, req *dto.PsreBackendUpdateClientStatusRequest) ([]byte, int, error)
}

type backendService struct {
	// Add necessary dependencies here
}

func NewBackendService() BackendService {
	return &backendService{}
}
func (s *backendService) ListClient(token, filter, page, limit string) ([]byte, int, error) {

	query := fmt.Sprintf("/backend/client?page=%s&limit=%s", page, limit)
	if filter != "" {
		query += "&filter=" + url.QueryEscape(filter)
	}

	data, status, err := utils.PsreRequest("GET", query, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre get users failed: %s", string(data))
	}

	return data, status, nil
}
func (s *backendService) UpdateClient(id, token string, req *dto.PsreBackendCreateClientRequest) ([]byte, int, error) {
	path := fmt.Sprintf("/backend/client/update/%s", id)
	data, status, err := utils.PsreRequest("POST", path, req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre update client failed: %s", string(data))
	}

	return data, status, nil
}
func (s *backendService) UpdateClientStatus(id, token string, req *dto.PsreBackendUpdateClientStatusRequest) ([]byte, int, error) {
	path := fmt.Sprintf("/backend/client/update_status/%s", id)
	data, status, err := utils.PsreRequest("POST", path, req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre update client status failed: %s", string(data))
	}

	return data, status, nil
}
func (s *backendService) CreateClient(token string, req *dto.PsreBackendCreateClientRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/backend/client/create", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre phone activation failed: %s", string(data))
	}

	return data, status, nil
}

func (s *backendService) LoginBackend(req dto.PsreBackendLoginRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/backend/login", req, "", nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre login failed: %s", string(data))
	}
	return data, status, nil
}
