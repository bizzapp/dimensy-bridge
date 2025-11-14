package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type MasterProductService interface {
	GetProducts(page, limit int, filters map[string]interface{}) ([]model.MasterProduct, int64, error)
	GetProductByID(id int64) (*model.MasterProduct, error)
	CreateProduct(product *model.MasterProduct) error
	UpdateProduct(product *model.MasterProduct) error
	DeleteProduct(id int64) error
	GetHistory(masterProductID *int64, limit int) ([]dto.MasterProductHistoryItem, error)
}

type masterProductService struct {
	repo repository.MasterProductRepository
}

func NewMasterProductService(repo repository.MasterProductRepository) MasterProductService {
	return &masterProductService{repo}
}

func (s *masterProductService) GetHistory(masterProductID *int64, limit int) ([]dto.MasterProductHistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetHistory(masterProductID, limit)
}

func (s *masterProductService) GetProducts(page, limit int, filters map[string]interface{}) ([]model.MasterProduct, int64, error) {
	offset := (page - 1) * limit
	return s.repo.FindAll(limit, offset, filters)
}

func (s *masterProductService) GetProductByID(id int64) (*model.MasterProduct, error) {
	return s.repo.FindByID(id)
}

func (s *masterProductService) CreateProduct(product *model.MasterProduct) error {
	return s.repo.Create(product)
}

func (s *masterProductService) UpdateProduct(product *model.MasterProduct) error {
	return s.repo.Update(product)
}

func (s *masterProductService) DeleteProduct(id int64) error {
	return s.repo.Delete(id)
}
