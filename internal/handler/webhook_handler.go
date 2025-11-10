package handler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	"fmt"
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
	var req dto.WebhookDocumentNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request payload",
		})
		return
	}
	err := h.webhookSvc.SendDocumentNotification(req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to process webhook",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "webhook received",
	})
}
