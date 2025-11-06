package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
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
	token := c.Request.Header.Get("Authorization")

	// Implementation for creating a client in PSRE backend
	var req dto.PsreBackendCreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	respBody, status, err := h.backendSvc.CreateClient(token, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)

}
func (h *PsreBackendHandler) ListClient(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	filter := c.Query("filter") // <--- ambil filter di sini
	page := c.Query("page")
	limit := c.Query("limit")

	data, status, err := h.backendSvc.ListClient(token, filter, page, limit)
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
	token := c.Request.Header.Get("Authorization")

	id := c.Param("id")

	respBody, status, err := h.backendSvc.UpdateClient(id, token, &req)
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
	token := c.Request.Header.Get("Authorization")
	id := c.Param("id")

	respBody, status, err := h.backendSvc.UpdateClientStatus(id, token, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
	// Implementation for updating client status in PSRE backend
}
