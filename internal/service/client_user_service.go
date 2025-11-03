package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type ClientUserService interface {
	GetAll() ([]model.ClientUser, error)
	GetByID(id uint) (*model.ClientUser, error)
	Create(user *model.ClientUser) error
	Update(user *model.ClientUser) error
	Delete(id uint) error
}

type clientUserService struct {
	repo repository.ClientUserRepository
}

func NewClientUserService(repo repository.ClientUserRepository) ClientUserService {
	return &clientUserService{repo}
}

func (s *clientUserService) GetAll() ([]model.ClientUser, error) {
	return s.repo.FindAll()
}

func (s *clientUserService) GetByID(id uint) (*model.ClientUser, error) {
	return s.repo.FindByID(id)
}

func (s *clientUserService) Create(user *model.ClientUser) error {
	return s.repo.Create(user)
}

func (s *clientUserService) Update(user *model.ClientUser) error {
	return s.repo.Update(user)
}

func (s *clientUserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
