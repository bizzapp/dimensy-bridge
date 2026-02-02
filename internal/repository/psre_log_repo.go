package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type PsreLogRepository interface {
	Create(data *model.PsreLog) error
}

type psreLogRepository struct {
	db *gorm.DB
}

func NewPsreLogRepository(db *gorm.DB) PsreLogRepository {
	return &psreLogRepository{db: db}
}

func (r *psreLogRepository) Create(data *model.PsreLog) error {
	return r.db.Create(data).Error
}
