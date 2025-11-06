package config

import (
	"dimensy-bridge/internal/handler"
	psre_handler "dimensy-bridge/internal/handler/psre_handler"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	psre_service "dimensy-bridge/internal/service/psre_service"

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

	QuotaClientRepo repository.QuotaClientRepository
	QuotaClientSvc  service.QuotaClientService
	QuotaClientHdl  *handler.QuotaClientHandler

	QuotaClientAdditionRepo repository.QuotaClientAdditionRepository
	QuotaClientAdditionSvc  service.QuotaClientAdditionService
	QuotaClientAdditionHdl  *handler.QuotaClientAdditionHandler

	QuotaClientReductionRepo repository.QuotaClientReductionRepository
	// QuotaClientReductionSvc service.

	ClientPsreRepo repository.ClientPsreRepository
	ClientPsreSvc  service.ClientPsreService
	ClientPsreHdl  *handler.ClientPsreHandler

	ClientRequestLogRepo repository.ClientRequestLogRepository
	PsreSvc              service.PsreService

	ClientCompanyRepo repository.ClientCompanyRepository
	ClientCompanySvc  service.ClientCompanyService
	ClientCompanyHdl  *handler.ClientCompanyHandler

	ClientUserRepo repository.ClientUserRepository
	ClientUserSvc  service.ClientUserService
	ClientUserHdl  *handler.ClientUserHandler

	CertificateHdl *handler.CertificateHandler

	PsreCompanyHdl    *psre_handler.PsreCompanyHandler
	PsreClientHdl     *psre_handler.PsreClientHandler
	PsreClientUserHdl *psre_handler.PsreClientUserHandler

	PsreClientCompanySvc psre_service.ClientCompanyService
	PsreClientUserSvc    psre_service.ClientUserService

	PsreCertificateHdl *psre_handler.PsreCertificateHandler
	PsreCertificateSvc psre_service.CertificateService

	PsreDocumentHdl *psre_handler.PsreDocumentHandler
	PsreDocumentSvc psre_service.DocumentService

	PsreDashboardHdl *psre_handler.PsreDashboardHandler
	PsreDashboardSvc psre_service.DashboardService
}

func NewAppDependencies(db *gorm.DB) *AppDependencies {

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo)
	authHdl := handler.NewAuthHandler(authSvc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)

	clientRepo := repository.NewClientRepository(db)

	masterProductRepo := repository.NewMasterProductRepository(db)
	masterProductSvc := service.NewMasterProductService(masterProductRepo)
	masterProductHdl := handler.NewMasterProductHandler(masterProductSvc)

	quotaClientRepo := repository.NewQuotaClientRepository(db)
	quotaClientReductionRepo := repository.NewQuotaClientReductionRepository(db)
	quotaClientSvc := service.NewQuotaClientService(db, quotaClientRepo, quotaClientReductionRepo)
	quotaClientHdl := handler.NewQuotaClientHandler(quotaClientSvc)

	quotaClientAdditionRepo := repository.NewQuotaClientAdditionRepository(db)
	quotaClientAdditionSvc := service.NewQuotaClientAdditionService(db, quotaClientAdditionRepo, quotaClientRepo)
	quotaClientAdditionHdl := handler.NewQuotaClientAdditionHandler(quotaClientAdditionSvc)

	clientPsreRepo := repository.NewClientPsreRepository(db)
	clientPsreSvc := service.NewClientPsreService(clientPsreRepo, clientRepo)

	clientRequestLogRepo := repository.NewClientRequestLogRepository(db)

	clientCompanyRepo := repository.NewClientCompanyRepository(db)
	clientCompanySvc := service.NewClientCompanyService(clientCompanyRepo, quotaClientSvc)
	clientCompanyHdl := handler.NewClientCompanyHandler(clientCompanySvc)

	ClientUserRepo := repository.NewClientUserRepository(db)
	ClientUserSvc := service.NewClientUserService(ClientUserRepo)
	ClientUserHdl := handler.NewClientUserHandler(ClientUserSvc)

	psreClientSvc := psre_service.NewClientService(clientRequestLogRepo, userRepo, clientRepo, clientPsreRepo)
	clientPsreHdl := handler.NewClientPsreHandler(psreClientSvc)
	psreClientHdl := psre_handler.NewPsreClientHandler(psreClientSvc)

	certificateRepo := repository.NewCertificateRepository(db)
	certificateSvc := service.NewCertificateService(certificateRepo)
	certificateHdl := handler.NewCertificateHandler(certificateSvc)

	clientSvc := service.NewClientService(clientRepo, userRepo, quotaClientRepo, quotaClientAdditionRepo)
	psreClientCompanySvc := psre_service.NewClientCompanyService(db, clientSvc, clientCompanyRepo, quotaClientSvc)
	psreCompanyHdl := psre_handler.NewPsreCompanyHandler(clientSvc, clientCompanySvc, psreClientCompanySvc)

	psreDashboardSvc := psre_service.NewDashboardService()
	psreDashboardHdl := psre_handler.NewPsreDashboardHandler(psreDashboardSvc)

	psreDocumentSvc := psre_service.NewDocumentService()
	psreDocumentHdl := psre_handler.NewPsreDocumentHandler(psreDocumentSvc)
	psreSvc := service.NewPsreService(clientRequestLogRepo, userRepo, clientCompanyRepo)

	psreClientUserSvc := psre_service.NewClientUserService(db, clientPsreSvc, clientCompanySvc, ClientUserSvc, ClientUserRepo)
	psreClientUserHdl := psre_handler.NewPsreClientUserHandler(ClientUserSvc, clientPsreSvc, clientCompanySvc, psreClientUserSvc)

	psreCertificateSvc := psre_service.NewCertificateService(certificateRepo, clientSvc)

	psreCertificateHdl := psre_handler.NewPsreCertificateHandler(psreCertificateSvc)
	clientHdl := handler.NewClientHandler(clientSvc, quotaClientSvc, quotaClientAdditionSvc)
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

		QuotaClientRepo: quotaClientRepo,
		QuotaClientSvc:  quotaClientSvc,
		QuotaClientHdl:  quotaClientHdl,

		QuotaClientAdditionRepo: quotaClientAdditionRepo,
		QuotaClientAdditionSvc:  quotaClientAdditionSvc,
		QuotaClientAdditionHdl:  quotaClientAdditionHdl,

		ClientPsreRepo: clientPsreRepo,
		ClientPsreSvc:  clientPsreSvc,
		ClientPsreHdl:  clientPsreHdl,

		ClientUserRepo: ClientUserRepo,
		ClientUserSvc:  ClientUserSvc,
		ClientUserHdl:  ClientUserHdl,

		ClientRequestLogRepo: clientRequestLogRepo,
		PsreSvc:              psreSvc,
		CertificateHdl:       certificateHdl,

		ClientCompanyRepo: clientCompanyRepo,
		ClientCompanySvc:  clientCompanySvc,
		ClientCompanyHdl:  clientCompanyHdl,

		PsreClientHdl:        psreClientHdl,
		PsreCompanyHdl:       psreCompanyHdl,
		PsreClientUserHdl:    psreClientUserHdl,
		PsreClientCompanySvc: psreClientCompanySvc,
		PsreClientUserSvc:    psreClientUserSvc,

		PsreCertificateHdl: psreCertificateHdl,
		PsreCertificateSvc: psreCertificateSvc,

		PsreDocumentHdl:  psreDocumentHdl,
		PsreDocumentSvc:  psreDocumentSvc,
		PsreDashboardHdl: psreDashboardHdl,
		PsreDashboardSvc: psreDashboardSvc,
	}
}
