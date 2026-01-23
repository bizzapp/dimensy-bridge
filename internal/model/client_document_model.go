package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	DOCUMENT_STATUS_PENDING    = "PENDING"
	DOCUMENT_STATUS_SUCCESS    = "SUCCESS"
	DOCUMENT_STATUS_WAITING    = "WAITING"
	DOCUMENT_STATUS_ON_PROCESS = "ON_PROCCESS"
	DOCUMENT_STATUS_SIGNED     = "SIGNED"
)

type ClientDocument struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64          `gorm:"not null" json:"client_id"`
	ExternalID        *string        `gorm:"size:255" json:"external_id,omitempty"`
	FileName          string         `gorm:"size:255;not null" json:"file_name"`
	Type              string         `gorm:"size:100;not null" json:"type"`
	GroupExternalID   *string        `gorm:"size:255" json:"group_external_id,omitempty"`
	Status            string         `gorm:"size:50;default:'pending'" json:"status"`
	TotalParticipants *int           `json:"total_participants,omitempty"`
	FileSizeKB        *int64         `json:"file_size_kb,omitempty"`
	CallbackURL       *string        `json:"callback_url,omitempty"`
	ClientCallbackURL *string        `json:"client_callback_url,omitempty"`
	CreatedAt         *time.Time     `json:"created_at,omitempty"`
	UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty"`

	Client *Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`
}
