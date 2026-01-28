package repository

import (
	"context"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"time"

	"gorm.io/gorm"
)

type QuotaClientReductionRepository interface {
	Create(reduction *model.QuotaClientReduction) error
	FindByID(id int64) (*model.QuotaClientReduction, error)
	FindAll() ([]model.QuotaClientReduction, error)
	FindByQuotaClientID(quotaClientID int64) ([]model.QuotaClientReduction, error)
	Update(reduction *model.QuotaClientReduction) error
	Delete(id int64) error
	GetTimeSeries(
		ctx context.Context,
		startTime time.Time,
		rangeKey string,
		clientID *int64,
		masterProductID *int64,
	) ([]dto.ReductionTimeSeriesRow, error)
}

type quotaClientReductionRepository struct {
	db *gorm.DB
}

func NewQuotaClientReductionRepository(db *gorm.DB) QuotaClientReductionRepository {
	return &quotaClientReductionRepository{db: db}
}

func (r *quotaClientReductionRepository) Create(reduction *model.QuotaClientReduction) error {
	return r.db.Create(reduction).Error
}

func (r *quotaClientReductionRepository) FindByID(id int64) (*model.QuotaClientReduction, error) {
	var reduction model.QuotaClientReduction
	if err := r.db.First(&reduction, id).Error; err != nil {
		return nil, err
	}
	return &reduction, nil
}

func (r *quotaClientReductionRepository) FindAll() ([]model.QuotaClientReduction, error) {
	var reductions []model.QuotaClientReduction
	if err := r.db.Find(&reductions).Error; err != nil {
		return nil, err
	}
	return reductions, nil
}

func (r *quotaClientReductionRepository) FindByQuotaClientID(quotaClientID int64) ([]model.QuotaClientReduction, error) {
	var reductions []model.QuotaClientReduction
	if err := r.db.Where("quota_client_id = ?", quotaClientID).Find(&reductions).Error; err != nil {
		return nil, err
	}
	return reductions, nil
}

func (r *quotaClientReductionRepository) Update(reduction *model.QuotaClientReduction) error {
	return r.db.Save(reduction).Error
}

func (r *quotaClientReductionRepository) Delete(id int64) error {
	return r.db.Delete(&model.QuotaClientReduction{}, id).Error
}

func (r *quotaClientReductionRepository) GetTimeSeries(
	ctx context.Context,
	startTime time.Time,
	rangeKey string,
	clientID *int64,
	masterProductID *int64,
) ([]dto.ReductionTimeSeriesRow, error) {

	var rows []dto.ReductionTimeSeriesRow

	// default daily
	timeGroupExpr := "DATE_TRUNC('day', qcr.created_at)"
	timeLabelExpr := "TO_CHAR(DATE_TRUNC('day', qcr.created_at), 'YYYY-MM-DD')"

	// monthly
	if rangeKey == "last_12_months" {
		timeGroupExpr = "DATE_TRUNC('month', qcr.created_at)"
		timeLabelExpr = "TO_CHAR(DATE_TRUNC('month', qcr.created_at), 'YYYY-MM')"
	}

	q := r.db.WithContext(ctx).
		Table("quota_client_reductions qcr").
		Select(`
			`+timeLabelExpr+` AS label,
			mp.name AS master_product_name,
			COALESCE(SUM(qcr.quantity), 0) AS total_quantity
		`).
		Joins("JOIN quota_clients qc ON qc.id = qcr.quota_client_id").
		Joins("JOIN master_products mp ON mp.id = qc.master_product_id").
		Where("qcr.created_at >= ?", startTime).
		Group(timeGroupExpr + ", mp.name").
		Order(timeGroupExpr + " ASC")

	if clientID != nil {
		q = q.Where("qc.client_id = ?", *clientID)
	}
	if masterProductID != nil {
		q = q.Where("mp.id = ?", *masterProductID)
	}

	err := q.Scan(&rows).Error
	return rows, err
}
