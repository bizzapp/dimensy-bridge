package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.Client, int64, error)
	FindByID(id int64) (*model.Client, error)
	FindByExternalID(externalID string) (*model.Client, error)
	Create(client *model.Client) error
	Update(client *model.Client) error
	Delete(id int64) error
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db}
}
func (r *clientRepository) FindByExternalID(externalID string) (*model.Client, error) {
	var clientPsre model.ClientPsre

	// Cari di tabel client_psres dan preload Client (dan relasi lain bila perlu)
	if err := r.db.
		Preload("Client").
		Where("external_id = ?", externalID).
		First(&clientPsre).Error; err != nil {
		return nil, err
	}

	// Kembalikan client terkait
	return &clientPsre.Client, nil
}

func (r *clientRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]model.Client, int64, error) {
	var clients []model.Client
	var total int64

	query := r.db.Model(&model.Client{}).Preload("User").Preload("ClientPsre")

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

func (r *clientRepository) FindByID(id int64) (*model.Client, error) {
	var client model.Client
	if err := r.db.Preload("User").Preload("ClientPsre").First(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) Create(client *model.Client) error {
	return r.db.Create(client).Error
}

func (r *clientRepository) Update(client *model.Client) error {
	return r.db.Save(client).Error
}

func (r *clientRepository) Delete(id int64) error {
	return r.db.Delete(&model.Client{}, id).Error
}
