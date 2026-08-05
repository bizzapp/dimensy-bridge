package utils

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/model/seeder"
	"fmt"

	"gorm.io/gorm"
)

// QuotaUtils provides utility functions for quota management
type QuotaUtils struct{}

// NewQuotaUtils creates a new instance of QuotaUtils
func NewQuotaUtils() *QuotaUtils {
	return &QuotaUtils{}
}

// CreateOrUpdateQuota creates new quota or updates existing quota based on clientID and masterProductID
// This is a flexible utility that can be used by any service
func (qu *QuotaUtils) CreateOrUpdateQuota(tx *gorm.DB, req dto.AddQuotaClientWithApproveRequest) (*model.QuotaClient, error) {
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

			// Create initial quota addition record if quantity > 0
			if req.Quantity > 0 {
				addition := model.QuotaClientAddition{
					QuotaClientID: quota.ID,
					Quantity:      req.Quantity,
					LatestQuota:   quota.CurrentQuota,
					Type:          "initial",
					CreatedBy:     req.CreatedBy,
					ProcessBy:     req.ProcessBy,
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

		// Create quota addition record for the update if quantity > 0
		if req.Quantity > 0 {
			addition := model.QuotaClientAddition{
				QuotaClientID: quota.ID,
				Quantity:      req.Quantity,
				LatestQuota:   quota.CurrentQuota,
				Type:          "addition",
				CreatedBy:     req.CreatedBy,
				ProcessBy:     req.ProcessBy,
				IsProcess:     true,
			}

			if err := tx.Create(&addition).Error; err != nil {
				return nil, fmt.Errorf("failed to create quota addition record: %w", err)
			}
		}
	}

	return &quota, nil
}

// UseQuota reduces quota and creates reduction record
func (qu *QuotaUtils) UseQuota(tx *gorm.DB, req dto.UseQuotaClientRequest) (*model.QuotaClient, error) {
	var quota model.QuotaClient

	if err := tx.Where("client_id = ? AND master_product_id = ?", req.ClientID, req.MasterProductID).First(&quota).Error; err != nil {
		return nil, fmt.Errorf("quota not found: %w", err)
	}

	if quota.CurrentQuota < req.Quantity {
		return nil, fmt.Errorf("insufficient quota")
	}

	quota.CurrentQuota -= req.Quantity
	if err := tx.Save(&quota).Error; err != nil {
		return nil, fmt.Errorf("failed to update quota: %w", err)
	}

	typeReduction := "usage"
	if req.TypeReduction != nil {
		typeReduction = *req.TypeReduction
	}

	// Create a quota reduction record
	reduction := model.QuotaClientReduction{
		QuotaClientID: quota.ID,
		Quantity:      req.Quantity,
		LatestQuota:   quota.CurrentQuota,
		Type:          typeReduction,
		UsedBy:        req.UsedBy,
	}

	if err := tx.Create(&reduction).Error; err != nil {
		return nil, fmt.Errorf("failed to create quota reduction record: %w", err)
	}

	return &quota, nil
}

// UseQuotaEkyc reduces E-KYC quota specifically
func (qu *QuotaUtils) UseQuotaEkyc(tx *gorm.DB, req dto.UseQuotaClientRequest) (*model.QuotaClient, error) {
	req.MasterProductID = seeder.ID_PRODUCT_KYC
	return qu.UseQuota(tx, req)
}

func (qu *QuotaUtils) QuotaLimit(tx *gorm.DB, clientID int64, masterProductID int64) (*model.QuotaClient, error) {
	var quota model.QuotaClient

	if err := tx.Where("client_id = ? AND master_product_id = ?", clientID, masterProductID).First(&quota).Error; err != nil {
		return nil, fmt.Errorf("quota not found: %w", err)
	}
	return &quota, nil
}

// CreateInitialQuota creates initial quota from subscription plan
func (qu *QuotaUtils) CreateInitialQuota(tx *gorm.DB, clientID int64, masterProductID int64, quantity int64, isUnlimited bool, createdBy int64) (*model.QuotaClient, error) {
	// Create initial quota
	quota := model.QuotaClient{
		MasterProductID: masterProductID,
		Quantity:        quantity,
		CurrentQuota:    quantity,
		ClientID:        clientID,
		IsUnlimited:     isUnlimited,
	}

	if err := tx.Create(&quota).Error; err != nil {
		return nil, fmt.Errorf("failed to create initial quota: %w", err)
	}

	// Create initial quota addition record if quantity > 0
	if quantity > 0 {
		addition := model.QuotaClientAddition{
			QuotaClientID: quota.ID,
			Quantity:      quantity,
			LatestQuota:   quota.CurrentQuota,
			Type:          "initial",
			CreatedBy:     createdBy,
			IsProcess:     true,
		}

		if err := tx.Create(&addition).Error; err != nil {
			return nil, fmt.Errorf("failed to create initial quota addition record: %w", err)
		}
	}

	return &quota, nil
}

// AddQuotaToExisting adds quota to existing client quota
func (qu *QuotaUtils) AddQuotaToExisting(tx *gorm.DB, quotaClientID int64, quantity int64, createdBy int64, processBy *int64) (*model.QuotaClient, error) {
	var quota model.QuotaClient

	if err := tx.First(&quota, quotaClientID).Error; err != nil {
		return nil, fmt.Errorf("quota not found: %w", err)
	}

	// Update current quota
	quota.CurrentQuota += quantity
	quota.Quantity += quantity

	if err := tx.Save(&quota).Error; err != nil {
		return nil, fmt.Errorf("failed to update quota: %w", err)
	}

	// Create quota addition record
	addition := model.QuotaClientAddition{
		QuotaClientID: quota.ID,
		Quantity:      quantity,
		LatestQuota:   quota.CurrentQuota,
		Type:          "addition",
		CreatedBy:     createdBy,
		ProcessBy:     processBy,
		IsProcess:     true,
	}

	if err := tx.Create(&addition).Error; err != nil {
		return nil, fmt.Errorf("failed to create quota addition record: %w", err)
	}

	return &quota, nil
}
