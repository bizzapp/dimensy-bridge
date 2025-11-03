package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
	// "project/model"
)

type ClientUserRepository interface {
	FindAll() ([]model.ClientUser, error)
	FindByID(id uint) (*model.ClientUser, error)
	Create(user *model.ClientUser) error
	Update(user *model.ClientUser) error
	Delete(id uint) error
	WithTx(tx *gorm.DB) ClientUserRepository
	DB() *gorm.DB
}

type clientUserRepository struct {
	db *gorm.DB
}

func NewClientUserRepository(db *gorm.DB) ClientUserRepository {
	return &clientUserRepository{db}
}

func (r *clientUserRepository) FindAll() ([]model.ClientUser, error) {
	var users []model.ClientUser
	err := r.db.Preload("ClientCompany").Preload("Client").Find(&users).Error
	return users, err
}

func (r *clientUserRepository) WithTx(tx *gorm.DB) ClientUserRepository {
	return &clientUserRepository{db: tx}
}
func (r *clientUserRepository) DB() *gorm.DB { return r.db }

func (r *clientUserRepository) FindByID(id uint) (*model.ClientUser, error) {
	var user model.ClientUser
	err := r.db.Preload("Client").Preload("ClientCompany").First(&user, id).Error
	return &user, err
}

func (r *clientUserRepository) Create(user *model.ClientUser) error {
	return r.db.Create(user).Error
}

func (r *clientUserRepository) Update(user *model.ClientUser) error {
	return r.db.Save(user).Error
}

func (r *clientUserRepository) Delete(id uint) error {
	return r.db.Delete(&model.ClientUser{}, id).Error
}
