package seeder

import (
	"dimensy-bridge/internal/model"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedUsers(db *gorm.DB) error {
	now := time.Now()

	// Hash password P@ssw0rd
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("P@ssw0rd"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}
	hashedPasswordStr := string(hashedPassword)

	users := []model.User{
		{
			Name:            "Administrator",
			Email:           stringPtr("admin@dimensy.com"),
			Username:        stringPtr("administrator"),
			EmailVerifiedAt: &now,
			Password:        &hashedPasswordStr,
			Role:            "administrator",
			CreatedAt:       &now,
			UpdatedAt:       &now,
		},
	}

	for _, user := range users {
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "email", "username", "email_verified_at",
				"password", "role", "updated_at",
			}),
		}).Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Name, err)
			return err
		}
	}

	log.Println("✅ Users seeded successfully")
	return nil
}

func stringPtr(s string) *string {
	return &s
}
