package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
)

type ClientService interface {
	GetClients(page, limit int, filters map[string]interface{}) ([]model.Client, int64, error)
	GetClientByID(id int64) (*model.Client, error)
	CreateClient(companyName, picName, email string) (*model.Client, error)
	UpdateClient(client *model.Client) error
	DeleteClient(id int64) error
}

type clientService struct {
	clientRepo repository.ClientRepository
	userRepo   repository.UserRepository
}

func NewClientService(clientRepo repository.ClientRepository, userRepo repository.UserRepository) ClientService {
	return &clientService{clientRepo, userRepo}
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
