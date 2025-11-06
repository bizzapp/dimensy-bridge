package dto

type PsreBackendCreateClientRequest struct {
	ClientName  string `json:"clientName" binding:"required"`
	PicName     string `json:"picName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=10"`
	ExpiredDate string `json:"expiredDate" binding:"required,datetime=2006-01-02"`
}

type PsreBackendUpdateClientStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE BANNED"`
}
