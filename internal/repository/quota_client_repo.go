package repository

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"sort"

	"gorm.io/gorm"
)

type QuotaClientRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.QuotaClient, int64, error)
	FindByID(id int64) (*model.QuotaClient, error)
	Create(quota *model.QuotaClient) error
	Update(quota *model.QuotaClient) error
	Delete(id int64) error
	FindByClientProduct(req dto.FindQuotaClientByClientProductRequest) (*model.QuotaClient, error)

	GetHistory(clientID int64, limit int) ([]dto.QuotaHistoryItem, error)
}

type quotaClientRepository struct {
	db *gorm.DB
}

func NewQuotaClientRepository(db *gorm.DB) QuotaClientRepository {
	return &quotaClientRepository{db}
}

func (r *quotaClientRepository) GetHistory(clientID int64, limit int) ([]dto.QuotaHistoryItem, error) {
	var additions []dto.QuotaHistoryItem
	var reductions []dto.QuotaHistoryItem

	// ============================
	// 📌 Query Addition History
	// ============================
	if err := r.db.Model(&model.QuotaClientAddition{}).
		Select(`
			quota_client_additions.id,
			quota_clients.master_product_id,
			master_products.name AS master_product_name,
			quota_client_additions.type,
			quota_client_additions.quantity,
			quota_client_additions.created_at
		`).
		Joins("JOIN quota_clients ON quota_clients.id = quota_client_additions.quota_client_id").
		Joins("JOIN master_products ON master_products.id = quota_clients.master_product_id").
		Where("quota_clients.client_id = ?", clientID).
		Limit(limit).
		Find(&additions).Error; err != nil {
		return nil, err
	}

	// Tambahkan jenis direction
	for i := range additions {
		additions[i].Direction = "ADDITION"
	}

	// ============================
	// 📌 Query Reduction History
	// ============================
	if err := r.db.Model(&model.QuotaClientReduction{}).
		Select(`
			quota_client_reductions.id,
			quota_clients.master_product_id,
			master_products.name AS master_product_name,
			quota_client_reductions.type,

			(quota_client_reductions.quantity * -1) AS quantity,
			quota_client_reductions.created_at
		`).
		Joins("JOIN quota_clients ON quota_clients.id = quota_client_reductions.quota_client_id").
		Joins("JOIN master_products ON master_products.id = quota_clients.master_product_id").
		Where("quota_clients.client_id = ?", clientID).
		Limit(limit).
		Find(&reductions).Error; err != nil {
		return nil, err
	}

	for i := range reductions {
		reductions[i].Direction = "REDUCTION"
	}

	// ============================
	// 📌 Gabungkan & Sort DESC
	// ============================

	all := append(additions, reductions...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(*all[j].CreatedAt)
	})

	// Limit global
	if len(all) > limit {
		all = all[:limit]
	}

	return all, nil
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
