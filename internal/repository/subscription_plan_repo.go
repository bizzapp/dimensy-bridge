package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type SubscriptionPlanRepository interface {
	FindAll() ([]model.SubscriptionPlan, error)
	FindByID(id int64) (*model.SubscriptionPlan, error)
	Create(plan *model.SubscriptionPlan) error
}

type subscriptionPlanRepository struct {
	db *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{db: db}
}

func (r *subscriptionPlanRepository) FindAll() ([]model.SubscriptionPlan, error) {
	var plans []model.SubscriptionPlan
	if err := r.db.Preload("Details").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *subscriptionPlanRepository) FindByID(id int64) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan
	if err := r.db.Preload("Details").First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *subscriptionPlanRepository) Create(plan *model.SubscriptionPlan) error {
	return r.db.Create(plan).Error
}
