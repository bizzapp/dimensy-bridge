package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientPsreRepository interface {
	FindByID(id int64) (*model.ClientPsre, error)
	FindByClientID(clientID int64) (*model.ClientPsre, error)
	Create(psre *model.ClientPsre) error
	Update(psre *model.ClientPsre) error
	Delete(id int64) error
}

type clientPsreRepository struct {
	db *gorm.DB
}

func NewClientPsreRepository(db *gorm.DB) ClientPsreRepository {
	return &clientPsreRepository{db}
}

func (r *clientPsreRepository) FindByID(id int64) (*model.ClientPsre, error) {
	var psre model.ClientPsre
	if err := r.db.Preload("Client").First(&psre, id).Error; err != nil {
		return nil, err
	}
	return &psre, nil
}

func (r *clientPsreRepository) FindByClientID(clientID int64) (*model.ClientPsre, error) {
	var psre model.ClientPsre
	if err := r.db.Preload("Client").Where("client_id = ?", clientID).First(&psre).Error; err != nil {
		return nil, err
	}
	return &psre, nil
}

func (r *clientPsreRepository) Create(psre *model.ClientPsre) error {
	return r.db.Create(psre).Error
}

func (r *clientPsreRepository) Update(psre *model.ClientPsre) error {
	return r.db.Save(psre).Error
}

func (r *clientPsreRepository) Delete(id int64) error {
	return r.db.Delete(&model.ClientPsre{}, id).Error
}
