package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type CertificateService interface {
	Issue(token, externalID string, req *dto.CertificateIssueRequest) ([]byte, int, error)
	Active(token, externalID string, req *dto.CertificateActiveRequest) ([]byte, int, error)
	RevokeRequest(token, externalID string, req *dto.CertificateRevokeRequest) ([]byte, int, error)
	Revoke(token, externalID string, req *dto.CertificateRevokeValidateRequest) ([]byte, int, error)
}

type certificateService struct {
	certificateRepo  repository.CertificateRepository
	clientSvc        service.ClientService
	userSvc          service.UserService
	clientCompanySvc service.ClientCompanyService
}

func NewCertificateService(certificateRepo repository.CertificateRepository, clientSvc service.ClientService) CertificateService {
	return &certificateService{certificateRepo: certificateRepo, clientSvc: clientSvc}
}

func (s *certificateService) Issue(token, externalID string, req *dto.CertificateIssueRequest) ([]byte, int, error) {

	client, err := s.clientSvc.GetClientByExternalId(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	data, status, err := utils.PsreRequest("POST", "/certificate/issue", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre phone activation failed: %s", string(data))
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return data, status, fmt.Errorf("failed to parse psre response: %w", err)
	}
	if resp.Code == 0 {
		// Get user if external ID provided
		var user *model.User
		if req.UserID != "" {
			u, err := s.userSvc.GetUserByExternalID(req.UserID)
			if err != nil {
				return data, status, fmt.Errorf("failed to get user by external id: %w", err)
			}
			user = u
		}

		// Get company if external ID provided
		var company *model.ClientCompany
		if req.CompanyID != "" {
			c, err := s.clientCompanySvc.GetByExternalID(req.CompanyID)
			if err != nil {
				return data, status, fmt.Errorf("failed to get company by external id: %w", err)
			}
			company = c
		}

		// Build certificate entity
		cert := &model.Certificate{
			ClientID: client.ID, // pastikan variabel client sudah diisi sebelumnya
			Status:   "PROCESSED",
		}

		// Only assign optional IDs when available
		if user != nil {
			cert.UserID = &user.ID
			cert.ExternalUserID = &req.UserID
		}
		if company != nil {
			cert.CompanyID = &company.ID
			cert.ExternalCompanyID = &req.CompanyID
		}

		// Simpan ke database
		if err := s.certificateRepo.Create(cert); err != nil {
			return data, status, fmt.Errorf("failed to create certificate: %w", err)
		}
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
