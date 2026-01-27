package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type QuotaClientAdditionRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error)
	FindByID(id int64) (*model.QuotaClientAddition, error)
	Create(addition *model.QuotaClientAddition) error
	Update(addition *model.QuotaClientAddition) error
	Delete(id int64) error
}

type quotaClientAdditionRepository struct {
	db *gorm.DB
}

func NewQuotaClientAdditionRepository(db *gorm.DB) QuotaClientAdditionRepository {
	return &quotaClientAdditionRepository{db}
}

func (r *quotaClientAdditionRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]model.QuotaClientAddition, int64, error) {
	var additions []model.QuotaClientAddition
	var total int64

	query := r.db.Model(&model.QuotaClientAddition{}).Preload("QuotaClient").Preload("QuotaClient.MasterProduct")

	// Handle client_id filter with JOIN
	if clientID, ok := filters["client_id"]; ok {
		query = query.Joins("JOIN quota_clients ON quota_clients.id = quota_client_additions.quota_client_id").
			Where("quota_clients.client_id = ?", clientID)
		delete(filters, "client_id")
	}

	// Apply other filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&additions).Error; err != nil {
		return nil, 0, err
	}

	return additions, total, nil
}

func (r *quotaClientAdditionRepository) FindByID(id int64) (*model.QuotaClientAddition, error) {
	var addition model.QuotaClientAddition
	if err := r.db.Preload("QuotaClient").First(&addition, id).Error; err != nil {
		return nil, err
	}
	return &addition, nil
}

func (r *quotaClientAdditionRepository) Create(addition *model.QuotaClientAddition) error {
	return r.db.Create(addition).Error
}

func (r *quotaClientAdditionRepository) Update(addition *model.QuotaClientAddition) error {
	return r.db.Save(addition).Error
}

func (r *quotaClientAdditionRepository) Delete(id int64) error {
	return r.db.Delete(&model.QuotaClientAddition{}, id).Error
}
