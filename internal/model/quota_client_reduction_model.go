package model

import "time"

type QuotaClientReduction struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	QuotaClientID int64  `gorm:"not null;index" json:"quota_client_id"`
	Quantity      int    `json:"quantity"`
	LatestQuota   int    `json:"latest_quota"`
	Type          string `gorm:"size:50" json:"type"`
	UsedBy        *int64 `json:"used_by,omitempty"` // bisa FK ke user admin

	// Relasi
	QuotaClient QuotaClient `gorm:"foreignKey:QuotaClientID;references:ID;constraint:OnDelete:CASCADE" json:"quota_client"`
	UsedByUser  *User       `gorm:"foreignKey:UsedBy;references:ID;constraint:OnDelete:CASCADE" json:"used_by_user,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
