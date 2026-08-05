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
)

func SetupRoutes(deps *config.AppDependencies) *gin.Engine {
	r := gin.Default()

	// Index route - Hello World
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Hello World",
			"service":     "Dimensy Bridge API",
			"version":     "v1.0",
			"description": "Welcome to Dimensy Bridge - Document Management & Digital Signature API",
		})
	})

	// Health check endpoint for Docker health checks
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Service is running",
		})
	})

	r.POST("/api/v1/webhook-notification-document", deps.WebhookHdl.HandlePSRENotification)
	r.POST("/api/v1/webhook-notification-certificate", deps.WebhookHdl.HandlePSRENotification)
	// =======================

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

	rl := middleware.NewRedisRateLimiter(deps.RedisClient, 200, 60*time.Second)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/login", deps.AuthHdl.Login)

		auth.POST("/logout", middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo), deps.AuthHdl.Logout)
	}

	user := api.Group("/users")
	{
		user.GET("", deps.UserHdl.List)
		user.GET("/:id", deps.UserHdl.Get)
		user.POST("", deps.UserHdl.Create)
		user.PUT("/:id", deps.UserHdl.Update)
		user.DELETE("/:id", deps.UserHdl.Delete)
	}

	clients := api.Group("/clients")
	{
		clients.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))

		clients.GET("/profile-psre/:id", deps.ClientPsreHdl.Profile)
		clients.GET("/fill-external-id/:id", deps.ClientPsreHdl.FillExternalId)
		clients.GET("", deps.ClientHdl.List)
		clients.GET("/:id", deps.ClientHdl.Get)
		clients.POST("", deps.ClientHdl.Create)
		clients.PUT("/:id", deps.ClientHdl.Update)
		clients.DELETE("/:id", deps.ClientHdl.Delete)
		clients.POST("/add-quota", deps.ClientHdl.AddQuota)
		clients.POST("/approve-add-quota", deps.ClientHdl.ApproveAddQuota)
	}
	clientSubs := api.Group("/client-subscriptions")
	{
		clientSubs.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		clientSubs.GET("", deps.ClientHasSubscriptionHdl.GetAll)
		clientSubs.GET("/:id", deps.ClientHasSubscriptionHdl.GetByID)
		clientSubs.POST("", deps.ClientHasSubscriptionHdl.Create)
		clientSubs.POST("/:id/process", deps.ClientHasSubscriptionHdl.Process)
		clientSubs.DELETE("/:id", deps.ClientHasSubscriptionHdl.Delete)
	}

	products := api.Group("/products")
	{
		products.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		products.GET("/history", deps.MasterProductHdl.GetHistory)
		products.GET("", deps.MasterProductHdl.List)
		products.GET("/:id", deps.MasterProductHdl.Get)
		products.POST("", deps.MasterProductHdl.Create)
		products.PUT("/:id", deps.MasterProductHdl.Update)
		products.DELETE("/:id", deps.MasterProductHdl.Delete)
	}

	quotas := api.Group("/quotas")
	{
		quotas.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		quotas.GET("/history", deps.QuotaClientHdl.GetHistory)
		quotas.GET("", deps.QuotaClientHdl.List)
		quotas.GET("/:id", deps.QuotaClientHdl.Get)
		quotas.POST("", deps.QuotaClientHdl.Create)
		quotas.PUT("/:id", deps.QuotaClientHdl.Update)
		quotas.DELETE("/:id", deps.QuotaClientHdl.Delete)
	}

	additions := api.Group("/quota-additions")
	{
		additions.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		additions.GET("", deps.QuotaClientAdditionHdl.List)
		additions.GET("/:id", deps.QuotaClientAdditionHdl.Get)
		additions.POST("", deps.QuotaClientAdditionHdl.Create)
		additions.PUT("/:id", deps.QuotaClientAdditionHdl.Update)
		additions.DELETE("/:id", deps.QuotaClientAdditionHdl.Delete)
	}

	reductions := api.Group("/quota-reductions")
	{
		reductions.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		reductions.GET("/chart", deps.QuotaClientReductionHdl.GetChart)
	}

	clientPsre := api.Group("/client-psre")
	{
		clientPsre.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		clientPsre.POST("/register", deps.ClientPsreHdl.Register)
		clientPsre.POST("/register-with-fill-external-id", deps.ClientPsreHdl.RegisterWithFillExternalId)
	}

	clientCompany := api.Group("/client-companies")
	{
		clientCompany.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		clientCompany.GET("", deps.ClientCompanyHdl.List)
		clientCompany.GET("/:id", deps.ClientCompanyHdl.Get)
		clientCompany.POST("", deps.ClientCompanyHdl.Create)
		clientCompany.PUT("/:id", deps.ClientCompanyHdl.Update)
		clientCompany.DELETE("/:id", deps.ClientCompanyHdl.Delete)
	}

	group := api.Group("/client_users")
	{
		group.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		group.GET("", deps.ClientUserHdl.GetAll)
		group.GET("/:id", deps.ClientUserHdl.GetByID)
		group.POST("", deps.ClientUserHdl.Create)
		group.PUT("", deps.ClientUserHdl.Update)
		group.DELETE("/:id", deps.ClientUserHdl.Delete)
	}

	// Client IP Whitelist routes
	clientIPWhitelist := api.Group("/client-ip-whitelist")
	{
		clientIPWhitelist.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		clientIPWhitelist.GET("/client/:client_id", deps.ClientIPWhitelistHdl.GetByClientID)
		clientIPWhitelist.GET("/:id", deps.ClientIPWhitelistHdl.GetByID)
		clientIPWhitelist.POST("", deps.ClientIPWhitelistHdl.Create)
		clientIPWhitelist.PUT("/:id", deps.ClientIPWhitelistHdl.Update)
		clientIPWhitelist.DELETE("/:id", deps.ClientIPWhitelistHdl.Delete)
	}
	// Inventory Master Product routes
	inventory := api.Group("/inventory_master_product")
	{
		inventory.Use(middleware.JWTAuthMiddlewareWithBlacklist(deps.TokenBlacklistRepo))
		inventory.GET("", deps.InventoryMasterProductHdl.Index)
		inventory.GET("/list", deps.InventoryMasterProductHdl.List)
		inventory.POST("/store_or_update", deps.InventoryMasterProductHdl.StoreOrUpdate)
		inventory.GET("/:id/show", deps.InventoryMasterProductHdl.Show)
		inventory.POST("/:id/mark_processed", deps.InventoryMasterProductHdl.MarkAsProcessed)
		inventory.POST("/:id/adjust_stock", deps.InventoryMasterProductHdl.AdjustStock)
		inventory.GET("/logs", deps.InventoryMasterProductHdl.GetLogs)
		inventory.POST("/:id/toggle_priority", deps.InventoryMasterProductHdl.TogglePriority)
		inventory.DELETE("/:id/delete", deps.InventoryMasterProductHdl.Delete)
		inventory.GET("/low_stock/items", deps.InventoryMasterProductHdl.GetLowStockItems)
		inventory.GET("/total/value", deps.InventoryMasterProductHdl.GetTotalValue)
	}

	// Setup PSRE routes
	SetupPsreRoutes(api, deps, rl)
	return r
}
