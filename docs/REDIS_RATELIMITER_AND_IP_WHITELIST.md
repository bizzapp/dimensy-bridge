# Redis Rate Limiter & IP Whitelist Setup Guide

## Overview

Sistem ini menyediakan 2 solusi utama:

1. **Client IP Whitelist** - Setiap client bisa punya lebih dari 1 IP yang di-whitelist
2. **Redis Rate Limiter** - Distributed rate limiting menggunakan Redis

## 1. Client IP Whitelist Module

### Database Model

```go
type ClientIPWhitelist struct {
    ID          int64
    ClientID    int64
    IPAddress   string      // IPv4 atau IPv6
    Description string      // Deskripsi IP
    IsActive    bool        // Status aktif/tidak
    CreatedAt   *time.Time
    UpdatedAt   *time.Time
    DeletedAt   *time.Time
}
```

### Files Created

- `internal/model/client_ip_whitelist_model.go` - Data model
- `internal/dto/client_ip_whitelist-dto.go` - Request/Response DTOs
- `internal/repository/client_ip_whitelist_repo.go` - Database operations
- `internal/service/client_ip_whitelist_service.go` - Business logic
- `internal/handler/client_ip_whitelist_handler.go` - HTTP handlers

### API Endpoints

```
POST   /api/v1/client-ip-whitelist              - Create whitelist entry
GET    /api/v1/client-ip-whitelist/:id          - Get by ID
GET    /api/v1/client-ip-whitelist/client/:client_id  - Get all IPs for client
PUT    /api/v1/client-ip-whitelist/:id          - Update whitelist
DELETE /api/v1/client-ip-whitelist/:id          - Delete whitelist
```

### Usage Example

#### Create Whitelist Entry
```bash
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Office Server",
    "is_active": true
  }'
```

#### Get All IPs for Client
```bash
curl http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json"
```

### Database Migration

Run this SQL to create the table:

```sql
CREATE TABLE client_ip_whitelists (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_client_ip_whitelists_client_id 
        FOREIGN KEY (client_id) REFERENCES clients(id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX idx_client_ip_whitelists_client_id ON client_ip_whitelists(client_id);
CREATE INDEX idx_client_ip_whitelists_ip_address ON client_ip_whitelists(ip_address);
CREATE INDEX idx_client_ip_whitelists_deleted_at ON client_ip_whitelists(deleted_at);
```

---

## 2. Redis Rate Limiter

### Architecture

**Token Bucket Algorithm** dengan Redis sebagai storage:
- Setiap client IP mendapat token bucket
- Tokens direset setiap detik
- Burst capacity memungkinkan beberapa request sekaligus

### Configuration

Update `go.mod`:
```go
require (
    github.com/redis/go-redis/v9 v9.0.0+
)
```

### Setup Redis Connection

Update `internal/config/db.go` atau file config lainnya:

```go
import "github.com/redis/go-redis/v9"

func InitRedis() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
    
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        panic("Failed to connect to Redis: " + err.Error())
    }
    
    return client
}
```

### Update AppDependencies

```go
// Di AppDependencies struct
type AppDependencies struct {
    RedisClient *redis.Client
    RedisRateLimiter *middleware.RedisRateLimiter
    // ... existing fields
}

// Di NewAppDependencies
func NewAppDependencies(db *gorm.DB, redisClient *redis.Client) *AppDependencies {
    redisRL := middleware.NewRedisRateLimiter(redisClient, 5, 10)
    // ... rest of initialization
}
```

### Update Main Application

```go
// cmd/server/main.go
func main() {
    // ... existing code
    
    // Initialize Redis
    redisClient := config.InitRedis()
    defer redisClient.Close()
    
    // Initialize dependencies
    deps := config.NewAppDependencies(db, redisClient)
    
    // Setup routes
    r := routes.SetupRoutes(deps)
    
    // ... rest of code
}
```

### Middleware Usage

#### Global Rate Limiting

```go
// di routes.go
func SetupRoutes(deps *config.AppDependencies) *gin.Engine {
    r := gin.Default()
    
    // Apply global rate limiting
    r.Use(deps.RedisRateLimiter.Middleware())
    
    // ... rest of routes
}
```

#### Per-Route Rate Limiting

```go
// Untuk specific routes yang perlu rate limiting
api := r.Group("/api/v1")
api.Use(deps.RedisRateLimiter.Middleware())

// Or untuk route-specific
clients := api.Group("/clients")
clients.Use(deps.RedisRateLimiter.Middleware())
```

### Rate Limiter Parameters

Current configuration (dari routes.go):
- **RPS (Requests Per Second)**: 5 
- **Burst Capacity**: 10
- **Window Size**: 1 second

Modifiable via:
```go
redisRL := middleware.NewRedisRateLimiter(
    redisClient,
    5,   // rps
    10,  // burst
)
```

---

## 3. IP Whitelist Middleware

### Two Types of Middleware

#### 1. Strict IP Whitelist (Blocks non-whitelisted IPs)

```go
import "dimensy-bridge/internal/middleware"

// di routes.go
clients := api.Group("/clients")
clients.Use(middleware.IPWhitelistMiddleware(deps.ClientIPWhitelistRepo))
clients.GET("", deps.ClientHdl.List)
```

#### 2. Optional IP Whitelist (Only checks if whitelist exists)

```go
// di routes.go
clients := api.Group("/clients")
clients.Use(middleware.OptionalIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))
clients.GET("", deps.ClientHdl.List)
```

### Middleware Flow

1. Extract `client_id` dari request (query param, URL param, atau token)
2. Get active IPs dari database untuk client tersebut
3. Check apakah current request IP ada dalam whitelist
4. Allow/Block request sesuai middleware type

---

## 4. Environment Variables

Add to `.env`:

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

---

## 5. Docker Compose Setup (Optional)

```yaml
version: '3.9'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    environment:
      - REDIS_PASSWORD=
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

volumes:
  redis_data:
```

Run: `docker-compose up -d`

---

## 6. Testing

### Test IP Whitelist

```bash
# Create whitelist entry
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"client_id": 1, "ip_address": "127.0.0.1", "is_active": true}'

# Get all whitelisted IPs for client
curl http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer {token}"
```

### Test Rate Limiter

```bash
# Make rapid requests to trigger rate limit
for i in {1..15}; do
  curl http://localhost:8080/api/v1/clients
done
# After 5 requests per second, others should return 429 Too Many Requests
```

---

## 7. Best Practices

1. **Use Optional Whitelist by default** - Lebih fleksibel
2. **Combine Rate Limiting + Whitelist** - Defense in depth
3. **Monitor Redis Memory** - Rate limiter entries bisa accumulate
4. **Backup Redis Data** - Use persistence (`appendonly yes`)
5. **Set TTL on Rate Limit Keys** - Automatic cleanup setiap 2 detik

---

## 8. Troubleshooting

### Redis Connection Failed
```
Error: redis: connection refused
Solution: Pastikan Redis running di port 6379
```

### Rate Limit Not Working
```
Check:
- Redis connection active
- Middleware properly registered
- Correct RPS/Burst values
```

### IP Whitelist Not Validating
```
Check:
- client_id parameter present di request
- IP address format valid (IPv4/IPv6)
- is_active = true di database
```

