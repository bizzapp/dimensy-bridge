package handler

import (
	"dimensy-bridge/internal/service"
	psreintegration "dimensy-bridge/internal/service/psre_integration"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PsreHandler struct {
	service service.PsreService

	clientService psreintegration.ClientService
}

func NewPsreHandler(clientService psreintegration.ClientService) *PsreHandler {
	return &PsreHandler{clientService: clientService}
}
func (h *PsreHandler) Login(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)

	resp, err := h.clientService.Login(body)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}


func (h *PsreHandler) CreateClientCompany(c *gin.Context) {

	// 🔹 Ambil token dari header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
		return
	}

	// Format header: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
		return
	}

	token := parts[1]

	// 🔹 Verifikasi JWE token
	data, err := utils.VerifyJWE(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	response.JSON(c, http.StatusOK, "Client company berhasil dibuat", data, nil)

	return

	// ambil id dari query param atau body (misalnya pakai query param ?id=)
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return
	}

	resp, err := h.service.CreateClientCompany(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}
