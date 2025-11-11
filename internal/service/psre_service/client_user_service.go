package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/model/seeder"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

type ClientUserService interface {
	GetUsers(token, externalID, filter, page, limit string) ([]byte, int, error)
	GetUserDetail(token, externalID, id string) ([]byte, int, error)
	RegisterUser(token, externalID string, req *dto.ClientUserRequest) ([]byte, int, error)
	ActivateUser(token, externalID string, req *dto.ClientUserActivateRequest) ([]byte, int, error)
	ResendActivationUser(token, externalID string, req *dto.ClientUserResendActivationRequest) ([]byte, int, error)
	RequestPhoneActivation(token, externalID string, req *dto.ClientUserRequestPhoneActivationRequest) ([]byte, int, error)
	PhoneActivation(token, externalID string, req *dto.ClientUserPhoneActivationRequest) ([]byte, int, error)
	RequestKYC(token, externalID string, req *dto.ClientUserKYCRequest) ([]byte, int, error)
	VerifyKYC(token, externalID string, req *dto.ClientUserVerifyKYCRequest) ([]byte, int, error)
	SyncUsers(token, externalID string) ([]byte, int, error)
}

type clientUserService struct {
	db                   *gorm.DB
	clientSvc            service.ClientService
	clientPsreSvc        service.ClientPsreService
	clientCompanySvc     service.ClientCompanyService
	clientUserSvc        service.ClientUserService
	clientUserRepo       repository.ClientUserRepository
	clientKYCHistorySvc  service.ClientKYCHistoryService
	clientKYCHistoryRepo repository.ClientKYCHistoryRepository
}

func NewClientUserService(
	db *gorm.DB,
	clientPsreSvc service.ClientPsreService,
	clientCompanySvc service.ClientCompanyService,
	clientUserSvc service.ClientUserService,
	clientUserRepo repository.ClientUserRepository,
	clientKYCHistorySvc service.ClientKYCHistoryService,
	clientSvc service.ClientService,
	clientKYCHistoryRepo repository.ClientKYCHistoryRepository,
) ClientUserService {
	return &clientUserService{
		db:                   db,
		clientPsreSvc:        clientPsreSvc,
		clientCompanySvc:     clientCompanySvc,
		clientUserSvc:        clientUserSvc,
		clientUserRepo:       clientUserRepo,
		clientKYCHistorySvc:  clientKYCHistorySvc,
		clientSvc:            clientSvc,
		clientKYCHistoryRepo: clientKYCHistoryRepo,
	}
}

func (s *clientUserService) GetUsers(token, externalID, filter, page, limit string) ([]byte, int, error) {
	query := fmt.Sprintf("/user?page=%s&limit=%s", page, limit)
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

func (s *clientUserService) GetUserDetail(token, externalID, id string) ([]byte, int, error) {
	path := fmt.Sprintf("/user/%s", id)

	data, status, err := utils.PsreRequest("GET", path, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre get user detail failed: %s", string(data))
	}

	return data, status, nil
}

// 🔹 REGISTER USER + TRANSACTION
func (s *clientUserService) RegisterUser(token, externalID string, req *dto.ClientUserRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		user := model.ClientUser{
			NIK:       &req.NIK,
			Name:      &req.FullName,
			Birthdate: &req.BirthDate.Time,
			Email:     &req.Email,
			Phone:     &req.Phone,
			IsWNI:     &req.IsWNI,
			ClientID:  client.ID,
		}

		// jika ada company
		if req.CompanyID != nil {
			clientCompany, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
			if err != nil {
				return fmt.Errorf("failed get client company: %w", err)
			}
			user.ClientCompanyID = &clientCompany.ID
		}

		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed create user: %w", err)
		}
		// utils.NewQuotaUtils().UseQuota(tx,)
		quantity := int64(1)
		_, err := utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: seeder.ID_PRODUCT_USER_PERSONAL,
			ClientID:        client.ID,
			Quantity:        quantity,
		})
		if err != nil {
			return fmt.Errorf("failed use quota: %w", err)
		}

		// 🔹 Call PSrE Register
		data, st, err := utils.PsreRequest("POST", "/user/register", req, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		var psreResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			UserID  string `json:"userId"`
		}
		if err := json.Unmarshal(data, &psreResp); err != nil {
			return errors.New("invalid psre response format")
		}

		if psreResp.Code != 0 {
			return fmt.Errorf("psre register failed: %s", psreResp.Message)
		}

		user.ExternalID = &psreResp.UserID
		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("failed update external_id: %w", err)
		}

		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}
