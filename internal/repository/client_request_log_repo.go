package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientRequestLogRepository interface {
	Create(log *model.ClientRequestLog) error
}

type clientRequestLogRepository struct {
	db *gorm.DB
}

func NewClientRequestLogRepository(db *gorm.DB) ClientRequestLogRepository {
	return &clientRequestLogRepository{db}
}

func (r *clientRequestLogRepository) Create(log *model.ClientRequestLog) error {
	return r.db.Create(log).Error
}
