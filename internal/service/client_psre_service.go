package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type ClientPsreService interface {
	GetByID(id int64) (*model.ClientPsre, error)
	GetByClientID(clientID int64) (*model.ClientPsre, error)
	UpdatePsre(psre *model.ClientPsre) error
	DeletePsre(id int64) error
}

type clientPsreService struct {
	psreRepo   repository.ClientPsreRepository
	clientRepo repository.ClientRepository
}

func NewClientPsreService(psreRepo repository.ClientPsreRepository, clientRepo repository.ClientRepository) ClientPsreService {
	return &clientPsreService{psreRepo, clientRepo}
}

func (s *clientPsreService) GetByID(id int64) (*model.ClientPsre, error) {
	return s.psreRepo.FindByID(id)
}

func (s *clientPsreService) GetByClientID(clientID int64) (*model.ClientPsre, error) {
	return s.psreRepo.FindByClientID(clientID)
}

func (s *clientPsreService) UpdatePsre(psre *model.ClientPsre) error {
	return s.psreRepo.Update(psre)
}

func (s *clientPsreService) DeletePsre(id int64) error {
	return s.psreRepo.Delete(id)
}
