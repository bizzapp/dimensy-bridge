package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreClientDocumentHandler struct {
	clientDocumentSvc psreservice.ClientDocumentService
}

func NewPsreClientDocumentHandler(clientDocumentSvc psreservice.ClientDocumentService) *PsreClientDocumentHandler {
	return &PsreClientDocumentHandler{
		clientDocumentSvc: clientDocumentSvc,
	}
}

func (h *PsreClientDocumentHandler) Upload(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusUnauthorized)
		c.Data(http.StatusUnauthorized, "application/json", message)
		return
	}
	var req dto.PsreDocumentSingleFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}
	respBody, status, err := h.clientDocumentSvc.UploadSingle(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreClientDocumentHandler) UploadBulk(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusUnauthorized)
		c.Data(http.StatusUnauthorized, "application/json", message)
		return
	}

	var req dto.PsreDocumentBulkFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}
	respBody, status, err := h.clientDocumentSvc.UploadBulk(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreClientDocumentHandler) RequestSign(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.PsreDocumentSignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.clientDocumentSvc.RequestSign(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreClientDocumentHandler) ProcessSign(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusUnauthorized)
		c.Data(http.StatusUnauthorized, "application/json", message)
		return
	}

	var req dto.PsreDocumentProcessSignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}
	respBody, status, err := h.clientDocumentSvc.ProcessSign(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreClientDocumentHandler) RequestStamp(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	var req dto.PsreDocumentStampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}
	respBody, status, err := h.clientDocumentSvc.RequestStamp(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreClientDocumentHandler) ProcessStamp(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusUnauthorized)
		c.Data(http.StatusUnauthorized, "application/json", message)
		return
	}

	var req dto.PsreDocumentProcessStampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// c.JSON(http.StatusBadRequest, err.Error())
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.clientDocumentSvc.ProcessStamp(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreClientDocumentHandler) RequestOtpSign(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusUnauthorized)
		c.Data(http.StatusUnauthorized, "application/json", message)
		return
	}

	var req dto.PsreDocumentOtpSignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.clientDocumentSvc.RequestOtpSign(token, externalID, req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreClientDocumentHandler) PreviewDocument(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		message := utils.ResponseError(err.Error(), http.StatusUnauthorized)
		c.Data(http.StatusUnauthorized, "application/json", message)
		return
	}

	documentID := c.Param("id")

	respBody, status, err := h.clientDocumentSvc.Preview(token, externalID, documentID)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
