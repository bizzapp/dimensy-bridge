package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type UserService interface {
	GetUsers(page, limit int, filters map[string]interface{}) ([]model.User, int64, error)
	GetUserByID(id int64) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUsers(page, limit int, filters map[string]interface{}) ([]model.User, int64, error) {
	offset := (page - 1) * limit
	return s.repo.FindAll(limit, offset, filters)
}

func (s *userService) GetUserByID(id int64) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) CreateUser(user *model.User) error {
	return s.repo.Create(user)
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}
