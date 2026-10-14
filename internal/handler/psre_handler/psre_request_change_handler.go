package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreRequestChangeHandler struct {
	requestChangeSvc psreservice.RequestChangeService
}

func NewPsreRequestChangeHandler(requestChangeSvc psreservice.RequestChangeService) *PsreRequestChangeHandler {
	return &PsreRequestChangeHandler{
		requestChangeSvc: requestChangeSvc,
	}
}

func (h *PsreRequestChangeHandler) RequestChangePhone(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.RequestChangePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.requestChangeSvc.RequestChangePhone(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreRequestChangeHandler) ResendPhone(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ResendPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.requestChangeSvc.ResendPhone(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreRequestChangeHandler) ValidatePhone(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ValidatePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.requestChangeSvc.ValidatePhone(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreRequestChangeHandler) RequestChangeEmail(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.RequestChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.requestChangeSvc.RequestChangeEmail(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreRequestChangeHandler) ResendEmail(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ResendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.requestChangeSvc.ResendEmail(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreRequestChangeHandler) ValidateEmail(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ValidateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := utils.ResponseError(err.Error(), http.StatusBadRequest)
		c.Data(http.StatusBadRequest, "application/json", message)
		return
	}

	respBody, status, err := h.requestChangeSvc.ValidateEmail(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
