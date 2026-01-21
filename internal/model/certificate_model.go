package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Certificate struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64          `gorm:"not null;index" json:"client_id"`
	ClienUserID       *int64         `gorm:"index" json:"client_user_id"`
	CompanyID         *int64         `gorm:"index" json:"company_id"`
	ExternalUserID    *uuid.UUID     `gorm:"type:uuid" json:"external_user_id"`
	ExternalCompanyID *uuid.UUID     `gorm:"type:uuid" json:"external_company_id"`
	Status            string         `gorm:"size:50;not null;default:'ACTIVE'" json:"status"`
	SerialNumber      *string        `gorm:"size:255;" json:"serial_number"`
	RevokeToken       *string        `gorm:"size:255" json:"revoke_token"`
	RevokeAt          *time.Time     `json:"revoke_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Client     Client         `gorm:"foreignKey:ClientID" json:"client"`
	ClientUser *ClientUser    `gorm:"foreignKey:ClienUserID" json:"client_user,omitempty"`
	Company    *ClientCompany `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
