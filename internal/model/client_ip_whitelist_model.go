package model

import "time"

type ClientIPWhitelist struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID    int64  `gorm:"not null;index" json:"client_id"`
	IPAddress   string `gorm:"size:45;not null;index" json:"ip_address"` // IPv4 atau IPv6
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Relasi ke Client
	Client *Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"client,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name
func (ClientIPWhitelist) TableName() string {
	return "client_ip_whitelists"
}
