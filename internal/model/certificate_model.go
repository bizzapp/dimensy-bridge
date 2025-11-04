package model

import "time"

type Certificate struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          uint64     `gorm:"not null;index" json:"client_id"`
	UserID            *uint64    `gorm:"index" json:"user_id"`
	CompanyID         *uint64    `gorm:"index" json:"company_id"`
	ExternalUserID    *string    `gorm:"size:255" json:"external_user_id"`
	ExternalCompanyID *string    `gorm:"size:255" json:"external_company_id"`
	Status            string     `gorm:"size:50;not null;default:'ACTIVE'" json:"status"`
	SerialNumber      string     `gorm:"size:255;uniqueIndex;not null" json:"serial_number"`
	RevokeToken       *string    `gorm:"size:255" json:"revoke_token"`
	RevokeAt          *time.Time `json:"revoke_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
