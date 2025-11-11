package dto

import "time"

type CertificateIssueActiveRequest struct {
	UserID    *string `json:"userId" binding:"omitempty"`
	CompanyID *string `json:"companyId" binding:"omitempty"`
}

type CertificateRevokeRequest struct {
	UserID    *string `json:"userId" binding:"omitempty"`
	CompanyID *string `json:"companyId" binding:"omitempty"`
}

type CertificateRevokeValidateRequest struct {
	UserID    *string `json:"userId" binding:"omitempty"`
	CompanyID *string `json:"companyId" binding:"omitempty"`
	OTP       string  `json:"otp" binding:"required"`
}
type CertificateActiveResponse struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    CertificateActiveResponseData `json:"data"`
}

type CertificateActiveResponseData struct {
	ID                 string     `json:"id"`
	SerialNumber       string     `json:"serialNumber"`
	Status             string     `json:"status"`
	ExpiredAt          *time.Time `json:"expiredAt"`
	RevokeToken        *string    `json:"revokeToken"`
	RevokeTokenExpired *time.Time `json:"revokeTokenExpired"`
	RevocationReasonID *string    `json:"revocationReasonId"`
	RevokedAt          *time.Time `json:"revokedAt"`
	CreatedAt          *time.Time `json:"createdAt"`
	UpdatedAt          *time.Time `json:"updatedAt"`
	UserID             *string    `json:"userId"`
	CompanyID          *string    `json:"companyId"`
}
