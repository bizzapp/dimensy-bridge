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
