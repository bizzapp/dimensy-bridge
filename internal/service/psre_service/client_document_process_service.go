package psreservice

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"fmt"
)

type ClientDocumentProcessService interface {
	CreateProcess(req *model.ClientDocumentProcess) error
	GetByExternalID(externalID string) (*model.ClientDocumentProcess, error)
	UpdateStatus(externalID, status string) error
	DeleteByExternalID(externalID string) error
}

type clientDocumentProcessService struct {
	repo repository.ClientDocumentProcessRepository
}

func NewClientDocumentProcessService(repo repository.ClientDocumentProcessRepository) ClientDocumentProcessService {
	return &clientDocumentProcessService{repo: repo}
}

func (s *clientDocumentProcessService) CreateProcess(req *model.ClientDocumentProcess) error {
	if req.ExternalID == "" {
		return fmt.Errorf("external_id is required")
	}
	return s.repo.Create(req)
}

func (s *clientDocumentProcessService) GetByExternalID(externalID string) (*model.ClientDocumentProcess, error) {
	return s.repo.FindByExternalID(externalID)
}

func (s *clientDocumentProcessService) UpdateStatus(externalID, status string) error {
	return s.repo.UpdateStatus(externalID, status)
}

func (s *clientDocumentProcessService) DeleteByExternalID(externalID string) error {
	return s.repo.DeleteByExternalID(externalID)
}
