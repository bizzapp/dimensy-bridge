package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type CertificateV2Service interface {
	RequestIssueV2(token, externalID string, req *dto.CertificateRequestIssueV2Request) ([]byte, int, error)
	RevokeRequestV2(token, externalID string, req *dto.CertificateRequestIssueV2Request) ([]byte, int, error)
	IssueV2(token, externalID string, req *dto.CertificateIssueV2Request) ([]byte, int, error)
	RevokeV2(token, externalID string, req *dto.CertificateIssueV2Request) ([]byte, int, error)
}

type certificateV2Service struct {
	db                   *gorm.DB
	certificateRepo      repository.CertificateRepository
	clientSvc            service.ClientService
	userSvc              service.UserService
	clientCompanySvc     service.ClientCompanyService
	clientUserSvc        service.ClientUserService
	clientKYCHistoryRepo repository.ClientKYCHistoryRepository
	certificateService   CertificateService
}

func NewCertificateV2Service(db *gorm.DB, certificateRepo repository.CertificateRepository, clientSvc service.ClientService, userSvc service.UserService, clientCompanySvc service.ClientCompanyService, clientUserSvc service.ClientUserService,
	clientKYCHistoryRepo repository.ClientKYCHistoryRepository, certificateService CertificateService) CertificateV2Service {
	return &certificateV2Service{
		db:                   db,
		certificateRepo:      certificateRepo,
		clientSvc:            clientSvc,
		userSvc:              userSvc,
		clientCompanySvc:     clientCompanySvc,
		clientUserSvc:        clientUserSvc,
		clientKYCHistoryRepo: clientKYCHistoryRepo,
		certificateService:   certificateService,
	}
}

