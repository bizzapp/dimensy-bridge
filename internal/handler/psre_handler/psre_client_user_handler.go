package psrehandler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreClientUserHandler struct {
	clientUserSvc    service.ClientUserService
	clientPsreSvc    service.ClientPsreService
	clientCompanySvc service.ClientCompanyService
}

func NewPsreClientUserHandler(clientUserSvc service.ClientUserService, clientPsreSvc service.ClientPsreService, clientCompanySvc service.ClientCompanyService) *PsreClientUserHandler {
	return &PsreClientUserHandler{
		clientUserSvc:    clientUserSvc,
		clientPsreSvc:    clientPsreSvc,
		clientCompanySvc: clientCompanySvc,
	}
}

func (h *PsreClientUserHandler) Register(c *gin.Context) {

	authData, _ := c.Get("authData")
	// token := c.Request.Header.Get("Authorization")

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

	if err := h.clientUserSvc.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusCreated, user)
	response.JSON(c, http.StatusCreated, "Client user berhasil didaftarkan", user, nil)

}
