package config

import (
	"dimensy-bridge/internal/model/seeder"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedMasterProducts(db *gorm.DB) error {
	products := seeder.SeedMasterProducts()
	for _, p := range products {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedUsers(db *gorm.DB) error {
	return seeder.SeedUsers(db)
}

func SeedSubscriptionPlans(db *gorm.DB) error {
	return seeder.SeedSubscriptionPlans(db)
}
