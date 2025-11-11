package dto

type CertificateIssueRequest struct {
	UserID    *string `json:"userId" binding:"omitempty"`
	CompanyID *string `json:"companyId" binding:"omitempty"`
}

type CertificateActiveRequest struct {
	UserID    string `json:"userId" binding:"omitempty"`
	CompanyID string `json:"companyId" binding:"omitempty"`
}

type CertificateRevokeRequest struct {
	UserID    string `json:"userId" binding:"omitempty"`
	CompanyID string `json:"companyId" binding:"omitempty"`
}

type CertificateRevokeValidateRequest struct {
	UserID    string `json:"userId" binding:"omitempty"`
	CompanyID string `json:"companyId" binding:"omitempty"`
	OTP       string `json:"otp" binding:"required"`
}
