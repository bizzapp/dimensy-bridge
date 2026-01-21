package dto

import (
	"time"

	"github.com/google/uuid"
)

type CertificateIssueActiveRequest struct {
	UserID    *uuid.UUID `json:"userId" binding:"required_without=CompanyID"`
	CompanyID *uuid.UUID `json:"companyId" binding:"required_without=UserID"`
}

type CertificateRevokeRequest struct {
	UserID    *uuid.UUID `json:"userId" binding:"omitempty" validate:"required_without=CompanyID"`
	CompanyID *uuid.UUID `json:"companyId" binding:"omitempty" validate:"required_without=UserID"`
}

type CertificateRevokeValidateRequest struct {
	UserID    *uuid.UUID `json:"userId" binding:"omitempty" validate:"required_without=CompanyID"`
	CompanyID *uuid.UUID `json:"companyId" binding:"omitempty" validate:"required_without=UserID"`
	OTP       string     `json:"otp" binding:"required"`
}

type CertificateRevokeValidateResponse struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
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
