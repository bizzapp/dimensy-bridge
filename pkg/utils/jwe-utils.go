package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
)

// CreateJWE membuat token terenkripsi JWE (setara jose.EncryptJWT)
func CreateJWE(payload map[string]interface{}) (string, error) {
	// Ambil secret dari environment variable
	secretB64 := os.Getenv("JWT_SECRET")
	if secretB64 == "" {
		secretB64 = "eWx0ZXpwNGpmMEFuUTBleWdSblc5WndWSTZ3U08wakE=" // default JS
	}

	// Decode base64
	secret, err := base64.StdEncoding.DecodeString(secretB64)
	if err != nil {
		return "", fmt.Errorf("invalid base64 secret: %w", err)
	}

	// Tambahkan claim standard
	payload["iss"] = "Dimensy"
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Add(24 * time.Hour).Unix()

	// Serialize ke JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	encrypted, err := jwe.Encrypt(
		payloadBytes,
		jwe.WithKey(jwa.DIRECT, secret),
		jwe.WithContentEncryption(jwa.A128CBC_HS256),
	)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}

	return string(encrypted), nil
}

// VerifyJWE mendekripsi dan memvalidasi JWE token
func VerifyJWE(token string) (map[string]interface{}, error) {
	secretB64 := os.Getenv("JWT_SECRET")
	if secretB64 == "" {
		secretB64 = "eWx0ZXpwNGpmMEFuUTBleWdSblc5WndWSTZ3U08wakE="
	}

	secret, err := base64.StdEncoding.DecodeString(secretB64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 secret: %w", err)
	}

	// Dekripsi token
	decrypted, err := jwe.Decrypt([]byte(token), jwe.WithKey(jwa.DIRECT, secret))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	// Parse payload JSON
	var payload map[string]interface{}
	if err := json.Unmarshal(decrypted, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// Validasi issuer
	if iss, ok := payload["iss"].(string); !ok || iss != "Dimensy" {
		return nil, errors.New("invalid issuer")
	}

	// Validasi exp
	if expVal, ok := payload["exp"].(float64); ok {
		exp := time.Unix(int64(expVal), 0)
		if time.Now().After(exp) {
			return nil, errors.New("token expired")
		}
	}

	return payload, nil
}
