package model

import "time"

type Client struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyName string `gorm:"size:255;not null" json:"company_name"`
	PicName     string `gorm:"size:255;not null" json:"pic_name"`
	UserID      int64  `gorm:"not null;index" json:"user_id"`

	// Relasi ke User
	User       User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	ClientPsre *ClientPsre `gorm:"foreignKey:ClientID;constraint:OnDelete:SET NULL;" json:"client_psre,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
