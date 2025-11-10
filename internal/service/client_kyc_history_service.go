package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"time"
)

type ClientKYCHistoryService interface {
	Create(data *model.ClientKYCHistory) (*model.ClientKYCHistory, error)
	VerifyKYC(id int64) (*model.ClientKYCHistory, error)
	GetAll() ([]model.ClientKYCHistory, error)
	GetByID(id int64) (*model.ClientKYCHistory, error)
	GetByClientUserID(clientUserID int64) ([]model.ClientKYCHistory, error)
	Delete(id int64) error
}

type clientKYCHistoryService struct {
	repo repository.ClientKYCHistoryRepository
}

func NewClientKYCHistoryService(repo repository.ClientKYCHistoryRepository) ClientKYCHistoryService {
	return &clientKYCHistoryService{repo: repo}
}

func (s *clientKYCHistoryService) Create(data *model.ClientKYCHistory) (*model.ClientKYCHistory, error) {
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *clientKYCHistoryService) VerifyKYC(id int64) (*model.ClientKYCHistory, error) {
	data, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	data.IsVerify = true
	data.VerifyTime = &now

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *clientKYCHistoryService) GetAll() ([]model.ClientKYCHistory, error) {
	return s.repo.FindAll()
}

func (s *clientKYCHistoryService) GetByID(id int64) (*model.ClientKYCHistory, error) {
	return s.repo.FindByID(id)
}

func (s *clientKYCHistoryService) GetByClientUserID(clientUserID int64) ([]model.ClientKYCHistory, error) {
	return s.repo.FindByClientUserID(clientUserID)
}

func (s *clientKYCHistoryService) Delete(id int64) error {
	return s.repo.Delete(id)
}
