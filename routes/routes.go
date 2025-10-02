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

	masterProducts := api.Group("/master-products")
	{
		masterProducts.GET("/", deps.MasterProductHdl.List)
		masterProducts.GET("/:id", deps.MasterProductHdl.Get)
		masterProducts.POST("/", deps.MasterProductHdl.Create)
		masterProducts.PUT("/:id", deps.MasterProductHdl.Update)
		masterProducts.DELETE("/:id", deps.MasterProductHdl.Delete)
	}
	return r
}
