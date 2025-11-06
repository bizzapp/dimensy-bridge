package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientDocumentProcessRepository interface {
	Create(data *model.ClientDocumentProcess) error
	FindByExternalID(externalID string) (*model.ClientDocumentProcess, error)
	UpdateStatus(externalID, status string) error
	DeleteByExternalID(externalID string) error
}

type clientDocumentProcessRepository struct {
	db *gorm.DB
}

func NewClientDocumentProcessRepository(db *gorm.DB) ClientDocumentProcessRepository {
	return &clientDocumentProcessRepository{db: db}
}

func (r *clientDocumentProcessRepository) Create(data *model.ClientDocumentProcess) error {
	return r.db.Create(data).Error
}

func (r *clientDocumentProcessRepository) FindByExternalID(externalID string) (*model.ClientDocumentProcess, error) {
	var process model.ClientDocumentProcess
	if err := r.db.Where("external_id = ?", externalID).First(&process).Error; err != nil {
		return nil, err
	}
	return &process, nil
}

func (r *clientDocumentProcessRepository) UpdateStatus(externalID, status string) error {
	return r.db.Model(&model.ClientDocumentProcess{}).
		Where("external_id = ?", externalID).
		Update("status", status).Error
}

func (r *clientDocumentProcessRepository) DeleteByExternalID(externalID string) error {
	return r.db.Where("external_id = ?", externalID).Delete(&model.ClientDocumentProcess{}).Error
}
