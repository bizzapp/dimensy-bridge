package seeder

import "dimensy-bridge/internal/model"

const (
	ID_PRODUCT_CA_COMPANY                  int64 = 1
	ID_PRODUCT_CA_PERSONAL                 int64 = 2
	ID_PRODUCT_CA_PERSONAL_COMPANY         int64 = 3
	ID_PRODUCT_METERAI                     int64 = 4
	ID_PRODUCT_SIGN                        int64 = 5
	ID_PRODUCT_STAMP                       int64 = 6
	ID_PRODUCT_OTP                         int64 = 7
	ID_PRODUCT_DOCUMENT_STORAGE_KB         int64 = 8
	ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT int64 = 9
	ID_PRODUCT_COMPANY                     int64 = 10
	ID_PRODUCT_USER_PERSONAL               int64 = 11
	ID_PRODUCT_USER_PERSONAL_COMPANY       int64 = 12
	ID_PRODUCT_UPLOAD_SINGLE_DOCUMENT      int64 = 13
	ID_PRODUCT_UPLOAD_BULK_DOCUMENT        int64 = 14
)

// SeedMasterProducts returns initial data for master_products table.
func SeedMasterProducts() []model.MasterProduct {

	return []model.MasterProduct{
		{
			ID:           ID_PRODUCT_CA_COMPANY,
			Name:         "User Certificate Company",
			Code:         "CA_COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA_COMPANY),
			Icon:         strPtr("https://example.com/icon.png"),
			Notes:        strPtr("User Certificate"),
		},
		{
			ID:           ID_PRODUCT_CA_PERSONAL,
			Name:         "User Certificate Personal",
			Code:         "CA_PERSONAL",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA_PERSONAL),
			Icon:         strPtr("https://example.com/icon-upload.png"),
			Notes:        strPtr("Upload Service"),
		},
		{
			ID:           ID_PRODUCT_CA_PERSONAL_COMPANY,
			Name:         "USER Certificate Personal Company",
			Code:         "CA_PERSONAL_COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA_PERSONAL_COMPANY),
			Icon:         strPtr("https://example.com/icon-storage.png"),
			Notes:        strPtr("Storage Service"),
		},
		{
			ID:           ID_PRODUCT_METERAI,
			Name:         "Meterai",
			Code:         "METERAI",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_METERAI),
			Icon:         strPtr("https://example.com/icon-estamp.png"),
			Notes:        strPtr("Electronic Stamp"),
		},
		{
			ID:           ID_PRODUCT_SIGN,
			Name:         "Sign",
			Code:         "SIGN",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_SIGN),
			Icon:         strPtr("https://example.com/icon-esign.png"),
			Notes:        strPtr("Electronic Signature"),
		},
		{
			ID:           ID_PRODUCT_STAMP,
			Name:         "Stamp",
			Code:         "STAMP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_STAMP),
			Icon:         strPtr("https://example.com/icon-emeterai.png"),
			Notes:        strPtr("Electronic Stamp Duty"),
		},
		{
			ID:           ID_PRODUCT_OTP,
			Name:         "OTP",
			Code:         "OTP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_OTP),
			Icon:         strPtr("https://example.com/icon-otp.png"),
			Notes:        strPtr("One Time Password Service"),
		},
		{
			ID:           ID_PRODUCT_DOCUMENT_STORAGE_KB,
			Name:         "Document Storage KB",
			Code:         "DOCUMENT_STORAGE_KB",
			Quantity:     1000000000000,
			CurrentStock: 1000000000000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_DOCUMENT_STORAGE_KB),
			Icon:         strPtr("https://example.com/icon-email.png"),
			Notes:        strPtr("Email Notification Service"),
		},
		{
			ID:           ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT,
			Name:         "Document Upload Limit",
			Code:         "DOCUMENT_UPLOAD_COUNT_LIMIT",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT),
			Icon:         strPtr("https://example.com/icon-keyla.png"),
			Notes:        strPtr("Keyla AI Assistant"),
		},
		{
			ID:           ID_PRODUCT_COMPANY,
			Name:         "Company",
			Code:         "COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_COMPANY),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Custom Company Stamp"),
		},
		{
			ID:           ID_PRODUCT_USER_PERSONAL,
			Name:         "User Personal",
			Code:         "USER_PERSONAL",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         11,
			Icon:         strPtr("https://example.com/icon-ekyc.png"),
			Notes:        strPtr("Electronic Know Your Customer"),
		},
		{
			ID:           ID_PRODUCT_USER_PERSONAL_COMPANY,
			Name:         "User Personal Company",
			Code:         "USER_PERSONAL_COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_USER_PERSONAL_COMPANY),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
		},
		{
			ID:           ID_PRODUCT_UPLOAD_SINGLE_DOCUMENT,
			Name:         "Upload Single Document",
			Code:         "UPLOAD_SINGLE_DOCUMENT",
			Quantity:     0,
			CurrentStock: 0,
			IsUnlimited:  true,
			Sort:         int(ID_PRODUCT_UPLOAD_SINGLE_DOCUMENT),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
		},
		{
			ID:           ID_PRODUCT_UPLOAD_BULK_DOCUMENT,
			Name:         "Upload Bulk Document",
			Code:         "UPLOAD_BULK_DOCUMENT",
			Quantity:     0,
			CurrentStock: 0,
			IsUnlimited:  true,
			Sort:         int(ID_PRODUCT_UPLOAD_BULK_DOCUMENT),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
		},
	}
}

func strPtr(s string) *string {
	return &s
}
