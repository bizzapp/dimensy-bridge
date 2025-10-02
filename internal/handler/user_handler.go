package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// ambil filter dari query
	filters := map[string]interface{}{}
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	users, total, err := h.service.GetUsers(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_LIST_ERROR", "Gagal mengambil data user", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	response.JSON(c, http.StatusOK, "List user berhasil diambil", users, meta)
}

func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.service.GetUserByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "User ditemukan", user, nil)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	// validasi password wajib
	if user.Password == nil || *user.Password == "" {
		response.Error(c, http.StatusBadRequest, "MISSING_PASSWORD", "Password wajib diisi", "")
		return
	}

	// hash password sebelum simpan
	hashed, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "HASH_ERROR", "Gagal meng-hash password", err.Error())
		return
	}
	hashedStr := string(hashed)
	user.Password = &hashedStr

	if err := h.service.CreateUser(&user); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_CREATE_ERROR", "Gagal membuat user", err.Error())
		return
	}

	// jangan return password ke response
	user.Password = nil

	response.JSON(c, http.StatusCreated, "User berhasil dibuat", user, nil)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}
	user.ID = id

	if err := h.service.UpdateUser(&user); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_UPDATE_ERROR", "Gagal update user", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "User berhasil diupdate", user, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.service.DeleteUser(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_DELETE_ERROR", "Gagal hapus user", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "User berhasil dihapus", nil, nil)
}
