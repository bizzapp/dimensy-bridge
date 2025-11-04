package dto

type CertificateIssueRequest struct {
	UserID    string `json:"userId" binding:"required"`
	CompanyID string `json:"companyId" binding:"omitempty"`
}
