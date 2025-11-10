package service

import (
	"bytes"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type WebhookService interface {
	// Define webhook-related service methods here
	SendDocumentNotification(req dto.WebhookDocumentNotificationRequest) error
}

type webhookService struct {
	db                   *gorm.DB
	clientDocumentRepo   repository.ClientDocumentRepository
	clientRequestLogRepo repository.ClientRequestLogRepository
}

func NewWebhookService(db *gorm.DB, clientDocumentRepo repository.ClientDocumentRepository, clientRequestLogRepo repository.ClientRequestLogRepository) WebhookService {
	return &webhookService{
		db:                   db,
		clientDocumentRepo:   clientDocumentRepo,
		clientRequestLogRepo: clientRequestLogRepo,
	}
}

func (s *webhookService) SendDocumentNotification(req dto.WebhookDocumentNotificationRequest) error {
	clientDocument, err := s.clientDocumentRepo.FindByExternalID(req.DocumentID)
	if err != nil {
		return fmt.Errorf("failed to find client document: %w", err)
	}

	// Check if client has callback URL
	if clientDocument.ClientCallbackURL == nil || *clientDocument.ClientCallbackURL == "" {
		return fmt.Errorf("client callback URL is not set for document ID: %s", req.DocumentID)
	}

	clientCallbackURL := *clientDocument.ClientCallbackURL

	// Prepare request payload
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", clientCallbackURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "Dimensy-Bridge-Webhook/1.0")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send the request
	resp, err := client.Do(httpReq)
	if err != nil {
		// Log failed request
		s.logRequest(clientDocument.ClientID, clientCallbackURL, "WEBHOOK_NOTIFICATION", string(payload), "", fmt.Sprintf("HTTP Error: %v", err))
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	// Log the request and response
	s.logRequest(clientDocument.ClientID, clientCallbackURL, "WEBHOOK_NOTIFICATION", string(payload), responseBody.String(), fmt.Sprintf("Status: %d", resp.StatusCode))

	// Update document status if needed
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Update document status based on webhook status
		s.updateDocumentStatus(clientDocument, req.Status)
		return nil
	}

	return fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
}

// Helper method to log requests
func (s *webhookService) logRequest(clientID int64, url, requestType, body, response, status string) {
	requestLog := &model.ClientRequestLog{
		URL:      url,
		Type:     requestType,
		ClientID: &clientID,
		Body:     body,
		Header:   "Content-Type: application/json",
		Response: fmt.Sprintf("%s | %s", status, response),
	}

	// Save log to database (ignore errors for logging)
	s.clientRequestLogRepo.Create(requestLog)
}

// Helper method to update document status
func (s *webhookService) updateDocumentStatus(clientDocument *model.ClientDocument, status string) {
	switch status {
	case "SIGNED":
		clientDocument.Status = model.DOCUMENT_STATUS_SIGNED
	case "ON_PROCESS":
		clientDocument.Status = model.DOCUMENT_STATUS_ON_PROCESS
	case "WAITING":
		clientDocument.Status = model.DOCUMENT_STATUS_WAITING
	default:
		// Keep existing status if unknown status received
		return
	}

	// Update document status (ignore errors for status update)
	s.db.Save(clientDocument)
}
