package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientUserHandler struct {
	service service.ClientUserService
}

func NewClientUserHandler(s service.ClientUserService) *ClientUserHandler {
	return &ClientUserHandler{s}
}

func (h *ClientUserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.JSON(c, http.StatusOK, "List user berhasil diambil", users, nil)
}

func (h *ClientUserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "User ditemukan", user, nil)
}

func (h *ClientUserHandler) Create(c *gin.Context) {
	var user model.ClientUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.JSON(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", nil)
		return
	}
	createdUser, err := h.service.Create(&user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_CREATE_ERROR", "Gagal buat user", err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, "User berhasil dibuat", createdUser, nil)
}

func (h *ClientUserHandler) Update(c *gin.Context) {
	var user model.ClientUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.JSON(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", nil)
		return
	}
	if err := h.service.Update(&user); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_UPDATE_ERROR", "Gagal update user", err.Error())
		return
	}
	// c.JSON(http.StatusOK, user)
	response.JSON(c, http.StatusOK, "User berhasil diupdate", user, nil)
}

func (h *ClientUserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_DELETE_ERROR", "Gagal hapus user", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "User berhasil dihapus", nil, nil)
}
