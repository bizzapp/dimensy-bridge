package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type CertificateService interface {
	Issue(token, externalID string, req *dto.CertificateIssueRequest) ([]byte, int, error)
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
