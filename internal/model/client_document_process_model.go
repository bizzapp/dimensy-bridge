package model

import "time"

const (
	ClientDocumentProcessStatusWaiting   = "WAITING"
	ClientDocumentProcessStatusOnProcess = "ON_PROCESS"
	ClientDocumentProcessStatusSigned    = "SIGNED"
)

type ClientDocumentProcess struct {
	ID                int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64   `gorm:"not null;index" json:"client_id"`
	ExternalID        string  `gorm:"size:255;not null;uniqueIndex" json:"external_id"`
	ExternalUserID    *string `gorm:"size:255" json:"external_user_id,omitempty"`
	ExternalCompanyID *string `gorm:"size:255" json:"external_company_id,omitempty"`
	Status            string  `gorm:"size:100;not null" json:"status"`

	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"client,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
