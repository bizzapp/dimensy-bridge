package routes

import (
	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPsreRoutes(api *gin.RouterGroup, deps *config.AppDependencies, rl *middleware.RedisRateLimiter) {
	psre := api.Group("/psre")
	{
		psre.Use(rl.Middleware()) // pasang rate limiter di group ini

		// Apply JWE-based IP whitelist middleware untuk PSRE routes
		psre.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(
			deps.ClientIPWhitelistRepo,
			deps.ClientPsreRepo,
		))

		company := psre.Group("/company")
		company.Use(middleware.AuthJWE())
		company.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(deps.ClientIPWhitelistRepo, deps.ClientPsreRepo)) // Strict IP whitelist

		user := psre.Group("/user")
		user.Use(middleware.AuthJWE())
		user.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(deps.ClientIPWhitelistRepo, deps.ClientPsreRepo)) // Strict IP whitelist
		user.GET("/", deps.PsreClientUserHdl.List)
		user.GET("/sync", deps.PsreClientUserHdl.Sync)
		user.GET("/:id", deps.PsreClientUserHdl.Detail)
		user.POST("/register", deps.PsreClientUserHdl.Register)
		user.POST("/activate", deps.PsreClientUserHdl.Activate)
		user.POST("/resend-activation", deps.PsreClientUserHdl.ResendActivationUser)
		user.POST("/request-phone-activation", deps.PsreClientUserHdl.RequestPhoneActivation)
		user.POST("/phone-activation", deps.PsreClientUserHdl.PhoneActivation)
		user.POST("/request-kyc", deps.PsreClientUserHdl.RequestKYC)
		user.POST("/verify-kyc", deps.PsreClientUserHdl.VerifyKYC)

		certificate := psre.Group("/certificate")
		certificate.Use(middleware.AuthJWE())
		certificate.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(deps.ClientIPWhitelistRepo, deps.ClientPsreRepo)) // Strict IP whitelist
		certificate.POST("/issue", deps.PsreCertificateHdl.Issue)
		certificate.POST("/active", deps.PsreCertificateHdl.Active)
		certificate.POST("/revoke-request", deps.PsreCertificateHdl.RevokeRequest)
		certificate.POST("/revoke", deps.PsreCertificateHdl.Revoke)
		certificate.POST("/resync-certificates", deps.PsreCertificateHdl.ResyncCertificates)

		certificateV2 := certificate.Group("/v2")
		certificateV2.POST("/request-issue", deps.PsreCertificateHdl.RequestIssueV2)
		certificateV2.POST("/issue", deps.PsreCertificateHdl.IssueV2)
		certificateV2.POST("/revoke-request", deps.PsreCertificateHdl.RevokeRequestV2)
		certificateV2.POST("/revoke", deps.PsreCertificateHdl.RevokeV2)
		certificateV2.POST("/revoke-ra", deps.PsreCertificateHdl.RevokeRA)

		document := psre.Group("/document")
		document.Use(middleware.AuthJWE())
		document.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(deps.ClientIPWhitelistRepo, deps.ClientPsreRepo)) // Strict IP whitelist
		document.POST("/upload", deps.PsreClientDocumentHdl.Upload)
		document.POST("/upload-bulk", deps.PsreClientDocumentHdl.UploadBulk)
		document.POST("/request-sign", deps.PsreClientDocumentHdl.RequestSign)
		document.POST("/proccess-sign", deps.PsreClientDocumentHdl.ProcessSign)
		document.POST("/request-stamp", deps.PsreClientDocumentHdl.RequestStamp)
		document.POST("/proccess-stamp", deps.PsreClientDocumentHdl.ProcessStamp)
		document.POST("/request-otp-sign", deps.PsreClientDocumentHdl.RequestOtpSign)
		document.GET("/preview/:id", deps.PsreClientDocumentHdl.PreviewDocument)
		document.POST("/retry-process", deps.PsreClientDocumentHdl.RetryProcess)

		client := psre.Group("/client")
		client.POST("/login", deps.PsreClientHdl.Login) // Login tanpa IP whitelist
		client.Use(middleware.AuthJWE())
		client.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(deps.ClientIPWhitelistRepo, deps.ClientPsreRepo)) // Strict IP whitelist untuk endpoint setelah login
		client.POST("/company/create", deps.PsreCompanyHdl.CreateClientCompany)
		client.POST("/company/invite", deps.PsreCompanyHdl.InviteClientCompany)
		client.GET("/company", deps.PsreCompanyHdl.GetClientCompany)
		client.GET("/company/detail/:id", deps.PsreCompanyHdl.DetailClientCompany)
		client.POST("/users/accept-invitation", deps.PsreCompanyHdl.AcceptInvitationClientUser)
		client.GET("/documents", deps.PsreClientHdl.Documents)
		client.GET("/documents/:id", deps.PsreClientHdl.DocumentDetail)

		backend := psre.Group("/backend")
		backend.POST("/login", deps.PsreBackendHdl.Login) // Login tanpa IP whitelist
		backend.Use(middleware.AuthJWE())
		backend.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(deps.ClientIPWhitelistRepo, deps.ClientPsreRepo)) // Strict IP whitelist untuk endpoint setelah login
		backend.POST("/client/create", deps.PsreBackendHdl.CreateClient)
		backend.GET("/client", deps.PsreBackendHdl.ListClient)
		backend.POST("/client/update/:id", deps.PsreBackendHdl.UpdateClient)
		backend.POST("/client/update_status/:id", deps.PsreBackendHdl.UpdateClientStatus)
		backend.GET("/dashboard/certificate", deps.PsreDashboardHdl.Certificate)
		backend.GET("/dashboard/document", deps.PsreDashboardHdl.Document)
	}
}
