package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ClientUserService interface {
	GetAll() ([]model.ClientUser, error)
	GetByID(id uint) (*model.ClientUser, error)
	Create(user *model.ClientUser) (*model.ClientUser, error)
	Update(user *model.ClientUser) error
	Delete(id uint) error
	GetByExternalID(externalID uuid.UUID) (*model.ClientUser, error)
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

func (s *clientUserService) GetByExternalID(externalID uuid.UUID) (*model.ClientUser, error) {
	return s.repo.FindByExternalID(externalID)
}

func (s *clientUserService) GetByID(id uint) (*model.ClientUser, error) {
	return s.repo.FindByID(id)
}

func (s *clientUserService) Create(user *model.ClientUser) (*model.ClientUser, error) {
	// Validasi input
	if err := s.validateClientUser(user); err != nil {
		return nil, err
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Simpan ke database
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *clientUserService) Update(user *model.ClientUser) error {
	return s.repo.Update(user)
}

func (s *clientUserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// validateClientUser melakukan validasi data ClientUser
func (s *clientUserService) validateClientUser(user *model.ClientUser) error {
	// Validasi NIK
	if strings.TrimSpace(*user.NIK) == "" {
		return errors.New("NIK is required")
	}

	// Validasi Name
	if strings.TrimSpace(*user.Name) == "" {
		return errors.New("name is required")
	}

	// Validasi Email
	if strings.TrimSpace(*user.Email) == "" {
		return errors.New("email is required")
	}

	// Validasi Phone
	if strings.TrimSpace(*user.Phone) == "" {
		return errors.New("phone is required")
	}

	// Validasi ClientID
	if user.ClientID <= 0 {
		return errors.New("client_id is required and must be greater than 0")
	}

	// Validasi Birthdate (tidak boleh zero time)
	if user.Birthdate.IsZero() {
		return errors.New("birthdate is required")
	}

	return nil
}
