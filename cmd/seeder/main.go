package main

import (
	"dimensy-bridge/internal/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	db := config.ConnectDatabase()

	// 2️⃣ Jalankan seeder dari internal/seeder.go
	if err := config.SeedUsers(db); err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}

	if err := config.SeedMasterProducts(db); err != nil {
		log.Fatalf("failed to seed master products: %v", err)
	}

	if err := config.SeedSubscriptionPlans(db); err != nil {
		log.Fatalf("failed to seed subscription plans: %v", err)
	}

	log.Println("🌱 Seeder selesai dijalankan.")
}
