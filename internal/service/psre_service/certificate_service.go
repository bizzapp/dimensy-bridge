package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type CertificateService interface {
	Issue(token, externalID string, req *dto.CertificateIssueRequest) ([]byte, int, error)
	Active(token, externalID string, req *dto.CertificateActiveRequest) ([]byte, int, error)
	RevokeRequest(token, externalID string, req *dto.CertificateRevokeRequest) ([]byte, int, error)
	Revoke(token, externalID string, req *dto.CertificateRevokeValidateRequest) ([]byte, int, error)
}

type certificateService struct {
}

func NewCertificateService() CertificateService {
	return &certificateService{}
}

func (s *certificateService) Issue(token, externalID string, req *dto.CertificateIssueRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/certificate/issue", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre phone activation failed: %s", string(data))
	}

	return data, status, nil
}

func (s *certificateService) Active(token, externalID string, req *dto.CertificateActiveRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/certificate/active", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre phone activation failed: %s", string(data))
	}

	return data, status, nil
}

func (s *certificateService) RevokeRequest(token, externalID string, req *dto.CertificateRevokeRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/certificate/revoke-request", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre phone activation failed: %s", string(data))
	}

	return data, status, nil
}
func (s *certificateService) Revoke(token, externalID string, req *dto.CertificateRevokeValidateRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/certificate/revoke", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre phone activation failed: %s", string(data))
	}

	return data, status, nil
}
