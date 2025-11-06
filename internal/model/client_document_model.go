package model

import "time"

const (
	DOCUMENT_STATUS_PENDING    = "PENDING"
	DOCUMENT_STATUS_APPROVED   = "APPROVED"
	DOCUMENT_STATUS_REJECTED   = "REJECTED"
	DOCUMENT_STATUS_PROCESSING = "PROCESSING"
)

type ClientDocument struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64      `gorm:"not null" json:"client_id"`
	FileName          string     `gorm:"size:255;not null" json:"file_name"`
	Type              string     `gorm:"size:100;not null" json:"type"`
	ExternalID        *string    `gorm:"size:255" json:"external_id,omitempty"`
	GroupExternalID   *string    `gorm:"size:255" json:"group_external_id,omitempty"`
	Status            string     `gorm:"size:50;default:'pending'" json:"status"`
	TotalParticipants *int       `json:"total_participants,omitempty"`
	FileSizeKB        *int64     `json:"file_size_kb,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`

	Client *Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`
}
