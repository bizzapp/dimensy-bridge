package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type ClientCompanyService interface {
	GetAll() ([]model.ClientCompany, error)
	GetByID(id int64) (*model.ClientCompany, error)
	Create(company *model.ClientCompany) error
	Update(company *model.ClientCompany) error
	Delete(id int64) error
}

type clientCompanyService struct {
	repo repository.ClientCompanyRepository
}

func NewClientCompanyService(repo repository.ClientCompanyRepository) ClientCompanyService {
	return &clientCompanyService{repo}
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
