package model

import "time"

type QuotaClient struct {
	ID              int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MasterProductID int64 `gorm:"not null;index" json:"master_product_id"`
	Quantity        int   `json:"quantity"`
	CurrentQuota    int   `json:"current_quota"`
	ClientID        int64 `gorm:"not null;index" json:"client_id"`

	// Relasi
	MasterProduct MasterProduct `gorm:"foreignKey:MasterProductID" json:"master_product"`
	Client        Client        `gorm:"foreignKey:ClientID" json:"client"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
