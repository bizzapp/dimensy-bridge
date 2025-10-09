package model

import "time"

type ClientCompany struct {
	ID         int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID   int64   `gorm:"index;not null" json:"client_id"`
	Name       string  `gorm:"size:255;not null" json:"name"`
	Address    string  `gorm:"size:500" json:"address"`
	Industry   string  `gorm:"size:255" json:"industry"`
	NPWP       string  `gorm:"size:100" json:"npwp"`
	NIB        string  `gorm:"size:100" json:"nib"`
	PICName    string  `gorm:"size:255" json:"pic_name"`
	PICEmail   string  `gorm:"size:255" json:"pic_email"`
	ExternalID *string `gorm:"size:255" json:"external_id,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
