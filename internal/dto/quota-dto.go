package dto

type ReductionTimeSeriesRow struct {
	Label             string `json:"label"`
	MasterProductName string `json:"master_product_name"`
	TotalQuantity     int64  `json:"total_quantity"`
}
