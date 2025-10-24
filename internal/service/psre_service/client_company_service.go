package psreservice

import (
	"dimensy-bridge/pkg/utils"
)

type ClientCompanyService interface {
	RegisterCompany(token string, body interface{}) ([]byte, int, error)
	GetCompany(token string, params map[string]string) ([]byte, int, error)
}
type clientCompanyService struct {
}

func NewClientCompanyService() ClientCompanyService {
	return &clientCompanyService{}
}

func (s *clientCompanyService) GetCompany(token string, params map[string]string) ([]byte, int, error) {
	return utils.PsreRequest("GET", "/client/company", nil, token, params)
}

func (s *clientCompanyService) RegisterCompany(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/client/company/create", body, token, nil)
}
