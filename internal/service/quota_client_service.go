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
	AddQuota(tx *gorm.DB, req dto.AddQuotaClientRequest) (*model.QuotaClient, error)
}

type quotaClientService struct {
	db                      *gorm.DB
	quotaClientRepo         repository.QuotaClientRepository
	quotaClientReduction    repository.QuotaClientReductionRepository
	quotaClientAdditionRepo repository.QuotaClientAdditionRepository
}

func NewQuotaClientService(db *gorm.DB, quotaClientRepo repository.QuotaClientRepository, quotaClientReduction repository.QuotaClientReductionRepository, quotaClientAdditionRepo repository.QuotaClientAdditionRepository) QuotaClientService {
	return &quotaClientService{
		db:                      db,
		quotaClientRepo:         quotaClientRepo,
		quotaClientReduction:    quotaClientReduction,
		quotaClientAdditionRepo: quotaClientAdditionRepo,
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

// CreateOrUpdateQuota creates new quota or updates existing quota based on clientID and masterProductID
func (s *quotaClientService) AddQuota(tx *gorm.DB, req dto.AddQuotaClientRequest) (*model.QuotaClient, error) {
	var quota model.QuotaClient

	// Try to find existing quota
	err := tx.Where("client_id = ? AND master_product_id = ?", req.ClientID, req.MasterProductID).First(&quota).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new quota if not exists
			quota = model.QuotaClient{
				MasterProductID:       req.MasterProductID,
				Quantity:              req.Quantity,
				CurrentQuota:          req.Quantity,
				ClientID:              req.ClientID,
				IsUnlimited:           req.IsUnlimited,
				MaxSingleUpload:       req.MaxSingleUpload,
				MaxBulkUploadLimitPcs: req.MaxBulkUploadLimitPcs,
				MaxBulkUploadLimitAll: req.MaxBulkUploadLimitAll,
				MaxBulkUploadCount:    req.MaxBulkUploadCount,
			}

			if err := tx.Create(&quota).Error; err != nil {
				return nil, fmt.Errorf("failed to create quota: %w", err)
			}

			if req.Quantity > 0 {
				addition := model.QuotaClientAddition{
					QuotaClientID: quota.ID,
					Quantity:      req.Quantity,
					LatestQuota:   quota.CurrentQuota,
					Type:          "initial",
					CreatedBy:     req.CreatedBy,
					IsProcess:     true,
				}

				if err := tx.Create(&addition).Error; err != nil {
					return nil, fmt.Errorf("failed to create initial quota addition record: %w", err)
				}
			}
		} else {
			return nil, fmt.Errorf("failed to query quota: %w", err)
		}
	} else {
		// Update existing quota - add to current values
		quota.Quantity += req.Quantity
		quota.CurrentQuota += req.Quantity
		quota.IsUnlimited = req.IsUnlimited // Update unlimited status
		quota.MaxSingleUpload = req.MaxSingleUpload
		quota.MaxBulkUploadLimitPcs = req.MaxBulkUploadLimitPcs
		quota.MaxBulkUploadLimitAll = req.MaxBulkUploadLimitAll
		quota.MaxBulkUploadCount = req.MaxBulkUploadCount

		if err := tx.Save(&quota).Error; err != nil {
			return nil, fmt.Errorf("failed to update quota: %w", err)
		}

		// Create quota addition record for the update
		if req.Quantity > 0 {
			addition := model.QuotaClientAddition{
				QuotaClientID: quota.ID,
				Quantity:      req.Quantity,
				LatestQuota:   quota.CurrentQuota,
				Type:          "addition",
				CreatedBy:     req.CreatedBy,
				IsProcess:     true,
			}

			if err := tx.Create(&addition).Error; err != nil {
				return nil, fmt.Errorf("failed to create quota addition record: %w", err)
			}
		}
	}

	return &quota, nil
}
