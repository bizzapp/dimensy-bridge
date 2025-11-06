package utils

import (
	"encoding/base64"
	"fmt"
	"time"
)

func NowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// calculateBase64FileSize menghitung ukuran file dari base64 string
func CalculateBase64FileSize(base64Data string) (int64, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return 0, fmt.Errorf("invalid base64 data: %w", err)
	}

	// Hitung ukuran dalam KB (1 KB = 1024 bytes)
	sizeKB := int64(len(decodedData)) / 1024

	return sizeKB, nil
}
