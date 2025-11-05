package seeder

import "dimensy-bridge/internal/model"

const (
	ID_PRODUCT_CA                  int64 = 1
	ID_PRODUCT_UPLOAD              int64 = 2
	ID_PRODUCT_STORAGE             int64 = 3
	ID_PRODUCT_ESTAMP              int64 = 4
	ID_PRODUCT_ESIGN               int64 = 5
	ID_PRODUCT_EMETERAI            int64 = 6
	ID_PRODUCT_OTP                 int64 = 7
	ID_PRODUCT_EMAIL               int64 = 8
	ID_PRODUCT_KEYLA               int64 = 9
	ID_PRODUCT_CSTAMP              int64 = 10
	ID_PRODUCT_EKYC                int64 = 11
	ID_PRODUCT_COMPANY             int64 = 12
	ID_PRODUCT_CA_PERSONAL         int64 = 13
	ID_PRODUCT_CA_COMPANY          int64 = 14
	ID_PRODUCT_CA_PERSONAL_COMPANY int64 = 15
)

// SeedMasterProducts returns initial data for master_products table.
func SeedMasterProducts() []model.MasterProduct {

	return []model.MasterProduct{
		{
			ID:           ID_PRODUCT_CA,
			Name:         "User Certificate",
			Code:         "CA",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA),
			Icon:         strPtr("https://example.com/icon.png"),
			Notes:        strPtr("User Certificate"),
		},
		{
			ID:           ID_PRODUCT_UPLOAD,
			Name:         "Upload",
			Code:         "UPLOAD",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_UPLOAD),
			Icon:         strPtr("https://example.com/icon-upload.png"),
			Notes:        strPtr("Upload Service"),
		},
		{
			ID:           ID_PRODUCT_STORAGE,
			Name:         "Storage",
			Code:         "STORAGE",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_STORAGE),
			Icon:         strPtr("https://example.com/icon-storage.png"),
			Notes:        strPtr("Storage Service"),
		},
		{
			ID:           ID_PRODUCT_ESTAMP,
			Name:         "eStamp",
			Code:         "ESTAMP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_ESTAMP),
			Icon:         strPtr("https://example.com/icon-estamp.png"),
			Notes:        strPtr("Electronic Stamp"),
		},
		{
			ID:           ID_PRODUCT_ESIGN,
			Name:         "eSign",
			Code:         "ESIGN",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_ESIGN),
			Icon:         strPtr("https://example.com/icon-esign.png"),
			Notes:        strPtr("Electronic Signature"),
		},
		{
			ID:           ID_PRODUCT_EMETERAI,
			Name:         "eMeterai",
			Code:         "EMETERAI",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_EMETERAI),
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
			ID:           ID_PRODUCT_EMAIL,
			Name:         "Email",
			Code:         "EMAIL",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_EMAIL),
			Icon:         strPtr("https://example.com/icon-email.png"),
			Notes:        strPtr("Email Notification Service"),
		},
		{
			ID:           ID_PRODUCT_KEYLA,
			Name:         "Keyla",
			Code:         "KEYLA",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_KEYLA),
			Icon:         strPtr("https://example.com/icon-keyla.png"),
			Notes:        strPtr("Keyla AI Assistant"),
		},
		{
			ID:           ID_PRODUCT_CSTAMP,
			Name:         "Custom Stamp",
			Code:         "CSTAMP",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CSTAMP),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Custom Company Stamp"),
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
			Notes:        strPtr("Company"),
		},
		{
			ID:           ID_PRODUCT_CA_PERSONAL,
			Name:         "CA Personal",
			Code:         "CA_PERSONAL",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA_PERSONAL),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
		},
		{
			ID:           ID_PRODUCT_CA_COMPANY,
			Name:         "CA Company",
			Code:         "CA_COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA_COMPANY),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
		},
		{
			ID:           ID_PRODUCT_CA_PERSONAL_COMPANY,
			Name:         "CA Personal Company",
			Code:         "CA_PERSONAL_COMPANY",
			Quantity:     100000,
			CurrentStock: 100000,
			IsUnlimited:  false,
			Sort:         int(ID_PRODUCT_CA_PERSONAL_COMPANY),
			Icon:         strPtr("https://example.com/icon-customstamp.png"),
			Notes:        strPtr("Company"),
		},
	}
}

func strPtr(s string) *string {
	return &s
}
