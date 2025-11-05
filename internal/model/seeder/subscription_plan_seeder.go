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
			IsDefault:   true,
			IsUnlimited: true,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			Details: []model.SubscriptionPlanDetail{
				// Basic plan details
				{
					MasterProductID: ID_PRODUCT_CA,
					Quantity:        1,
				},
				{
					MasterProductID: ID_PRODUCT_COMPANY,
					Quantity:        1,
				},
				{
					MasterProductID: ID_PRODUCT_ESIGN,
					Quantity:        5,
				},
				{
					MasterProductID: ID_PRODUCT_EMETERAI,
					Quantity:        5,
				},
				{
					MasterProductID: ID_PRODUCT_OTP,
					Quantity:        5,
				},
				{
					MasterProductID: ID_PRODUCT_EKYC,
					Quantity:        1,
				},
				{
					MasterProductID: ID_PRODUCT_UPLOAD, // upload single
					MaxSingleUpload: intPtr(20),
				},
				{
					MasterProductID:       ID_PRODUCT_STORAGE, // upload bulk
					MaxBulkUploadLimitPcs: intPtr(1),
					MaxBulkUploadLimitAll: intPtr(5),
					MaxBulkUploadCount:    intPtr(5),
				},
				{
					MasterProductID: ID_PRODUCT_CSTAMP, // personal company user
					Quantity:        20,
				},
			},
		},
		{
			ID:          2,
			Name:        "Client Pro",
			Description: strPtr("Full production plan for premium clients"),
			IsDefault:   false,
			IsUnlimited: false,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			Details: []model.SubscriptionPlanDetail{
				// Pro plan details
				{
					MasterProductID: ID_PRODUCT_CA,
					IsUnlimited:     true,
				},
				{
					MasterProductID: ID_PRODUCT_COMPANY,
					IsUnlimited:     true,
				},
				{
					MasterProductID: ID_PRODUCT_ESIGN,
					IsUnlimited:     true,
				},
				{
					MasterProductID: ID_PRODUCT_EMETERAI,
					Quantity:        100,
				},
				{
					MasterProductID: ID_PRODUCT_OTP,
					IsUnlimited:     true,
				},
				{
					MasterProductID: ID_PRODUCT_EKYC,
					IsUnlimited:     true,
				},
				{
					MasterProductID: ID_PRODUCT_UPLOAD, // upload single
					MaxSingleUpload: intPtr(5),
				},
				{
					MasterProductID:       ID_PRODUCT_STORAGE, // upload bulk
					MaxBulkUploadLimitPcs: intPtr(3),
					MaxBulkUploadLimitAll: intPtr(20),
					MaxBulkUploadCount:    intPtr(20),
				},
				{
					MasterProductID: ID_PRODUCT_CSTAMP, // personal company user
					IsUnlimited:     true,
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
