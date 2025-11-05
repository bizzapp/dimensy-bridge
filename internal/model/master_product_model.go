package model

import "time"

type MasterProduct struct {
	ID           int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"size:255;not null" json:"name"`
	Code         string  `gorm:"size:100;unique;not null" json:"code"`
	Quantity     int     `json:"quantity"`
	CurrentStock int     `json:"current_stock"`
	IsUnlimited  bool    `gorm:"default:false" json:"is_unlimited"`
	Sort         int     `gorm:"default:0" json:"sort"`
	Icon         *string `gorm:"size:255" json:"icon,omitempty"`
	Notes        *string `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
