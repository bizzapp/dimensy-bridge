package model

import (
	"time"
)

type User struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	Email           *string    `gorm:"size:255;unique;" json:"email"`
	Username        *string    `gorm:"size:255;unique;" json:"username"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Password        *string    `gorm:"size:255;" json:"password,omitempty"`
	RememberToken   *string    `gorm:"size:100" json:"remember_token,omitempty"`
	Role            string     `gorm:"size:255;default:admin;not null" json:"role"`
	Image           *string    `gorm:"size:255;index" json:"image,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
