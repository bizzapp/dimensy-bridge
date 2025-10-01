package main

import (
	"log"
	"os"
	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/model"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	// Connect DB
	db := config.ConnectDatabase()

	// Auto migrate tabel User
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ gagal migrate tabel: %v", err)
	}

	log.Printf("🚀 Server running on port %s\n", os.Getenv("PORT"))
}
