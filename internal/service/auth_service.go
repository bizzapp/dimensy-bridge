package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// import "dimensy-bridge/internal/model"

type AuthService interface {
	Login(email *string, password *string) (string, *model.User, error)
	Logout(token string) error
}

type authService struct {
	authRepo repository.AuthRepository

	blacklistRepo repository.TokenBlacklistRepository
}

func NewAuthService(authRepo repository.AuthRepository, blacklistRepo repository.TokenBlacklistRepository) AuthService {
	return &authService{
		authRepo: authRepo, blacklistRepo: blacklistRepo,
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

	// pastikan password tidak nil
	if user.Password == nil || *user.Password == "" {
		return "", nil, errors.New("password not set for this user")
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

func (s *authService) Logout(tokenString string) error {
	// Parse token to get expiration time
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Get expiration time from token
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid token expiration")
	}

	expiresAt := time.Unix(int64(exp), 0)

	// Add token to blacklist
	return s.blacklistRepo.Create(tokenString, expiresAt)
}
