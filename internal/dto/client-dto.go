package dto

type AddQuotaClientRequest struct {
	MasterProductID int64 `json:"master_product_id" binding:"required"`
	ClientID        int64 `json:"client_id" binding:"required"`
	CreatedBy       int64 `json:"created_by" binding:"required"`
	Quantity        int64 `json:"quantity" binding:"required,gt=0"`
}

type ApproveAddQuotaClientRequest struct {
	QuotaAdditionID int64  `json:"quota_addition_id" binding:"required"`
	ProcessBy       *int64 `json:"process_by" binding:"required"`
}

type UseQuotaClientRequest struct {
	MasterProductID int64  `json:"master_product_id" binding:"required"`
	ClientID        int64  `json:"client_id" binding:"required"`
	Quantity        int64  `json:"quantity" binding:"required,gt=0"`
	UsedBy          *int64 `json:"used_by" binding:"omitempty"`
}

type FindQuotaClientByClientProductRequest struct {
	MasterProductID int64 `json:"master_product_id" binding:"required"`
	ClientID        int64 `json:"client_id" binding:"required"`
}
