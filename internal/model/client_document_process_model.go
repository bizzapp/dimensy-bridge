package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ClientDocumentProcessStatusWaiting   = "WAITING"
	ClientDocumentProcessStatusOnProcess = "ON_PROCESS"
	ClientDocumentProcessStatusSigned    = "SIGNED"

	DocumentProcessExpiredHour int = 24
)

type ClientDocumentProcess struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID           int64      `gorm:"not null;index" json:"client_id"`
	ExternalID         string     `gorm:"size:255;not null;uniqueIndex" json:"external_id"`
	ExternalUserID     *string    `gorm:"size:255" json:"external_user_id,omitempty"`
	ExternalCompanyID  *string    `gorm:"size:255" json:"external_company_id,omitempty"`
	Status             string     `gorm:"size:100;not null" json:"status"`
	IsDone             bool       `gorm:"default:false" json:"is_done"`
	IsProcess          bool       `gorm:"default:false" json:"is_process"`
	ExpireTime         *time.Time `json:"expire_time,omitempty"`
	IsNeedReverseStock bool       `gorm:"default:true" json:"is_need_reverse_stock"`

	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"client,omitempty"`

	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
