package main

import (
	"context"
	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/jobs"
	"dimensy-bridge/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// Connect database
	db := config.ConnectDatabase()
	log.Println("✅ Database connected successfully!")

	// Dependencies
	deps := config.NewAppDependencies(db)

	// Setup HTTP server
	r := routes.SetupRoutes(deps)
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}
	// 🚀 RUN SERVER

	sched := jobs.NewScheduler(deps)
	if err := sched.Register(); err != nil {
		log.Fatalf("❌ Failed to register cron jobs: %v", err)
	}
	go sched.Start()
	log.Println("⏰ Cron scheduler started")
	go func() {
		port := os.Getenv("PORT") // <- ini dibaca dulu
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("❌ Failed to run server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	log.Println("🛑 Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop cron jobs
	if err := sched.Stop(shutdownCtx); err != nil {
		log.Printf("Error stopping scheduler: %v", err)
	}
	log.Println("✅ Cron scheduler stopped")

}
