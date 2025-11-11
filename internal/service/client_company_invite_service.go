package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
)

type ClientCompanyInviteService interface {
	CreateInvite(req *model.ClientCompanyInvite) (*model.ClientCompanyInvite, error)
	GetInviteByExternal(userID, companyID string) (*model.ClientCompanyInvite, error)
	VerifyInvite(id int64) error
}

type clientCompanyInviteService struct {
	repo repository.ClientCompanyInviteRepository
}

func NewClientCompanyInviteService(repo repository.ClientCompanyInviteRepository) ClientCompanyInviteService {
	return &clientCompanyInviteService{repo}
}

func (s *clientCompanyInviteService) CreateInvite(req *model.ClientCompanyInvite) (*model.ClientCompanyInvite, error) {
	existing, _ := s.repo.FindByExternal(req.ExternalUserID, req.ExternalCompanyID)
	if existing != nil {
		return nil, errors.New("invite already exists")
	}
	err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *clientCompanyInviteService) GetInviteByExternal(userID, companyID string) (*model.ClientCompanyInvite, error) {
	return s.repo.FindByExternal(userID, companyID)
}

func (s *clientCompanyInviteService) VerifyInvite(id int64) error {
	return s.repo.VerifyInvite(id)
}
