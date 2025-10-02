package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type MasterProductRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.MasterProduct, int64, error)
	FindByID(id int64) (*model.MasterProduct, error)
	Create(product *model.MasterProduct) error
	Update(product *model.MasterProduct) error
	Delete(id int64) error
}

type masterProductRepository struct {
	db *gorm.DB
}

func NewMasterProductRepository(db *gorm.DB) MasterProductRepository {
	return &masterProductRepository{db}
}

func (r *masterProductRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]model.MasterProduct, int64, error) {
	var products []model.MasterProduct
	var total int64

	query := r.db.Model(&model.MasterProduct{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *masterProductRepository) FindByID(id int64) (*model.MasterProduct, error) {
	var product model.MasterProduct
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *masterProductRepository) Create(product *model.MasterProduct) error {
	return r.db.Create(product).Error
}

func (r *masterProductRepository) Update(product *model.MasterProduct) error {
	return r.db.Save(product).Error
}

func (r *masterProductRepository) Delete(id int64) error {
	return r.db.Delete(&model.MasterProduct{}, id).Error
}
