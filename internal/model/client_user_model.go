package model

import "time"

type ClientUser struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	ClientID        int64     `json:"client_id"`
	ClientCompanyID *int64    `json:"client_company_id"`
	ExternalID      *string   `json:"external_id" gorm:"uniqueIndex"`
	NIK             string    `json:"nik"`
	Name            string    `json:"name"`
	Birthdate       time.Time `json:"birthdate"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	IsWNI           bool      `json:"is_wni"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Client        Client        `json:"client" gorm:"foreignKey:ClientID"`
	ClientCompany ClientCompany `json:"client_company" gorm:"foreignKey:ClientCompanyID"`
}