func (s *certificateV2Service) RequestIssueV2(token, externalID string, req *dto.CertificateRequestIssueV2Request) ([]byte, int, error) {
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
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/v2/request-issue", req, token, nil)
		respBody, status = data, psreStatus

		userID := req.UserID
		if userID == nil {
			userID = req.UserPicID
		}

		var requestIssueReq = dto.CertificateIssueActiveRequest{
			UserID:    userID,
			CompanyID: req.CompanyID,
		}

		if err != nil {
			// Try Active as fallback
			fallbackData, fallbackStatus, fallbackErr := s.certificateService.HandleActiveAsFallbackWithResponse(tx, token, &requestIssueReq, client.ID)
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
			Code int `json:"code"`
			Data struct {
				SignatureID string `json:"signatureId"`
				URL         string `json:"url"`
			}
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to parse psre response: %w", err)
		}

		if resp.Code == 0 {

			clientUser, err := s.clientUserSvc.GetByExternalID(*userID)
			if err != nil {
				return fmt.Errorf("failed to get client user by id: %w", err)
			}
			kycHistory := model.ClientKYCHistory{
				ExternalUserID: *userID,
				Signature:      resp.Data.SignatureID,
				IsVerify:       false,
				ClientID:       client.ID,
				ClientUserID:   clientUser.ID,
				Type:           model.TypeLivenessGenerateCertificate,
			}
			// Check if KYC history already exists
			_, err = s.clientKYCHistoryRepo.GetBySignatureID(resp.Data.SignatureID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to get client kyc history: %w", err)
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new KYC history record
				if err := s.clientKYCHistoryRepo.CreateTx(tx, &kycHistory); err != nil {
					return fmt.Errorf("failed to save kyc history: %w", err)
				}

				// Use quota for KYC request
				quantity := int64(1)
				_, err = utils.NewQuotaUtils().UseQuotaEkyc(tx, dto.UseQuotaClientRequest{
					ClientID: client.ID,
					Quantity: quantity,
				})
				if err != nil {
					status = http.StatusBadRequest
					respBody = utils.ResponseError(err.Error(), status)
					return fmt.Errorf("failed use quota: %w", err)
				}
			} else {
				// Update existing KYC history record
				if err := s.clientKYCHistoryRepo.UpdateTx(tx, &kycHistory); err != nil {
					return fmt.Errorf("failed to update kyc history: %w", err)
				}
			}
			// Get user if external ID provided
			var userID *int64
			if req.UserID != nil {
				u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
				if err != nil {
					return fmt.Errorf("failed to get user by external id: %w", err)
				}
				userID = &u.ID
			}

			// Get company if external ID provided
			var companyID *int64
			if req.CompanyID != nil {
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
			if err := s.certificateService.CreateOrUpdateCertificateWithTx(tx, client.ID, userID, companyID, &requestIssueReq, dataResp); err != nil {
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

func (s *certificateV2Service) IssueV2(token, externalID string, req *dto.CertificateIssueV2Request) ([]byte, int, error) {

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client info
		_, err := s.clientSvc.GetClientByExternalId(externalID)
		if err != nil {
			message := fmt.Sprintf("Unauthorized: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client: %w", err)
		}

		kycHistory, err := s.clientKYCHistoryRepo.GetBySignatureID(req.SignatureID)
		if err != nil {
			message := fmt.Sprintf("Unauthorized: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client: %w", err)
		}

		clientCallbackURL := req.CallbackURL
		req.CallbackURL = os.Getenv("APP_URL_WEBHOOK_CERTIFICATE")

		kycHistory.CallbackURL = &req.CallbackURL
		kycHistory.ClientCallbackURL = &clientCallbackURL

		if err := s.clientKYCHistoryRepo.UpdateTx(tx, kycHistory); err != nil {
			return fmt.Errorf("failed to update client kyc history: %w", err)
		}

		// Call PSrE API
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/v2/issue", req, token, nil)
		respBody, status = data, psreStatus

		return err
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *certificateV2Service) RevokeRequestV2(token, externalID string, req *dto.CertificateRequestIssueV2Request) ([]byte, int, error) {
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
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/v2/revoke-request", req, token, nil)
		respBody, status = data, psreStatus

		userID := req.UserID
		if userID == nil {
			userID = req.UserPicID
		}

		var requestIssueReq = dto.CertificateIssueActiveRequest{
			UserID:    userID,
			CompanyID: req.CompanyID,
		}

		if err != nil {
			// Try Active as fallback
			fallbackData, fallbackStatus, fallbackErr := s.certificateService.HandleActiveAsFallbackWithResponse(tx, token, &requestIssueReq, client.ID)
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
			Code int `json:"code"`
			Data struct {
				SignatureID string `json:"signatureId"`
				URL         string `json:"url"`
			}
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to parse psre response: %w", err)
		}

		if resp.Code == 0 {

			clientUser, err := s.clientUserSvc.GetByExternalID(*req.UserID)
			if err != nil {
				return fmt.Errorf("failed to get client user by id: %w", err)
			}
			kycHistory := model.ClientKYCHistory{
				ExternalUserID: *userID,
				Signature:      resp.Data.SignatureID,
				IsVerify:       false,
				ClientID:       client.ID,
				ClientUserID:   clientUser.ID,
				Type:           model.TypeLivenessRevokeCertificate,
			}
			// Check if KYC history already exists
			_, err = s.clientKYCHistoryRepo.GetBySignatureID(resp.Data.SignatureID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to get client kyc history: %w", err)
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new KYC history record
				if err := s.clientKYCHistoryRepo.CreateTx(tx, &kycHistory); err != nil {
					return fmt.Errorf("failed to save kyc history: %w", err)
				}

				// Use quota for KYC request
				quantity := int64(1)
				_, err = utils.NewQuotaUtils().UseQuotaEkyc(tx, dto.UseQuotaClientRequest{
					ClientID: client.ID,
					Quantity: quantity,
				})
				if err != nil {
					status = http.StatusBadRequest
					respBody = utils.ResponseError(err.Error(), status)
					return fmt.Errorf("failed use quota: %w", err)
				}
			} else {
				// Update existing KYC history record
				if err := s.clientKYCHistoryRepo.UpdateTx(tx, &kycHistory); err != nil {
					return fmt.Errorf("failed to update kyc history: %w", err)
				}
			}
			// Get user if external ID provided
			var userID *int64
			if req.UserID != nil {
				u, err := s.clientUserSvc.GetByExternalID(*req.UserID)
				if err != nil {
					return fmt.Errorf("failed to get user by external id: %w", err)
				}
				userID = &u.ID
			}

			// Get company if external ID provided
			var companyID *int64
			if req.CompanyID != nil {
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
			if err := s.certificateService.CreateOrUpdateCertificateWithTx(tx, client.ID, userID, companyID, &requestIssueReq, dataResp); err != nil {
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

func (s *certificateV2Service) RevokeV2(token, externalID string, req *dto.CertificateIssueV2Request) ([]byte, int, error) {

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client info
		_, err := s.clientSvc.GetClientByExternalId(externalID)
		if err != nil {
			message := fmt.Sprintf("Unauthorized: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client: %w", err)
		}

		// Call PSrE API
		data, psreStatus, err := utils.PsreRequest("POST", "/certificate/v2/revoke", req, token, nil)
		respBody, status = data, psreStatus

		return err
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}
