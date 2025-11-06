package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientDocumentProcessDetailRepository interface {
	Create(data *model.ClientDocumentProcessDetail) error
	FindByClientID(clientID int64) ([]model.ClientDocumentProcessDetail, error)
	DeleteByClientID(clientID int64) error
}

type clientDocumentProcessDetailRepository struct {
	db *gorm.DB
}

func NewClientDocumentProcessDetailRepository(db *gorm.DB) ClientDocumentProcessDetailRepository {
	return &clientDocumentProcessDetailRepository{db: db}
}

func (r *clientDocumentProcessDetailRepository) Create(data *model.ClientDocumentProcessDetail) error {
	return r.db.Create(data).Error
}

func (r *clientDocumentProcessDetailRepository) FindByClientID(clientID int64) ([]model.ClientDocumentProcessDetail, error) {
	var details []model.ClientDocumentProcessDetail
	if err := r.db.Where("client_id = ?", clientID).Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *clientDocumentProcessDetailRepository) DeleteByClientID(clientID int64) error {
	return r.db.Where("client_id = ?", clientID).Delete(&model.ClientDocumentProcessDetail{}).Error
}
