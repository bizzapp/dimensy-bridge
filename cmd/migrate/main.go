package main

import (
	"log"
	"os"
	"reflect"

	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/model"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	db := config.ConnectDatabase()

	// daftar semua model
	tables := []interface{}{
		&model.User{},
		&model.Client{},
		&model.MasterProduct{},
		&model.QuotaClient{},
		&model.QuotaClientAddition{},
		&model.QuotaClientReduction{},
		&model.ClientIPWhitelist{},
		&model.ClientPsre{},
		&model.ClientRequestLog{},
		&model.ClientCompany{},
		&model.ClientUser{},
		&model.Certificate{},
		&model.SubscriptionPlan{},
		&model.SubscriptionPlanDetail{},
		&model.ClientDocument{},
		&model.ClientDocumentProcess{},
		&model.ClientDocumentProcessDetail{},
		&model.ClientDocumentResendOtp{},
		&model.ClientHasSubscriptionPlan{},
		&model.ClientKYCHistory{},
		&model.ClientCompanyInvite{},
		&model.MasterProductAddition{},
		&model.MasterProductReduction{},
		&model.InventoryMasterProduct{},
		&model.InventoryMasterProductLog{},
		&model.TokenBlacklist{},
	}

	refreshMode := len(os.Args) > 1 && os.Args[1] == "refresh"

	if refreshMode {
		runMigrateRefresh(db, tables)
	} else {
		runMigrate(db, tables)
	}
}

// 🔹 migrate biasa
func runMigrate(db *gorm.DB, tables []interface{}) {
	log.Println("🚀 Running migrations...")
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("❌ Gagal migrate database: %v", err)
	}
	log.Println("✅ Database migration berhasil!")
}

// 🔹 drop semua tabel + migrate ulang
func runMigrateRefresh(db *gorm.DB, tables []interface{}) {
	migrator := db.Migrator()

	log.Println("🧹 Dropping all tables...")
	for _, table := range tables {
		tableName := db.NamingStrategy.TableName(reflect.TypeOf(table).Elem().Name())
		if err := migrator.DropTable(table); err != nil {
			log.Fatalf("❌ Gagal drop table %s: %v", tableName, err)
		}
		log.Printf("🗑️  Dropped table: %s", tableName)
	}

	log.Println("🚀 Running fresh migrations...")
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("❌ Gagal migrate database: %v", err)
	}

	log.Println("✅ Database migration refresh berhasil!")
}
