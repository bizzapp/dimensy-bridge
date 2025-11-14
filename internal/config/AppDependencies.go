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

	SubscriptionPlanRepo repository.SubscriptionPlanRepository
	SubscriptionPlanSvc  service.SubscriptionPlanService
	SubscriptionPlanHdl  *handler.SubscriptionPlanHandler

	ClientHasSubscriptionRepo repository.ClientHasSubscriptionPlanRepository
	ClientHasSubscriptionSvc  service.ClientHasSubscriptionPlanService
	ClientHasSubscriptionHdl  *handler.ClientHasSubscriptionPlanHandler

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

	// Client KYC History Module
	ClientKYCHistoryRepo repository.ClientKYCHistoryRepository
	ClientKYCHistorySvc  service.ClientKYCHistoryService
	ClientKYCHistoryHdl  *handler.ClientKYCHistoryHandler

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
	WebhookSvc service.WebhookService  // 👈 tambahkan ini

	ClientCompanyInviteRepo repository.ClientCompanyInviteRepository
	ClientCompanyInviteSvc  service.ClientCompanyInviteService

	ClientDocumentProcessRepo repository.ClientDocumentProcessRepository

	TokenBlacklistRepo repository.TokenBlacklistRepository
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
	clientHasSubscriptionPlanRepo := repository.NewClientHasSubscriptionPlanRepository(db)
	subscriptionPlanRepo := repository.NewSubscriptionPlanRepository(db)
	clientKYCHistoryRepo := repository.NewClientKYCHistoryRepository(db)
	clientCompanyInviteRepo := repository.NewClientCompanyInviteRepository(db)
	clientDocumentProcessRepo := repository.NewClientDocumentProcessRepository(db)

	tokenBlacklistRepo := repository.NewTokenBlacklistRepository(db)

	// === CORE SERVICES ===
	authSvc := service.NewAuthService(authRepo, tokenBlacklistRepo)
	userSvc := service.NewUserService(userRepo)
	quotaClientSvc := service.NewQuotaClientService(db, quotaClientRepo, quotaClientReductionRepo, quotaClientAdditionRepo)
	clientHasSubscriptionPlanSvc := service.NewClientHasSubscriptionPlanService(clientHasSubscriptionPlanRepo, quotaClientRepo, subscriptionPlanRepo, quotaClientSvc)
	clientCompanyInviteSvc := service.NewClientCompanyInviteService(clientCompanyInviteRepo)

	quotaClientAdditionSvc := service.NewQuotaClientAdditionService(db, quotaClientAdditionRepo, quotaClientRepo)
	clientSvc := service.NewClientService(clientRepo, userRepo, quotaClientRepo, quotaClientAdditionRepo)
	clientPsreSvc := service.NewClientPsreService(clientPsreRepo, clientRepo)
	clientCompanySvc := service.NewClientCompanyService(clientCompanyRepo, quotaClientSvc)
	clientUserSvc := service.NewClientUserService(clientUserRepo)
	masterProductSvc := service.NewMasterProductService(masterProductRepo)
	certificateSvc := service.NewCertificateService(certificateRepo)
	clientKYCHistorySvc := service.NewClientKYCHistoryService(clientKYCHistoryRepo)
	psreSvc := service.NewPsreService(clientRequestLogRepo, userRepo, clientCompanyRepo)
	webhookSvc := service.NewWebhookService(db, clientDocumentRepo, clientRequestLogRepo)

	// === PSRE SERVICES ===
	psreClientSvc := psre_service.NewClientService(clientRequestLogRepo, userRepo, clientRepo, clientPsreRepo)
	psreClientCompanySvc := psre_service.NewClientCompanyService(db, clientSvc, clientCompanyRepo, quotaClientSvc, clientCompanyInviteSvc, clientCompanyInviteRepo, clientUserRepo)
	psreClientUserSvc := psre_service.NewClientUserService(db, clientPsreSvc, clientCompanySvc, clientUserSvc, clientUserRepo, clientKYCHistorySvc, clientSvc, clientKYCHistoryRepo)
	psreCertificateSvc := psre_service.NewCertificateService(db, certificateRepo, clientSvc, userSvc, clientCompanySvc, clientUserSvc)
	psreClientDocumentSvc := psre_service.NewClientDocumentService(db, clientPsreSvc, clientDocumentRepo, clientDocumentProcessRepo)
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
	clientHasSubscriptionPlanHdl := handler.NewClientHasSubscriptionPlanHandler(clientHasSubscriptionPlanSvc)

	// === PSRE HANDLERS ===
	psreClientHdl := psre_handler.NewPsreClientHandler(psreClientSvc)
	psreCompanyHdl := psre_handler.NewPsreCompanyHandler(clientSvc, clientCompanySvc, psreClientCompanySvc)
	psreClientUserHdl := psre_handler.NewPsreClientUserHandler(clientUserSvc, clientPsreSvc, clientCompanySvc, psreClientUserSvc)
	psreCertificateHdl := psre_handler.NewPsreCertificateHandler(psreCertificateSvc)
	psreClientDocumentHdl := psre_handler.NewPsreClientDocumentHandler(psreClientDocumentSvc)
	psreDashboardHdl := psre_handler.NewPsreDashboardHandler(psreDashboardSvc)
	psreBackendHdl := psre_handler.NewPsreBackendHandler(psreBackendSvc)
	webhookHdl := handler.NewWebhookHandler(webhookSvc)
	return &AppDependencies{
		DB:                 db,
		AuthRepo:           authRepo,
		AuthSvc:            authSvc,
		AuthHdl:            authHdl,
		TokenBlacklistRepo: tokenBlacklistRepo,

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

		ClientKYCHistoryRepo: clientKYCHistoryRepo,
		ClientKYCHistorySvc:  clientKYCHistorySvc,

		ClientRequestLogRepo: clientRequestLogRepo,
		PsreSvc:              psreSvc,
		CertificateHdl:       certificateHdl,

		ClientCompanyRepo: clientCompanyRepo,
		ClientCompanySvc:  clientCompanySvc,
		ClientCompanyHdl:  clientCompanyHdl,

		ClientHasSubscriptionHdl: clientHasSubscriptionPlanHdl,

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
		WebhookHdl:            webhookHdl,
		WebhookSvc:            webhookSvc,
	}
}
