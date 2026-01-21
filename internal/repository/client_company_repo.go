package repository

import (
	"dimensy-bridge/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientCompanyRepository interface {
	FindAll() ([]model.ClientCompany, error)
	FindByID(id int64) (*model.ClientCompany, error)
	Create(company *model.ClientCompany) error
	Update(company *model.ClientCompany) error
	Delete(id int64) error
	UpdateExternalID(id int64, externalID uuid.UUID) error
	FindByExternalID(externalID uuid.UUID) (*model.ClientCompany, error)
	CreateTx(tx *gorm.DB, c *model.ClientCompany) error
	UpdateExternalIDTx(tx *gorm.DB, id int64, externalID uuid.UUID) error
	UpdateTx(tx *gorm.DB, c *model.ClientCompany) error
}

type clientCompanyRepository struct {
	db *gorm.DB
}

func NewClientCompanyRepository(db *gorm.DB) ClientCompanyRepository {
	return &clientCompanyRepository{db}
}

func (r *clientCompanyRepository) CreateTx(tx *gorm.DB, c *model.ClientCompany) error {
	return tx.Create(c).Error
}

func (r *clientCompanyRepository) UpdateExternalIDTx(tx *gorm.DB, id int64, externalID uuid.UUID) error {
	return tx.Model(&model.ClientCompany{}).Where("id = ?", id).Update("external_id", externalID).Error
}

func (r *clientCompanyRepository) UpdateTx(tx *gorm.DB, c *model.ClientCompany) error {
	return tx.Save(c).Error
}

func (r *clientCompanyRepository) FindByExternalID(externalID uuid.UUID) (*model.ClientCompany, error) {
	var company model.ClientCompany
	if err := r.db.Where("external_id = ?", externalID).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *clientCompanyRepository) FindAll() ([]model.ClientCompany, error) {
	var companies []model.ClientCompany
	if err := r.db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *clientCompanyRepository) FindByID(id int64) (*model.ClientCompany, error) {
	var company model.ClientCompany
	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *clientCompanyRepository) Create(company *model.ClientCompany) error {
	return r.db.Create(company).Error
}

func (r *clientCompanyRepository) Update(company *model.ClientCompany) error {
	return r.db.Save(company).Error
}

func (r *clientCompanyRepository) Delete(id int64) error {
	return r.db.Delete(&model.ClientCompany{}, id).Error
}

func (r *clientCompanyRepository) UpdateExternalID(id int64, externalID uuid.UUID) error {
	return r.db.Model(&model.ClientCompany{}).
		Where("id = ?", id).
		Update("external_id", externalID).Error
}
