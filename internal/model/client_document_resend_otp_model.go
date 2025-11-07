package model

import (
	"time"

	"gorm.io/gorm"
)

type ClientDocumentResendOtp struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64          `gorm:"not null" json:"client_id"`
	ExternalID        string         `gorm:"size:255;not null" json:"external_id"`
	ExternalUserID    *string        `gorm:"size:255" json:"external_user_id,omitempty"`
	ExternalCompanyID *string        `gorm:"size:255" json:"external_company_id,omitempty"`
	Type              string         `gorm:"size:50" json:"type" binding:"required"`
	CreatedAt         *time.Time     `json:"created_at,omitempty"`
	UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty"`

	Client *Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`
}
