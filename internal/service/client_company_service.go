package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	// "dimensy-bridge/internal/service"
)

type ClientCompanyService interface {
	GetAll() ([]model.ClientCompany, error)
	GetByID(id int64) (*model.ClientCompany, error)
	Create(company *model.ClientCompany) error
	Update(company *model.ClientCompany) error
	Delete(id int64) error
	UpdateExternalID(id int64, externalID string) error

	GetByExternalID(externalID string) (*model.ClientCompany, error)
}

type clientCompanyService struct {
	repo           repository.ClientCompanyRepository
	quotaClientSvc QuotaClientService
}

func NewClientCompanyService(repo repository.ClientCompanyRepository, quotaClientSvc QuotaClientService) ClientCompanyService {
	return &clientCompanyService{repo, quotaClientSvc}
}

func (s *clientCompanyService) GetByExternalID(externalID string) (*model.ClientCompany, error) {
	return s.repo.FindByExternalID(externalID)
}

func (s *clientCompanyService) GetAll() ([]model.ClientCompany, error) {
	return s.repo.FindAll()
}

func (s *clientCompanyService) GetByID(id int64) (*model.ClientCompany, error) {
	return s.repo.FindByID(id)
}

func (s *clientCompanyService) Create(company *model.ClientCompany) error {
	return s.repo.Create(company)
}

func (s *clientCompanyService) Update(company *model.ClientCompany) error {
	return s.repo.Update(company)
}

func (s *clientCompanyService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *clientCompanyService) UpdateExternalID(id int64, externalID string) error {
	return s.repo.UpdateExternalID(id, externalID)
}
