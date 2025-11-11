package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"

	"gorm.io/gorm"
)

type QuotaClientService interface {
	GetQuotas(page, limit int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	GetQuotaByID(id int64) (*model.QuotaClient, error)
	CreateQuota(quota *model.QuotaClient) error
	UpdateQuota(quota *model.QuotaClient) error
	DeleteQuota(id int64) error
	// UseQuota(dto.UseQuotaClientRequest) (*model.QuotaClient, error)
	AddQuotaWithApprove(tx *gorm.DB, req dto.AddQuotaClientWithApproveRequest) (*model.QuotaClient, error)
}

type quotaClientService struct {
	db                      *gorm.DB
	quotaClientRepo         repository.QuotaClientRepository
	quotaClientReduction    repository.QuotaClientReductionRepository
	quotaClientAdditionRepo repository.QuotaClientAdditionRepository
	quotaUtils              *utils.QuotaUtils
}

func NewQuotaClientService(db *gorm.DB, quotaClientRepo repository.QuotaClientRepository, quotaClientReduction repository.QuotaClientReductionRepository, quotaClientAdditionRepo repository.QuotaClientAdditionRepository) QuotaClientService {
	return &quotaClientService{
		db:                      db,
		quotaClientRepo:         quotaClientRepo,
		quotaClientReduction:    quotaClientReduction,
		quotaClientAdditionRepo: quotaClientAdditionRepo,
		quotaUtils:              utils.NewQuotaUtils(),
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

func (s *quotaClientService) AddQuotaWithApprove(tx *gorm.DB, req dto.AddQuotaClientWithApproveRequest) (*model.QuotaClient, error) {
	return s.quotaUtils.CreateOrUpdateQuota(tx, req)
}
