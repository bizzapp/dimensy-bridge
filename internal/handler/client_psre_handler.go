package handler

import (
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientPsreHandler struct {
	service service.ClientPsreService
}

func NewClientPsreHandler(s service.ClientPsreService) *ClientPsreHandler {
	return &ClientPsreHandler{s}
}

func (h *ClientPsreHandler) Register(c *gin.Context) {
	var req struct {
		ClientID int64 `json:"client_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	// ambil password default dari ENV
	password := os.Getenv("PSRE_DEFAULT_PASSWORD")
	if password == "" {
		password = "DefaultP@ssw0rd!" // fallback default
	}

	// hitung expired date berdasarkan ENV
	expDate := time.Now()

	// bisa override dengan satuan hari/bulan/tahun
	if days, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_DAYS")); days > 0 {
		expDate = expDate.AddDate(0, 0, days)
	} else if months, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_MONTHS")); months > 0 {
		expDate = expDate.AddDate(0, months, 0)
	} else if years, _ := strconv.Atoi(os.Getenv("PSRE_EXPIRE_YEARS")); years > 0 {
		expDate = expDate.AddDate(years, 0, 0)
	} else {
		// default 1 tahun
		expDate = expDate.AddDate(1, 0, 0)
	}

	// kirim ke service
	psre, err := h.service.RegisterPsre(req.ClientID, expDate, "", password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PSRE_REGISTER_ERROR", "Gagal register client ke PSRE", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Client PSRE berhasil dibuat", psre, nil)
}

func (h *ClientPsreHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	psre, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "PSRE_NOT_FOUND", "Client PSRE tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Data ditemukan", psre, nil)
}
