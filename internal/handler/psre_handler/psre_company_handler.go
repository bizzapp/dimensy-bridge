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
func (h *PsreCompanyHandler) GetClientCompany(c *gin.Context) {
	// Ambil data hasil verifikasi dari middleware
	authData, _ := c.Get("authData")
	token := c.Request.Header.Get("Authorization")

	_, err := utils.ExtractExternalID(authData)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	// Ambil semua query params dan ubah jadi map[string]string
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// Forward ke PSRE
	respBody, status, err := h.psreClientCompanySvc.GetCompany(token, params)
	if err != nil {
		var psreResp map[string]interface{}
		if jsonErr := json.Unmarshal(respBody, &psreResp); jsonErr == nil {
			c.JSON(status, psreResp)
			return
		}

		c.JSON(status, gin.H{
			"code":    status,
			"message": string(respBody),
		})
		return
	}

	// Sukses → langsung teruskan response dari PSRE
	c.JSON(status, json.RawMessage(respBody))
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

	var psreSuccessResp struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		CompanyID string `json:"companyId"`
	}
	if err := json.Unmarshal(respBody, &psreSuccessResp); err != nil {
		response.JSON(c, http.StatusInternalServerError, "Failed to parse PSRE response", nil, nil)
		return
	}

	// Update external_id
	if psreSuccessResp.CompanyID != "" {
		if err := h.clientCompanySvc.UpdateExternalID(reqLocal.ID, psreSuccessResp.CompanyID); err != nil {
			response.JSON(c, http.StatusInternalServerError, "Failed to update external ID", nil, nil)
			return
		}
	}

	// sukses
	c.JSON(status, json.RawMessage(respBody))

}
