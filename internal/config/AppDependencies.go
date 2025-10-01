package config

import (
	"dimensy-bridge/internal/handler"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"

	"gorm.io/gorm"
)

type AppDependencies struct {
	DB *gorm.DB

	AuthRepo repository.AuthRepository
	AuthSvc  service.AuthService
	AuthHdl  *handler.AuthHandler
}

func NewAppDependencies(db *gorm.DB) *AppDependencies {

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo)
	authHdl := handler.NewAuthHandler(authSvc)
	return &AppDependencies{
		DB:       db,
		AuthRepo: authRepo,
		AuthSvc:  authSvc,
		AuthHdl:  authHdl,
	}
}
