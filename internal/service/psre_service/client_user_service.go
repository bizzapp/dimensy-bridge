package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ClientUserService interface {
	Register(token string, body interface{}) ([]byte, int, error)
	RegisterUser(token, externalID string, req *dto.ClientUserRequest) ([]byte, int, error)
}

type clientUserService struct {
	db               *gorm.DB
	clientPsreSvc    service.ClientPsreService
	clientCompanySvc service.ClientCompanyService
	clientUserSvc    service.ClientUserService
}

func NewClientUserService(
	db *gorm.DB,
	clientPsreSvc service.ClientPsreService,
	clientCompanySvc service.ClientCompanyService,
	clientUserSvc service.ClientUserService,
) ClientUserService {
	return &clientUserService{
		db:               db,
		clientPsreSvc:    clientPsreSvc,
		clientCompanySvc: clientCompanySvc,
		clientUserSvc:    clientUserSvc,
	}
}

// --- panggil endpoint PSrE langsung ---
func (s *clientUserService) Register(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/user/register", body, token, nil)
}

// logic utama: create client_user + register ke PSrE dalam satu transaction
func (s *clientUserService) RegisterUser(token, externalID string, req *dto.ClientUserRequest) ([]byte, int, error) {
	// ambil data client PSrE
	clientPsre, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	// mulai transaction
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// mapping request ke model
		user := model.ClientUser{
			NIK:       req.NIK,
			Name:      req.FullName,
			Birthdate: req.BirthDate.Time,
			Email:     req.Email,
			Phone:     req.Phone,
			IsWNI:     req.IsWNI,
			ClientID:  clientPsre.ClientID,
		}

		// jika ada company
		if req.CompanyID != nil {
			clientCompany, err := s.clientCompanySvc.GetByExternalID(*req.CompanyID)
			if err != nil {
				return fmt.Errorf("failed get client company: %w", err)
			}
			user.ClientCompanyID = &clientCompany.ID
		}

		// create user di DB
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed create user: %w", err)
		}

		// panggil PSrE API
		data, st, err := s.Register(token, req)
		respBody, status = data, st // simpan agar bisa dikembalikan nanti
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		// kalau PSrE return error HTTP (400/500)
		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		// parse response PSrE
		var psreResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			UserID  string `json:"userId"`
		}
		if err := json.Unmarshal(data, &psreResp); err != nil {
			return errors.New("invalid psre response format")
		}

		// kalau PSrE gagal (code != 0)
		if psreResp.Code != 0 {
			return fmt.Errorf("psre register failed: %s", psreResp.Message)
		}

		// update external_id
		user.ExternalID = &psreResp.UserID
		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("failed update external_id: %w", err)
		}

		// commit otomatis kalau return nil
		return nil
	})

	// --- transaksi gagal (rollback otomatis oleh GORM) ---
	if txErr != nil {
		if respBody != nil {
			// Kalau PSrE mengembalikan body JSON (misalnya {"code":400,...})
			return respBody, status, txErr
		}
		// Kalau PSrE tidak merespons (network error, parsing error, dll)
		errorJSON, _ := json.Marshal(map[string]interface{}{
			"code":    400,
			"message": txErr.Error(),
		})
		return errorJSON, http.StatusBadRequest, txErr
	}

	// --- transaksi sukses ---
	return respBody, status, nil
}
