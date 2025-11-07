package model

import "time"

type SubscriptionPlan struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	ClientID    *int64     `json:"client_id,omitempty"`
	IsDefault   bool       `gorm:"default:false" json:"is_default"`
	IsUnlimited bool       `gorm:"default:false" json:"is_unlimited"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`

	Client  *Client                  `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	Details []SubscriptionPlanDetail `gorm:"foreignKey:SubscriptionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"details,omitempty"`
}
