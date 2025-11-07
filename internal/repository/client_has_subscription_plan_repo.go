package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientHasSubscriptionPlanRepository interface {
	Create(data *model.ClientHasSubscriptionPlan) error
	Update(data *model.ClientHasSubscriptionPlan) error
	FindByID(id int64) (*model.ClientHasSubscriptionPlan, error)
	FindAll() ([]model.ClientHasSubscriptionPlan, error)
	Delete(id int64) error
	WithTransaction(fn func(tx *gorm.DB) error) error
	FindByIDTx(tx *gorm.DB, id int64) (*model.ClientHasSubscriptionPlan, error)
}

type clientHasSubscriptionPlanRepository struct {
	db *gorm.DB
}

func NewClientHasSubscriptionPlanRepository(db *gorm.DB) ClientHasSubscriptionPlanRepository {
	return &clientHasSubscriptionPlanRepository{db: db}
}
func (r *clientHasSubscriptionPlanRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *clientHasSubscriptionPlanRepository) FindByIDTx(tx *gorm.DB, id int64) (*model.ClientHasSubscriptionPlan, error) {
	var sub model.ClientHasSubscriptionPlan
	err := tx.Preload("SubscriptionPlan.Details").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&sub, id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *clientHasSubscriptionPlanRepository) Create(data *model.ClientHasSubscriptionPlan) error {
	return r.db.Create(data).Error
}

func (r *clientHasSubscriptionPlanRepository) Update(data *model.ClientHasSubscriptionPlan) error {
	return r.db.Save(data).Error
}

func (r *clientHasSubscriptionPlanRepository) FindByID(id int64) (*model.ClientHasSubscriptionPlan, error) {
	var result model.ClientHasSubscriptionPlan
	if err := r.db.Preload("SubscriptionPlan.Details").First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *clientHasSubscriptionPlanRepository) FindAll() ([]model.ClientHasSubscriptionPlan, error) {
	var results []model.ClientHasSubscriptionPlan
	if err := r.db.Preload("SubscriptionPlan.Details").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *clientHasSubscriptionPlanRepository) Delete(id int64) error {
	return r.db.Delete(&model.ClientHasSubscriptionPlan{}, id).Error
}
