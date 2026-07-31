package dto

import "github.com/google/uuid"

type PsreCreateClientCompanyRequest struct {
	CompanyName     string `json:"companyName" validate:"required"`
	CompanyAddress  string `json:"companyAddress" validate:"required"`
	CompanyIndustry string `json:"companyIndustry" validate:"required"`
	NPWP            string `json:"npwp" validate:"required"`
	NIB             string `json:"nib" validate:"required"`
	PICPhone        string `json:"picPhone" validate:"required"`
	PICName         string `json:"picName" validate:"required"`
	PICEmail        string `json:"picEmail" validate:"required,email"`
}

type PsreRegisterCompanyResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	CompanyID uuid.UUID `json:"companyId"`
}

type PsreAcceptInvitationResponse struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Data    PsreAcceptInvitationResponseData `json:"data"`
}

type PsreAcceptInvitationResponseData struct {
	UserID    uuid.UUID `json:"userId"`
	CompanyID uuid.UUID `json:"companyId"`
}

type PsreInviteClientCompanyRequest struct {
	UserID    uuid.UUID `json:"userId" validate:"required"`
	CompanyID uuid.UUID `json:"companyId" validate:"required"`
	URL       string    `json:"url" validate:"required,url"`
}

type PsreAcceptInvitationClientUserRequest struct {
	Token string `json:"token" binding:"required"`
}
