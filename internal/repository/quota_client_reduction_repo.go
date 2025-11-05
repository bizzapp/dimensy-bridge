package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type QuotaClientReductionRepository interface {
	Create(reduction *model.QuotaClientReduction) error
	FindByID(id int64) (*model.QuotaClientReduction, error)
	FindAll() ([]model.QuotaClientReduction, error)
	FindByQuotaClientID(quotaClientID int64) ([]model.QuotaClientReduction, error)
	Update(reduction *model.QuotaClientReduction) error
	Delete(id int64) error
}

type quotaClientReductionRepository struct {
	db *gorm.DB
}

func NewQuotaClientReductionRepository(db *gorm.DB) QuotaClientReductionRepository {
	return &quotaClientReductionRepository{db: db}
}

func (r *quotaClientReductionRepository) Create(reduction *model.QuotaClientReduction) error {
	return r.db.Create(reduction).Error
}

func (r *quotaClientReductionRepository) FindByID(id int64) (*model.QuotaClientReduction, error) {
	var reduction model.QuotaClientReduction
	if err := r.db.First(&reduction, id).Error; err != nil {
		return nil, err
	}
	return &reduction, nil
}

func (r *quotaClientReductionRepository) FindAll() ([]model.QuotaClientReduction, error) {
	var reductions []model.QuotaClientReduction
	if err := r.db.Find(&reductions).Error; err != nil {
		return nil, err
	}
	return reductions, nil
}

func (r *quotaClientReductionRepository) FindByQuotaClientID(quotaClientID int64) ([]model.QuotaClientReduction, error) {
	var reductions []model.QuotaClientReduction
	if err := r.db.Where("quota_client_id = ?", quotaClientID).Find(&reductions).Error; err != nil {
		return nil, err
	}
	return reductions, nil
}

func (r *quotaClientReductionRepository) Update(reduction *model.QuotaClientReduction) error {
	return r.db.Save(reduction).Error
}

func (r *quotaClientReductionRepository) Delete(id int64) error {
	return r.db.Delete(&model.QuotaClientReduction{}, id).Error
}
