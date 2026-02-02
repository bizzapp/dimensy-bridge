package dto

import "github.com/google/uuid"

// PsreDocumentSingleFileRequest untuk upload single document
type PsreDocumentSingleFileRequest struct {
	CallbackURL string `json:"callbackUrl" binding:"required,url"`
	FileName    string `json:"fileName" binding:"required"`
	Document    string `json:"document" binding:"required,base64"` // Base64 encoded document
}

// PsreDocumentBulkFileRequest untuk upload multiple documents

type PsreDocumentBulkFileRequest struct {
	CallbackURL string                     `json:"callbackUrl" binding:"required,url"`
	Documents   []PsreDocumentBulkFileItem `json:"documents" binding:"required,min=1,dive"`
}

type PsreDocumentBulkFileItem struct {
	FileName string `json:"fileName" binding:"required"`
	Document string `json:"document" binding:"required"` // base64 string
}

// PsreDocumentFileRequest untuk individual document dalam bulk upload
type PsreDocumentFileRequest struct {
	FileName string `json:"fileName" binding:"required"`
	Document string `json:"document" binding:"required,base64"` // Base64 encoded document
}

// PsreDocumentSignRequest untuk request digital signature

type PsreDocumentSignRequest struct {
	DocumentOrGroupID uuid.UUID                  `json:"documentOrGroupId" binding:"required"`
	UserID            *string                    `json:"userId,omitempty"`
	CompanyID         *string                    `json:"companyId,omitempty"`
	Positions         []PsreDocumentSignPosition `json:"positions" binding:"required,min=1,dive"`
}

type PsreDocumentSignPosition struct {
	StampType string   `json:"stampType" binding:"required"` // SIGN / EMETERAI
	Reason    string   `json:"reason" binding:"required"`
	Location  string   `json:"location" binding:"required"`
	X         *float64 `json:"x"` // gunakan pointer agar 0 tidak dianggap "kosong"
	Y         *float64 `json:"y"`
	W         *float64 `json:"w"`
	H         *float64 `json:"h"`
	Page      int      `json:"page" binding:"required"`
	Image     string   `json:"image" binding:"required"` // base64 string
}

// PsreDocumentProcessSignRequest untuk process signature dengan OTP
type PsreDocumentProcessSignRequest struct {
	IsApprove         bool   `json:"isApprove" binding:"required"`
	DocumentOrGroupID string `json:"documentOrGroupId" binding:"required"`
	Otp               string `json:"otp" binding:"required"`
	IP                string `json:"ip" binding:"required,ip"` // validasi IP otomatis
}

// PsreDocumentStampRequest untuk request digital stamp

type PsreDocumentStampRequest struct {
	DocumentOrGroupID uuid.UUID                   `json:"documentOrGroupId" binding:"required"`
	CompanyID         uuid.UUID                   `json:"companyId" binding:"required"`
	UserID            uuid.UUID                   `json:"userId" binding:"required"`
	Positions         []PsreDocumentStampPosition `json:"positions" binding:"required,dive"`
}

// PsreDocumentStampPosition mendefinisikan posisi tanda tangan/stempel dalam dokumen
type PsreDocumentStampPosition struct {
	Reason   string   `json:"reason,omitempty"`
	Location string   `json:"location,omitempty"`
	X        *float64 `json:"x"` // gunakan pointer agar 0 tidak dianggap "kosong"
	Y        *float64 `json:"y"`
	W        *float64 `json:"w"`
	H        *float64 `json:"h"`
	Page     int      `json:"page" binding:"required"`
	Image    string   `json:"image" binding:"required,base64"` // base64 encoded image
}

// PsreDocumentProcessStampRequest untuk process stamping
type PsreDocumentProcessStampRequest struct {
	IsApprove         bool      `json:"isApprove" binding:"required"`
	DocumentOrGroupID uuid.UUID `json:"documentOrGroupId" binding:"required"`
	Otp               string    `json:"otp" binding:"required"`
	IP                string    `json:"ip" binding:"required,ip"` // validasi IP otomatis
}

// PsreDocumentOtpSignRequest untuk request OTP untuk signing
type PsreDocumentOtpSignRequest struct {
	DocumentOrGroupID uuid.UUID `json:"documentOrGroupId" binding:"required"`
	UserID            *string   `json:"userId,omitempty"`    // boleh null
	CompanyID         *string   `json:"companyId,omitempty"` // boleh null
	DocumentType      string    `json:"documentType" binding:"required,oneof=sign stamp"`
}
