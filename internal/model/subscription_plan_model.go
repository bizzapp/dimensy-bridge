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

type SubscriptionPlanDetail struct {
	ID                    int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubscriptionPlanID    int64      `gorm:"not null" json:"subscription_plan_id"`
	MasterProductID       int64      `gorm:"not null" json:"master_product_id"`
	IsUnlimited           bool       `gorm:"default:false" json:"is_unlimited"`
	MaxSingleUpload       *int       `json:"max_single_upload,omitempty"`
	MaxBulkUploadLimitPcs *int       `json:"max_bulk_upload_limit_pcs,omitempty"`
	MaxBulkUploadLimitAll *int       `json:"max_bulk_upload_limit_all,omitempty"`
	MaxBulkUploadCount    *int       `json:"max_bulk_upload_count,omitempty"`
	Quantity              int        `gorm:"default:0" json:"quantity"`
	CreatedAt             *time.Time `json:"created_at,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`

	MasterProduct    MasterProduct    `gorm:"foreignKey:MasterProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"master_product,omitempty"`
	SubscriptionPlan SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscription_plan,omitempty"`
}
