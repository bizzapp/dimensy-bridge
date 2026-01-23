package dto

import "github.com/google/uuid"

type WebhookDocumentNotificationRequest struct {
	DocumentID   uuid.UUID  `json:"documentId" binding:"required"`
	UserID       *uuid.UUID `json:"userId" binding:"omitempty"`
	CompanyID    *uuid.UUID `json:"companyId" binding:"omitempty"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
	Status       string     `json:"status" binding:"required"`
	SignedAT     string     `json:"signedAt,omitempty"`
}
