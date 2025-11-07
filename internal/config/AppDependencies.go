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

	// Auth Module
	AuthRepo repository.AuthRepository
	AuthSvc  service.AuthService
	AuthHdl  *handler.AuthHandler

	// User Module
	UserRepo repository.UserRepository
	UserSvc  service.UserService
	UserHdl  *handler.UserHandler

	// Client Module
	ClientRepo repository.ClientRepository
	ClientSvc  service.ClientService
	ClientHdl  *handler.ClientHandler

	// Client PSRE Module
	ClientPsreRepo repository.ClientPsreRepository
	ClientPsreSvc  service.ClientPsreService
	ClientPsreHdl  *handler.ClientPsreHandler

	// Client Company Module
	ClientCompanyRepo repository.ClientCompanyRepository
	ClientCompanySvc  service.ClientCompanyService
	ClientCompanyHdl  *handler.ClientCompanyHandler

	// Client User Module
	ClientUserRepo repository.ClientUserRepository
	ClientUserSvc  service.ClientUserService
	ClientUserHdl  *handler.ClientUserHandler

	// Master Product Module
	MasterProductRepo repository.MasterProductRepository
	MasterProductSvc  service.MasterProductService
	MasterProductHdl  *handler.MasterProductHandler

	// Quota Client Module
	QuotaClientRepo          repository.QuotaClientRepository
	QuotaClientReductionRepo repository.QuotaClientReductionRepository
	QuotaClientSvc           service.QuotaClientService
	QuotaClientHdl           *handler.QuotaClientHandler

	// Quota Client Addition Module
	QuotaClientAdditionRepo repository.QuotaClientAdditionRepository
	QuotaClientAdditionSvc  service.QuotaClientAdditionService
	QuotaClientAdditionHdl  *handler.QuotaClientAdditionHandler

	ClientDocumentRepo repository.ClientDocumentRepository
	ClientDocumentSvc  psre_service.ClientDocumentService

	// Certificate Module
	CertificateHdl *handler.CertificateHandler

	// Client Request Log Module
	ClientRequestLogRepo repository.ClientRequestLogRepository

	// PSRE Service Module
	PsreSvc service.PsreService

	// PSRE Client Module
	PsreClientHdl *psre_handler.PsreClientHandler

	// PSRE Company Module
	PsreCompanyHdl       *psre_handler.PsreCompanyHandler
	PsreClientCompanySvc psre_service.ClientCompanyService

	// PSRE User Module
	PsreClientUserHdl *psre_handler.PsreClientUserHandler
	PsreClientUserSvc psre_service.ClientUserService

	// PSRE Certificate Module
	PsreCertificateHdl *psre_handler.PsreCertificateHandler
	PsreCertificateSvc psre_service.CertificateService

	// PSRE Document Module
	PsreClientDocumentHdl *psre_handler.PsreClientDocumentHandler
	PsreDocumentSvc       psre_service.ClientDocumentService

	// PSRE Dashboard Module
	PsreDashboardHdl *psre_handler.PsreDashboardHandler
	PsreDashboardSvc psre_service.DashboardService

	// PSRE Backend Module
	PsreBackendHdl *psre_handler.PsreBackendHandler
	PsreBackendSvc psre_service.BackendService

	WebhookHdl *handler.WebhookHandler // 👈 tambahkan ini

}

