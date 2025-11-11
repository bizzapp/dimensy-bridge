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
	Issue(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error)
	Active(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error)
	RevokeRequest(token, externalID string, req *dto.CertificateRevokeRequest) ([]byte, int, error)
	Revoke(token, externalID string, req *dto.CertificateRevokeValidateRequest) ([]byte, int, error)
}

type certificateService struct {
	certificateRepo  repository.CertificateRepository
	clientSvc        service.ClientService
	userSvc          service.UserService
	clientCompanySvc service.ClientCompanyService
	clientUserSvc    service.ClientUserService
}

func NewCertificateService(certificateRepo repository.CertificateRepository, clientSvc service.ClientService, userSvc service.UserService, clientCompanySvc service.ClientCompanyService, clientUserSvc service.ClientUserService) CertificateService {
	return &certificateService{
		certificateRepo:  certificateRepo,
		clientSvc:        clientSvc,
		userSvc:          userSvc,
		clientCompanySvc: clientCompanySvc,
		clientUserSvc:    clientUserSvc,
	}
}

func (s *certificateService) Issue(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error) {

	client, err := s.clientSvc.GetClientByExternalId(externalID)
	if err != nil {
		fmt.Printf("[Certificate Issue] Failed to get client: %v\n", err)
		return nil, http.StatusBadRequest, err
	}

	data, status, err := utils.PsreRequest("POST", "/certificate/issue", req, token, nil)
	if err != nil {
		s.Active(token, externalID, req)
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
		var user *model.ClientUser
		var userID *int64
		if req.UserID != nil && *req.UserID != "" {
			u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
			if err != nil {
				return data, status, fmt.Errorf("failed to get user by external id: %w", err)
			}
			user = u
			userID = &user.ID
		}

		// Get company if external ID provided
		var company *model.ClientCompany
		var companyID *int64
		if req.CompanyID != nil && *req.CompanyID != "" {
			c, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
			if err != nil {
				return data, status, fmt.Errorf("failed to get company by external id: %w", err)
			}
			company = c
			companyID = &company.ID
		}

		// CreateOrUpdate certificate berdasarkan clientID, userID, dan companyID
		if err := s.createOrUpdateCertificate(client.ID, userID, companyID, req, "PROCESSED"); err != nil {
			return data, status, fmt.Errorf("failed to create or update certificate: %w", err)
		}
	}

	return data, status, nil
}

func (s *certificateService) Active(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error) {
	client, err := s.clientSvc.GetClientByExternalId(externalID)
	if err != nil {
		fmt.Printf("[Certificate Active] Failed to get client: %v\n", err)
		return nil, http.StatusBadRequest, err
	}

	data, status, err := utils.PsreRequest("POST", "/certificate/active", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre certificate activation failed: %s", string(data))
	}

	// Parse response to check if activation was successful
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return data, status, fmt.Errorf("failed to parse psre response: %w", err)
	}

	if resp.Code == 0 {
		// Get user if external ID provided
		var userID *int64
		if req.UserID != nil && *req.UserID != "" {
			u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
			if err != nil {
				return data, status, fmt.Errorf("failed to get user by external id: %w", err)
			}
			userID = &u.ID
		}

		// Get company if external ID provided
		var companyID *int64
		if req.CompanyID != nil && *req.CompanyID != "" {
			c, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
			if err != nil {
				return data, status, fmt.Errorf("failed to get company by external id: %w", err)
			}
			companyID = &c.ID
		}

		// Convert DTO to CertificateIssueActiveRequest format for createOrUpdate
		issueReq := &dto.CertificateIssueActiveRequest{}
		if req.UserID != nil && *req.UserID != "" {
			issueReq.UserID = req.UserID
		}
		if req.CompanyID != nil && *req.CompanyID != "" {
			issueReq.CompanyID = req.CompanyID
		}

		// CreateOrUpdate certificate berdasarkan clientID, userID, dan companyID
		if err := s.createOrUpdateCertificate(client.ID, userID, companyID, issueReq, "ACTIVE"); err != nil {
			return data, status, fmt.Errorf("failed to create or update certificate: %w", err)
		}
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

// createOrUpdateCertificate implements create or update logic based on clientID, userID, and companyID combination
func (s *certificateService) createOrUpdateCertificate(clientID int64, userID *int64, companyID *int64, req *dto.CertificateIssueActiveRequest, status string) error {
	// Try to find existing certificate
	existingCert, err := s.certificateRepo.FindByClientUserAndCompany(clientID, userID, companyID)

	if err != nil {
		// Certificate doesn't exist, create new one
		fmt.Printf("[Certificate] Creating new certificate for clientID=%d, userID=%v, companyID=%v\n", clientID, userID, companyID)

		cert := &model.Certificate{
			ClientID: clientID,
			Status:   status,
		}

		// Set optional IDs
		if userID != nil {
			cert.ClienUserID = userID
			cert.ExternalUserID = req.UserID
		}
		if companyID != nil {
			cert.CompanyID = companyID
			cert.ExternalCompanyID = req.CompanyID
		}

		return s.certificateRepo.Create(cert)
	}

	// Certificate exists, update it
	fmt.Printf("[Certificate] Updating existing certificate ID=%d for clientID=%d, userID=%v, companyID=%v\n", existingCert.ID, clientID, userID, companyID)

	existingCert.Status = status

	// Update external IDs if provided
	if req.UserID != nil {
		existingCert.ExternalUserID = req.UserID
	}
	if req.CompanyID != nil {
		existingCert.ExternalCompanyID = req.CompanyID
	}

	return s.certificateRepo.Update(existingCert)
}
