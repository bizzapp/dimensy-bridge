package dto

type PsreCreateClientCompanyRequest struct {
	CompanyName     string `json:"companyName" validate:"required"`
	CompanyAddress  string `json:"companyAddress" validate:"required"`
	CompanyIndustry string `json:"companyIndustry" validate:"required"`
	NPWP            string `json:"npwp" validate:"required"`
	NIB             string `json:"nib" validate:"required"`
	PICName         string `json:"picName" validate:"required"`
	PICEmail        string `json:"picEmail" validate:"required,email"`
}

type PsreRegisterCompanyResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	CompanyID string `json:"companyId"`
}

type PsreInviteClientCompanyRequest struct {
	UserID    string `json:"userId" validate:"required"`
	CompanyID string `json:"companyId" validate:"required"`
	URL       string `json:"url" validate:"required,url"`
}

type PsreAcceptInvitationClientUserRequest struct {
	Token string `json:"token" binding:"required"`
}
