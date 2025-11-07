package model

import (
	"time"

	"gorm.io/gorm"
)

type ClientDocumentProcessDetail struct {
	ID                      int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientDocumentProcessID int64   `gorm:"not null;index" json:"client_document_process_id"`
	ClientID                int64   `gorm:"not null;index" json:"client_id"`
	Type                    string  `gorm:"size:100;not null" json:"type"`
	Reason                  string  `gorm:"size:255" json:"reason,omitempty"`
	Location                string  `gorm:"size:255" json:"location,omitempty"`
	X                       float64 `json:"x"`
	Y                       float64 `json:"y"`
	W                       float64 `json:"w"`
	H                       float64 `json:"h"`
	Page                    int     `json:"page"`
	ImageFileSizeKB         *int64  `json:"image_file_size_kb,omitempty"`

	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"client,omitempty"`

	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
