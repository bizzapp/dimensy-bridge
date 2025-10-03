package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientCompanyRepository interface {
	FindAll() ([]model.ClientCompany, error)
	FindByID(id int64) (*model.ClientCompany, error)
	Create(company *model.ClientCompany) error
	Update(company *model.ClientCompany) error
	Delete(id int64) error
}

type clientCompanyRepository struct {
	db *gorm.DB
}

func NewClientCompanyRepository(db *gorm.DB) ClientCompanyRepository {
	return &clientCompanyRepository{db}
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
