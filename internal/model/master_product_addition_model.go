package model

import "time"

type MasterProductAddition struct {
	ID              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	MasterProductID int64  `gorm:"not null;index" json:"master_product_id"`
	Quantity        int64  `json:"quantity"`
	LatestQuota     int64  `json:"latest_quota"`
	Type            string `gorm:"size:50" json:"type"`
	CreatedBy       int64  `json:"created_by"`
	ProcessBy       *int64 `json:"process_by,omitempty"`

	IsProcess bool `gorm:"default:false" json:"is_process"`

	// Relasi
	CreatedByUser *User         `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:CASCADE" json:"created_by_user,omitempty"`
	MasterProduct MasterProduct `gorm:"foreignKey:MasterProductID;references:ID;constraint:OnDelete:CASCADE" json:"master_product"`
	ProcessByUser *User         `gorm:"foreignKey:ProcessBy;references:ID;constraint:OnDelete:CASCADE" json:"process_by_user,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
