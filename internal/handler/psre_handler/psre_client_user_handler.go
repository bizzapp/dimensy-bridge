package psrehandler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreClientUserHandler struct {
	clientUserSvc     service.ClientUserService
	psreClientUserSvc psreservice.ClientUserService
	clientPsreSvc     service.ClientPsreService
	clientCompanySvc  service.ClientCompanyService
}

func NewPsreClientUserHandler(clientUserSvc service.ClientUserService, clientPsreSvc service.ClientPsreService, clientCompanySvc service.ClientCompanyService, psreClientUserSvc psreservice.ClientUserService) *PsreClientUserHandler {
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

	clientPsre, err := h.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	user := model.ClientUser{
		NIK:       req.NIK,
		Name:      req.FullName,
		Birthdate: req.BirthDate.Time, // tinggal ambil time.Time-nya
		Email:     req.Email,
		Phone:     req.Phone,
		IsWNI:     req.IsWNI,
		ClientID:  clientPsre.ClientID,
	}

	user.ClientCompanyID = nil
	if req.CompanyID != nil {
		clientCompany, err := h.clientCompanySvc.GetByExternalID(*req.CompanyID)
		if err != nil {
			response.JSON(c, http.StatusInternalServerError, err.Error(), nil, nil)
			return
		}
		user.ClientCompanyID = &clientCompany.ID
	}

	userPtr, err := h.clientUserSvc.Create(&user)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	data, _, err := h.psreClientUserSvc.Register(token, req)
	fmt.Println("data psre register user:", string(data))
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// --- parse JSON PSrE response ---
	var psreResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		UserID  string `json:"userId"`
	}
	if err := json.Unmarshal(data, &psreResp); err != nil {
		response.JSON(c, http.StatusInternalServerError, "invalid psre response", nil, nil)
		return
	}

	// --- kalau sukses, update external_id ---
	if psreResp.Code == 0 && psreResp.UserID != "" {
		userPtr.ExternalID = &psreResp.UserID

		if err := h.clientUserSvc.Update(userPtr); err != nil {
			response.JSON(c, http.StatusInternalServerError, "failed to update external_id", nil, nil)
			return
		}
	}

	// c.JSON(http.StatusCreated, user)
	response.JSON(c, http.StatusCreated, "Client user berhasil didaftarkan", userPtr, nil)

}
