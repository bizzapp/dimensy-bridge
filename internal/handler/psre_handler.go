package handler

import (
	"dimensy-bridge/internal/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreHandler struct {
	service service.PsreService
}

func NewPsreHandler(s service.PsreService) *PsreHandler {
	return &PsreHandler{s}
}
func (h *PsreHandler) Login(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)

	resp, err := h.service.Login(body, map[string]string{})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

func (h *PsreHandler) ClientCompany(c *gin.Context) {
	// ambil id dari query param atau body (misalnya pakai query param ?id=)
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return
	}

	resp, err := h.service.RegisterClientCompany(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}
