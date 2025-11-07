package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type QuotaClientService interface {
	GetQuotas(page, limit int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	GetQuotaByID(id int64) (*model.QuotaClient, error)
	CreateQuota(quota *model.QuotaClient) error
	UpdateQuota(quota *model.QuotaClient) error
	DeleteQuota(id int64) error
	UseQuota(dto.UseQuotaClientRequest) (*model.QuotaClient, error)
}

type quotaClientService struct {
	db                   *gorm.DB
	quotaClientRepo      repository.QuotaClientRepository
	quotaClientReduction repository.QuotaClientReductionRepository
}

func NewQuotaClientService(db *gorm.DB, quotaClientRepo repository.QuotaClientRepository, quotaClientReduction repository.QuotaClientReductionRepository) QuotaClientService {
	return &quotaClientService{
		db:                   db,
		quotaClientRepo:      quotaClientRepo,
		quotaClientReduction: quotaClientReduction,
	}
}

func (s *quotaClientService) GetQuotas(page, limit int, filters map[string]interface{}) ([]model.QuotaClient, int64, error) {
	offset := (page - 1) * limit
	return s.quotaClientRepo.FindAll(limit, offset, filters)
}

func (s *quotaClientService) GetQuotaByID(id int64) (*model.QuotaClient, error) {
	return s.quotaClientRepo.FindByID(id)
}

func (s *quotaClientService) CreateQuota(quota *model.QuotaClient) error {
	return s.quotaClientRepo.Create(quota)
}

func (s *quotaClientService) UpdateQuota(quota *model.QuotaClient) error {
	return s.quotaClientRepo.Update(quota)
}

func (s *quotaClientService) DeleteQuota(id int64) error {
	return s.quotaClientRepo.Delete(id)
}

func (s *quotaClientService) UseQuota(req dto.UseQuotaClientRequest) (*model.QuotaClient, error) {
	var quota model.QuotaClient
	if err := s.db.First(&quota, req.ClientID).Error; err != nil {
		return nil, fmt.Errorf("quota not found: %w", err)
	}

	if quota.CurrentQuota < req.Quantity {
		return nil, fmt.Errorf("insufficient quota")
	}

	quota.CurrentQuota -= req.Quantity
	if err := s.db.Save(&quota).Error; err != nil {
		return nil, fmt.Errorf("failed to update quota: %w", err)
	}

	// Create a quota reduction record
	reduction := model.QuotaClientReduction{
		QuotaClientID: quota.ID,
		Quantity:      req.Quantity,
		LatestQuota:   quota.CurrentQuota,
		Type:          "usage",
		UsedBy:        req.UsedBy,
	}
	if err := s.quotaClientReduction.Create(&reduction); err != nil {
		return nil, fmt.Errorf("failed to create quota reduction record: %w", err)
	}
	return &quota, nil
}
