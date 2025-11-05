package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type SubscriptionPlanService interface {
	GetAll() ([]model.SubscriptionPlan, error)
	GetByID(id int64) (*model.SubscriptionPlan, error)
	Create(plan *model.SubscriptionPlan) error
}

type subscriptionPlanService struct {
	repo repository.SubscriptionPlanRepository
}

func NewSubscriptionPlanService(repo repository.SubscriptionPlanRepository) SubscriptionPlanService {
	return &subscriptionPlanService{repo: repo}
}

func (s *subscriptionPlanService) GetAll() ([]model.SubscriptionPlan, error) {
	return s.repo.FindAll()
}

func (s *subscriptionPlanService) GetByID(id int64) (*model.SubscriptionPlan, error) {
	return s.repo.FindByID(id)
}

func (s *subscriptionPlanService) Create(plan *model.SubscriptionPlan) error {
	return s.repo.Create(plan)
}
