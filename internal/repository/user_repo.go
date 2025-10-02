package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(limit, offset int, filters map[string]interface{}) ([]model.User, int64, error)
	FindByID(id int64) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

	// apply filters (example: name, email, role)
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Delete(&model.User{}, id).Error
}
