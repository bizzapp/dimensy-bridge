package psrehandler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreClientUserHandler struct {
	clientUserSvc     service.ClientUserService
	clientPsreSvc     service.ClientPsreService
	clientCompanySvc  service.ClientCompanyService
	psreClientUserSvc psreservice.ClientUserService
}

func NewPsreClientUserHandler(
	clientUserSvc service.ClientUserService,
	clientPsreSvc service.ClientPsreService,
	clientCompanySvc service.ClientCompanyService,
	psreClientUserSvc psreservice.ClientUserService,
) *PsreClientUserHandler {
	return &PsreClientUserHandler{
		clientUserSvc:     clientUserSvc,
		clientPsreSvc:     clientPsreSvc,
		clientCompanySvc:  clientCompanySvc,
		psreClientUserSvc: psreClientUserSvc,
	}
}

func (h *PsreClientUserHandler) Register(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	var req dto.ClientUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	respBody, status, err := h.psreClientUserSvc.RegisterUser(token, externalID, &req)
	if err != nil {
		// langsung kirim JSON dari PSrE bila ada
		c.Data(status, "application/json", respBody)
		return
	}

	response.JSON(c, http.StatusCreated, "Client user berhasil didaftarkan", respBody, nil)
}
