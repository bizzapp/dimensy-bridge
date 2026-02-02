package main

import (
	"dimensy-bridge/internal/config"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	db := config.ConnectDatabase()

	// Check if specific seeder is requested via command line argument
	if len(os.Args) > 1 {
		seederName := os.Args[1]
		runSpecificSeeder(db, seederName)
		return
	}

	// Run all seeders if no specific seeder is requested
	runAllSeeders(db)
}

func runSpecificSeeder(db *gorm.DB, seederName string) {
	switch seederName {
	case "users":
		log.Println("🌱 Running users seeder...")
		if err := config.SeedUsers(db); err != nil {
			log.Fatalf("failed to seed users: %v", err)
		}
		log.Println("✅ Users seeder completed")

	case "master-products", "products":
		log.Println("🌱 Running master products seeder...")
		if err := config.SeedMasterProducts(db); err != nil {
			log.Fatalf("failed to seed master products: %v", err)
		}
		log.Println("✅ Master products seeder completed")

	case "subscription-plans", "plans":
		log.Println("🌱 Running subscription plans seeder...")
		if err := config.SeedSubscriptionPlans(db); err != nil {
			log.Fatalf("failed to seed subscription plans: %v", err)
		}
		log.Println("✅ Subscription plans seeder completed")

	case "client-users":
		log.Println("🌱 Running client users seeder...")
		if err := config.SeedClientUsersFromUserActivations(db); err != nil {
			log.Fatalf("failed to seed client users: %v", err)
		}
		log.Println("✅ Client users seeder completed")

	default:
		fmt.Printf("❌ Unknown seeder: %s\n", seederName)
		fmt.Println("\nAvailable seeders:")
		fmt.Println("  users              - Seed admin users")
		fmt.Println("  master-products    - Seed master products")
		fmt.Println("  products           - Alias for master-products")
		fmt.Println("  subscription-plans - Seed subscription plans")
		fmt.Println("  plans              - Alias for subscription-plans")
		fmt.Println("  client-users       - Seed client users from user activations")
		fmt.Println("\nUsage:")
		fmt.Println("  go run cmd/seeder/main.go [seeder-name]")
		fmt.Println("  go run cmd/seeder/main.go                    # Run all seeders")
		fmt.Println("  go run cmd/seeder/main.go master-products    # Run specific seeder")
		os.Exit(1)
	}
}

func runAllSeeders(db *gorm.DB) {
	log.Println("🌱 Running all seeders...")

	if err := config.SeedUsers(db); err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}
	log.Println("✅ Users seeded")

	if err := config.SeedMasterProducts(db); err != nil {
		log.Fatalf("failed to seed master products: %v", err)
	}
	log.Println("✅ Master products seeded")

	if err := config.SeedSubscriptionPlans(db); err != nil {
		log.Fatalf("failed to seed subscription plans: %v", err)
	}
	log.Println("✅ Subscription plans seeded")

	if err := config.SeedClientUsersFromUserActivations(db); err != nil {
		log.Fatalf("failed to seed client users: %v", err)
	}
	log.Println("✅ Client users seeded")

	log.Println("🌱 All seeders completed successfully!")
}
