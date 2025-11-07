package model

import (
	"time"

	"gorm.io/gorm"
)

type ClientHasSubscriptionPlan struct {
	ID                 int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID           int64          `gorm:"not null" json:"client_id"`
	SubscriptionPlanID int64          `gorm:"not null" json:"subscription_plan_id"`
	CreatedBy          int64          `json:"created_by" binding:"required"`
	ProcessBy          *int64         `json:"process_by,omitempty"`
	ProcessTime        *time.Time     `json:"process_time,omitempty"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	ExpiredDate        *time.Time     `json:"expired_date,omitempty"`
	CreatedAt          *time.Time     `json:"created_at,omitempty"`
	UpdatedAt          *time.Time     `json:"updated_at,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Client           Client           `gorm:"foreignKey:ClientID" json:"client"`
	SubscriptionPlan SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID" json:"subscription_plan"`
}
