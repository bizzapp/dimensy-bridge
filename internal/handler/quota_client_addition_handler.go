package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuotaClientAdditionHandler struct {
	service service.QuotaClientAdditionService
}

func NewQuotaClientAdditionHandler(s service.QuotaClientAdditionService) *QuotaClientAdditionHandler {
	return &QuotaClientAdditionHandler{s}
}

func (h *QuotaClientAdditionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	additions, total, err := h.service.GetAdditions(page, limit, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "ADDITION_LIST_ERROR", "Gagal mengambil data quota addition", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	response.JSON(c, http.StatusOK, "List quota addition berhasil diambil", additions, meta)
}

func (h *QuotaClientAdditionHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	addition, err := h.service.GetAdditionByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "ADDITION_NOT_FOUND", "Data tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Data ditemukan", addition, nil)
}

func (h *QuotaClientAdditionHandler) Create(c *gin.Context) {
	var addition model.QuotaClientAddition
	if err := c.ShouldBindJSON(&addition); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	if err := h.service.CreateAddition(&addition); err != nil {
		response.Error(c, http.StatusInternalServerError, "ADDITION_CREATE_ERROR", "Gagal membuat quota addition", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Quota addition berhasil dibuat", addition, nil)
}

func (h *QuotaClientAdditionHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var addition model.QuotaClientAddition
	if err := c.ShouldBindJSON(&addition); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}
	addition.ID = id

	if err := h.service.UpdateAddition(&addition); err != nil {
		response.Error(c, http.StatusInternalServerError, "ADDITION_UPDATE_ERROR", "Gagal update quota addition", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Quota addition berhasil diupdate", addition, nil)
}

func (h *QuotaClientAdditionHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.service.DeleteAddition(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "ADDITION_DELETE_ERROR", "Gagal hapus quota addition", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Quota addition berhasil dihapus", nil, nil)
}

func (h *QuotaClientAdditionHandler) Process(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	processBy, _ := strconv.ParseInt(c.DefaultQuery("process_by", "0"), 10, 64)

	addition, err := h.service.ProcessAddition(id, processBy)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ADDITION_PROCESS_ERROR", "Gagal memproses quota addition", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Quota addition berhasil diproses", addition, nil)
}
