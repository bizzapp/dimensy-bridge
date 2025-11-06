package psrehandler

import (
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

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
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	respBody, status, err := h.dashboardService.GetCertificateDashboard(token, externalID)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreDashboardHandler) Document(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	respBody, status, err := h.dashboardService.GetDocumentDashboard(token, externalID)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
