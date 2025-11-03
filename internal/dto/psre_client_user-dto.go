package dto

import (
	"encoding/json"
	"time"
)

// ClientUserRequest digunakan untuk Create & Update
type ClientUserRequest struct {
	NIK       string     `json:"nik" binding:"required"`
	FullName  string     `json:"fullName" binding:"required"`
	BirthDate CustomDate `json:"birthDate" binding:"required"`
	Email     string     `json:"email" binding:"required,email"`
	Phone     string     `json:"phone" binding:"required"`
	IsWNI     bool       `json:"isWni" binding:"required"`
	URL       string     `json:"url" binding:"omitempty,url"`
	CompanyID *string    `json:"companyId,omitempty"` // nullable
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
