package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuotaClientAdditionService interface {
	GetAdditions(page, limit int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error)
	GetAdditionByID(id int64) (*model.QuotaClientAddition, error)
	CreateAddition(addition *model.QuotaClientAddition) error
	UpdateAddition(addition *model.QuotaClientAddition) error
	DeleteAddition(id int64) error
	ApproveAddQuota(dto.ApproveAddQuotaClientRequest) error
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
func (s *quotaClientAdditionService) ApproveAddQuota(req dto.ApproveAddQuotaClientRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		// ============================================
		// LOCK addition row (mencegah double approve)
		// ============================================
		var addition model.QuotaClientAddition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("QuotaClient").
			First(&addition, req.QuotaAdditionID).Error; err != nil {
			return fmt.Errorf("addition not found: %w", err)
		}

		if addition.IsProcess {
			return errors.New("addition already approved")
		}

		qc := addition.QuotaClient

		// ============================================
		// LOCK master product row
		// ============================================
		var mp model.MasterProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&mp, qc.MasterProductID).Error; err != nil {
			return fmt.Errorf("master product not found: %w", err)
		}

		// Cek stok latest (sudah locked)
		if !mp.IsUnlimited && mp.CurrentStock < addition.Quantity {
			return fmt.Errorf("not enough stock in master product")
		}

		// ============================================
		// UPDATE QuotaClient secara atomik
		// ============================================
		if err := tx.Model(&qc).
			Where("id = ?", qc.ID).
			Updates(map[string]interface{}{
				"quantity":      gorm.Expr("quantity + ?", addition.Quantity),
				"current_quota": gorm.Expr("current_quota + ?", addition.Quantity),
			}).Error; err != nil {
			return fmt.Errorf("failed to update quota client: %w", err)
		}

		// ============================================
		// Kurangi stok master product (atomic)
		// ============================================
		var latestStock int64

		if !mp.IsUnlimited {
			if err := tx.Model(&mp).
				Where("id = ?", mp.ID).
				Update("current_stock", gorm.Expr("current_stock - ?", addition.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to update master product stock: %w", err)
			}

			// Baca ulang stok terbaru setelah update (agar akurat)
			if err := tx.Model(&mp).Select("current_stock").First(&mp).Error; err != nil {
				return fmt.Errorf("failed to reload master product stock: %w", err)
			}

			latestStock = mp.CurrentStock

			// Simpan record reduction
			reduction := &model.MasterProductReduction{
				MasterProductID: mp.ID,
				Quantity:        addition.Quantity,
				LatestQuota:     latestStock,
				Type:            "ADDITION",
				UsedBy:          req.ProcessBy,
			}

			if err := tx.Create(reduction).Error; err != nil {
				return fmt.Errorf("failed to create reduction record: %w", err)
			}
		}

		// ============================================
		// Tandai addition sudah diproses
		// ============================================
		if err := tx.Model(&addition).Updates(map[string]interface{}{
			"latest_quota": gorm.Expr("latest_quota + ?", addition.Quantity),
			"is_process":   true,
			"process_by":   req.ProcessBy,
			"updated_at":   time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("failed to update addition: %w", err)
		}

		return nil
	})
}
