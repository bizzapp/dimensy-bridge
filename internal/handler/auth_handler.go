package handler

import (
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/response"
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
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, http.StatusBadRequest, "TOKEN_REQUIRED", "Authorization header required", "")
		return
	}

	// Remove "Bearer " prefix
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		response.Error(c, http.StatusBadRequest, "INVALID_TOKEN_FORMAT", "Invalid token format", "")
		return
	}

	// Add token to blacklist
	err := h.authSvc.Logout(tokenString)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Logout successful", nil, nil)
}
