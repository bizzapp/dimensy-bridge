package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuotaClientHandler struct {
	service service.QuotaClientService
}

func NewQuotaClientHandler(s service.QuotaClientService) *QuotaClientHandler {
	return &QuotaClientHandler{s}
}

func (h *QuotaClientHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filters := map[string]interface{}{}
	if clientID := c.Query("client_id"); clientID != "" {
		filters["client_id"] = clientID
	}
	if mpID := c.Query("master_product_id"); mpID != "" {
		filters["master_product_id"] = mpID
	}

	quotas, total, err := h.service.GetQuotas(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "QUOTA_LIST_ERROR", "Gagal mengambil data quota", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	response.JSON(c, http.StatusOK, "List quota berhasil diambil", quotas, meta)
}

func (h *QuotaClientHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	quota, err := h.service.GetQuotaByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "QUOTA_NOT_FOUND", "Quota tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Quota ditemukan", quota, nil)
}

func (h *QuotaClientHandler) Create(c *gin.Context) {
	var quota model.QuotaClient
	if err := c.ShouldBindJSON(&quota); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	if err := h.service.CreateQuota(&quota); err != nil {
		response.Error(c, http.StatusInternalServerError, "QUOTA_CREATE_ERROR", "Gagal membuat quota", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Quota berhasil dibuat", quota, nil)
}

func (h *QuotaClientHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var quota model.QuotaClient
	if err := c.ShouldBindJSON(&quota); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}
	quota.ID = id

	if err := h.service.UpdateQuota(&quota); err != nil {
		response.Error(c, http.StatusInternalServerError, "QUOTA_UPDATE_ERROR", "Gagal update quota", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Quota berhasil diupdate", quota, nil)
}

func (h *QuotaClientHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.service.DeleteQuota(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "QUOTA_DELETE_ERROR", "Gagal hapus quota", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Quota berhasil dihapus", nil, nil)
}
