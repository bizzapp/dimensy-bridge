package repository

import (
"dimensy-bridge/internal/model"

"gorm.io/gorm"
)

type ClientIPWhitelistRepository interface {
	Create(ipWhitelist *model.ClientIPWhitelist) (*model.ClientIPWhitelist, error)
	GetByID(id int64) (*model.ClientIPWhitelist, error)
	GetByClientID(clientID int64, page, limit int) ([]model.ClientIPWhitelist, int64, error)
	GetActiveIPsByClientID(clientID int64) ([]string, error)
	Update(id int64, ipWhitelist *model.ClientIPWhitelist) (*model.ClientIPWhitelist, error)
	Delete(id int64) error
	IsIPWhitelisted(clientID int64, ipAddress string) (bool, error)
}

type clientIPWhitelistRepository struct {
	db *gorm.DB
}

func NewClientIPWhitelistRepository(db *gorm.DB) ClientIPWhitelistRepository {
	return &clientIPWhitelistRepository{db: db}
}

func (r *clientIPWhitelistRepository) Create(ipWhitelist *model.ClientIPWhitelist) (*model.ClientIPWhitelist, error) {
	if err := r.db.Create(ipWhitelist).Error; err != nil {
		return nil, err
	}
	return ipWhitelist, nil
}

func (r *clientIPWhitelistRepository) GetByID(id int64) (*model.ClientIPWhitelist, error) {
	var ipWhitelist model.ClientIPWhitelist
	if err := r.db.Where("id = ?", id).First(&ipWhitelist).Error; err != nil {
		return nil, err
	}
	return &ipWhitelist, nil
}

func (r *clientIPWhitelistRepository) GetByClientID(clientID int64, page, limit int) ([]model.ClientIPWhitelist, int64, error) {
	var ipWhitelists []model.ClientIPWhitelist
	var total int64

	query := r.db.Where("client_id = ?", clientID)

	if err := query.Model(&model.ClientIPWhitelist{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&ipWhitelists).Error; err != nil {
		return nil, 0, err
	}

	return ipWhitelists, total, nil
}

func (r *clientIPWhitelistRepository) GetActiveIPsByClientID(clientID int64) ([]string, error) {
	var ips []string
	if err := r.db.Model(&model.ClientIPWhitelist{}).
		Where("client_id = ? AND is_active = ?", clientID, true).
		Pluck("ip_address", &ips).Error; err != nil {
		return nil, err
	}
	return ips, nil
}

func (r *clientIPWhitelistRepository) Update(id int64, ipWhitelist *model.ClientIPWhitelist) (*model.ClientIPWhitelist, error) {
	if err := r.db.Model(&model.ClientIPWhitelist{}).Where("id = ?", id).Updates(ipWhitelist).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *clientIPWhitelistRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.ClientIPWhitelist{}).Error
}

func (r *clientIPWhitelistRepository) IsIPWhitelisted(clientID int64, ipAddress string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.ClientIPWhitelist{}).
		Where("client_id = ? AND ip_address = ? AND is_active = ?", clientID, ipAddress, true).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
