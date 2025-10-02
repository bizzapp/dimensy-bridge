package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MasterProductHandler struct {
	service service.MasterProductService
}

func NewMasterProductHandler(s service.MasterProductService) *MasterProductHandler {
	return &MasterProductHandler{s}
}

func (h *MasterProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filters := map[string]interface{}{}
	if code := c.Query("code"); code != "" {
		filters["code"] = code
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	products, total, err := h.service.GetProducts(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PRODUCT_LIST_ERROR", "Gagal mengambil data produk", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	response.JSON(c, http.StatusOK, "List produk berhasil diambil", products, meta)
}

func (h *MasterProductHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product, err := h.service.GetProductByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Produk tidak ditemukan", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Produk ditemukan", product, nil)
}

func (h *MasterProductHandler) Create(c *gin.Context) {
	var product model.MasterProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}

	if err := h.service.CreateProduct(&product); err != nil {
		response.Error(c, http.StatusInternalServerError, "PRODUCT_CREATE_ERROR", "Gagal membuat produk", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Produk berhasil dibuat", product, nil)
}

func (h *MasterProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var product model.MasterProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Input tidak valid", err.Error())
		return
	}
	product.ID = id

	if err := h.service.UpdateProduct(&product); err != nil {
		response.Error(c, http.StatusInternalServerError, "PRODUCT_UPDATE_ERROR", "Gagal update produk", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Produk berhasil diupdate", product, nil)
}

func (h *MasterProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.service.DeleteProduct(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "PRODUCT_DELETE_ERROR", "Gagal hapus produk", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Produk berhasil dihapus", nil, nil)
}
