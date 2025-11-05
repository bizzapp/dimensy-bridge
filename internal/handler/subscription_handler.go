package handler

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscriptionPlanHandler struct {
	svc service.SubscriptionPlanService
}

func NewSubscriptionPlanHandler(router *gin.Engine, svc service.SubscriptionPlanService) {
	h := &SubscriptionPlanHandler{svc: svc}
	group := router.Group("/subscription-plans")
	{
		group.GET("", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.POST("", h.Create)
	}
}

func (h *SubscriptionPlanHandler) GetAll(c *gin.Context) {
	plans, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *SubscriptionPlanHandler) GetByID(c *gin.Context) {
	var plan model.SubscriptionPlan
	if err := c.BindUri(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.GetByID(plan.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SubscriptionPlanHandler) Create(c *gin.Context) {
	var plan model.SubscriptionPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}
