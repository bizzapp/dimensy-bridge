package seeder

import (
	"dimensy-bridge/internal/model"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedSubscriptionPlans(db *gorm.DB) error {
	now := time.Now()

	plans := []model.SubscriptionPlan{
		{
			ID:          1,
			Name:        "Client Basic",
			Description: strPtr("Default basic plan for new clients"),
			IsDefault:   false,
			IsUnlimited: false,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			Details: []model.SubscriptionPlanDetail{
				// Basic plan details - 10k quantity each
				{
					MasterProductID: ID_PRODUCT_CA_COMPANY,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_CA_PERSONAL,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_CA_PERSONAL_COMPANY,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_METERAI,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_SIGN,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_STAMP,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_OTP,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_DOCUMENT_STORAGE_KB,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_COMPANY,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_USER_PERSONAL,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_USER_PERSONAL_COMPANY,
					Quantity:        10000,
				},
				{
					MasterProductID: ID_PRODUCT_UPLOAD_SINGLE_DOCUMENT,
					Quantity:        10000,
					MaxSingleUpload: intPtr(20),
				},
				{
					MasterProductID:       ID_PRODUCT_UPLOAD_BULK_DOCUMENT,
					Quantity:              10000,
					MaxBulkUploadLimitPcs: intPtr(1),
					MaxBulkUploadLimitAll: intPtr(5),
					MaxBulkUploadCount:    intPtr(5),
				},
			},
		},
		{
			ID:          2,
			Name:        "Client Pro",
			Description: strPtr("Full production plan for premium clients"),
			IsDefault:   true,
			IsUnlimited: false,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			Details: []model.SubscriptionPlanDetail{
				// Pro plan details - 20k quantity each
				{
					MasterProductID: ID_PRODUCT_CA_COMPANY,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_CA_PERSONAL,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_CA_PERSONAL_COMPANY,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_METERAI,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_SIGN,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_STAMP,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_OTP,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_DOCUMENT_STORAGE_KB,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_COMPANY,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_USER_PERSONAL,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_USER_PERSONAL_COMPANY,
					Quantity:        20000,
				},
				{
					MasterProductID: ID_PRODUCT_UPLOAD_SINGLE_DOCUMENT,
					Quantity:        20000,
					MaxSingleUpload: intPtr(100),
				},
				{
					MasterProductID:       ID_PRODUCT_UPLOAD_BULK_DOCUMENT,
					Quantity:              20000,
					MaxBulkUploadLimitPcs: intPtr(10),
					MaxBulkUploadLimitAll: intPtr(50),
					MaxBulkUploadCount:    intPtr(50),
				},
			},
		},
	}

	for _, plan := range plans {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&plan).Error; err != nil {
			return err
		}
	}
	log.Println("✅ Subscription plans (Basic & Prod) seeded successfully!")
	return nil
}

func intPtr(i int) *int {
	return &i
}
