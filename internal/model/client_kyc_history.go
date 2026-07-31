package model

import (
	"time"

	"github.com/google/uuid"
)

type TypeLiveness string

const (
	TypeLivenessGenerateCertificate TypeLiveness = "GENERATE_CERTIFICATE"
	TypeLivenessRevokeCertificate   TypeLiveness = "REVOKE_CERTIFICATE"
	TypeLivenessUserValidation      TypeLiveness = "USER_VALIDATION"
)

type ClientKYCHistory struct {
	ID                int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID          int64        `gorm:"not null;index" json:"client_id"`
	ClientUserID      int64        `gorm:"not null;index" json:"client_user_id"`
	ExternalUserID    uuid.UUID    `gorm:"type:uuid;not null" json:"external_user_id"`
	Signature         string       `gorm:"size:255;unique;not null" json:"signature"` // 🔹 unique & required string
	IsVerify          bool         `gorm:"default:false" json:"is_verify"`
	IsReject          bool         `gorm:"default:false" json:"is_reject"`
	RejectTime        *time.Time   `json:"reject_time"`
	Count             int          `gorm:"default:0" json:"count"`
	VerifyTime        *time.Time   `json:"verify_time,omitempty"`
	Type              TypeLiveness `gorm:"size:50;not null;default:'USER_VALIDATION'" json:"type"`
	CallbackURL       *string      `json:"callback_url,omitempty"`
	ClientCallbackURL *string      `json:"client_callback_url,omitempty"`
	CreatedAt         *time.Time   `json:"created_at,omitempty"`
	UpdatedAt         *time.Time   `json:"updated_at,omitempty"`

	Client     *Client     `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	ClientUser *ClientUser `gorm:"foreignKey:ClientUserID" json:"client_user,omitempty"`
}
