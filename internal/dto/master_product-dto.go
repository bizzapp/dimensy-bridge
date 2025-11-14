package dto

import "time"

type MasterProductHistoryItem struct {
	ID              int64      `json:"id"`
	Quantity        int64      `json:"quantity"`
	MasterProductID int64      `json:"master_product_id"`
	Type            string     `json:"type"`
	Direction       string     `json:"direction"` // ADDITION / REDUCTION
	CreatedAt       *time.Time `json:"created_at"`
}
