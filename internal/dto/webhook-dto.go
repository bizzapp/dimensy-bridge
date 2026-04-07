package dto

import "github.com/google/uuid"

type WebhookDocumentNotificationRequest struct {
	DocumentID   uuid.UUID `json:"documentId" binding:"required"`
	GroupID      uuid.UUID `json:"groupId,omitempty"`
	UserID       *string   `json:"userId" binding:"omitempty"`
	CompanyID    *string   `json:"companyId" binding:"omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
	Status       *string   `json:"status"`
	SignedAT     *string   `json:"signedAt,omitempty"`
}
