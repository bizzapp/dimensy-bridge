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

	return r
}
