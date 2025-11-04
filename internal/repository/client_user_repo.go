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
	UpdateActiveStatus(externalID string, active bool) error
	UpdateVerifyPhoneStatus(externalID string, verify bool) error
	UpdateVerifyKYCStatus(externalID string, verify bool) error
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

func (r *clientUserRepository) FindByID(id uint) (*model.ClientUser, error) {
	var user model.ClientUser
	err := r.db.Preload("Client").Preload("ClientCompany").First(&user, id).Error
	return &user, err
}
func (r *clientUserRepository) UpdateActiveStatus(externalID string, active bool) error {
	return r.db.Model(&model.ClientUser{}).
		Where("external_id = ?", externalID).
		Update("is_active", active).Error
}
func (r *clientUserRepository) UpdateVerifyPhoneStatus(externalID string, verify bool) error {
	return r.db.Model(&model.ClientUser{}).
		Where("external_id = ?", externalID).
		Update("is_verify_phone", verify).Error
}

func (r *clientUserRepository) UpdateVerifyKYCStatus(externalID string, verify bool) error {
	return r.db.Model(&model.ClientUser{}).
		Where("external_id = ?", externalID).
		Update("is_verify_kyc", verify).Error
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
