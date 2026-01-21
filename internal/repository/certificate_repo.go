package repository

import (
	"dimensy-bridge/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificateRepository interface {
	Create(cert *model.Certificate) error
	FindByID(id uint64) (*model.Certificate, error)
	FindBySerialNumber(serial string) (*model.Certificate, error)
	FindAllByClientID(clientID uint64) ([]model.Certificate, error)
	FindByClientUserAndCompany(clientID int64, userID *int64, companyID *int64) (*model.Certificate, error)
	Update(cert *model.Certificate) error
	Delete(id uint64) error
	FindByExternal(externalUserID, externalCompanyID *uuid.UUID) (*model.Certificate, error)
}

type certificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) CertificateRepository {
	return &certificateRepository{db: db}
}

func (r *certificateRepository) Create(cert *model.Certificate) error {
	return r.db.Create(cert).Error
}
func (r *certificateRepository) FindByExternal(externalUserID, externalCompanyID *uuid.UUID) (*model.Certificate, error) {
	var cert model.Certificate
	query := r.db.Model(&model.Certificate{})

	// Add externalUserID condition - handle both nil and value cases
	if externalUserID != nil {
		query = query.Where("external_user_id = ?", *externalUserID)
	} else {
		query = query.Where("external_user_id IS NULL")
	}

	// Add externalCompanyID condition - handle both nil and value cases
	if externalCompanyID != nil {
		query = query.Where("external_company_id = ?", *externalCompanyID)
	} else {
		query = query.Where("external_company_id IS NULL")
	}

	err := query.First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) FindByID(id uint64) (*model.Certificate, error) {
	var cert model.Certificate
	err := r.db.First(&cert, id).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) FindBySerialNumber(serial string) (*model.Certificate, error) {
	var cert model.Certificate
	err := r.db.Where("serial_number = ?", serial).First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) FindAllByClientID(clientID uint64) ([]model.Certificate, error) {
	var certs []model.Certificate
	err := r.db.Where("client_id = ?", clientID).Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindByClientUserAndCompany(clientID int64, userID *int64, companyID *int64) (*model.Certificate, error) {
	var cert model.Certificate
	query := r.db.Where("client_id = ?", clientID)

	// Add userID condition - handle both nil and value cases
	if userID != nil {
		query = query.Where("clien_user_id = ?", *userID)
	} else {
		query = query.Where("clien_user_id IS NULL")
	}

	// Add companyID condition - handle both nil and value cases
	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) Update(cert *model.Certificate) error {
	return r.db.Save(cert).Error
}

func (r *certificateRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Certificate{}, id).Error
}
