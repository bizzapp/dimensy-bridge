package main

import (
	"log"

	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/model"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	db := config.ConnectDatabase()

	if err := db.AutoMigrate(
		&model.User{},
	); err != nil {
		log.Fatalf("❌ Gagal migrate database: %v", err)
	}

	log.Println("✅ Database migration berhasil!")
}
