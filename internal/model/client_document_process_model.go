package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ClientDocumentProcessStatusWaiting   = "WAITING"
	ClientDocumentProcessStatusOnProcess = "ON_PROCESS"
	ClientDocumentProcessStatusSigned    = "SIGNED"

	TypeSignMeterai = "SIGN_METERAI"
	TypeStamp       = "STAMP"

	StampTypeMeterai = "EMETERAI"
	StampTypeSign    = "SIGN"

	DocumentProcessExpiredHour int = 24
)

type ClientDocumentProcess struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientDocumentID   int64      `gorm:"not null;index" json:"client_document_id"`
	ClientID           int64      `gorm:"not null;index" json:"client_id"`
	Type               string     `gorm:"size:50;not null;default:'SIGN_METERAI'" json:"type"`
	ExternalID         uuid.UUID  `gorm:"type:uuid;not null" json:"external_id"`
	ExternalGroupID    *uuid.UUID `gorm:"type:uuid" json:"group_id,omitempty"`
	ExternalUserID     *uuid.UUID `gorm:"type:uuid" json:"external_user_id,omitempty"`
	ClientUserID       *int64     `gorm:"index" json:"client_user_id,omitempty"`
	ExternalCompanyID  *uuid.UUID `gorm:"type:uuid" json:"external_company_id,omitempty"`
	ClientCompanyID    *int64     `gorm:"index" json:"client_company_id,omitempty"`
	Status             string     `gorm:"size:100;not null" json:"status"`
	IsDone             bool       `gorm:"default:false" json:"is_done"`
	IsProcess          bool       `gorm:"default:false" json:"is_process"`
	ExpireTime         *time.Time `json:"expire_time,omitempty"`
	IsNeedReverseStock bool       `gorm:"default:true" json:"is_need_reverse_stock"`

	Client         *Client         `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"client,omitempty"`
	ClientDocument *ClientDocument `gorm:"foreignKey:ClientDocumentID;constraint:OnDelete:CASCADE;" json:"client_document,omitempty"`
	ClientUser     *ClientUser     `gorm:"foreignKey:ClientUserID;constraint:OnDelete:SET NULL;" json:"client_user,omitempty"`
	ClientCompany  *ClientCompany  `gorm:"foreignKey:ClientCompanyID;constraint:OnDelete:SET NULL;" json:"client_company,omitempty"`

	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
