package handler

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientIPWhitelistHandler struct {
	svc service.ClientIPWhitelistService
}

func NewClientIPWhitelistHandler(svc service.ClientIPWhitelistService) *ClientIPWhitelistHandler {
	return &ClientIPWhitelistHandler{svc: svc}
}

// Create creates a new IP whitelist entry
func (h *ClientIPWhitelistHandler) Create(c *gin.Context) {
	var req dto.CreateClientIPWhitelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.svc.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create IP whitelist",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "IP whitelist created successfully",
		"data":    result,
	})
}

// GetByID gets IP whitelist by ID
func (h *ClientIPWhitelistHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	result, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetByClientID gets all IP whitelists for a client
func (h *ClientIPWhitelistHandler) GetByClientID(c *gin.Context) {
	clientID, err := strconv.ParseInt(c.Param("client_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid client ID",
		})
		return
	}

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	result, err := h.svc.GetByClientID(clientID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch IP whitelists",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// Update updates an IP whitelist entry
func (h *ClientIPWhitelistHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	var req dto.UpdateClientIPWhitelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.svc.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update IP whitelist",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "IP whitelist updated successfully",
		"data":    result,
	})
}

// Delete deletes an IP whitelist entry
func (h *ClientIPWhitelistHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete IP whitelist",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "IP whitelist deleted successfully",
	})
}
