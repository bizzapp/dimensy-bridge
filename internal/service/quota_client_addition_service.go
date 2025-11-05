package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type QuotaClientAdditionService interface {
	GetAdditions(page, limit int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error)
	GetAdditionByID(id int64) (*model.QuotaClientAddition, error)
	CreateAddition(addition *model.QuotaClientAddition) error
	UpdateAddition(addition *model.QuotaClientAddition) error
	DeleteAddition(id int64) error
	ApproveAddition(additionID int64, approvedBy *int64) error
}

type quotaClientAdditionService struct {
	db           *gorm.DB
	additionRepo repository.QuotaClientAdditionRepository
	quotaRepo    repository.QuotaClientRepository
}

func NewQuotaClientAdditionService(db *gorm.DB, additionRepo repository.QuotaClientAdditionRepository, quotaRepo repository.QuotaClientRepository) QuotaClientAdditionService {
	return &quotaClientAdditionService{db, additionRepo, quotaRepo}
}

func (s *quotaClientAdditionService) GetAdditions(page, limit int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error) {
	offset := (page - 1) * limit
	return s.additionRepo.FindAll(limit, offset, filters)
}

func (s *quotaClientAdditionService) GetAdditionByID(id int64) (*model.QuotaClientAddition, error) {
	return s.additionRepo.FindByID(id)
}

func (s *quotaClientAdditionService) CreateAddition(addition *model.QuotaClientAddition) error {
	if addition == nil {
		return fmt.Errorf("addition data cannot be nil")
	}

	if addition.QuotaClientID == 0 {
		return fmt.Errorf("quota_client_id is required")
	}

	if addition.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	// Default value jika belum diisi
	// if addition.Status == "" {
	// addition.Status = "PENDING"
	// }
	if addition.Type == "" {
		addition.Type = "manual"
	}

	if err := s.db.Create(addition).Error; err != nil {
		return fmt.Errorf("failed to create quota client addition: %w", err)
	}
	return nil
}

func (s *quotaClientAdditionService) UpdateAddition(addition *model.QuotaClientAddition) error {
	return s.additionRepo.Update(addition)
}

func (s *quotaClientAdditionService) DeleteAddition(id int64) error {
	return s.additionRepo.Delete(id)
}

// 🔹 Step 2: Approve addition (apply quota & reduce stock)
func (s *quotaClientAdditionService) ApproveAddition(additionID int64, approvedBy *int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var addition model.QuotaClientAddition
		if err := tx.Preload("QuotaClient").First(&addition, additionID).Error; err != nil {
			return fmt.Errorf("addition not found: %w", err)
		}

		if addition.IsProcess {
			return errors.New("addition already approved")
		}

		var qc = addition.QuotaClient

		// Ambil MasterProduct
		var mp model.MasterProduct
		if err := tx.First(&mp, qc.MasterProductID).Error; err != nil {
			return fmt.Errorf("master product not found: %w", err)
		}

		// Cek stok
		if !mp.IsUnlimited && mp.CurrentStock < addition.Quantity {
			return fmt.Errorf("not enough stock in master product")
		}

		// Tambah quota client
		newQuota := qc.CurrentQuota + addition.Quantity
		if err := tx.Model(&qc).Updates(map[string]interface{}{
			"quantity":      gorm.Expr("quantity + ?", addition.Quantity),
			"current_quota": newQuota,
		}).Error; err != nil {
			return fmt.Errorf("failed to update quota client: %w", err)
		}

		// Kurangi stok master product
		if !mp.IsUnlimited {
			if err := tx.Model(&mp).Update("current_stock", gorm.Expr("current_stock - ?", addition.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to update master product stock: %w", err)
			}
		}

		// Update addition record
		if err := tx.Model(&addition).Updates(map[string]interface{}{
			"latest_quota": newQuota,
			"status":       "APPROVED",
			"is_process":   true,
			"process_by":   approvedBy,
			"updated_at":   time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("failed to update addition: %w", err)
		}

		return nil
	})
}
