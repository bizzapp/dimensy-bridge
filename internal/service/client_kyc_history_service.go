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
	CreateOrUpdate(data *model.ClientKYCHistory) (*model.ClientKYCHistory, error)
	Update(data *model.ClientKYCHistory) (*model.ClientKYCHistory, error)
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
func (s *clientKYCHistoryService) Update(data *model.ClientKYCHistory) (*model.ClientKYCHistory, error) {
	if err := s.repo.Update(data); err != nil {
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

// CreateOrUpdate creates new KYC history or updates existing based on ExternalID and Signature
func (s *clientKYCHistoryService) CreateOrUpdate(data *model.ClientKYCHistory) (*model.ClientKYCHistory, error) {
	// Try to find existing KYC history by external_id and signature combination
	existing, err := s.repo.FindByExternalUserIDAndSignature(data.ExternalUserID, data.Signature)

	if err != nil {
		// If record not found, create new
		if err.Error() == "record not found" {
			if err := s.repo.Create(data); err != nil {
				return nil, err
			}
			return data, nil
		}
		// Other error occurred
		return nil, err
	}

	// Record exists, update it
	existing.IsVerify = data.IsVerify
	existing.VerifyTime = data.VerifyTime
	existing.ClientUserID = data.ClientUserID
	existing.Signature = data.Signature
	existing.Count = existing.Count + 1

	// Update other fields if provided
	if data.ClientID != 0 {
		existing.ClientID = data.ClientID
	}
	if data.ClientUserID != 0 {
		existing.ClientUserID = data.ClientUserID
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}
