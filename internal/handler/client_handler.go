package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service service.ClientService
}

func NewClientHandler(s service.ClientService) *ClientHandler {
	return &ClientHandler{s}
}

func (h *ClientHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filters := map[string]interface{}{}
	if companyName := c.Query("company_name"); companyName != "" {
		filters["company_name"] = companyName
	}

	clients, total, err := h.service.GetClients(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "CLIENT_LIST_ERROR", "Gagal mengambil data client", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	response.JSON(c, http.StatusOK, "List client berhasil diambil", clients, meta)
}

func (h *ClientHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	client, err := h.service.GetClientByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "CLIENT_NOT_FOUND", "Client tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Client ditemukan", client, nil)
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req struct {
		CompanyName string `json:"company_name" binding:"required"`
		PicName     string `json:"pic_name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	client, err := h.service.CreateClient(req.CompanyName, req.PicName, req.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "CLIENT_CREATE_ERROR", "Gagal membuat client", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Client berhasil dibuat", client, nil)
}

func (h *ClientHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var client model.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}
	client.ID = id

	if err := h.service.UpdateClient(&client); err != nil {
		response.Error(c, http.StatusInternalServerError, "CLIENT_UPDATE_ERROR", "Gagal update client", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Client berhasil diupdate", client, nil)
}

func (h *ClientHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.service.DeleteClient(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "CLIENT_DELETE_ERROR", "Gagal hapus client", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Client berhasil dihapus", nil, nil)
}
