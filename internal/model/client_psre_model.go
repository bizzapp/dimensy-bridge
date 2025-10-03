package model

import "time"

type ClientPsre struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID   int64     `gorm:"not null;index;unique" json:"client_id"`
	ExternalID string    `gorm:"size:255;not null" json:"external_id"`
	ExpireDate time.Time `json:"expire_date"`

	// relasi
	Client Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"client"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
