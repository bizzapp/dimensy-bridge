package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// Get database URL from environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Printf("Database connection error: %v", err)
		panic("❌ Failed to connect to database!")
	}
	return db
}

func ConnectClientDatabase() *gorm.DB {
	// Get client database URL from environment variable
	// If not set, fallback to main DATABASE_URL
	clientDatabaseURL := os.Getenv("CLIENT_DATABASE_URL")
	if clientDatabaseURL == "" {
		clientDatabaseURL = os.Getenv("DATABASE_URL")
		log.Println("⚠️ CLIENT_DATABASE_URL not set, using DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(clientDatabaseURL), &gorm.Config{})
	if err != nil {
		log.Printf("Client database connection error: %v", err)
		panic("❌ Failed to connect to client database!")
	}
	return db
}
