package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientCompanyInvite struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64          `gorm:"not null" json:"client_id"`
	ClientUserID      int64          `gorm:"not null" json:"client_user_id"`
	ClientCompanyID   int64          `gorm:"not null" json:"client_company_id"`
	ExternalUserID    uuid.UUID      `gorm:"type:uuid;not null" json:"external_user_id"`
	ExternalCompanyID uuid.UUID      `gorm:"type:uuid;not null" json:"external_company_id"`
	IsVerify          bool           `gorm:"default:false" json:"is_verify"`
	VerifyTime        *time.Time     `json:"verify_time,omitempty"`
	CreatedAt         *time.Time     `json:"created_at,omitempty"`
	UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Client        Client         `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	ClientUser    *ClientUser    `gorm:"foreignKey:ClientUserID" json:"client_user,omitempty"`
	ClientCompany *ClientCompany `gorm:"foreignKey:ClientCompanyID" json:"client_company,omitempty"`
}
