package psrehandler

import (
	psreService "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreClientHandler struct {
	psreClientSvc psreService.ClientService
}

func NewPsreClientHandler(psreClientSvc psreService.ClientService) *PsreClientHandler {
	return &PsreClientHandler{psreClientSvc: psreClientSvc}
}

func (h *PsreClientHandler) Login(c *gin.Context) {

	body, _ := io.ReadAll(c.Request.Body)

	resp, err := h.psreClientSvc.Login(body)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

func (h *PsreClientHandler) Documents(c *gin.Context) {
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

	data, status, err := h.psreClientSvc.GetDocuments(token, externalID, filter, page, limit)
	if err != nil {
		c.Data(status, "application/json", data)
		return
	}
	c.Data(status, "application/json", data)
}

func (h *PsreClientHandler) DocumentDetail(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	documentID := c.Param("id")

	data, status, err := h.psreClientSvc.GetDocumentDetail(token, externalID, documentID)
	if err != nil {
		c.Data(status, "application/json", data)
		return
	}
	c.Data(status, "application/json", data)
}
