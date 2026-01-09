package dto

type CreateInventoryMasterProductRequest struct {
	ID              *int64  `json:"id,omitempty"` // Optional: if null = create, if provided = update
	MasterProductID int64   `json:"master_product_id" binding:"required"`
	VendorName      string  `json:"vendor_name" binding:"required"`
	Price           float64 `json:"price" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required"`
	IsPriorityUse   bool    `json:"is_priority_use"`
}

type UpdateInventoryMasterProductRequest struct {
	VendorName    string  `json:"vendor_name" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	IsPriorityUse bool    `json:"is_priority_use"`
}

type AdjustStockRequest struct {
	Adjustment int    `json:"adjustment" binding:"required"` // positive = stock in (debit), negative = stock out (credit)
	Notes      string `json:"notes"`
}

type InventoryMasterProductResponse struct {
	ID              int64   `json:"id"`
	MasterProductID int64   `json:"master_product_id"`
	VendorName      string  `json:"vendor_name"`
	Price           float64 `json:"price"`
	Quantity        int     `json:"quantity"`
	CurrentStock    int     `json:"current_stock"`
	IsProcessed     bool    `json:"is_processed"`
	IsPriorityUse   bool    `json:"is_priority_use"`
	CreatedAt       string  `json:"created_at,omitempty"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
	DeletedAt       string  `json:"deleted_at,omitempty"`
}

type InventoryMasterProductLogResponse struct {
	ID                       int64  `json:"id"`
	InventoryMasterProductID int64  `json:"inventory_master_product_id"`
	MasterProductID          int64  `json:"master_product_id"`
	Debit                    int    `json:"debit"`
	Credit                   int    `json:"credit"`
	PreviousStock            int    `json:"previous_stock"`
	CurrentStock             int    `json:"current_stock"`
	Time                     string `json:"time"`
	Notes                    string `json:"notes,omitempty"`
	CreatedAt                string `json:"created_at,omitempty"`
	UpdatedAt                string `json:"updated_at,omitempty"`
}

type LowStockItemResponse struct {
	ID              int64   `json:"id"`
	MasterProductID int64   `json:"master_product_id"`
	VendorName      string  `json:"vendor_name"`
	CurrentStock    int     `json:"current_stock"`
	Price           float64 `json:"price"`
	TotalValue      float64 `json:"total_value"`
}

type TotalInventoryValueResponse struct {
	TotalValue float64 `json:"total_value"`
	Currency   string  `json:"currency"`
}
