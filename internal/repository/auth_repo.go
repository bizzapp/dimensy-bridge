package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email *string) (*model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) FindByEmail(email *string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
