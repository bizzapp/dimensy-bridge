# Quick Reference Guide

## 🎯 Checklist for Integration

### Phase 1: Setup (15 min)
- [ ] Run `go get github.com/redis/go-redis/v9`
- [ ] Start Redis service: `redis-server` or `docker run -d -p 6379:6379 redis:7-alpine`
- [ ] Update `.env` with Redis config
- [ ] Run database migrations

### Phase 2: Code Updates (10 min)
- [ ] Update `main.go` to initialize Redis
- [ ] Update `AppDependencies.go` (already done ✓)
- [ ] Update `routes.go` (already done ✓)

### Phase 3: Testing (10 min)
- [ ] Test IP whitelist CRUD
- [ ] Test rate limiter
- [ ] Verify Redis connection

---

## 📋 Files at a Glance

### Must Read First
```
docs/SUMMARY.md                                    ← Start here
docs/IMPLEMENTATION_GUIDE.md                       ← Step-by-step
docs/REDIS_RATELIMITER_AND_IP_WHITELIST.md        ← Detailed docs
```

### Reference
```
docs/ARCHITECTURE_DIAGRAMS.md                      ← How it works
docs/DATABASE_MIGRATION.md                         ← DB setup
docs/TESTING_GUIDE.md                              ← Test examples
docs/CLIENT_IP_WHITELIST_README.md                 ← Feature overview
```

### Code
```
internal/model/client_ip_whitelist_model.go        ← Data model
internal/dto/client_ip_whitelist-dto.go            ← Request/Response
internal/repository/client_ip_whitelist_repo.go    ← DB operations
internal/service/client_ip_whitelist_service.go    ← Business logic
internal/handler/client_ip_whitelist_handler.go    ← HTTP handlers
internal/middleware/redis_ratelimiter.go           ← Rate limiting
internal/middleware/ip_whitelist_middleware.go     ← IP validation
internal/config/redis.go                           ← Redis setup
```

---

## 🔥 Most Common Commands

### Create IP Whitelist
```bash
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"client_id": 1, "ip_address": "192.168.1.100", "is_active": true}'
```

### Get All IPs for Client
```bash
curl http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Check Redis Connection
```bash
redis-cli ping
```

### Check Rate Limit Keys
```bash
redis-cli KEYS "rate_limit:*"
```

### Check Database Table
```bash
psql -c "SELECT * FROM client_ip_whitelists LIMIT 5;"
```

---

## 🔧 Configuration Reference

### Rate Limiter Defaults
```
RPS: 5 (requests per second)
Burst: 10 (max concurrent)
Window: 1 second
Expiry: 2 seconds
```

### Adjust Rate Limiter
```go
// In routes.go or where middleware is created
rl := middleware.NewRedisRateLimiter(
    redisClient,
    5,   // RPS (requests per second)
    10,  // Burst capacity
)
```

### Environment Variables
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## 📊 API Response Examples

### Success (201 Created)
```json
{
  "success": true,
  "message": "IP whitelist created successfully",
  "data": {
    "id": 1,
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "is_active": true,
    "created_at": "2024-12-29T10:00:00Z"
  }
}
```

### Success (200 OK - List)
```json
{
  "success": true,
  "data": {
    "data": [
      {"id": 1, "client_id": 1, "ip_address": "192.168.1.100", "is_active": true},
      {"id": 2, "client_id": 1, "ip_address": "10.0.0.5", "is_active": true}
    ],
    "total": 2,
    "page": 1,
    "limit": 10
  }
}
```

### Rate Limit (429)
```json
{
  "success": false,
  "message": "Too many requests, please try again later."
}
```

### IP Not Whitelisted (403)
```json
{
  "success": false,
  "message": "Your IP address is not whitelisted",
  "ip": "192.168.1.50"
}
```

### Invalid IP (400)
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Key: 'CreateClientIPWhitelistRequest.IPAddress' Error:Field validation for 'IPAddress' failed on the 'ip' tag"
}
```

---

