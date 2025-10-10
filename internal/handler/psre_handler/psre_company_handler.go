package psrehandler

import (
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreCompanyHandler struct {
	clientSvc        service.ClientService
	clientCompanySvc service.ClientCompanyService
}

func NewPsreCompanyHandler(clientSvc service.ClientService, clientCompanySvc service.ClientCompanyService) *PsreCompanyHandler {
	return &PsreCompanyHandler{
		clientSvc:        clientSvc,
		clientCompanySvc: clientCompanySvc,
	}
}

func (h *PsreCompanyHandler) CreateClientCompany(c *gin.Context) {

	// Ambil data hasil verifikasi dari middleware
	authData, _ := c.Get("authData")

	fmt.Println("authData:", authData)
	externalID, err := utils.ExtractExternalID(authData)
	fmt.Println("External ID:", externalID)

	if err != nil {
		response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	// }

	client, err := h.clientSvc.GetClientByExternalId(externalID)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	response.JSON(c, http.StatusOK, "Client company berhasil dibuat", client, nil)
	// if err != nil {

	// // Bind body JSON
	// var req model.ClientCompany
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	response.JSON(c, http.StatusBadRequest, "Invalid request body", nil, err)
	// 	return
	// }

	// // Simpan ke DB via service
	// created, err := h.clientCompanySvc.Create(&req)
	// if err != nil {
	// 	response.JSON(c, http.StatusInternalServerError, "Failed to create client company", nil, err)
	// 	return
	// }

	// response.JSON(c, http.StatusOK, "Client company berhasil dibuat", created, nil)

}
