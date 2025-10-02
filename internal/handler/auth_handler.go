package handler

import (
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
	"dimensy-bridge/pkg/utils/jwtutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: s,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(c, http.StatusBadRequest, "MISSING_CREDENTIALS", "Email and password are required", "")
		return
	}

	token, user, err := h.authSvc.Login(&req.Email, &req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "LOGIN_FAILED", err.Error(), "")
		return
	}

	data := gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	response.JSON(c, http.StatusOK, "Login successful", data, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Ambil token dari header Authorization
	token, err := jwtutil.GetTokenFromContext(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Token not provided",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get token"})
		return
	}

	_, _, err = jwtutil.GetUserIDAndRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Blacklist token (opsional, tergantung implementasi)
	if err := h.authSvc.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to logout",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
