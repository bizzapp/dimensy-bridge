package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	CtxUserIDKey = "user_id"
	CtxRoleKey   = "role"
	CtxTokenKey  = "token"
)

// constants.go
const (
	RoleAdmin         = "admin"
	RoleAdministrator = "administrator"
	RoleUser          = "user"
	RoleGuest         = "guest"
)

// claims.go
type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware() gin.HandlerFunc {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	tokenCookieName := os.Getenv("COOKIE_NAME")
	if tokenCookieName == "" {
		tokenCookieName = "ACCESS_TOKEN"
	}

	return func(c *gin.Context) {
		var tokenString string

		// 1) Coba dari cookie
		if cookie, err := c.Cookie(tokenCookieName); err == nil && cookie != "" {
			fmt.Println("[AUTH] Token dari cookie")
			tokenString = cookie
		} else {
			// 2) Fallback ke Authorization header
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				fmt.Println("[AUTH] Token dari Authorization header")
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			fmt.Println("[AUTH] Token tidak ditemukan")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not provided"})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			// Validasi algoritma
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtKey, nil
		}, jwt.WithLeeway(30*time.Second)) // toleransi clock skew optional

		if err != nil || !token.Valid {
			// prioritas: signature invalid itu fatal
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
				return
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Simpan ke context (konsisten pakai konstanta)
		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRoleKey, claims.Role)
		c.Set(CtxTokenKey, tokenString)

		c.Next()
	}
}

// middleware/role_required.go
func RoleRequired(allowed ...string) gin.HandlerFunc {
	norm := func(s string) string { return strings.TrimSpace(strings.ToLower(s)) }
	allowedSet := map[string]struct{}{}
	for _, a := range allowed {
		allowedSet[norm(a)] = struct{}{}
	}

	return func(c *gin.Context) {
		v, ok := c.Get(CtxRoleKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no role in context"})
			return
		}
		role := norm(fmt.Sprintf("%v", v))
		if _, ok := allowedSet[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden", "role": role, "need": allowed})
			return
		}
		c.Next()
	}
}
