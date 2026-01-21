package dto

import "github.com/google/uuid"

type WebhookDocumentNotificationRequest struct {
	DocumentID uuid.UUID `json:"documentId" binding:"required"`
	Status     string    `json:"status" binding:"required"`
	SignedAT   string    `json:"signedAt,omitempty"`
}