func (s *clientUserService) ActivateUser(token, externalID string, req *dto.ClientUserActivateRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/activate", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre activate failed: %s", string(data))
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			UserID string `json:"userId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return data, status, fmt.Errorf("failed to parse psre response: %w", err)
	}

	// ✅ Panggil repository untuk update DB
	if resp.Code == 0 {
		if err := s.clientUserRepo.UpdateActiveStatus(resp.Data.UserID, true); err != nil {
			return data, status, fmt.Errorf("failed to update user active status: %w", err)
		}
	}

	return data, status, nil
}

func (s *clientUserService) RequestKYC(token, externalID string, req *dto.ClientUserKYCRequest) ([]byte, int, error) {
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

		// Get client user info
		clientUser, err := s.clientUserSvc.GetByExternalID(req.UserID)
		if err != nil {
			message := fmt.Sprintf("Unauthorized client user: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client user: %w", err)
		}

		// Call PSrE API
		data, psreStatus, err := utils.PsreRequest("POST", "/user/request-kyc", req, token, nil)
		respBody, status = data, psreStatus

		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if psreStatus >= 400 {
			return fmt.Errorf("psre request kyc failed: %s", string(data))
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
			kycHistory := model.ClientKYCHistory{
				ExternalUserID: req.UserID,
				Signature:      resp.Data.SignatureID,
				IsVerify:       false,
				ClientID:       client.ID,
				ClientUserID:   clientUser.ID,
			}

			// Check if KYC history already exists
			_, err := s.clientKYCHistoryRepo.GetBySignatureID(resp.Data.SignatureID)
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
				_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
					MasterProductID: seeder.ID_PRODUCT_KYC,
					ClientID:        client.ID,
					Quantity:        quantity,
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

// 🔹 RESEND ACTIVATION
func (s *clientUserService) ResendActivationUser(token, externalID string, req *dto.ClientUserResendActivationRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/resend-activation", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre resend activation failed: %s", string(data))
	}

	return data, status, nil
}

func (s *clientUserService) RequestPhoneActivation(token, externalID string, req *dto.ClientUserRequestPhoneActivationRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-phone-activation", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre request phone activation failed: %s", string(data))
	}

	return data, status, nil
}

func (s *clientUserService) PhoneActivation(token, externalID string, req *dto.ClientUserPhoneActivationRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/phone-activation", req, token, nil)
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
		if err := s.clientUserRepo.UpdateVerifyPhoneStatus(req.UserID, true); err != nil {
			return data, status, fmt.Errorf("failed to update user verify phone status: %w", err)
		}
	}

	return data, status, nil
}

func (s *clientUserService) VerifyKYC(token, externalID string, req *dto.ClientUserVerifyKYCRequest) ([]byte, int, error) {
	data, status, _ := utils.PsreRequest("POST", "/user/verify-kyc", req, token, nil)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return data, status, fmt.Errorf("failed to parse psre response: %w", err)
	}

	// Is Verify
	if resp.Code == 0 {
		clientKycUser, err := s.clientKYCHistoryRepo.GetBySignatureID(req.SignatureID)
		if err != nil {
			return data, status, fmt.Errorf("failed to get client kyc history: %w", err)
		}

		if err := s.clientUserRepo.UpdateVerifyKYCStatus(clientKycUser.ExternalUserID, true); err != nil {
			return data, status, fmt.Errorf("failed to update user verify phone status: %w", err)
		}
		if err := s.clientKYCHistoryRepo.UpdateIsVerifyStatus(req.SignatureID, true); err != nil {
			return data, status, fmt.Errorf("failed to update kyc history is_verify status: %w", err)
		}
	}
	// Is Reject
	if resp.Code == 90 {
		if err := s.clientKYCHistoryRepo.UpdateIsRejectStatus(req.SignatureID, true); err != nil {
			return data, status, fmt.Errorf("failed to update kyc history is_reject status: %w", err)
		}
	}

	return data, status, nil
}

// SyncUsers syncs user data from PSrE to local database
func (s *clientUserService) SyncUsers(token, externalID string) ([]byte, int, error) {
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

		// Get users from PSrE API
		data, psreStatus, err := utils.PsreRequest("GET", "/user", nil, token, map[string]string{
			"limit": "1000", // Get all users, adjust if needed
		})

		if err != nil {
			message := fmt.Sprintf("Failed to fetch users from PSrE: %v", err)
			respBody = utils.ResponseError(message, 500)
			status = 500
			return fmt.Errorf("PSrE request failed: %w", err)
		}

		if psreStatus >= 400 {
			respBody = data
			status = psreStatus
			return fmt.Errorf("PSrE returned HTTP %d: %s", psreStatus, string(data))
		}

		// Parse response from PSrE
		var psreResponse dto.PsreClientUserSyncResponse
		if err := json.Unmarshal(data, &psreResponse); err != nil {
			message := fmt.Sprintf("Failed to parse PSrE response: %v", err)
			respBody = utils.ResponseError(message, 500)
			status = 500
			return fmt.Errorf("failed to parse PSrE response: %w", err)
		}

		if psreResponse.Code != 0 {
			message := fmt.Sprintf("PSrE API error: %s", psreResponse.Message)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("PSrE API error: %s", psreResponse.Message)
		}

		// Process each user from PSrE
		syncedUsers := 0
		createdUsers := 0
		updatedUsers := 0

		for _, userData := range psreResponse.Data {
			// Parse birthdate
			birthdate, err := utils.ParseBirthDate(userData.BirthDate)
			if err != nil {
				fmt.Printf("Warning: Failed to parse birthdate for user %s: %v\n", userData.ID, err)
				continue
			}

			// Create ClientUser model
			clientUser := &model.ClientUser{
				ClientID:        client.ID,
				ClientCompanyID: nil, // Will be updated based on userCompany data if needed
				ExternalID:      &userData.ID,
				NIK:             &userData.NIK,
				Name:            &userData.FullName,
				Birthdate:       birthdate,
				Email:           &userData.Email,
				Phone:           &userData.Phone,
				IsWNI:           nil,   // Not provided in response, keep as nil
				IsActive:        true,  // Assume active if returned by API
				IsVerifyPhone:   false, // Default value
				IsVerifyKYC:     false, // Default value
				ParentID:        nil,
			}

			// Check if user already exists
			existingUser, err := s.clientUserRepo.FindByExternalID(userData.ID)
			isNewUser := err == gorm.ErrRecordNotFound

			if isNewUser {
				// Create new user
				if err := s.clientUserRepo.Create(clientUser); err != nil {
					fmt.Printf("Warning: Failed to create user %s: %v\n", userData.ID, err)
					continue
				}
				createdUsers++
			} else if err == nil {
				// Update existing user
				clientUser.ID = existingUser.ID
				clientUser.CreatedAt = existingUser.CreatedAt // Keep original creation time

				if err := s.clientUserRepo.Update(clientUser); err != nil {
					fmt.Printf("Warning: Failed to update user %s: %v\n", userData.ID, err)
					continue
				}
				updatedUsers++
			} else {
				// Other error
				fmt.Printf("Warning: Error checking user %s: %v\n", userData.ID, err)
				continue
			}

			syncedUsers++
		}

		// Create success response
		successMessage := fmt.Sprintf("Successfully synced %d users (created: %d, updated: %d)",
			syncedUsers, createdUsers, updatedUsers)

		respBody, _ = json.Marshal(map[string]interface{}{
			"code":    0,
			"message": successMessage,
			"data": map[string]interface{}{
				"total_synced": syncedUsers,
				"created":      createdUsers,
				"updated":      updatedUsers,
			},
		})
		status = 200

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
