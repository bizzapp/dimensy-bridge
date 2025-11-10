package psrehandler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"fmt"
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

func (h *PsreClientUserHandler) Activate(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ClientUserActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.psreClientUserSvc.ActivateUser(token, externalID, &req)
	if err != nil {
		fmt.Println(respBody, "respBody")
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
func (h *PsreClientUserHandler) ResendActivationUser(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ClientUserResendActivationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.psreClientUserSvc.ResendActivationUser(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreClientUserHandler) RequestPhoneActivation(c *gin.Context) {

	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ClientUserRequestPhoneActivationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.psreClientUserSvc.RequestPhoneActivation(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)

}

func (h *PsreClientUserHandler) PhoneActivation(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var req dto.ClientUserPhoneActivationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	respBody, status, err := h.psreClientUserSvc.PhoneActivation(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreClientUserHandler) List(c *gin.Context) {
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

	data, status, err := h.psreClientUserSvc.GetUsers(token, externalID, filter, page, limit)
	if err != nil {
		c.Data(status, "application/json", data)
		return
	}
	c.Data(status, "application/json", data)
}
func (h *PsreClientUserHandler) Detail(c *gin.Context) {
	token := c.GetHeader("Authorization")
	authData, _ := c.Get("authData")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id := c.Param("id")

	data, status, err := h.psreClientUserSvc.GetUserDetail(token, externalID, id)
	if err != nil {
		c.Data(status, "application/json", data)
		return
	}
	c.Data(status, "application/json", data)
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
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreClientUserHandler) RequestKYC(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	var req dto.ClientUserKYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	respBody, status, err := h.psreClientUserSvc.RequestKYC(token, externalID, &req)
	if err != nil {
		// fmt.Println(string(respBody), "respBody")
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *PsreClientUserHandler) VerifyKYC(c *gin.Context) {
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	externalID, err := utils.ExtractExternalID(authData)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	var req dto.ClientUserVerifyKYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	respBody, status, err := h.psreClientUserSvc.VerifyKYC(token, externalID, &req)
	if err != nil {
		c.Data(status, "application/json", respBody)
		return
	}
	c.Data(status, "application/json", respBody)
}
