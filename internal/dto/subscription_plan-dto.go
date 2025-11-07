package dto

type CreateSubscriptionPlanRequest struct {
	ClientID           int64 `json:"client_id" binding:"required"`
	SubscriptionPlanID int64 `json:"subscription_plan_id" binding:"required"`
	CreatedBy          int64 `json:"created_by" binding:"required"`
}
