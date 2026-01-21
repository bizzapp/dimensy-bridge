package repository

import (
	"dimensy-bridge/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientKYCHistoryRepository interface {
	Create(data *model.ClientKYCHistory) error
	CreateTx(tx *gorm.DB, c *model.ClientKYCHistory) error
	UpdateTx(tx *gorm.DB, c *model.ClientKYCHistory) error
	Update(data *model.ClientKYCHistory) error
	FindByID(id int64) (*model.ClientKYCHistory, error)
	FindAll() ([]model.ClientKYCHistory, error)
	Delete(id int64) error
	FindByClientUserID(clientUserID int64) ([]model.ClientKYCHistory, error)
	FindByExternalUserID(externalUserID uuid.UUID) (*model.ClientKYCHistory, error)
	FindByExternalUserIDAndSignature(externalUserID uuid.UUID, signature string) (*model.ClientKYCHistory, error)
	UpdateIsRejectStatus(signatureID string, isReject bool) error
	UpdateIsVerifyStatus(signatureID string, isVerify bool) error
	GetBySignatureID(signatureID string) (*model.ClientKYCHistory, error)
}

type clientKYCHistoryRepository struct {
	db *gorm.DB
}

func NewClientKYCHistoryRepository(db *gorm.DB) ClientKYCHistoryRepository {
	return &clientKYCHistoryRepository{db: db}
}
func (r *clientKYCHistoryRepository) CreateTx(tx *gorm.DB, c *model.ClientKYCHistory) error {
	return tx.Create(c).Error
}

func (r *clientKYCHistoryRepository) UpdateTx(tx *gorm.DB, c *model.ClientKYCHistory) error {
	return tx.Save(c).Error
}
func (r *clientKYCHistoryRepository) GetBySignatureID(signatureID string) (*model.ClientKYCHistory, error) {
	var result model.ClientKYCHistory
	if err := r.db.Where("signature = ?", signatureID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
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

func (r *clientKYCHistoryRepository) FindByExternalUserID(externalUserID uuid.UUID) (*model.ClientKYCHistory, error) {
	var result model.ClientKYCHistory
	if err := r.db.Where("external_user_id = ?", externalUserID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *clientKYCHistoryRepository) FindByExternalUserIDAndSignature(externalUserID uuid.UUID, signature string) (*model.ClientKYCHistory, error) {
	var result model.ClientKYCHistory
	if err := r.db.Where("external_user_id = ? AND signature = ?", externalUserID, signature).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *clientKYCHistoryRepository) UpdateIsRejectStatus(signatureID string, isReject bool) error {
	return r.db.Model(&model.ClientKYCHistory{}).
		Where("signature = ?", signatureID).
		Where("reject_time IS NULL").
		Updates(map[string]interface{}{
			"is_reject":   isReject,
			"reject_time": time.Now(),
		}).Error
}

func (r *clientKYCHistoryRepository) UpdateIsVerifyStatus(signatureID string, isVerify bool) error {
	return r.db.Model(&model.ClientKYCHistory{}).
		Where("signature = ?", signatureID).
		Where("verify_time IS NULL").
		Updates(map[string]interface{}{
			"is_verify":   isVerify,
			"verify_time": time.Now(),
		}).Error
}
