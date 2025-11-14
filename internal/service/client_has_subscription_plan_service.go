package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientHasSubscriptionPlanService interface {
	Create(data *dto.CreateSubscriptionPlanRequest) (*model.ClientHasSubscriptionPlan, error)
	Process(id int64, processedBy int64) error
	GetAll() ([]model.ClientHasSubscriptionPlan, error)
	GetByID(id int64) (*model.ClientHasSubscriptionPlan, error)
	Delete(id int64) error
}

type clientHasSubscriptionPlanService struct {
	repo             repository.ClientHasSubscriptionPlanRepository
	quotaClientRepo  repository.QuotaClientRepository
	subscriptionRepo repository.SubscriptionPlanRepository
	quotaClientSvc   QuotaClientService
	quotaUtils       *utils.QuotaUtils
}

func NewClientHasSubscriptionPlanService(
	repo repository.ClientHasSubscriptionPlanRepository,
	quotaRepo repository.QuotaClientRepository,
	subRepo repository.SubscriptionPlanRepository,
	quotaClientSvc QuotaClientService,
) ClientHasSubscriptionPlanService {
	return &clientHasSubscriptionPlanService{
		repo:             repo,
		quotaClientRepo:  quotaRepo,
		subscriptionRepo: subRepo,
		quotaClientSvc:   quotaClientSvc,
		quotaUtils:       utils.NewQuotaUtils(),
	}
}

func (s *clientHasSubscriptionPlanService) Create(data *dto.CreateSubscriptionPlanRequest) (*model.ClientHasSubscriptionPlan, error) {
	plan := &model.ClientHasSubscriptionPlan{
		ClientID:           data.ClientID,
		SubscriptionPlanID: data.SubscriptionPlanID,
		CreatedBy:          data.CreatedBy,
	}
	if err := s.repo.Create(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

// Process: menandai subscription sudah aktif dan membuat quota_clients sesuai detail plan
func (s *clientHasSubscriptionPlanService) Process(id int64, processedBy int64) error {
	return s.repo.WithTransaction(func(tx *gorm.DB) error {
		// 1️⃣ Ambil data subscription dengan preload details, gunakan lock (FOR UPDATE)
		sub, err := s.repo.FindByIDTx(tx, id)
		if err != nil {
			return err
		}

		now := time.Now()
		sub.ProcessBy = &processedBy
		sub.ProcessTime = &now
		sub.IsActive = true

		// 2️⃣ Proses setiap detail dari subscription plan
		for _, detail := range sub.SubscriptionPlan.Details {

			// Ambil master product dengan lock FOR UPDATE untuk cegah race condition
			var master model.MasterProduct
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", detail.MasterProductID).
				First(&master).Error; err != nil {
				return fmt.Errorf("master product not found: %w", err)
			}

			// Validasi stok cukup
			if !master.IsUnlimited {
				if master.CurrentStock < int64(detail.Quantity) {
					return fmt.Errorf("insufficient stock for product ID %d", master.ID)
				}

				// Simpan current stock sebelum dikurangi

				// Kurangi stok master product
				master.CurrentStock -= int64(detail.Quantity)
				if err := tx.Save(&master).Error; err != nil {
					return fmt.Errorf("failed to update stock for product %d: %w", master.ID, err)
				}

				// Simpan record reduction
				reduction := &model.MasterProductReduction{
					MasterProductID: master.ID,
					Quantity:        int64(detail.Quantity),
					LatestQuota:     master.CurrentStock,
					Type:            "SUBSCRIPTION",
					UsedBy:          sub.ProcessBy,
				}
				if err := tx.Create(reduction).Error; err != nil {
					return fmt.Errorf("failed to create reduction record: %w", err)
				}
			}

			// Create or update quota untuk client dengan addition record menggunakan utils
			req := dto.AddQuotaClientWithApproveRequest{
				ClientID:              sub.ClientID,
				MasterProductID:       master.ID,
				Quantity:              int64(detail.Quantity),
				IsUnlimited:           detail.IsUnlimited,
				CreatedBy:             sub.CreatedBy, // Use createdBy as createdBy for addition record
				ProcessBy:             sub.ProcessBy,
				MaxSingleUpload:       detail.MaxSingleUpload,
				MaxBulkUploadLimitPcs: detail.MaxBulkUploadLimitPcs,
				MaxBulkUploadLimitAll: detail.MaxBulkUploadLimitAll,
				MaxBulkUploadCount:    detail.MaxBulkUploadCount,
			}
			_, err = s.quotaUtils.CreateOrUpdateQuota(tx, req)
			if err != nil {
				return fmt.Errorf("failed to create or update quota: %w", err)
			}
		}

		// 3️⃣ Update status subscription
		if err := tx.Save(sub).Error; err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}

		return nil
	})
}

func (s *clientHasSubscriptionPlanService) GetAll() ([]model.ClientHasSubscriptionPlan, error) {
	return s.repo.FindAll()
}

func (s *clientHasSubscriptionPlanService) GetByID(id int64) (*model.ClientHasSubscriptionPlan, error) {
	return s.repo.FindByID(id)
}

func (s *clientHasSubscriptionPlanService) Delete(id int64) error {
	return s.repo.Delete(id)
}
