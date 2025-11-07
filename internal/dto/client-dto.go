package dto

type AddQuotaClientRequest struct {
	MasterProductID       int64 `json:"master_product_id" binding:"required"`
	ClientID              int64 `json:"client_id" binding:"required"`
	CreatedBy             int64 `json:"created_by" binding:"required"`
	Quantity              int64 `json:"quantity" binding:"required"`
	IsUnlimited           bool  `json:"is_unlimited"`
	MaxSingleUpload       *int  `json:"max_single_upload,omitempty"`
	MaxBulkUploadLimitPcs *int  `json:"max_bulk_upload_limit_pcs,omitempty"`
	MaxBulkUploadLimitAll *int  `json:"max_bulk_upload_limit_all,omitempty"`
	MaxBulkUploadCount    *int  `json:"max_bulk_upload_count,omitempty"`
}

type AddQuotaClientWithApproveRequest struct {
	MasterProductID       int64  `json:"master_product_id" binding:"required"`
	ClientID              int64  `json:"client_id" binding:"required"`
	CreatedBy             int64  `json:"created_by" binding:"required"`
	Quantity              int64  `json:"quantity" binding:"required"`
	IsUnlimited           bool   `json:"is_unlimited"`
	ProcessBy             *int64 `json:"process_by" binding:"required"`
	MaxSingleUpload       *int   `json:"max_single_upload,omitempty"`
	MaxBulkUploadLimitPcs *int   `json:"max_bulk_upload_limit_pcs,omitempty"`
	MaxBulkUploadLimitAll *int   `json:"max_bulk_upload_limit_all,omitempty"`
	MaxBulkUploadCount    *int   `json:"max_bulk_upload_count,omitempty"`
}

type ApproveAddQuotaClientRequest struct {
	QuotaAdditionID int64  `json:"quota_addition_id" binding:"required"`
	ProcessBy       *int64 `json:"process_by" binding:"required"`
}

type UseQuotaClientRequest struct {
	MasterProductID int64  `json:"master_product_id" binding:"required"`
	ClientID        int64  `json:"client_id" binding:"required"`
	Quantity        int64  `json:"quantity" binding:"required"`
	UsedBy          *int64 `json:"used_by" binding:"omitempty"`
}

type FindQuotaClientByClientProductRequest struct {
	MasterProductID int64 `json:"master_product_id" binding:"required"`
	ClientID        int64 `json:"client_id" binding:"required"`
}
