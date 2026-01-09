package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"dimensy-bridge/internal/model"
)

type InventoryMasterProductRepository interface {
	Create(ctx context.Context, inventory *model.InventoryMasterProduct) error
	Update(ctx context.Context, inventory *model.InventoryMasterProduct) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.InventoryMasterProduct, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*model.InventoryMasterProduct, int64, error)
	FindByMasterProductIDAndVendor(ctx context.Context, masterProductID int64, vendorName string) (*model.InventoryMasterProduct, error)
	GetLowStockItems(ctx context.Context, threshold int) ([]*model.InventoryMasterProduct, error)
	GetTotalInventoryValue(ctx context.Context) (float64, error)
}

type inventoryMasterProductRepository struct {
	db *gorm.DB
}

func NewInventoryMasterProductRepository(db *gorm.DB) InventoryMasterProductRepository {
	return &inventoryMasterProductRepository{db: db}
}

func (r *inventoryMasterProductRepository) Create(ctx context.Context, inventory *model.InventoryMasterProduct) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *inventoryMasterProductRepository) Update(ctx context.Context, inventory *model.InventoryMasterProduct) error {
	return r.db.WithContext(ctx).Save(inventory).Error
}

func (r *inventoryMasterProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.InventoryMasterProduct{}, id).Error
}

func (r *inventoryMasterProductRepository) FindByID(ctx context.Context, id int64) (*model.InventoryMasterProduct, error) {
	var inventory *model.InventoryMasterProduct
	err := r.db.WithContext(ctx).
		Preload("MasterProduct").
		First(&inventory, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory master product not found")
		}
		return nil, err
	}

	return inventory, nil
}

func (r *inventoryMasterProductRepository) FindAll(ctx context.Context, page, pageSize int) ([]*model.InventoryMasterProduct, int64, error) {
	var inventories []*model.InventoryMasterProduct
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Preload("MasterProduct").
		Model(&model.InventoryMasterProduct{}).
		Count(&total).
		Offset(offset).
		Limit(pageSize).
		Find(&inventories).Error

	if err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}

func (r *inventoryMasterProductRepository) FindByMasterProductIDAndVendor(ctx context.Context, masterProductID int64, vendorName string) (*model.InventoryMasterProduct, error) {
	var inventory *model.InventoryMasterProduct
	err := r.db.WithContext(ctx).
		Where("master_product_id = ? AND vendor_name = ?", masterProductID, vendorName).
		Preload("MasterProduct").
		First(&inventory).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return inventory, nil
}

func (r *inventoryMasterProductRepository) GetLowStockItems(ctx context.Context, threshold int) ([]*model.InventoryMasterProduct, error) {
	var inventories []*model.InventoryMasterProduct
	err := r.db.WithContext(ctx).
		Where("current_stock <= ? AND is_processed = false", threshold).
		Preload("MasterProduct").
		Find(&inventories).Error

	return inventories, err
}

func (r *inventoryMasterProductRepository) GetTotalInventoryValue(ctx context.Context) (float64, error) {
	var totalValue float64
	err := r.db.WithContext(ctx).
		Model(&model.InventoryMasterProduct{}).
		Select("COALESCE(SUM(price * current_stock), 0) as total").
		Row().
		Scan(&totalValue)

	return totalValue, err
}
