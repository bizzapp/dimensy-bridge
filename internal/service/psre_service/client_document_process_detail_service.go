package psreservice

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type ClientDocumentProcessDetailService interface {
	CreateDetail(req *model.ClientDocumentProcessDetail) error
	GetByClientID(clientID int64) ([]model.ClientDocumentProcessDetail, error)
	DeleteByClientID(clientID int64) error
}

type clientDocumentProcessDetailService struct {
	repo repository.ClientDocumentProcessDetailRepository
}

func NewClientDocumentProcessDetailService(repo repository.ClientDocumentProcessDetailRepository) ClientDocumentProcessDetailService {
	return &clientDocumentProcessDetailService{repo: repo}
}

func (s *clientDocumentProcessDetailService) CreateDetail(req *model.ClientDocumentProcessDetail) error {
	return s.repo.Create(req)
}

func (s *clientDocumentProcessDetailService) GetByClientID(clientID int64) ([]model.ClientDocumentProcessDetail, error) {
	return s.repo.FindByClientID(clientID)
}

func (s *clientDocumentProcessDetailService) DeleteByClientID(clientID int64) error {
	return s.repo.DeleteByClientID(clientID)
}
