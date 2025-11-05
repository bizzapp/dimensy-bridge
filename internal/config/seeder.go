package config

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedMasterProducts(db *gorm.DB) error {
	products := model.SeedMasterProducts()
	for _, p := range products {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p).Error; err != nil {
			return err
		}
	}
	return nil
}
