package model

import (
	"time"

	"gorm.io/gorm"
)

type InventoryMasterProduct struct {
	ID              int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	MasterProductID int64   `gorm:"not null;index" json:"master_product_id"`
	VendorName      string  `gorm:"size:255;not null" json:"vendor_name"`
	Price           float64 `gorm:"not null" json:"price"`
	Quantity        int     `gorm:"not null" json:"quantity"`
	CurrentStock    int     `gorm:"not null" json:"current_stock"`
	IsProcessed     bool    `gorm:"default:false;not null" json:"is_processed"`
	IsPriorityUse   bool    `gorm:"default:false;not null" json:"is_priority_use"`

	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign key relationship
	MasterProduct *MasterProduct `gorm:"foreignKey:MasterProductID;constraint:OnDelete:CASCADE" json:"master_product,omitempty"`
}
