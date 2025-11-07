package handler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils/jwtutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHasSubscriptionPlanHandler struct {
	svc service.ClientHasSubscriptionPlanService
}

func NewClientHasSubscriptionPlanHandler(svc service.ClientHasSubscriptionPlanService) *ClientHasSubscriptionPlanHandler {
	return &ClientHasSubscriptionPlanHandler{svc: svc}

}

// GetAll : ambil semua data client_has_subscriptions
func (h *ClientHasSubscriptionPlanHandler) GetAll(c *gin.Context) {
	data, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.JSON(c, http.StatusOK, "success", data, nil)
}

// GetByID : ambil data berdasarkan ID
func (h *ClientHasSubscriptionPlanHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	data, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// Create : buat client_has_subscription baru
func (h *ClientHasSubscriptionPlanHandler) Create(c *gin.Context) {
	var req dto.CreateSubscriptionPlanRequest

	userId, err := jwtutil.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	req.CreatedBy = userId
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.svc.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, data)
}

// Process : aktifkan subscription dan generate quota_clients
func (h *ClientHasSubscriptionPlanHandler) Process(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	userId, err := jwtutil.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.svc.Process(id, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscription processed successfully",
	})
}

// Delete : hapus client_has_subscription
func (h *ClientHasSubscriptionPlanHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
