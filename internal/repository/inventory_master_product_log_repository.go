package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"dimensy-bridge/internal/model"
)

type InventoryMasterProductLogRepository interface {
	Create(ctx context.Context, log *model.InventoryMasterProductLog) error
	FindByInventoryID(ctx context.Context, inventoryID int64, page, pageSize int) ([]*model.InventoryMasterProductLog, int64, error)
	FindByMasterProductID(ctx context.Context, masterProductID int64, page, pageSize int) ([]*model.InventoryMasterProductLog, int64, error)
	GetLogsInRange(ctx context.Context, startTime, endTime time.Time) ([]*model.InventoryMasterProductLog, error)
	FindInventoryLogs(
		ctx context.Context,
		inventoryID *int64,
		page, pageSize int,
	) ([]*model.InventoryMasterProductLog, int64, error)
}

type inventoryMasterProductLogRepository struct {
	db *gorm.DB
}

func NewInventoryMasterProductLogRepository(db *gorm.DB) InventoryMasterProductLogRepository {
	return &inventoryMasterProductLogRepository{db: db}
}

func (r *inventoryMasterProductLogRepository) Create(ctx context.Context, log *model.InventoryMasterProductLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *inventoryMasterProductLogRepository) FindByInventoryID(ctx context.Context, inventoryID int64, page, pageSize int) ([]*model.InventoryMasterProductLog, int64, error) {
	var logs []*model.InventoryMasterProductLog
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Where("inventory_master_product_id = ?", inventoryID).
		Preload("MasterProduct").
		Model(&model.InventoryMasterProductLog{}).
		Count(&total).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *inventoryMasterProductLogRepository) FindInventoryLogs(
	ctx context.Context,
	inventoryID *int64,
	page, pageSize int,
) ([]*model.InventoryMasterProductLog, int64, error) {
	var logs []*model.InventoryMasterProductLog
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).
		Model(&model.InventoryMasterProductLog{}).
		Preload("MasterProduct").
		Preload("InventoryMasterProduct")

	if inventoryID != nil {
		query = query.Where("inventory_master_product_id = ?", *inventoryID)
	}

	err := query.
		Count(&total).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *inventoryMasterProductLogRepository) FindByMasterProductID(ctx context.Context, masterProductID int64, page, pageSize int) ([]*model.InventoryMasterProductLog, int64, error) {
	var logs []*model.InventoryMasterProductLog
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Where("master_product_id = ?", masterProductID).
		Preload("InventoryMasterProduct").
		Preload("MasterProduct").
		Model(&model.InventoryMasterProductLog{}).
		Count(&total).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *inventoryMasterProductLogRepository) GetLogsInRange(ctx context.Context, startTime, endTime time.Time) ([]*model.InventoryMasterProductLog, error) {
	var logs []*model.InventoryMasterProductLog
	err := r.db.WithContext(ctx).
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Preload("InventoryMasterProduct").
		Preload("MasterProduct").
		Order("time DESC").
		Find(&logs).Error

	return logs, err
}
