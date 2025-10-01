package response

import (
	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
}

// Success response dengan metadata
func JSON(c *gin.Context, status int, message string, data interface{}, meta *Meta) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error response
func Error(c *gin.Context, status int, code string, message string, detail string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
		Error: &ErrorDetail{
			Code:    code,
			Details: detail,
		},
	})
}
