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
)

type ClientUserService interface {
	Register(token string, body interface{}) ([]byte, int, error)
	RegisterUser(token, externalID string, req *dto.ClientUserRequest) ([]byte, int, error)
}

type clientUserService struct {
	clientPsreSvc    service.ClientPsreService
	clientCompanySvc service.ClientCompanyService
	clientUserSvc    service.ClientUserService
}

func NewClientUserService(
	clientPsreSvc service.ClientPsreService,
	clientCompanySvc service.ClientCompanyService,
	clientUserSvc service.ClientUserService,
) ClientUserService {
	return &clientUserService{
		clientPsreSvc:    clientPsreSvc,
		clientCompanySvc: clientCompanySvc,
		clientUserSvc:    clientUserSvc,
	}
}

// --- panggil endpoint PSrE langsung ---
func (s *clientUserService) Register(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/user/register", body, token, nil)
}

// --- logic utama register user ke PSrE dan simpan DB ---
func (s *clientUserService) RegisterUser(token, externalID string, req *dto.ClientUserRequest) ([]byte, int, error) {
	// ambil data client PSrE
	clientPsre, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	// map request ke model
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
			return nil, http.StatusBadRequest, fmt.Errorf("failed get client company: %w", err)
		}
		user.ClientCompanyID = &clientCompany.ID
	}

	// simpan user lokal
	userPtr, err := s.clientUserSvc.Create(&user)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed create user: %w", err)
	}

	// panggil PSrE
	data, status, err := s.Register(token, req)
	if err != nil {
		return data, status, err
	}

	// jika PSrE balas error (status >= 400), kembalikan body apa adanya
	if status >= 400 {
		return data, status, fmt.Errorf("psre error: %s", string(data))
	}

	// parsing response PSrE
	var psreResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		UserID  string `json:"userId"`
	}
	if err := json.Unmarshal(data, &psreResp); err != nil {
		return nil, http.StatusBadRequest, errors.New("invalid psre response format")
	}

	// update external_id kalau sukses
	if psreResp.Code == 0 && psreResp.UserID != "" {
		userPtr.ExternalID = &psreResp.UserID
		if err := s.clientUserSvc.Update(userPtr); err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("failed update external_id: %w", err)
		}
	}

	return data, http.StatusOK, nil
}
