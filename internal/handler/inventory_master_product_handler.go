package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
)

type InventoryMasterProductHandler struct {
	inventoryService service.InventoryMasterProductService
}

func NewInventoryMasterProductHandler(inventoryService service.InventoryMasterProductService) *InventoryMasterProductHandler {
	return &InventoryMasterProductHandler{
		inventoryService: inventoryService,
	}
}

// Index - Get all inventories with pagination
// GET /inventory_master_product
func (h *InventoryMasterProductHandler) Index(c *gin.Context) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	inventories, total, err := h.inventoryService.GetAllInventories(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch inventories", err.Error()))
		return
	}
	response.JSON(c, http.StatusOK, "Inventories fetched successfully", inventories, &response.Meta{
		Page:       page,
		Limit:      pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})

	// c.JSON(http.StatusOK, response.SuccessResponse("Inventories fetched successfully", gin.H{
	// 	"data":      inventories,
	// 	"total":     total,
	// 	"page":      page,
	// 	"page_size": pageSize,
	// }))
}

// List - Alternative endpoint to list inventories
// GET /inventory_master_product/list
func (h *InventoryMasterProductHandler) List(c *gin.Context) {
	h.Index(c)
}

// StoreOrUpdate - Create or update inventory
// POST /inventory_master_product/store_or_update
// If ID is null, create new inventory
// If ID exists, update only if is_processed is false
func (h *InventoryMasterProductHandler) StoreOrUpdate(c *gin.Context) {
	var req dto.CreateInventoryMasterProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", err.Error()))
		return
	}

	// Check if ID is provided (update) or null (create)
	if req.ID == nil {
		// CREATE: ID is null
		inventory, err := h.inventoryService.CreateInventory(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create inventory", err.Error()))
			return
		}

		c.JSON(http.StatusCreated, response.SuccessResponse("Inventory created successfully", inventory))
	} else {
		// UPDATE: ID exists
		updateReq := &dto.UpdateInventoryMasterProductRequest{
			VendorName:    req.VendorName,
			Price:         req.Price,
			IsPriorityUse: req.IsPriorityUse,
		}

		inventory, err := h.inventoryService.UpdateInventory(c.Request.Context(), *req.ID, updateReq)
		if err != nil {
			if err.Error() == "cannot update inventory that has been processed" {
				c.JSON(http.StatusBadRequest, response.ErrorResponse("Cannot update", "Inventory has been processed and cannot be updated"))
			} else {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update inventory", err.Error()))
			}
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse("Inventory updated successfully", inventory))
	}
}

// Show - Get inventory by ID
// GET /inventory_master_product/{id}/show
func (h *InventoryMasterProductHandler) Show(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid inventory ID", err.Error()))
		return
	}

	inventory, err := h.inventoryService.GetInventoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Inventory not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Inventory fetched successfully", inventory))
}

// AdjustStock - Adjust inventory stock (debit for stock in, credit for stock out)
// POST /inventory_master_product/{id}/adjust_stock
func (h *InventoryMasterProductHandler) AdjustStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid inventory ID", err.Error()))
		return
	}

	var req dto.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", err.Error()))
		return
	}

	inventory, err := h.inventoryService.AdjustStock(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to adjust stock", err.Error()))
		return
	}

	response.JSON(c, http.StatusOK, "Stock adjusted successfully", inventory, nil)
}

// GetLogs - Get inventory transaction logs
// GET /inventory_master_product/{id}/logs
func (h *InventoryMasterProductHandler) GetLogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid inventory ID", err.Error()))
		return
	}

	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	logs, total, err := h.inventoryService.GetInventoryLogs(c.Request.Context(), id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch logs", err.Error()))
		return
	}

	// c.JSON(http.StatusOK,)

	response.JSON(c, http.StatusOK, "Logs fetched successfully", logs, &response.Meta{
		Page:       page,
		Limit:      pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}

// MarkAsProcessed - Mark inventory as processed
// POST /inventory_master_product/{id}/mark_processed
func (h *InventoryMasterProductHandler) MarkAsProcessed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid inventory ID", err.Error()))
		return
	}

	inventory, err := h.inventoryService.MarkAsProcessed(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to mark as processed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Inventory marked as processed", inventory))
}

// TogglePriority - Toggle priority use status
// POST /inventory_master_product/{id}/toggle_priority
func (h *InventoryMasterProductHandler) TogglePriority(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid inventory ID", err.Error()))
		return
	}

	inventory, err := h.inventoryService.TogglePriority(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to toggle priority", err.Error()))
		return
	}

	response.JSON(c, http.StatusOK, "Priority toggled successfully", inventory, nil)
}

// Delete - Delete inventory
// DELETE /inventory_master_product/{id}/delete
func (h *InventoryMasterProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid inventory ID", err.Error()))
		return
	}

	if err := h.inventoryService.DeleteInventory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to delete inventory", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Inventory deleted successfully", nil))
}

// GetLowStockItems - Get items with low stock
// GET /inventory_master_product/low_stock/items
func (h *InventoryMasterProductHandler) GetLowStockItems(c *gin.Context) {
	threshold := 10

	if t := c.Query("threshold"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil && parsed > 0 {
			threshold = parsed
		}
	}

	items, err := h.inventoryService.GetLowStockItems(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch low stock items", err.Error()))
		return
	}

	response.JSON(c, http.StatusOK, "Low stock items fetched successfully", items, nil)
}

// GetTotalValue - Get total inventory value
// GET /inventory_master_product/total/value
func (h *InventoryMasterProductHandler) GetTotalValue(c *gin.Context) {
	totalValue, err := h.inventoryService.GetTotalInventoryValue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to calculate total value", err.Error()))
		return
	}

	response.JSON(c, http.StatusOK, "Total value calculated successfully", gin.H{
		"total_value": totalValue,
		"currency":    "IDR",
	}, nil)
}
