package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreCertificateHandler struct {
	certificateService psreservice.CertificateService
}

func NewPsreCertificateHandler(certificateService psreservice.CertificateService) *PsreCertificateHandler {
	return &PsreCertificateHandler{
		certificateService: certificateService,
	}
}

func (h *PsreCertificateHandler) Issue(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.certificateService.Issue(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreCertificateHandler) Active(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.certificateService.Active(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreCertificateHandler) RevokeRequest(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateRevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.certificateService.RevokeRequest(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreCertificateHandler) Revoke(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateRevokeValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.certificateService.Revoke(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
