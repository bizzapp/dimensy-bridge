package psrehandler

import (
	psreService "dimensy-bridge/internal/service/psre_service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreClientHandler struct {
	psreClientSvc psreService.ClientService
}

func NewPsreClientHandler(psreClientSvc psreService.ClientService) *PsreClientHandler {
	return &PsreClientHandler{psreClientSvc: psreClientSvc}
}

func (h *PsreClientHandler) Login(c *gin.Context) {

	body, _ := io.ReadAll(c.Request.Body)

	resp, err := h.psreClientSvc.Login(body)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}
