package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type CertificateRepository interface {
	Create(cert *model.Certificate) error
	FindByID(id uint64) (*model.Certificate, error)
	FindBySerialNumber(serial string) (*model.Certificate, error)
	FindAllByClientID(clientID uint64) ([]model.Certificate, error)
	Update(cert *model.Certificate) error
	Delete(id uint64) error
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

func (r *certificateRepository) Update(cert *model.Certificate) error {
	return r.db.Save(cert).Error
}

func (r *certificateRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Certificate{}, id).Error
}
