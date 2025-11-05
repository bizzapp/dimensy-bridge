package dto

type AddQuotaClientRequest struct {
	MasterProductId int64   `json:"master_product_id" binding:"required"`
	ClientID        int64   `json:"client_id" binding:"required"`
	CreatedBy       int64   `json:"created_by" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
}

type ApproveAddQuotaClientRequest struct {
	QuotaAdditionID int64  `json:"quota_addition_id" binding:"required"`
	ProcessBy       *int64 `json:"process_by" binding:"required"`
}

type UseQuotaClientRequest struct {
	MasterProductId int64   `json:"master_product_id" binding:"required"`
	ClientID        int64   `json:"client_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
}

type FindQuotaClientByClientProductRequest struct {
	MasterProductId int64 `json:"master_product_id" binding:"required"`
	ClientID        int64 `json:"client_id" binding:"required"`
}
