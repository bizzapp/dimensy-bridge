package repository

import (
	"context"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

type MasterProductRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.MasterProduct, int64, error)
	FindByID(id int64) (*model.MasterProduct, error)
	Create(product *model.MasterProduct) error
	Update(product *model.MasterProduct) error
	Delete(id int64) error
	GetHistory(masterProductID *int64, limit int) ([]dto.MasterProductHistoryItem, error)
	FindForChart(
		ctx context.Context,
		masterProductID *int64,
	) ([]string, error)
}

type masterProductRepository struct {
	db *gorm.DB
}

func NewMasterProductRepository(db *gorm.DB) MasterProductRepository {
	return &masterProductRepository{db}
}

func (r *masterProductRepository) GetHistory(masterProductID *int64, limit int) ([]dto.MasterProductHistoryItem, error) {
	var adds []dto.MasterProductHistoryItem
	var reducs []dto.MasterProductHistoryItem

	addQuery := r.db.Model(&model.MasterProductAddition{}).
		Select(`
			id,
			master_product_id,
			quantity,
			type,
			created_at
		`).
		Where("is_process = ?", true)

	reduceQuery := r.db.Model(&model.MasterProductReduction{}).
		Select(`
			id,
			master_product_id,
			(quantity * -1) AS quantity,
			type,
			created_at
		`)

	// Jika ada filter masterProductID
	if masterProductID != nil {
		addQuery = addQuery.Where("master_product_id = ?", *masterProductID)
		reduceQuery = reduceQuery.Where("master_product_id = ?", *masterProductID)
	}

	if err := addQuery.Limit(limit).Find(&adds).Error; err != nil {
		return nil, err
	}
	for i := range adds {
		adds[i].Direction = "ADDITION"
	}

	if err := reduceQuery.Limit(limit).Find(&reducs).Error; err != nil {
		return nil, err
	}
	for i := range reducs {
		reducs[i].Direction = "REDUCTION"
	}

	// Gabungkan
	all := append(adds, reducs...)

	// Sort berdasarkan waktu (DESC)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(*all[j].CreatedAt)
	})

	if len(all) > limit {
		all = all[:limit]
	}

	return all, nil
}
func (r *masterProductRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]model.MasterProduct, int64, error) {
	var products []model.MasterProduct
	var total int64

	query := r.db.Model(&model.MasterProduct{})

	for key, value := range filters {
		query = query.Where(key+" ILIKE ?", "%"+fmt.Sprint(value)+"%")
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

func (r *masterProductRepository) FindForChart(
	ctx context.Context,
	masterProductID *int64,
) ([]string, error) {

	var names []string

	q := r.db.WithContext(ctx).
		Model(&model.MasterProduct{}).
		Select("name").
		Where("deleted_at IS NULL")

	q = q.Where(`
    name ILIKE ? 
    OR name ILIKE ? 
    OR name ILIKE ?
    OR name ILIKE ?
    OR name ILIKE ?
	`,
		"%kyc%",
		"%meterai%",
		"%otp%",
		"%sign%",
		"%stamp%",
	)

	if masterProductID != nil {
		q = q.Where("id = ?", *masterProductID)
	}

	err := q.Order("name ASC").Pluck("name", &names).Error
	return names, err
}
