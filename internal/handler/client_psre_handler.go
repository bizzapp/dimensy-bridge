package handler

import (
	psreintegration "dimensy-bridge/internal/service/psre_integration"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientPsreHandler struct {
	psreClientSvc psreintegration.ClientService
}

func NewClientPsreHandler(psreClientSvc psreintegration.ClientService) *ClientPsreHandler {
	return &ClientPsreHandler{psreClientSvc: psreClientSvc}
}

func (h *ClientPsreHandler) Register(c *gin.Context) {
	var req struct {
		ClientID int64 `json:"client_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	_, err := h.psreClientSvc.Register(req.ClientID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PSRE_REGISTER_ERROR", "Gagal register client ke PSRE", err.Error())
		return
	}
	psreClient, err := h.psreClientSvc.FillExternalId(req.ClientID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PSRE_REGISTER_ERROR", "Gagal register client ke PSRE", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Client PSRE berhasil dibuat", psreClient, nil)
}

func (h *ClientPsreHandler) Profile(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	profile, err := h.psreClientSvc.Profile(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PSRE_REGISTER_ERROR", "Gagal register client ke PSRE", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Client PSRE berhasil dibuat", profile, nil)
}
func (h *ClientPsreHandler) FillExternalId(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	profile, err := h.psreClientSvc.FillExternalId(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PSRE_REGISTER_ERROR", "Gagal register client ke PSRE", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Client PSRE berhasil dibuat", profile, nil)
}
