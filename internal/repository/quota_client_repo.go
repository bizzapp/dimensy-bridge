package repository

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type QuotaClientRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	FindByID(id int64) (*model.QuotaClient, error)
	Create(quota *model.QuotaClient) error
	Update(quota *model.QuotaClient) error
	Delete(id int64) error
	FindByClientProduct(req dto.FindQuotaClientByClientProductRequest) (*model.QuotaClient, error)

	GetAdditionHistory(limit, offset int, filters map[string]interface{}) ([]dto.QuotaHistoryItem, error)
	GetReductionHistory(limit, offset int, filters map[string]interface{}) ([]dto.QuotaHistoryItem, error)
	CountAdditionHistory(filters map[string]interface{}) (int64, error)
	CountReductionHistory(filters map[string]interface{}) (int64, error)
}

type quotaClientRepository struct {
	db *gorm.DB
}

func NewQuotaClientRepository(db *gorm.DB) QuotaClientRepository {
	return &quotaClientRepository{db}
}

func (r *quotaClientRepository) GetAdditionHistory(limit, offset int, filters map[string]interface{}) ([]dto.QuotaHistoryItem, error) {
	var additions []dto.QuotaHistoryItem

	query := r.db.Model(&model.QuotaClientAddition{}).
		Select(`
        quota_client_additions.id,
        quota_clients.master_product_id,
        master_products.name AS master_product_name,
        quota_client_additions.type,
        quota_client_additions.quantity,
        quota_client_additions.created_at,
		clients.company_name as client_name
    `).
		Joins("JOIN quota_clients ON quota_clients.id = quota_client_additions.quota_client_id").
		Joins("JOIN master_products ON master_products.id = quota_clients.master_product_id").
		Order("quota_client_additions.created_at DESC").
		Joins("JOIN clients ON quota_clients.client_id = clients.id").
		Limit(limit).
		Offset(offset)

	if clientID, ok := filters["client_id"]; ok {
		query = query.Where("quota_clients.client_id = ?", clientID)
	}

	if err := query.Find(&additions).Error; err != nil {
		return nil, err
	}

	return additions, nil
}

func (r *quotaClientRepository) GetReductionHistory(limit, offset int, filters map[string]interface{}) ([]dto.QuotaHistoryItem, error) {
	var reductions []dto.QuotaHistoryItem

	query := r.db.Model(&model.QuotaClientReduction{}).
		Select(`
			quota_client_reductions.id,
			quota_clients.master_product_id,
			master_products.name AS master_product_name,
			quota_client_reductions.type,
			quota_client_reductions.quantity,
			quota_client_reductions.created_at,
			clients.company_name as client_name
		`).
		Joins("JOIN quota_clients ON quota_clients.id = quota_client_reductions.quota_client_id").
		Joins("JOIN master_products ON master_products.id = quota_clients.master_product_id").
		Joins("JOIN clients ON clients.id = quota_clients.client_id").
		Order("quota_client_reductions.created_at DESC").
		Limit(limit).
		Offset(offset)

	if clientID, ok := filters["client_id"]; ok {
		query = query.Where("quota_clients.client_id = ?", clientID)
	}

	if err := query.Find(&reductions).Error; err != nil {
		return nil, err
	}

	return reductions, nil
}

func (r *quotaClientRepository) CountAdditionHistory(filters map[string]interface{}) (int64, error) {
	var count int64

	query := r.db.Model(&model.QuotaClientAddition{}).
		Joins("JOIN quota_clients ON quota_clients.id = quota_client_additions.quota_client_id")

	if clientID, ok := filters["client_id"]; ok {
		query = query.Where("quota_clients.client_id = ?", clientID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *quotaClientRepository) CountReductionHistory(filters map[string]interface{}) (int64, error) {
	var count int64

	query := r.db.Model(&model.QuotaClientReduction{}).
		Joins("JOIN quota_clients ON quota_clients.id = quota_client_reductions.quota_client_id")

	if clientID, ok := filters["client_id"]; ok {
		query = query.Where("quota_clients.client_id = ?", clientID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
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

func (r *quotaClientRepository) FindByClientProduct(req dto.FindQuotaClientByClientProductRequest) (*model.QuotaClient, error) {
	var quota model.QuotaClient
	if err := r.db.Preload("MasterProduct").Preload("Client").
		Where("client_id = ? AND master_product_id = ?", req.ClientID, req.MasterProductID).
		First(&quota).Error; err != nil {
		return nil, err
	}
	return &quota, nil
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
