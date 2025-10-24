package psrehandler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	psreservice "dimensy-bridge/internal/service/psre_service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PsreCompanyHandler struct {
	clientSvc            service.ClientService
	clientCompanySvc     service.ClientCompanyService
	psreClientCompanySvc psreservice.ClientCompanyService
}

func NewPsreCompanyHandler(clientSvc service.ClientService, clientCompanySvc service.ClientCompanyService, psreClientCompanySvc psreservice.ClientCompanyService) *PsreCompanyHandler {
	return &PsreCompanyHandler{
		clientSvc:            clientSvc,
		clientCompanySvc:     clientCompanySvc,
		psreClientCompanySvc: psreClientCompanySvc,
	}
}

func (h *PsreCompanyHandler) CreateClientCompany(c *gin.Context) {

	// Ambil data hasil verifikasi dari middleware
	authData, _ := c.Get("authData")

	token := c.Request.Header.Get("Authorization")

	// fmt.Println("authData:", authData)
	externalID, err := utils.ExtractExternalID(authData)
	// fmt.Println("External ID:", externalID)

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

	// Bind body JSON
	var req dto.PsreCreateClientCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "Invalid request body", nil, nil)
		return
	}

	reqLocal := model.ClientCompany{
		ClientID: client.ID,
		Name:     req.CompanyName,
		Address:  req.CompanyAddress,
		Industry: req.CompanyIndustry,
		NPWP:     req.NPWP,
		NIB:      req.NIB,
		PICName:  req.PICName,
		PICEmail: req.PICEmail,
	}

	err = h.clientCompanySvc.Create(&reqLocal)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Failed to create client company", nil, nil)
		return
	}
	respBody, status, err := h.psreClientCompanySvc.RegisterCompany(token, req)
	if err != nil {
		// kalau error dari PSRE, respBody biasanya udah berisi JSON {code, message}
		var psreResp map[string]interface{}
		if jsonErr := json.Unmarshal(respBody, &psreResp); jsonErr == nil {
			c.JSON(status, psreResp)
			return
		}

		// fallback kalau bukan JSON valid
		c.JSON(status, gin.H{
			"code":    status,
			"message": string(respBody),
		})
		return
	}

	// sukses
	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "Company registered successfully",
		"data":    json.RawMessage(respBody),
	})

}
