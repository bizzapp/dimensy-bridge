package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type QuotaClientRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	FindByID(id int64) (*model.QuotaClient, error)
	Create(quota *model.QuotaClient) error
	Update(quota *model.QuotaClient) error
	Delete(id int64) error
}

type quotaClientRepository struct {
	db *gorm.DB
}

func NewQuotaClientRepository(db *gorm.DB) QuotaClientRepository {
	return &quotaClientRepository{db}
}

func (r *quotaClientRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]model.QuotaClient, int64, error) {
	var quotas []model.QuotaClient
	var total int64

	query := r.db.Model(&model.QuotaClient{}).Preload("MasterProduct").Preload("Client")

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	return quotas, total, nil
}

func (r *quotaClientRepository) FindByID(id int64) (*model.QuotaClient, error) {
	var quota model.QuotaClient
	if err := r.db.Preload("MasterProduct").Preload("Client").First(&quota, id).Error; err != nil {
		return nil, err
	}
	return &quota, nil
}

func (r *quotaClientRepository) Create(quota *model.QuotaClient) error {
	return r.db.Create(quota).Error
}

func (r *quotaClientRepository) Update(quota *model.QuotaClient) error {
	return r.db.Save(quota).Error
}

func (r *quotaClientRepository) Delete(id int64) error {
	return r.db.Delete(&model.QuotaClient{}, id).Error
}
