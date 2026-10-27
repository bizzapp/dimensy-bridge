package psrehandler

import (
	"dimensy-bridge/internal/dto"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreKycHandler struct {
	kycSvc psreservice.KycService
}

func NewPsreKycHandler(kycSvc psreservice.KycService) *PsreKycHandler {
	return &PsreKycHandler{
		kycSvc: kycSvc,
	}
}

func (h *PsreKycHandler) IdValidationFull(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	id := c.Param("id")
	respBody, status, err := h.kycSvc.IdValidationFull(token, externalID, id)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreKycHandler) LivenessFull(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	id := c.Param("id")

	respBody, status, err := h.kycSvc.LivenessFull(token, externalID, id)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreKycHandler) FailedKyc(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	id := c.Param("id")

	respBody, status, err := h.kycSvc.FailedKyc(token, externalID, id)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreKycHandler) FailedKycByBody(c *gin.Context) {
	externalID, token, err := utils.ValidateExternalID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.UserCompanyPicDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(http.StatusBadRequest, "application/json", utils.ResponseError(err.Error(), http.StatusBadRequest))
		return
	}

	respBody, status, err := h.kycSvc.FailedKycByBody(token, externalID, &req)
	if err != nil {
		if len(respBody) == 0 {
			respBody = utils.ResponseError(err.Error(), status)
		}
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
