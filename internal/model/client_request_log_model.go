package model

import "time"

type ClientRequestLog struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	URL      string `gorm:"size:512" json:"url"`
	Type     string `gorm:"size:50" json:"type"`
	ClientID *int64 `gorm:"index" json:"client_id,omitempty"`
	Body     string `gorm:"type:text" json:"body"`
	Header   string `gorm:"type:text" json:"header"`
	Response string `gorm:"type:text" json:"response"`

	Client *Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
