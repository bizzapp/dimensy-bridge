package psreservice

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

type ClientDocumentProcessService interface {
	CreateProcess(req *model.ClientDocumentProcess) error
	GetByExternalID(externalID uuid.UUID) (*model.ClientDocumentProcess, error)
	UpdateStatus(externalID uuid.UUID, status string) error
	DeleteByExternalID(externalID uuid.UUID) error
}

type clientDocumentProcessService struct {
	repo repository.ClientDocumentProcessRepository
}

func NewClientDocumentProcessService(repo repository.ClientDocumentProcessRepository) ClientDocumentProcessService {
	return &clientDocumentProcessService{repo: repo}
}

func (s *clientDocumentProcessService) CreateProcess(req *model.ClientDocumentProcess) error {
	if req.ExternalID == uuid.Nil {
		return fmt.Errorf("external_id is required")
	}
	return s.repo.Create(req)
}

func (s *clientDocumentProcessService) GetByExternalID(externalID uuid.UUID) (*model.ClientDocumentProcess, error) {
	return s.repo.FindByExternalID(externalID)
}

func (s *clientDocumentProcessService) UpdateStatus(externalID uuid.UUID, status string) error {
	return s.repo.UpdateStatus(externalID, status)
}

func (s *clientDocumentProcessService) DeleteByExternalID(externalID uuid.UUID) error {
	return s.repo.DeleteByExternalID(externalID)
}
