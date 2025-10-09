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
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     clientOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		clients.GET("/", deps.ClientHdl.List)
		clients.GET("/:id", deps.ClientHdl.Get)
		clients.POST("/", deps.ClientHdl.Create)
		clients.PUT("/:id", deps.ClientHdl.Update)
		clients.DELETE("/:id", deps.ClientHdl.Delete)
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
		additions.POST("/:id/process", deps.QuotaClientAdditionHdl.Process)
	}

	clientPsre := api.Group("/client-psre")
	{
		clientPsre.POST("/register", deps.ClientPsreHdl.Register)
		clientPsre.GET("/:id", deps.ClientPsreHdl.Get)
	}

	clientCompany := api.Group("/client-companies")
	{
		clientCompany.GET("/", deps.ClientCompanyHdl.List)
		clientCompany.GET("/:id", deps.ClientCompanyHdl.Get)
		clientCompany.POST("/", deps.ClientCompanyHdl.Create)
		clientCompany.PUT("/:id", deps.ClientCompanyHdl.Update)
		clientCompany.DELETE("/:id", deps.ClientCompanyHdl.Delete)
	}
	psre := api.Group("/psre")
	{
		psre.Use(rl.Middleware()) // pasang rate limiter di group ini
		psre.POST("/login", deps.PsreHdl.Login)

		company := psre.Group("/company")
		company.POST("/create", deps.PsreHdl.CreateClientCompany)
	}
	return r
}
