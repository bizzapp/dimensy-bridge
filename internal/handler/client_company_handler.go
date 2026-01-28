package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientCompanyHandler struct {
	service service.ClientCompanyService
}

func NewClientCompanyHandler(s service.ClientCompanyService) *ClientCompanyHandler {
	return &ClientCompanyHandler{s}
}

// GET /companies
func (h *ClientCompanyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filters := map[string]interface{}{}
	if clientID := c.Query("client_id"); clientID != "" {
		if id, err := strconv.ParseInt(clientID, 10, 64); err == nil {
			filters["client_id"] = id
		}
	}

	companies, total, err := h.service.GetCompanies(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COMPANY_LIST_ERROR", "Gagal ambil data company", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	response.JSON(c, http.StatusOK, "List company berhasil diambil", companies, meta)
}

// GET /companies/:id
func (h *ClientCompanyHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	company, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Company ditemukan", company, nil)
}

// POST /companies
func (h *ClientCompanyHandler) Create(c *gin.Context) {
	var company model.ClientCompany
	if err := c.ShouldBindJSON(&company); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	if err := h.service.Create(&company); err != nil {
		response.Error(c, http.StatusInternalServerError, "COMPANY_CREATE_ERROR", "Gagal buat company", err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, "Company berhasil dibuat", company, nil)
}

// PUT /companies/:id
func (h *ClientCompanyHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var company model.ClientCompany
	if err := c.ShouldBindJSON(&company); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}
	company.ID = id

	if err := h.service.Update(&company); err != nil {
		response.Error(c, http.StatusInternalServerError, "COMPANY_UPDATE_ERROR", "Gagal update company", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Company berhasil diupdate", company, nil)
}

// DELETE /companies/:id
func (h *ClientCompanyHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "COMPANY_DELETE_ERROR", "Gagal hapus company", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Company berhasil dihapus", nil, nil)
}
