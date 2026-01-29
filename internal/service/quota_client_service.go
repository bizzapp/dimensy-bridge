package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"sort"

	"gorm.io/gorm"
)

type QuotaClientService interface {
	GetQuotas(page, limit int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	GetQuotaByID(id int64) (*model.QuotaClient, error)
	CreateQuota(quota *model.QuotaClient) error
	UpdateQuota(quota *model.QuotaClient) error
	DeleteQuota(id int64) error
	GetHistory(page, limit int, filters map[string]interface{}) ([]dto.QuotaHistoryItem, int64, error)
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

func (s *quotaClientService) GetHistory(page, limit int, filters map[string]interface{}) ([]dto.QuotaHistoryItem, int64, error) {
	// Set default limit
	if limit <= 0 {
		limit = 20
	}

	// Set default page
	if page <= 0 {
		page = 1
	}

	// Calculate offset
	offset := (page - 1) * limit

	// ============================
	// 📌 Fetch Addition History
	// ============================
	additions, err := s.quotaClientRepo.GetAdditionHistory(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	// Set direction for additions
	for i := range additions {
		additions[i].Direction = "ADDITION"
	}

	// ============================
	// 📌 Fetch Reduction History
	// ============================
	reductions, err := s.quotaClientRepo.GetReductionHistory(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	// Set direction and negate quantity for reductions
	for i := range reductions {
		reductions[i].Direction = "REDUCTION"
		reductions[i].Quantity = reductions[i].Quantity * -1
	}

	// ============================
	// 📌 Get Total Count
	// ============================
	additionCount, err := s.quotaClientRepo.CountAdditionHistory(filters)
	if err != nil {
		return nil, 0, err
	}

	reductionCount, err := s.quotaClientRepo.CountReductionHistory(filters)
	if err != nil {
		return nil, 0, err
	}

	totalCount := additionCount + reductionCount

	// ============================
	// 📌 Merge & Sort DESC
	// ============================
	all := append(additions, reductions...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(*all[j].CreatedAt)
	})

	// Apply global limit
	if len(all) > limit {
		all = all[:limit]
	}

	return all, totalCount, nil
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
