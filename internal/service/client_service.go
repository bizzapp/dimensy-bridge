package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type ClientService interface {
	GetClients(page, limit int, filters map[string]interface{}) ([]model.Client, int64, error)
	GetClientByID(id int64) (*model.Client, error)
	CreateClient(companyName, picName, email string) (*model.Client, error)
	UpdateClient(client *model.Client) error
	DeleteClient(id int64) error
	GetClientByExternalId(externalID string) (*model.Client, error)

	AddQuota(dto.AddQuotaClientRequest) error
}

type clientService struct {
	tx                      *gorm.DB
	clientRepo              repository.ClientRepository
	userRepo                repository.UserRepository
	quotaClientRepo         repository.QuotaClientRepository
	quotaClientAdditionRepo repository.QuotaClientAdditionRepository
	quotaClientSvc          QuotaClientService
}

func NewClientService(clientRepo repository.ClientRepository, userRepo repository.UserRepository, quotaClientRepo repository.QuotaClientRepository, quotaClientAdditionRepo repository.QuotaClientAdditionRepository, quotaClientSvc QuotaClientService) ClientService {
	return &clientService{
		clientRepo:              clientRepo,
		userRepo:                userRepo,
		quotaClientRepo:         quotaClientRepo,
		quotaClientAdditionRepo: quotaClientAdditionRepo,
		quotaClientSvc:          quotaClientSvc,
	}
}

func (s *clientService) GetClientByExternalId(externalID string) (*model.Client, error) {
	client, err := s.clientRepo.FindByExternalID(externalID)
	if err != nil {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (s *clientService) AddQuota(req dto.AddQuotaClientRequest) error {
	// Use quota client service to handle both quota update and addition record
	_, err := s.quotaClientSvc.AddQuota(s.tx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *clientService) GetClients(page, limit int, filters map[string]interface{}) ([]model.Client, int64, error) {
	offset := (page - 1) * limit
	return s.clientRepo.FindAll(limit, offset, filters)
}

func (s *clientService) GetClientByID(id int64) (*model.Client, error) {
	return s.clientRepo.FindByID(id)
}

func (s *clientService) CreateClient(companyName, picName, email string) (*model.Client, error) {
	if email == "" {
		return nil, errors.New("email wajib diisi")
	}

	// buat user untuk client
	user := model.User{
		Name:  companyName,
		Email: &email,
		Role:  "client",
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	client := model.Client{
		CompanyName: companyName,
		PicName:     picName,
		UserID:      user.ID,
	}
	if err := s.clientRepo.Create(&client); err != nil {
		return nil, err
	}

	return &client, nil
}

func (s *clientService) UpdateClient(client *model.Client) error {
	return s.clientRepo.Update(client)
}

func (s *clientService) DeleteClient(id int64) error {
	return s.clientRepo.Delete(id)
}
