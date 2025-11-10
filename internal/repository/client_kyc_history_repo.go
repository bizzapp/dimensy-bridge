package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientKYCHistoryRepository interface {
	Create(data *model.ClientKYCHistory) error
	Update(data *model.ClientKYCHistory) error
	FindByID(id int64) (*model.ClientKYCHistory, error)
	FindAll() ([]model.ClientKYCHistory, error)
	Delete(id int64) error
	FindByClientUserID(clientUserID int64) ([]model.ClientKYCHistory, error)
}

type clientKYCHistoryRepository struct {
	db *gorm.DB
}

func NewClientKYCHistoryRepository(db *gorm.DB) ClientKYCHistoryRepository {
	return &clientKYCHistoryRepository{db: db}
}

func (r *clientKYCHistoryRepository) Create(data *model.ClientKYCHistory) error {
	return r.db.Create(data).Error
}

func (r *clientKYCHistoryRepository) Update(data *model.ClientKYCHistory) error {
	return r.db.Save(data).Error
}

func (r *clientKYCHistoryRepository) FindByID(id int64) (*model.ClientKYCHistory, error) {
	var result model.ClientKYCHistory
	if err := r.db.First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *clientKYCHistoryRepository) FindAll() ([]model.ClientKYCHistory, error) {
	var results []model.ClientKYCHistory
	if err := r.db.Preload("Client").Preload("ClientUser").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *clientKYCHistoryRepository) Delete(id int64) error {
	return r.db.Delete(&model.ClientKYCHistory{}, id).Error
}

func (r *clientKYCHistoryRepository) FindByClientUserID(clientUserID int64) ([]model.ClientKYCHistory, error) {
	var results []model.ClientKYCHistory
	if err := r.db.Where("client_user_id = ?", clientUserID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
