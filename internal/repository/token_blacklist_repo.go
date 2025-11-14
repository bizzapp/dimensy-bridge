package repository

import (
	"dimensy-bridge/internal/model"
	"time"

	"gorm.io/gorm"
)

type TokenBlacklistRepository interface {
	Create(token string, expiresAt time.Time) error
	IsBlacklisted(token string) (bool, error)
	CleanupExpired() error
}

type tokenBlacklistRepository struct {
	db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) TokenBlacklistRepository {
	return &tokenBlacklistRepository{db: db}
}

func (r *tokenBlacklistRepository) Create(token string, expiresAt time.Time) error {
	blacklist := &model.TokenBlacklist{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return r.db.Create(blacklist).Error
}

func (r *tokenBlacklistRepository) IsBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&model.TokenBlacklist{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *tokenBlacklistRepository) CleanupExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).
		Delete(&model.TokenBlacklist{}).Error
}
