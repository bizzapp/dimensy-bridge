package model

import "time"

type Certificate struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64      `gorm:"not null;index" json:"client_id"`
	UserID            *int64     `gorm:"index" json:"user_id"`
	CompanyID         *int64     `gorm:"index" json:"company_id"`
	ExternalUserID    *string    `gorm:"size:255" json:"external_user_id"`
	ExternalCompanyID *string    `gorm:"size:255" json:"external_company_id"`
	Status            string     `gorm:"size:50;not null;default:'ACTIVE'" json:"status"`
	SerialNumber      *string    `gorm:"size:255;" json:"serial_number"`
	RevokeToken       *string    `gorm:"size:255" json:"revoke_token"`
	RevokeAt          *time.Time `json:"revoke_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Client  Client         `gorm:"foreignKey:ClientID" json:"client"`
	User    *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *ClientCompany `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
