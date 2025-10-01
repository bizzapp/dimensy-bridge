package jwtutil

import (
	"dimensy-bridge/internal/middleware"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTokenFromContext(c *gin.Context) (string, error) {
	val, exists := c.Get(middleware.CtxTokenKey)
	if !exists {
		return "", errors.New("token not found in context")
	}
	fmt.Println(val, "VAL")
	token, ok := val.(string)
	if !ok {
		return "", errors.New("invalid user id type")
	}
	return token, nil
}
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	val, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		return 0, errors.New("user id not found in context")
	}
	fmt.Println(val, "VAL")
	userID, ok := val.(int64)
	if !ok {
		return 0, errors.New("invalid user id type")
	}
	return userID, nil
}

func GetUserIDAndRoleFromContext(c *gin.Context) (int64, string, error) {
	// Ambil UserID
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return 0, "", err
	}

	// Ambil Role
	val, exists := c.Get(middleware.CtxRoleKey)
	if !exists {
		return 0, "", errors.New("role not found in context")
	}
	role, ok := val.(string)
	if !ok {
		return 0, "", errors.New("invalid role type")
	}

	return userID, strings.ToLower(strings.TrimSpace(role)), nil
}
