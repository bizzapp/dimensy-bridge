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
	FindByExternalID(externalID string) (*model.ClientUser, error)
	CreateOrUpdate(user *model.ClientUser) error
	FindByExternalIDs(externalIDs []string) ([]model.ClientUser, error)
}

type clientUserRepository struct {
	db *gorm.DB
}

func NewClientUserRepository(db *gorm.DB) ClientUserRepository {
	return &clientUserRepository{db}
}
func (r *clientUserRepository) FindByExternalID(externalID string) (*model.ClientUser, error) {
	var user model.ClientUser
	err := r.db.Preload("Client").Preload("ClientCompany").Where("external_id = ?", externalID).First(&user).Error
	return &user, err
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

// CreateOrUpdate creates a new user if not exists, otherwise updates existing user
func (r *clientUserRepository) CreateOrUpdate(user *model.ClientUser) error {
	// Check if user exists by external_id
	var existingUser model.ClientUser
	err := r.db.Where("external_id = ?", user.ExternalID).First(&existingUser).Error

	if err == gorm.ErrRecordNotFound {
		// User doesn't exist, create new
		return r.db.Create(user).Error
	} else if err != nil {
		// Other error occurred
		return err
	}

	// User exists, update existing record
	user.ID = existingUser.ID // Keep the existing ID
	return r.db.Save(user).Error
}

// FindByExternalIDs finds multiple users by their external IDs
func (r *clientUserRepository) FindByExternalIDs(externalIDs []string) ([]model.ClientUser, error) {
	var users []model.ClientUser
	err := r.db.Preload("Client").Preload("ClientCompany").
		Where("external_id IN ?", externalIDs).Find(&users).Error
	return users, err
}
