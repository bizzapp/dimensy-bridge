package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreCertificateHandler struct {
	certificateService   psreservice.CertificateService
	certificateV2Service psreservice.CertificateV2Service
}

func NewPsreCertificateHandler(certificateService psreservice.CertificateService, certificateV2Service psreservice.CertificateV2Service) *PsreCertificateHandler {
	return &PsreCertificateHandler{
		certificateService:   certificateService,
		certificateV2Service: certificateV2Service,
	}
}

func (h *PsreCertificateHandler) Issue(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateIssueActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
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
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateIssueActiveRequest
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
	externalID, token, err := utils.ValidateExternalID(c)
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
	externalID, token, err := utils.ValidateExternalID(c)
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

func (h *PsreCertificateHandler) RequestIssueV2(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateRequestIssueV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.certificateV2Service.RequestIssueV2(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreCertificateHandler) IssueV2(c *gin.Context) {

	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.CertificateIssueV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.certificateV2Service.IssueV2(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreCertificateHandler) RevokeRequestV2(c *gin.Context) {
}
func (h *PsreCertificateHandler) RevokeV2(c *gin.Context) {
}
