package psrehandler

import (
	psreservice "dimensy-bridge/internal/service/psre_service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PsreDashboardHandler struct {
	dashboardService psreservice.DashboardService
}

func NewPsreDashboardHandler(dashboardService psreservice.DashboardService) *PsreDashboardHandler {
	return &PsreDashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *PsreDashboardHandler) Certificate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	respBody, status, err := h.dashboardService.GetCertificateDashboard(token)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreDashboardHandler) Document(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	respBody, status, err := h.dashboardService.GetDocumentDashboard(token)
	fmt.Println(respBody)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
