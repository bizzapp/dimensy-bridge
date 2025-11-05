package model

import "time"

const (
	ID_PRODUCT_CA       int64 = 1
	ID_PRODUCT_UPLOAD   int64 = 2
	ID_PRODUCT_STORAGE  int64 = 3
	ID_PRODUCT_ESTAMP   int64 = 4
	ID_PRODUCT_ESIGN    int64 = 5
	ID_PRODUCT_EMETERAI int64 = 6
	ID_PRODUCT_OTP      int64 = 7
	ID_PRODUCT_EMAIL    int64 = 8
	ID_PRODUCT_KEYLA    int64 = 9
	ID_PRODUCT_CSTAMP   int64 = 10
	ID_PRODUCT_EKYC     int64 = 11
	ID_PRODUCT_COMPANY  int64 = 12
)

// SeedMasterProducts returns initial data for master_products table.
func SeedMasterProducts() []MasterProduct {
	now := time.Now()

	return []MasterProduct{
		{
			ID:           ID_PRODUCT_CA,
			Name:         "User Certificate",
			Code:         "CA",
			Quantity:     100000,
			CurrentStock: 99900,
			IsUnlimited:  false,
			Sort:         1,
			Icon:         strPtr("https://example.com/icon.png"),
			Notes:        strPtr("User Certificate"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_UPLOAD,
			Name:         "Upload",
			Code:         "UPLOAD",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         2,
			Icon:         strPtr("https://example.com/icon-upload.png"),
			Notes:        strPtr("Upload Service"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_STORAGE,
			Name:         "Storage",
			Code:         "STORAGE",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         3,
			Icon:         strPtr("https://example.com/icon-storage.png"),
			Notes:        strPtr("Storage Service"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_ESTAMP,
			Name:         "eStamp",
			Code:         "ESTAMP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         4,
			Icon:         strPtr("https://example.com/icon-estamp.png"),
			Notes:        strPtr("Electronic Stamp"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_ESIGN,
			Name:         "eSign",
			Code:         "ESIGN",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         5,
			Icon:         strPtr("https://example.com/icon-esign.png"),
			Notes:        strPtr("Electronic Signature"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_EMETERAI,
			Name:         "eMeterai",
			Code:         "EMETERAI",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         6,
			Icon:         strPtr("https://example.com/icon-emeterai.png"),
			Notes:        strPtr("Electronic Stamp Duty"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_OTP,
			Name:         "OTP",
			Code:         "OTP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         7,
			Icon:         strPtr("https://example.com/icon-otp.png"),
			Notes:        strPtr("One Time Password Service"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_EMAIL,
			Name:         "Email",
			Code:         "EMAIL",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         8,
			Icon:         strPtr("https://example.com/icon-email.png"),
			Notes:        strPtr("Email Notification Service"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_KEYLA,
			Name:         "Keyla",
			Code:         "KEYLA",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         9,
			Icon:         strPtr("https://example.com/icon-keyla.png"),
			Notes:        strPtr("Keyla AI Assistant"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_CSTAMP,
			Name:         "Custom Stamp",
			Code:         "CSTAMP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         10,
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Custom Company Stamp"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_EKYC,
			Name:         "Ekyc",
			Code:         "EKYC",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         11,
			Icon:         strPtr("https://example.com/icon-ekyc.png"),
			Notes:        strPtr("Electronic Know Your Customer"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
		{
			ID:           ID_PRODUCT_COMPANY,
			Name:         "Company",
			Code:         "COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         12,
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		},
	}
}

func strPtr(s string) *string {
	return &s
}
