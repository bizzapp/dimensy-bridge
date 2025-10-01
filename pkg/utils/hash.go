package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	TokenExpire = 24 * 365 * time.Hour
) // 1 tahun

func HashString(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func GenerateJWT(userID int64, email *string, role string, name string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(TokenExpire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
