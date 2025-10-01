package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// import "dimensy-bridge/internal/model"

type AuthService interface {
	Login(email *string, password *string) (string, *model.User, error)
	Logout(token string) error
}

type authService struct {
	authRepo  repository.AuthRepository
	blacklist map[string]bool
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo:  authRepo,
		blacklist: make(map[string]bool),
	}
}

func (s *authService) Login(email *string, password *string) (string, *model.User, error) {
	if email == nil || password == nil || *email == "" || *password == "" {
		return "", nil, errors.New("email and password are required")
	}

	user, err := s.authRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("email not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*password)); err != nil {
		return "", nil, errors.New("invalid password")
	}
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, user.Name)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", nil, errors.New("JWT_SECRET not configured")
	}

	return token, user, nil
}

func (s *authService) Logout(token string) error {
	if token == "" {
		return errors.New("empty token")
	}

	// simpan ke blacklist (contoh memory, production lebih baik pakai Redis)
	s.blacklist[token] = true
	return nil
}

func (s *authService) IsBlacklisted(token string) bool {
	return s.blacklist[token]
}