## 🚨 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Redis connection refused | Check Redis running: `redis-cli ping` |
| Rate limiter not working | Verify middleware registered on routes |
| IP whitelist not blocking | Check `is_active=true` in database |
| 401 Unauthorized | Verify JWT token valid |
| Database table not found | Run migrations: See DATABASE_MIGRATION.md |
| Duplicate IP error | Use unique constraint: See DATABASE_MIGRATION.md |

---

## 🧪 Quick Test

```bash
# 1. Create whitelist
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"client_id": 1, "ip_address": "192.168.1.100", "is_active": true}')
echo $RESPONSE | jq .

# 2. Get all
curl -s http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer $TOKEN" | jq .

# 3. Test rate limiter
for i in {1..15}; do
  echo "Request $i:"
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/api/v1/clients \
    -H "Authorization: Bearer $TOKEN"
done
```

---

## 📍 Service Methods Cheat Sheet

### ClientIPWhitelistService
```go
Create(req *CreateClientIPWhitelistRequest) (*ClientIPWhitelistResponse, error)
GetByID(id int64) (*ClientIPWhitelistResponse, error)
GetByClientID(clientID int64, page, limit int) (*ListClientIPWhitelistResponse, error)
Update(id int64, req *UpdateClientIPWhitelistRequest) (*ClientIPWhitelistResponse, error)
Delete(id int64) error
IsIPWhitelisted(clientID int64, ipAddress string) (bool, error)
GetActiveIPsByClientID(clientID int64) ([]string, error)
```

### RedisRateLimiter
```go
NewRedisRateLimiter(client *redis.Client, rps int, burst int) *RedisRateLimiter
rl.Middleware() gin.HandlerFunc
rl.isAllowed(c *gin.Context, key string) bool
```

---

## 🏃 Integration Checklist

```markdown
# Step-by-step integration:

## 1. Dependencies ✓
- [x] redis/go-redis imported
- [x] Model created
- [x] DTO created
- [x] Repository created
- [x] Service created
- [x] Handler created
- [x] Middleware created

## 2. AppDependencies ✓
- [x] Repo added
- [x] Service added
- [x] Handler added
- [x] Initialized in NewAppDependencies

## 3. Routes ✓
- [x] Routes added
- [x] Middleware protection added
- [x] CRUD endpoints configured

## 4. Configuration
- [ ] .env updated with REDIS_*
- [ ] main.go updated (need to do)
- [ ] Migrations run (need to do)
- [ ] Redis started (need to do)

## 5. Testing
- [ ] CREATE endpoint tested
- [ ] GET endpoint tested
- [ ] List endpoint tested
- [ ] UPDATE endpoint tested
- [ ] DELETE endpoint tested
- [ ] Rate limiter tested
```

---

## 🎓 Learning Path

1. **Start here**: `docs/SUMMARY.md`
2. **How it works**: `docs/ARCHITECTURE_DIAGRAMS.md`
3. **Setup**: `docs/IMPLEMENTATION_GUIDE.md`
4. **Details**: `docs/REDIS_RATELIMITER_AND_IP_WHITELIST.md`
5. **Test**: `docs/TESTING_GUIDE.md`
6. **Database**: `docs/DATABASE_MIGRATION.md`

---

## 💡 Pro Tips

1. **Use Optional Whitelist** - More flexible than strict mode
2. **Monitor Redis memory** - Keys expire after 2 seconds but check `MEMORY USAGE`
3. **Index optimization** - Ensure `client_id` & `ip_address` indexed for speed
4. **Rate limit tiers** - Different RPS for different endpoint groups
5. **Test with JMeter** - For load testing beyond cURL

---

## 📞 Support

- See **docs/** folder for detailed documentation
- Check **TESTING_GUIDE.md** for test examples
- Review **IMPLEMENTATION_GUIDE.md** for setup help

---

## 🚀 Ready to Go!

All files created and integrated. Just need to:

1. Update `main.go` to initialize Redis
2. Run database migrations
3. Test endpoints
4. Deploy!

See **IMPLEMENTATION_GUIDE.md** for detailed instructions.

