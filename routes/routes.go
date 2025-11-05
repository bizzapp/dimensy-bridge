package routes

import (
	"dimensy-bridge/internal/config"
	"dimensy-bridge/internal/middleware"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRoutes(deps *config.AppDependencies) *gin.Engine {
	r := gin.Default()

	// Health check endpoint for Docker health checks
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Service is running",
		})
	})

	clientOriginsRaw := strings.Split(os.Getenv("CLIENT_ORIGIN"), ",")
	clientOrigins := make([]string, len(clientOriginsRaw))
	for i, origin := range clientOriginsRaw {
		clientOrigins[i] = strings.TrimSpace(origin)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     clientOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	rl := middleware.NewRateLimiter(rate.Every(200*time.Millisecond), 10)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/login", deps.AuthHdl.Login)

		auth.Use(middleware.JWTAuthMiddleware())
		auth.POST("/logout", deps.AuthHdl.Logout)
	}

	user := api.Group("/users")
	{
		user.GET("/", deps.UserHdl.List)
		user.GET("/:id", deps.UserHdl.Get)
		user.POST("/", deps.UserHdl.Create)
		user.PUT("/:id", deps.UserHdl.Update)
		user.DELETE("/:id", deps.UserHdl.Delete)
	}

	clients := api.Group("/clients")
	{
		clients.Use(middleware.JWTAuthMiddleware())
		clients.GET("/profile-psre/:id", deps.ClientPsreHdl.Profile)
		clients.GET("/fill-external-id/:id", deps.ClientPsreHdl.FillExternalId)
		clients.GET("/", deps.ClientHdl.List)
		clients.GET("/:id", deps.ClientHdl.Get)
		clients.POST("/", deps.ClientHdl.Create)
		clients.PUT("/:id", deps.ClientHdl.Update)
		clients.DELETE("/:id", deps.ClientHdl.Delete)
		clients.POST("/add-quota", deps.ClientHdl.AddQuota)
		clients.POST("/approve-add-quota", deps.ClientHdl.ApproveAddQuota)

	}

	products := api.Group("/products")
	{
		products.GET("/", deps.MasterProductHdl.List)
		products.GET("/:id", deps.MasterProductHdl.Get)
		products.POST("/", deps.MasterProductHdl.Create)
		products.PUT("/:id", deps.MasterProductHdl.Update)
		products.DELETE("/:id", deps.MasterProductHdl.Delete)
	}

	quotas := api.Group("/quotas")
	{
		quotas.GET("/", deps.QuotaClientHdl.List)
		quotas.GET("/:id", deps.QuotaClientHdl.Get)
		quotas.POST("/", deps.QuotaClientHdl.Create)
		quotas.PUT("/:id", deps.QuotaClientHdl.Update)
		quotas.DELETE("/:id", deps.QuotaClientHdl.Delete)
	}

	additions := api.Group("/quota-additions")
	{
		additions.GET("/", deps.QuotaClientAdditionHdl.List)
		additions.GET("/:id", deps.QuotaClientAdditionHdl.Get)
		additions.POST("/", deps.QuotaClientAdditionHdl.Create)
		additions.PUT("/:id", deps.QuotaClientAdditionHdl.Update)
		additions.DELETE("/:id", deps.QuotaClientAdditionHdl.Delete)
	}

	clientPsre := api.Group("/client-psre")
	{
		clientPsre.POST("/register", deps.ClientPsreHdl.Register)
		clientPsre.POST("/register-with-fill-external-id", deps.ClientPsreHdl.RegisterWithFillExternalId)
	}

	clientCompany := api.Group("/client-companies")
	{
		clientCompany.GET("/", deps.ClientCompanyHdl.List)
		clientCompany.GET("/:id", deps.ClientCompanyHdl.Get)
		clientCompany.POST("/", deps.ClientCompanyHdl.Create)
		clientCompany.PUT("/:id", deps.ClientCompanyHdl.Update)
		clientCompany.DELETE("/:id", deps.ClientCompanyHdl.Delete)
	}

	group := api.Group("/client_users")
	{
		group.GET("/", deps.ClientUserHdl.GetAll)
		group.GET("/:id", deps.ClientUserHdl.GetByID)
		group.POST("/", deps.ClientUserHdl.Create)
		group.PUT("/", deps.ClientUserHdl.Update)
		group.DELETE("/:id", deps.ClientUserHdl.Delete)
	}

	psre := api.Group("/psre")
	{
		psre.Use(rl.Middleware()) // pasang rate limiter di group ini
		psre.POST("/login", deps.PsreClientHdl.Login)

		company := psre.Group("/company")
		company.Use(middleware.AuthJWE())
		company.GET("/", deps.PsreCompanyHdl.GetClientCompany)
		company.GET("/detail/:id", deps.PsreCompanyHdl.DetailClientCompany)
		company.POST("/create", deps.PsreCompanyHdl.CreateClientCompany)
		company.POST("/invite", deps.PsreCompanyHdl.InviteClientCompany)

		user := psre.Group("/user")
		user.Use(middleware.AuthJWE())
		user.GET("/", deps.PsreClientUserHdl.List)
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
		certificate.POST("/issue", deps.PsreCertificateHdl.Issue)
		certificate.POST("/active", deps.PsreCertificateHdl.Active)
		certificate.POST("/revoke-request", deps.PsreCertificateHdl.RevokeRequest)
		certificate.POST("/revoke", deps.PsreCertificateHdl.Revoke)
	}
	return r
}