func NewAppDependencies(db *gorm.DB) *AppDependencies {
	// === REPOSITORIES ===
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	clientRepo := repository.NewClientRepository(db)
	clientPsreRepo := repository.NewClientPsreRepository(db)
	clientCompanyRepo := repository.NewClientCompanyRepository(db)
	clientUserRepo := repository.NewClientUserRepository(db)
	masterProductRepo := repository.NewMasterProductRepository(db)
	quotaClientRepo := repository.NewQuotaClientRepository(db)
	quotaClientReductionRepo := repository.NewQuotaClientReductionRepository(db)
	quotaClientAdditionRepo := repository.NewQuotaClientAdditionRepository(db)
	clientRequestLogRepo := repository.NewClientRequestLogRepository(db)
	certificateRepo := repository.NewCertificateRepository(db)
	clientDocumentRepo := repository.NewClientDocumentRepository(db)

	// === CORE SERVICES ===
	authSvc := service.NewAuthService(authRepo)
	userSvc := service.NewUserService(userRepo)
	quotaClientSvc := service.NewQuotaClientService(db, quotaClientRepo, quotaClientReductionRepo)
	quotaClientAdditionSvc := service.NewQuotaClientAdditionService(db, quotaClientAdditionRepo, quotaClientRepo)
	clientSvc := service.NewClientService(clientRepo, userRepo, quotaClientRepo, quotaClientAdditionRepo)
	clientPsreSvc := service.NewClientPsreService(clientPsreRepo, clientRepo)
	clientCompanySvc := service.NewClientCompanyService(clientCompanyRepo, quotaClientSvc)
	clientUserSvc := service.NewClientUserService(clientUserRepo)
	masterProductSvc := service.NewMasterProductService(masterProductRepo)
	certificateSvc := service.NewCertificateService(certificateRepo)
	psreSvc := service.NewPsreService(clientRequestLogRepo, userRepo, clientCompanyRepo)

	// === PSRE SERVICES ===
	psreClientSvc := psre_service.NewClientService(clientRequestLogRepo, userRepo, clientRepo, clientPsreRepo)
	psreClientCompanySvc := psre_service.NewClientCompanyService(db, clientSvc, clientCompanyRepo, quotaClientSvc)
	psreClientUserSvc := psre_service.NewClientUserService(db, clientPsreSvc, clientCompanySvc, clientUserSvc, clientUserRepo)
	psreCertificateSvc := psre_service.NewCertificateService(certificateRepo, clientSvc)
	psreClientDocumentSvc := psre_service.NewClientDocumentService(db, clientPsreSvc, clientDocumentRepo)
	psreDashboardSvc := psre_service.NewDashboardService()
	psreBackendSvc := psre_service.NewBackendService()

	// === CORE HANDLERS ===
	authHdl := handler.NewAuthHandler(authSvc)
	userHdl := handler.NewUserHandler(userSvc)
	clientHdl := handler.NewClientHandler(clientSvc, quotaClientSvc, quotaClientAdditionSvc)
	clientPsreHdl := handler.NewClientPsreHandler(psreClientSvc)
	clientCompanyHdl := handler.NewClientCompanyHandler(clientCompanySvc)
	clientUserHdl := handler.NewClientUserHandler(clientUserSvc)
	masterProductHdl := handler.NewMasterProductHandler(masterProductSvc)
	quotaClientHdl := handler.NewQuotaClientHandler(quotaClientSvc)
	quotaClientAdditionHdl := handler.NewQuotaClientAdditionHandler(quotaClientAdditionSvc)
	certificateHdl := handler.NewCertificateHandler(certificateSvc)

	// === PSRE HANDLERS ===
	psreClientHdl := psre_handler.NewPsreClientHandler(psreClientSvc)
	psreCompanyHdl := psre_handler.NewPsreCompanyHandler(clientSvc, clientCompanySvc, psreClientCompanySvc)
	psreClientUserHdl := psre_handler.NewPsreClientUserHandler(clientUserSvc, clientPsreSvc, clientCompanySvc, psreClientUserSvc)
	psreCertificateHdl := psre_handler.NewPsreCertificateHandler(psreCertificateSvc)
	psreClientDocumentHdl := psre_handler.NewPsreClientDocumentHandler(psreClientDocumentSvc)
	psreDashboardHdl := psre_handler.NewPsreDashboardHandler(psreDashboardSvc)
	psreBackendHdl := psre_handler.NewPsreBackendHandler(psreBackendSvc)
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

		QuotaClientRepo:          quotaClientRepo,
		QuotaClientReductionRepo: quotaClientReductionRepo,
		QuotaClientSvc:           quotaClientSvc,
		QuotaClientHdl:           quotaClientHdl,

		QuotaClientAdditionRepo: quotaClientAdditionRepo,
		QuotaClientAdditionSvc:  quotaClientAdditionSvc,
		QuotaClientAdditionHdl:  quotaClientAdditionHdl,

		ClientPsreRepo: clientPsreRepo,
		ClientPsreSvc:  clientPsreSvc,
		ClientPsreHdl:  clientPsreHdl,

		ClientUserRepo: clientUserRepo,
		ClientUserSvc:  clientUserSvc,
		ClientUserHdl:  clientUserHdl,

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

		PsreClientDocumentHdl: psreClientDocumentHdl,
		PsreDocumentSvc:       psreClientDocumentSvc,
		PsreDashboardHdl:      psreDashboardHdl,
		PsreDashboardSvc:      psreDashboardSvc,
		PsreBackendHdl:        psreBackendHdl,
		PsreBackendSvc:        psreBackendSvc,
	}
}
