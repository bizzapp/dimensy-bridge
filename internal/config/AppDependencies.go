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

	UserRepo repository.UserRepository
	UserSvc  service.UserService
	UserHdl  *handler.UserHandler

	ClientRepo repository.ClientRepository
	ClientSvc  service.ClientService
	ClientHdl  *handler.ClientHandler

	MasterProductRepo repository.MasterProductRepository
	MasterProductSvc  service.MasterProductService
	MasterProductHdl  *handler.MasterProductHandler
}

func NewAppDependencies(db *gorm.DB) *AppDependencies {

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo)
	authHdl := handler.NewAuthHandler(authSvc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)

	clientRepo := repository.NewClientRepository(db)
	clientSvc := service.NewClientService(clientRepo, userRepo)
	clientHdl := handler.NewClientHandler(clientSvc)

	masterProductRepo := repository.NewMasterProductRepository(db)
	masterProductSvc := service.NewMasterProductService(masterProductRepo)
	masterProductHdl := handler.NewMasterProductHandler(masterProductSvc)
	return &AppDependencies{
		DB:       db,
		AuthRepo: authRepo,
		AuthSvc:  authSvc,
		AuthHdl:  authHdl,

		UserRepo: userRepo,
		UserSvc:  userSvc,
		UserHdl:  userHdl,

		ClientRepo: clientRepo,
		ClientSvc:  clientSvc,
		ClientHdl:  clientHdl,

		MasterProductRepo: masterProductRepo,
		MasterProductSvc:  masterProductSvc,
		MasterProductHdl:  masterProductHdl,
	}
}
