package repository

import (
	"dimensy-bridge/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientDocumentResendOtpRepository interface {
	Create(data *model.ClientDocumentResendOtp) error
	FindByID(id int64) (*model.ClientDocumentResendOtp, error)
	FindAll() ([]model.ClientDocumentResendOtp, error)
	Delete(id int64) error
	FindByExternalID(externalID uuid.UUID) (*model.ClientDocumentResendOtp, error)
}

type clientDocumentResendOtpRepository struct {
	db *gorm.DB
}

func NewClientDocumentResendOtpRepository(db *gorm.DB) ClientDocumentResendOtpRepository {
	return &clientDocumentResendOtpRepository{db: db}
}

func (r *clientDocumentResendOtpRepository) Create(data *model.ClientDocumentResendOtp) error {
	return r.db.Create(data).Error
}

func (r *clientDocumentResendOtpRepository) FindByID(id int64) (*model.ClientDocumentResendOtp, error) {
	var result model.ClientDocumentResendOtp
	if err := r.db.First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *clientDocumentResendOtpRepository) FindAll() ([]model.ClientDocumentResendOtp, error) {
	var results []model.ClientDocumentResendOtp
	if err := r.db.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *clientDocumentResendOtpRepository) Delete(id int64) error {
	return r.db.Delete(&model.ClientDocumentResendOtp{}, id).Error
}

func (r *clientDocumentResendOtpRepository) FindByExternalID(externalID uuid.UUID) (*model.ClientDocumentResendOtp, error) {
	var result model.ClientDocumentResendOtp
	if err := r.db.Where("external_id = ?", externalID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
