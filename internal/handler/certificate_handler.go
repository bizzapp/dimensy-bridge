package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CertificateHandler struct {
	service service.CertificateService
}

func NewCertificateHandler(s service.CertificateService) *CertificateHandler {
	return &CertificateHandler{service: s}
}

func (h *CertificateHandler) RegisterRoutes(r *gin.RouterGroup) {
	certs := r.Group("/certificates")
	{
		certs.POST("", h.Create)
		certs.GET("/:id", h.GetByID)
		certs.GET("/client/:clientID", h.GetByClientID)
		certs.POST("/:id/revoke", h.Revoke)
	}
}

func (h *CertificateHandler) Create(c *gin.Context) {
	var req model.Certificate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCertificate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *CertificateHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	cert, err := h.service.GetCertificateByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cert)
}

func (h *CertificateHandler) GetByClientID(c *gin.Context) {
	clientID, _ := strconv.ParseUint(c.Param("clientID"), 10, 64)
	certs, err := h.service.GetCertificatesByClientID(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, certs)
}

func (h *CertificateHandler) Revoke(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.RevokeCertificate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "certificate revoked"})
}
