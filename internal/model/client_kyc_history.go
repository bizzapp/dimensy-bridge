package model

import "time"

type ClientKYCHistory struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID     int64      `gorm:"not null;index" json:"client_id"`
	ClientUserID int64      `gorm:"not null;index" json:"client_user_id"`
	ExternalID   string     `gorm:"size:255;not null" json:"external_id"`
	Signature    string     `gorm:"size:255;unique;not null" json:"signature"` // 🔹 unique & required string
	IsVerify     bool       `gorm:"default:false" json:"is_verify"`
	VerifyTime   *time.Time `json:"verify_time,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`

	Client     *Client     `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	ClientUser *ClientUser `gorm:"foreignKey:ClientUserID" json:"client_user,omitempty"`
}
