package main

import (
	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/model/seeder"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	db := config.ConnectDatabase()

	// 2️⃣ Jalankan seeder dari internal/seeder.go
	if err := config.SeedMasterProducts(db); err != nil {
		log.Fatalf("failed to seed master products: %v", err)
	}

	if err := seeder.SeedSubscriptionPlans(db); err != nil {
		log.Fatalf("failed to seed subscription plans: %v", err)
	}

	log.Println("🌱 Seeder selesai dijalankan.")
}
