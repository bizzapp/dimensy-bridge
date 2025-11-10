package middleware

import (
	"dimensy-bridge/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthJWE() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Missing Authorization header",
			})
			c.Abort()
			return
		}

		var token string

		// 🔹 Jika pakai format "Bearer <token>"
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		} else {
			// 🔹 Jika langsung token saja
			token = strings.TrimSpace(authHeader)
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Empty token",
			})
			c.Abort()
			return
		}

		// 🔹 Verifikasi token JWE
		data, err := utils.VerifyJWE(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// ✅ Simpan hasil verifikasi ke context
		c.Set("authData", data)
		c.Next()
	}
}
