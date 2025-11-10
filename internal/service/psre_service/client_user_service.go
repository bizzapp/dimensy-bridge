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
			NIK:       req.NIK,
			Name:      req.FullName,
			Birthdate: req.BirthDate.Time,
			Email:     req.Email,
			Phone:     req.Phone,
			IsWNI:     req.IsWNI,
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

	client, err := s.clientSvc.GetClientByExternalId(externalID)
	if err != nil {
		// return fmt.Errorf("unauthorized client: %w", err)
		return nil, http.StatusBadRequest, fmt.Errorf("unauthorized client: %w", err)
	}
	clientUser, err := s.clientUserSvc.GetByExternalID(req.UserID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unauthorized client user: %w", err)
	}

	data, status, err := utils.PsreRequest("POST", "/user/request-kyc", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre request kyc failed: %s", string(data))
	}

	var resp struct {
		Code int `json:"code"`
		Data struct {
			SignatureID string `json:"signatureId"`
			URL         string `json:"url"`
		}
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return data, status, fmt.Errorf("failed to parse psre response: %w", err)
	}
	if resp.Code == 0 {
		kycHistory := model.ClientKYCHistory{
			ExternalUserID: req.UserID,
			Signature:      resp.Data.SignatureID,
			IsVerify:       false,
			ClientID:       client.ID,
			ClientUserID:   clientUser.ID,
		}

		// Use CreateOrUpdate to handle existing records
		_, err := s.clientKYCHistorySvc.CreateOrUpdate(&kycHistory)
		if err != nil {
			return data, status, fmt.Errorf("failed to save kyc history: %w", err)
		}
	}

	return data, status, nil
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
