package psreservice

import (
	"context"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ClientCompanyService interface {
	RegisterCompany(token string, body interface{}) ([]byte, int, error)
	GetCompany(token string, params map[string]string) ([]byte, int, error)
	CreateClientCompany(ctx context.Context, authData interface{}, token string, req dto.PsreCreateClientCompanyRequest) ([]byte, int, error)
	DetailClientCompany(token string, id string) ([]byte, int, error)
	InviteClientCompany(token string, body interface{}) ([]byte, int, error)
}
type clientCompanyService struct {
	clientSvc         service.ClientService // gunakan interface
	clientCompanyRepo repository.ClientCompanyRepository
}

func NewClientCompanyService(clientSvc service.ClientService, clientCompanyRepo repository.ClientCompanyRepository) ClientCompanyService {
	return &clientCompanyService{
		clientSvc:         clientSvc,
		clientCompanyRepo: clientCompanyRepo,
	}
}
func (s *clientCompanyService) DetailClientCompany(token string, id string) ([]byte, int, error) {

	return utils.PsreRequest("GET", "/client/company/detail/"+id, nil, token, nil)
}

func (s *clientCompanyService) GetCompany(token string, params map[string]string) ([]byte, int, error) {
	return utils.PsreRequest("GET", "/client/company", nil, token, params)
}

func (s *clientCompanyService) RegisterCompany(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/client/company/create", body, token, nil)
}

func (s *clientCompanyService) InviteClientCompany(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/client/company/invite", body, token, nil)
}

func (s *clientCompanyService) CreateClientCompany(ctx context.Context, authData interface{}, token string, req dto.PsreCreateClientCompanyRequest) ([]byte, int, error) {
	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("unauthorized")
	}

	client, err := s.clientSvc.GetClientByExternalId(externalID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	clientCompany := model.ClientCompany{
		ClientID: client.ID,
		Name:     req.CompanyName,
		Address:  req.CompanyAddress,
		Industry: req.CompanyIndustry,
		NPWP:     req.NPWP,
		NIB:      req.NIB,
		PICName:  req.PICName,
		PICEmail: req.PICEmail,
	}

	if err := s.clientCompanyRepo.Create(&clientCompany); err != nil {
		return nil, http.StatusBadRequest, err
	}

	respBody, status, err := s.RegisterCompany(token, req)
	if err != nil {
		return respBody, status, fmt.Errorf("failed call psre api: %w", err)
	}

	var resp dto.PsreRegisterCompanyResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if resp.CompanyID != "" {
		_ = s.clientCompanyRepo.UpdateExternalID(clientCompany.ID, resp.CompanyID)
	}

	return respBody, status, nil
}
