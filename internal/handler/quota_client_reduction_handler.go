package handler

import (
	"dimensy-bridge/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuotaClientReductionHandler struct {
	service service.QuotaClientReductionService
}

func NewQuotaClientReductionHandler(service service.QuotaClientReductionService) *QuotaClientReductionHandler {
	return &QuotaClientReductionHandler{service: service}
}

func (h *QuotaClientReductionHandler) GetChart(c *gin.Context) {
	ctx := c.Request.Context()
	rangeKey := c.DefaultQuery("range", "last_7_days")

	var clientID *int64
	if v := c.Query("client_id"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)
		clientID = &id
	}

	var masterProductID *int64
	if v := c.Query("master_product_id"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)
		masterProductID = &id
	}

	resp, err := h.service.GetChart(ctx, rangeKey, clientID, masterProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
