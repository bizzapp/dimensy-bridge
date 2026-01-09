package model

import "time"

type InventoryMasterProductLog struct {
	ID                       int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	InventoryMasterProductID int64     `gorm:"not null;index" json:"inventory_master_product_id"`
	MasterProductID          int64     `gorm:"not null;index" json:"master_product_id"`
	Debit                    int       `gorm:"not null" json:"debit"`
	Credit                   int       `gorm:"not null" json:"credit"`
	PreviousStock            int       `gorm:"not null" json:"previous_stock"`
	CurrentStock             int       `gorm:"not null" json:"current_stock"`
	Time                     time.Time `gorm:"not null" json:"time"`
	Notes                    *string   `gorm:"size:255" json:"notes,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Foreign key relationships
	InventoryMasterProduct *InventoryMasterProduct `gorm:"foreignKey:InventoryMasterProductID;constraint:OnDelete:CASCADE" json:"inventory_master_product,omitempty"`
	MasterProduct          *MasterProduct          `gorm:"foreignKey:MasterProductID;constraint:OnDelete:CASCADE" json:"master_product,omitempty"`
}
