package repository

import (
	"dimensy-bridge/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientDocumentProcessRepository interface {
	Create(data *model.ClientDocumentProcess) error
	FindByExternalID(externalID uuid.UUID) (*model.ClientDocumentProcess, error)
	FindByExternalIDAndExternalUserID(externalID uuid.UUID, userID *uuid.UUID) (*model.ClientDocumentProcess, error)
	FindByExternalIDExternalUserIDExternalCompanyID(externalID *uuid.UUID, groupID *uuid.UUID, userID *uuid.UUID, companyID *uuid.UUID) (*model.ClientDocumentProcess, error)
	UpdateStatus(externalID uuid.UUID, status string) error
	DeleteByExternalID(externalID uuid.UUID) error
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
func (r *clientDocumentProcessRepository) FindByExternalIDExternalUserIDExternalCompanyID(externalID *uuid.UUID, groupID *uuid.UUID, userID *uuid.UUID, companyID *uuid.UUID) (*model.ClientDocumentProcess, error) {
	var process model.ClientDocumentProcess
	query := r.db
	if externalID != nil {
		query = query.Where("external_id = ?", *externalID)
	}

	// Handle groupID: if nil, check IS NULL; otherwise check equality
	if groupID != nil {
		query = query.Where("external_group_id = ?", *groupID)
	}

	// Handle userID: if nil, check IS NULL; otherwise check equality
	if userID != nil {
		query = query.Where("external_user_id = ?", *userID)
	}

	// Handle companyID: if nil, check IS NULL; otherwise check equality
	if companyID != nil {
		query = query.Where("external_company_id = ?", *companyID)
	}

	if err := query.First(&process).Error; err != nil {
		return nil, err
	}
	return &process, nil
}

func (r *clientDocumentProcessRepository) FindByExternalIDAndExternalUserID(externalID uuid.UUID, userID *uuid.UUID) (*model.ClientDocumentProcess, error) {
	var process model.ClientDocumentProcess
	if err := r.db.Where("external_id = ? AND external_user_id = ?", externalID, userID).First(&process).Error; err != nil {
		return nil, err
	}
	return &process, nil
}
func (r *clientDocumentProcessRepository) FindByExternalID(externalID uuid.UUID) (*model.ClientDocumentProcess, error) {
	var process model.ClientDocumentProcess
	if err := r.db.Where("external_id = ?", externalID).First(&process).Error; err != nil {
		return nil, err
	}
	return &process, nil
}

func (r *clientDocumentProcessRepository) UpdateStatus(externalID uuid.UUID, status string) error {
	return r.db.Model(&model.ClientDocumentProcess{}).
		Where("external_id = ?", externalID).
		Update("status", status).Error
}

func (r *clientDocumentProcessRepository) DeleteByExternalID(externalID uuid.UUID) error {
	return r.db.Where("external_id = ?", externalID).Delete(&model.ClientDocumentProcess{}).Error
}
