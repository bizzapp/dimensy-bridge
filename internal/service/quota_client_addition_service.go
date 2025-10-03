package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
)

type QuotaClientAdditionService interface {
	GetAdditions(page, limit int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error)
	GetAdditionByID(id int64) (*model.QuotaClientAddition, error)
	CreateAddition(addition *model.QuotaClientAddition) error
	UpdateAddition(addition *model.QuotaClientAddition) error
	DeleteAddition(id int64) error
	ProcessAddition(id int64, processBy int64) (*model.QuotaClientAddition, error)
}

type quotaClientAdditionService struct {
	additionRepo repository.QuotaClientAdditionRepository
	quotaRepo    repository.QuotaClientRepository
}

func NewQuotaClientAdditionService(additionRepo repository.QuotaClientAdditionRepository, quotaRepo repository.QuotaClientRepository) QuotaClientAdditionService {
	return &quotaClientAdditionService{additionRepo, quotaRepo}
}

func (s *quotaClientAdditionService) GetAdditions(page, limit int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error) {
	offset := (page - 1) * limit
	return s.additionRepo.FindAll(limit, offset, filters)
}

func (s *quotaClientAdditionService) GetAdditionByID(id int64) (*model.QuotaClientAddition, error) {
	return s.additionRepo.FindByID(id)
}

func (s *quotaClientAdditionService) CreateAddition(addition *model.QuotaClientAddition) error {
	addition.IsProcess = false
	return s.additionRepo.Create(addition)
}

func (s *quotaClientAdditionService) UpdateAddition(addition *model.QuotaClientAddition) error {
	return s.additionRepo.Update(addition)
}

func (s *quotaClientAdditionService) DeleteAddition(id int64) error {
	return s.additionRepo.Delete(id)
}

// Process menambah quota_client sesuai quantity
func (s *quotaClientAdditionService) ProcessAddition(id int64, processBy int64) (*model.QuotaClientAddition, error) {
	addition, err := s.additionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if addition.IsProcess {
		return nil, errors.New("quota addition sudah diproses")
	}

	quota, err := s.quotaRepo.FindByID(addition.QuotaClientID)
	if err != nil {
		return nil, err
	}

	// update quota_clients
	quota.Quantity += addition.Quantity
	quota.CurrentQuota += addition.Quantity
	if err := s.quotaRepo.Update(quota); err != nil {
		return nil, err
	}

	// update addition
	addition.IsProcess = true
	addition.ProcessBy = &processBy
	addition.LatestQuota = quota.CurrentQuota
	if err := s.additionRepo.Update(addition); err != nil {
		return nil, err
	}

	return addition, nil
}
