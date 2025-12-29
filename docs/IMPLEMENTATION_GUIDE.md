# Implementation Guide

## Step 1: Update go.mod

Add Redis dependency:

```bash
go get github.com/redis/go-redis/v9
```

Or manually add to go.mod:
```
require github.com/redis/go-redis/v9 v9.0.0
```

## Step 2: Update .env

```env
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Rate Limiter Config
RATE_LIMIT_RPS=5
RATE_LIMIT_BURST=10
```

## Step 3: Update AppDependencies

The file `internal/config/AppDependencies.go` has been updated with:
- ClientIPWhitelistRepo
- ClientIPWhitelistSvc
- ClientIPWhitelistHdl

## Step 4: Update main.go (Example)

```go
package main

import (
	"dimensy-bridge/internal/config"
	"dimensy-bridge/routes"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Database
	db := config.InitDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Initialize Redis (NEW)
	redisClient := config.InitRedis()
	defer config.CloseRedis(redisClient)

	// Run migrations
	config.MigrateDatabase(db)

	// Create dependencies (UPDATED - now includes redis)
	deps := config.NewAppDependencies(db, redisClient)

	// Setup routes
	r := routes.SetupRoutes(deps)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

## Step 5: Update Routes (Already Done)

The file `routes/routes.go` has been updated with:

```go
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
```

## Step 6: Update AppDependencies (if needed)

Current AppDependencies.go structure:

```go
type AppDependencies struct {
	DB *gorm.DB
	RedisClient *redis.Client  // Add this if not using separately
	
	// ... existing fields
	
	// Client IP Whitelist Module (ADDED)
	ClientIPWhitelistRepo repository.ClientIPWhitelistRepository
	ClientIPWhitelistSvc  service.ClientIPWhitelistService
	ClientIPWhitelistHdl  *handler.ClientIPWhitelistHandler
	
	// ... other fields
}
```

## Step 7: Database Migration

Run the migration:

```bash
# Option 1: Using GORM AutoMigrate (in main.go)
config.MigrateDatabase(db)

# Option 2: Manual SQL
psql -U postgres -d dimensy_bridge -f docs/DATABASE_MIGRATION.md
```

## Step 8: Optional - Add Rate Limiter to Routes

### Option A: Global Rate Limiting (All Routes)

```go
// di SetupRoutes
func SetupRoutes(deps *config.AppDependencies) *gin.Engine {
	r := gin.Default()
	
	// Apply rate limiter globally
	redisRL := middleware.NewRedisRateLimiter(
		deps.RedisClient,
		5,  // 5 requests per second
		10, // burst capacity
	)
	r.Use(redisRL.Middleware())
	
	// ... rest of routes
}
```

### Option B: Per-Route Rate Limiting

```go
// di SetupRoutes untuk specific routes
api := r.Group("/api/v1")

// High-frequency endpoints get stricter rate limit
clients := api.Group("/clients")
strictRL := middleware.NewRedisRateLimiter(deps.RedisClient, 3, 5)
clients.Use(strictRL.Middleware())

// Regular endpoints get normal rate limit
users := api.Group("/users")
normalRL := middleware.NewRedisRateLimiter(deps.RedisClient, 5, 10)
users.Use(normalRL.Middleware())
```

## Step 9: Testing

### Verify Redis Connection
```bash
redis-cli ping
# Should return: PONG
```

### Test IP Whitelist API

```bash
# 1. Create a whitelist entry
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Office Server",
    "is_active": true
  }'

# 2. Get all whitelisted IPs for a client
curl http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. Update whitelist entry
curl -X PUT http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_active": false, "description": "Disabled"}'

# 4. Delete whitelist entry
curl -X DELETE http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test Rate Limiter

```bash
# Make 15 rapid requests to trigger rate limit
for i in {1..15}; do
  echo "Request $i:"
  curl http://localhost:8080/api/v1/clients \
    -H "Authorization: Bearer YOUR_TOKEN"
  echo ""
done

# Expected: First 5 succeed, next 5 (burst) succeed, rest get 429
```

## Troubleshooting

### Redis Connection Failed
```
Error: Failed to connect to Redis
Solution: 
- Check Redis is running: redis-cli ping
- Verify REDIS_HOST and REDIS_PORT in .env
- Check firewall/network connectivity
```

### Module Not Found
```
Error: module dimensy-bridge/internal/handler: not found
Solution:
- Run: go mod tidy
- Run: go mod download
```

### Compilation Errors
```
Error: undefined: NewClientIPWhitelistRepository
Solution:
- Check repository/client_ip_whitelist_repo.go is created
- Run: go build ./...
```

## Files Created

1. ✅ `internal/model/client_ip_whitelist_model.go` - Data model
2. ✅ `internal/dto/client_ip_whitelist-dto.go` - DTOs
3. ✅ `internal/repository/client_ip_whitelist_repo.go` - Repository
4. ✅ `internal/service/client_ip_whitelist_service.go` - Service
5. ✅ `internal/handler/client_ip_whitelist_handler.go` - Handler
6. ✅ `internal/middleware/redis_ratelimiter.go` - Redis rate limiter
7. ✅ `internal/middleware/ip_whitelist_middleware.go` - IP whitelist middleware
8. ✅ `internal/config/redis.go` - Redis initialization
9. ✅ `internal/config/AppDependencies.go` - Updated with new modules
10. ✅ `routes/routes.go` - Updated with new routes
11. ✅ `docs/REDIS_RATELIMITER_AND_IP_WHITELIST.md` - Full documentation
12. ✅ `docs/DATABASE_MIGRATION.md` - Migration scripts

