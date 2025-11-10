package dto

type WebhookDocumentNotificationRequest struct {
	DocumentID string `json:"documentId" binding:"required"`
	Status     string `json:"status" binding:"required"`
	SignedAT   string `json:"signedAt,omitempty"`
}
