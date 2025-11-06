package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreBackendHandler struct {
	backendSvc psreservice.BackendService
}

func NewPsreBackendHandler(backendSvc psreservice.BackendService) *PsreBackendHandler {
	return &PsreBackendHandler{
		backendSvc: backendSvc,
	}
}
func (h *PsreBackendHandler) Login(c *gin.Context) {
	var req dto.PsreBackendLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	respBody, status, err := h.backendSvc.LoginBackend(req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreBackendHandler) CreateClient(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// Implementation for creating a client in PSRE backend
	var req dto.PsreBackendCreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	respBody, status, err := h.backendSvc.CreateClient(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)

}
func (h *PsreBackendHandler) ListClient(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	filter := c.Query("filter") // <--- ambil filter di sini
	page := c.Query("page")
	limit := c.Query("limit")

	data, status, err := h.backendSvc.ListClient(token, externalID, filter, page, limit)
	if err != nil {
		c.Data(status, "application/json", data)
		return
	}
	c.Data(status, "application/json", data)
}
func (h *PsreBackendHandler) UpdateClient(c *gin.Context) {
	// Implementation for updating a client in PSRE backend
	var req dto.PsreBackendCreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	id := c.Param("id")

	respBody, status, err := h.backendSvc.UpdateClient(id, token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreBackendHandler) UpdateClientStatus(c *gin.Context) {
	var req dto.PsreBackendUpdateClientStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	id := c.Param("id")

	respBody, status, err := h.backendSvc.UpdateClientStatus(id, token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
	// Implementation for updating client status in PSRE backend
}
