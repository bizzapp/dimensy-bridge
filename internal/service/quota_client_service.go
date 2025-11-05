package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"

	"gorm.io/gorm"
)

type QuotaClientService interface {
	GetQuotas(page, limit int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	GetQuotaByID(id int64) (*model.QuotaClient, error)
	CreateQuota(quota *model.QuotaClient) error
	UpdateQuota(quota *model.QuotaClient) error
	DeleteQuota(id int64) error
}

type quotaClientService struct {
	db   *gorm.DB
	repo repository.QuotaClientRepository
}

func NewQuotaClientService(db *gorm.DB, repo repository.QuotaClientRepository) QuotaClientService {
	return &quotaClientService{db, repo}
}

func (s *quotaClientService) GetQuotas(page, limit int, filters map[string]interface{}) ([]model.QuotaClient, int64, error) {
	offset := (page - 1) * limit
	return s.repo.FindAll(limit, offset, filters)
}

func (s *quotaClientService) GetQuotaByID(id int64) (*model.QuotaClient, error) {
	return s.repo.FindByID(id)
}

func (s *quotaClientService) CreateQuota(quota *model.QuotaClient) error {
	return s.repo.Create(quota)
}

func (s *quotaClientService) UpdateQuota(quota *model.QuotaClient) error {
	return s.repo.Update(quota)
}

func (s *quotaClientService) DeleteQuota(id int64) error {
	return s.repo.Delete(id)
}
