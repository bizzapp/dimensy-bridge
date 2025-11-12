package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/model/seeder"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type CertificateService interface {
	Issue(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error)
	Active(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error)
	RevokeRequest(token, externalID string, req *dto.CertificateRevokeRequest) ([]byte, int, error)
	Revoke(token, externalID string, req *dto.CertificateRevokeValidateRequest) ([]byte, int, error)
}

type certificateService struct {
	db               *gorm.DB
	certificateRepo  repository.CertificateRepository
	clientSvc        service.ClientService
	userSvc          service.UserService
	clientCompanySvc service.ClientCompanyService
	clientUserSvc    service.ClientUserService
}

func NewCertificateService(db *gorm.DB, certificateRepo repository.CertificateRepository, clientSvc service.ClientService, userSvc service.UserService, clientCompanySvc service.ClientCompanyService, clientUserSvc service.ClientUserService) CertificateService {
	return &certificateService{
		db:               db,
		certificateRepo:  certificateRepo,
		clientSvc:        clientSvc,
		userSvc:          userSvc,
		clientCompanySvc: clientCompanySvc,
		clientUserSvc:    clientUserSvc,
	}
}

func (s *certificateService) Issue(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error) {
	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client info
		client, err := s.clientSvc.GetClientByExternalId(externalID)
		if err != nil {
			message := fmt.Sprintf("Unauthorized: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client: %w", err)
		}

		// Call PSrE API
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/issue", req, token, nil)
		respBody, status = data, psreStatus

		if err != nil {
			// Try Active as fallback
			fallbackData, fallbackStatus, fallbackErr := s.handleActiveAsFallbackWithResponse(tx, token, req, client.ID)
			if fallbackErr != nil {
				return fallbackErr
			}
			respBody, status = fallbackData, fallbackStatus
			return nil // Success with fallback
		}

		if psreStatus >= 400 {
			return fmt.Errorf("PSrE returned HTTP %d: %s", psreStatus, string(data))
		}

		var resp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to parse psre response: %w", err)
		}

		if resp.Code == 0 {
			// Get user if external ID provided
			var userID *int64
			if req.UserID != nil && *req.UserID != "" {
				u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
				if err != nil {
					return fmt.Errorf("failed to get user by external id: %w", err)
				}
				userID = &u.ID
			}

			// Get company if external ID provided
			var companyID *int64
			if req.CompanyID != nil && *req.CompanyID != "" {
				c, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
				if err != nil {
					return fmt.Errorf("failed to get company by external id: %w", err)
				}
				companyID = &c.ID
			}

			// Parse response data
			dataResp := dto.CertificateActiveResponseData{
				Status:       "PROCESSED",
				SerialNumber: "",
			}
			if err := json.Unmarshal(data, &dataResp); err != nil {
				return fmt.Errorf("failed to parse psre response data: %w", err)
			}

			// CreateOrUpdate certificate within transaction
			if err := s.createOrUpdateCertificateWithTx(tx, client.ID, userID, companyID, req, dataResp); err != nil {
				return fmt.Errorf("failed to create or update certificate: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *certificateService) Active(token, externalID string, req *dto.CertificateIssueActiveRequest) ([]byte, int, error) {
	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client info
		client, err := s.clientSvc.GetClientByExternalId(externalID)
		if err != nil {
			fmt.Printf("[Certificate Active] Failed to get client: %v\n", err)
			message := fmt.Sprintf("Unauthorized: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client: %w", err)
		}

		// Call PSrE API
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/active", req, token, nil)
		respBody, status = data, psreStatus

		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if psreStatus >= 400 {
			return fmt.Errorf("psre certificate activation failed: %s", string(data))
		}

		var resp dto.CertificateActiveResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to parse psre response: %w", err)
		}

		if resp.Code == 0 {
			// Get user if external ID provided
			var userID *int64
			if req.UserID != nil && *req.UserID != "" {
				u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
				if err != nil {
					return fmt.Errorf("failed to get user by external id: %w", err)
				}
				userID = &u.ID
			}

			// Get company if external ID provided
			var companyID *int64
			if req.CompanyID != nil && *req.CompanyID != "" {
				c, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
				if err != nil {
					return fmt.Errorf("failed to get company by external id: %w", err)
				}
				companyID = &c.ID
			}

			// CreateOrUpdate certificate within transaction
			if err := s.createOrUpdateCertificateWithTx(tx, client.ID, userID, companyID, req, resp.Data); err != nil {
				return fmt.Errorf("failed to create or update certificate: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
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

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Call PSrE API
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/revoke", req, token, nil)
		respBody, status = data, psreStatus
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if psreStatus >= 400 {
			return fmt.Errorf("psre certificate revocation failed: %s", string(data))
		}

		var resp dto.CertificateRevokeValidateResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			message := fmt.Sprintf("Failed to parse PSrE response: %v", err)
			respBody = utils.ResponseError(message, 500)
			status = 400
			return fmt.Errorf("failed to parse psre response: %w", err)
		}

		if resp.Code == 0 {
			// Update certificate status to REVOKED
			cert, err := s.certificateRepo.FindByExternal(req.UserID, req.CompanyID)
			if err != nil {
				return fmt.Errorf("failed to find certificate by serial number: %w", err)
			}

			cert.Status = "REVOKED"
			// Set deleted timestamp for soft delete
			cert.DeletedAt = gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			}
			if err := tx.Save(cert).Error; err != nil {
				return fmt.Errorf("failed to update certificate status: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

// createOrUpdateCertificateWithTx implements create or update logic within a transaction
func (s *certificateService) createOrUpdateCertificateWithTx(tx *gorm.DB, clientID int64, userID *int64, companyID *int64, req *dto.CertificateIssueActiveRequest, dataResp dto.CertificateActiveResponseData) error {
	// Try to find existing certificate
	existingCert, err := s.certificateRepo.FindByClientUserAndCompany(clientID, userID, companyID)

	if err != nil {
		// Certificate doesn't exist, create new one

		cert := &model.Certificate{
			ClientID: clientID,
			Status:   dataResp.Status,
		}

		masterProductID := int64(seeder.ID_PRODUCT_CA_COMPANY)
		if userID != nil {
			cert.ClienUserID = userID
			cert.ExternalUserID = req.UserID
			masterProductID = int64(seeder.ID_PRODUCT_CA_PERSONAL)
		}
		if companyID != nil {
			cert.CompanyID = companyID
			cert.ExternalCompanyID = req.CompanyID
			masterProductID = int64(seeder.ID_PRODUCT_CA_COMPANY)
		}
		if userID != nil && companyID != nil {
			masterProductID = int64(seeder.ID_PRODUCT_CA_PERSONAL_COMPANY)
		}
		cert.SerialNumber = &dataResp.SerialNumber

		quantity := 1
		_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: masterProductID,
			ClientID:        clientID,
			Quantity:        int64(quantity),
		})
		if err != nil {
			return fmt.Errorf("failed to use quota: %w", err)
		}

		return tx.Create(cert).Error
	}

	// Certificate exists, update it

	existingCert.Status = dataResp.Status
	existingCert.SerialNumber = &dataResp.SerialNumber

	// Update external IDs if provided
	if req.UserID != nil {
		existingCert.ExternalUserID = req.UserID
	}
	if req.CompanyID != nil {
		existingCert.ExternalCompanyID = req.CompanyID
	}

	return tx.Save(existingCert).Error
}

// handleActiveAsFallbackWithResponse handles fallback to Active when Issue fails and returns response
func (s *certificateService) handleActiveAsFallbackWithResponse(tx *gorm.DB, token string, req *dto.CertificateIssueActiveRequest, clientID int64) ([]byte, int, error) {
	// fmt.Printf("[Certificate] Trying Active as fallback for clientID=%d\n", clientID)

	// Call PSrE Active API
	data, psreStatus, err := utils.PsreRequest("POST", "/certificate/active", req, token, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("fallback active also failed: %w", err)
	}

	if psreStatus >= 400 {
		return data, psreStatus, fmt.Errorf("fallback active returned HTTP %d: %s", psreStatus, string(data))
	}

	var resp dto.CertificateActiveResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return data, psreStatus, fmt.Errorf("failed to parse fallback active response: %w", err)
	}

	if resp.Code == 0 {
		// Get user if external ID provided
		var userID *int64
		if req.UserID != nil && *req.UserID != "" {
			u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
			if err != nil {
				return data, psreStatus, fmt.Errorf("failed to get user by external id in fallback: %w", err)
			}
			userID = &u.ID
		}

		// Get company if external ID provided
		var companyID *int64
		if req.CompanyID != nil && *req.CompanyID != "" {
			c, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
			if err != nil {
				return data, psreStatus, fmt.Errorf("failed to get company by external id in fallback: %w", err)
			}
			companyID = &c.ID
		}

		// CreateOrUpdate certificate within transaction
		if err := s.createOrUpdateCertificateWithTx(tx, clientID, userID, companyID, req, resp.Data); err != nil {
			return data, psreStatus, fmt.Errorf("failed to create or update certificate in fallback: %w", err)
		}
	}

	return data, psreStatus, nil
}
