package handler

import (
	"bytes"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookSvc service.WebhookService
}

func NewWebhookHandler(webhookSvc service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookSvc: webhookSvc,
	}
}

// HandlePSRENotification menerima webhook dari PSRE
func (h *WebhookHandler) HandlePSRENotification(c *gin.Context) {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	// convert ke string untuk disimpan
	bodyString := string(bodyBytes)

	// balikin lagi ke request body supaya bisa dipakai lagi (IMPORTANT)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// simpan log
	h.webhookSvc.StoreWebhookRequestLog(
		c.Request.URL.String(),
		"PSRE_NOTIFICATION",
		bodyString,
		"processing",
		"received",
	)

	var req dto.WebhookDocumentNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}
	err = h.webhookSvc.SendDocumentNotification(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to process webhook",
		})
		return
	}

	c.Status(http.StatusOK)
}
