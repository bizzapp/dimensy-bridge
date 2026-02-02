package model

import "time"

type PsreLog struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	Description *string    `json:"description" gorm:"column:description"`
	JsonHeader  *string    `json:"json_header" gorm:"column:json_header;type:text"`
	JsonContent string     `json:"json_content" gorm:"column:json_content;type:text"`
	FullURL     *string    `json:"full_url" gorm:"column:full_url;type:text"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName specifies the table name for PsreLog
func (PsreLog) TableName() string {
	return "psre_logs"
}
