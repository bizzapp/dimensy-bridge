package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientKYCHistoryHandler struct {
	svc service.ClientKYCHistoryService
}

func NewClientKYCHistoryHandler(r *gin.Engine, svc service.ClientKYCHistoryService) {
	h := &ClientKYCHistoryHandler{svc: svc}

	group := r.Group("/api/v1/client-kyc-histories")
	{
		group.GET("/", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.GET("/user/:clientUserID", h.GetByClientUserID)
		group.POST("/", h.Create)
		group.POST("/:id/verify", h.VerifyKYC)
		group.DELETE("/:id", h.Delete)
	}
}

func (h *ClientKYCHistoryHandler) GetAll(c *gin.Context) {
	data, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ClientKYCHistoryHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	data, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ClientKYCHistoryHandler) GetByClientUserID(c *gin.Context) {
	clientUserID, _ := strconv.ParseInt(c.Param("clientUserID"), 10, 64)
	data, err := h.svc.GetByClientUserID(clientUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ClientKYCHistoryHandler) Create(c *gin.Context) {
	var req model.ClientKYCHistory
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

func (h *ClientKYCHistoryHandler) VerifyKYC(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	data, err := h.svc.VerifyKYC(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "KYC verified successfully",
		"data":    data,
	})
}

func (h *ClientKYCHistoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
