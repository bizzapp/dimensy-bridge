package dto

import (
	"time"

	"github.com/google/uuid"
)

type CertificateIssueActiveRequest struct {
	UserID    *uuid.UUID `json:"userId" binding:"required_without=CompanyID"`
	CompanyID *uuid.UUID `json:"companyId" binding:"required_without=UserID"`
}

type CertificateRequestIssueV2Request struct {
	UserCompanyPicDTO
	LivenessCoreRequest
}

type CertificateIssueV2Request struct {
	SignatureID string `json:"signatureId" binding:"required"`
	Status      string `json:"status" binding:"required"`
	CallbackURL string `json:"callbackUrl" binding:"required"`
	RevokeType  int8   `json:"revokeType" binding:"required"`
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

type CertificateResyncRequest struct {
	SN string `json:"sn" binding:"required"`
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
