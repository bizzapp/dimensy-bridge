package model

import (
	"time"
)

type TokenBlacklist struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklists"
}
