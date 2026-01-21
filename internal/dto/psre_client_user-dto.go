package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ClientUserActivateRequest struct {
	ActivationToken string `json:"activationToken" binding:"required"`
}

// ClientUserRequest digunakan untuk Create & Update

type ClientUserResendActivationRequest struct {
	Email string `json:"email" binding:"required,email"`
	URL   string `json:"url" binding:"required,url"`
}

type ClientUserRequestPhoneActivationRequest struct {
	Phone  string `json:"phone" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}

type ClientUserPhoneActivationRequest struct {
	UserID string `json:"userId" binding:"required"`
	OTP    string `json:"otp" binding:"required"`
}

type ClientUserKYCRequest struct {
	UserID     uuid.UUID `json:"userId" binding:"required"`
	SuccessURL string    `json:"successUrl" binding:"required,url"`
	FailedURL  string    `json:"failedUrl" binding:"required,url"`
}

type ClientUserVerifyKYCRequest struct {
	SignatureID string `json:"signatureId" binding:"required"`
}
type ClientUserRequest struct {
	NIK       string     `json:"nik" binding:"required"`
	FullName  string     `json:"fullName" binding:"required"`
	BirthDate CustomDate `json:"birthDate" binding:"required"`
	Email     string     `json:"email" binding:"required,email"`
	Phone     string     `json:"phone" binding:"required"`
	IsWNI     bool       `json:"isWni" binding:"required"`
	URL       string     `json:"url" binding:"omitempty,url"`
	CompanyID *uuid.UUID `json:"companyId,omitempty"` // nullable
}

// CustomDate bisa mem-parse "1994-08-25" langsung ke time.Time
type CustomDate struct {
	time.Time
}

const dateLayout = "2006-01-02"

// UnmarshalJSON dipanggil otomatis ketika parsing JSON
func (d *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	// hapus tanda kutip dari JSON string
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// parse ke time.Time sesuai layout kita
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON agar saat return JSON, tetap "YYYY-MM-DD"
func (d CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateLayout))
}

// Response DTOs for Sync functionality
type PsreClientUserSyncResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    []PsreClientUserData `json:"data"`
}

type PsreClientUserData struct {
	ID          string                `json:"id"`
	DigitalID   string                `json:"digitalID"`
	NIK         string                `json:"nik"`
	FullName    string                `json:"fullName"`
	BirthDate   string                `json:"birthDate"`
	Email       string                `json:"email"`
	Phone       string                `json:"phone"`
	CreatedAt   string                `json:"createdAt"`
	UserCompany []PsreUserCompanyData `json:"userCompany"`
}

type PsreUserCompanyData struct {
	CompanyID string `json:"companyId"`
	// Add other company fields if needed
}
