package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

// HandlePSRENotification menerima webhook dari PSRE
func (h *WebhookHandler) HandlePSRENotification(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON payload",
		})
		return
	}
	// TODO: proses payload sesuai kebutuhan
	// contoh: simpan ke DB, kirim event, logging, dll.
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "webhook received",
	})
}
