package repository

import (
	"dimensy-bridge/internal/model"
	"time"

	"gorm.io/gorm"
)

type ClientCompanyInviteRepository interface {
	Create(invite *model.ClientCompanyInvite) error
	FindByID(id int64) (*model.ClientCompanyInvite, error)
	FindByExternal(userID, companyID string) (*model.ClientCompanyInvite, error)
	VerifyInvite(id int64) error

	CreateTx(tx *gorm.DB, c *model.ClientCompanyInvite) error
}

type clientCompanyInviteRepository struct {
	db *gorm.DB
}

func NewClientCompanyInviteRepository(db *gorm.DB) ClientCompanyInviteRepository {
	return &clientCompanyInviteRepository{db}
}

func (r *clientCompanyInviteRepository) CreateTx(tx *gorm.DB, c *model.ClientCompanyInvite) error {
	return tx.Create(c).Error
}
func (r *clientCompanyInviteRepository) Create(invite *model.ClientCompanyInvite) error {
	return r.db.Create(invite).Error
}

func (r *clientCompanyInviteRepository) FindByID(id int64) (*model.ClientCompanyInvite, error) {
	var invite model.ClientCompanyInvite
	if err := r.db.First(&invite, id).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *clientCompanyInviteRepository) FindByExternal(userID, companyID string) (*model.ClientCompanyInvite, error) {
	var invite model.ClientCompanyInvite
	err := r.db.Where("external_user_id = ? AND external_company_id = ?", userID, companyID).First(&invite).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *clientCompanyInviteRepository) VerifyInvite(id int64) error {
	return r.db.Model(&model.ClientCompanyInvite{}).
		Where("id = ?", id).
		UpdateColumns(map[string]interface{}{
			"is_verify":   true,
			"verify_time": time.Now(),
		}).Error
}
