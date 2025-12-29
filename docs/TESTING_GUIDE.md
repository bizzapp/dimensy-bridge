# Testing Guide

## Manual Testing with cURL

### 1. IP Whitelist CRUD Operations

#### Create Whitelist Entry
```bash
# Basic
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Office Server",
    "is_active": true
  }'

# With IPv6
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "2001:db8::1",
    "description": "IPv6 Server",
    "is_active": true
  }'

# Response (201 Created)
{
  "success": true,
  "message": "IP whitelist created successfully",
  "data": {
    "id": 1,
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Office Server",
    "is_active": true,
    "created_at": "2024-12-29T10:00:00Z",
    "updated_at": "2024-12-29T10:00:00Z"
  }
}
```

#### Get Whitelist by ID
```bash
curl http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response (200 OK)
{
  "success": true,
  "data": {
    "id": 1,
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Office Server",
    "is_active": true,
    "created_at": "2024-12-29T10:00:00Z"
  }
}
```

#### Get All IPs for Client
```bash
# Page 1, 10 items per page (default)
curl "http://localhost:8080/api/v1/client-ip-whitelist/client/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Page 2, 20 items per page
curl "http://localhost:8080/api/v1/client-ip-whitelist/client/1?page=2&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response (200 OK)
{
  "success": true,
  "data": {
    "data": [
      {
        "id": 1,
        "client_id": 1,
        "ip_address": "192.168.1.100",
        "description": "Office Server",
        "is_active": true,
        "created_at": "2024-12-29T10:00:00Z"
      },
      {
        "id": 2,
        "client_id": 1,
        "ip_address": "10.0.0.5",
        "description": "Dev Machine",
        "is_active": true,
        "created_at": "2024-12-29T10:01:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "limit": 10
  }
}
```

#### Update Whitelist
```bash
curl -X PUT http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated Description",
    "is_active": false
  }'

# Response (200 OK)
{
  "success": true,
  "message": "IP whitelist updated successfully",
  "data": {
    "id": 1,
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Updated Description",
    "is_active": false,
    "updated_at": "2024-12-29T10:05:00Z"
  }
}
```

#### Delete Whitelist
```bash
curl -X DELETE http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response (200 OK)
{
  "success": true,
  "message": "IP whitelist deleted successfully"
}
```

---

## Rate Limiter Testing

### Test Rate Limiting

```bash
# Bash script to test rate limiting
#!/bin/bash

echo "Testing Rate Limiter (5 RPS with 10 burst)"
echo "=============================================="

for i in {1..15}; do
  echo "Request $i:"
  curl -i http://localhost:8080/api/v1/clients \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    --silent | head -1
  sleep 0.05  # 50ms between requests = 20 RPS
done

echo ""
echo "Expected: First 10 requests succeed (5 per second + 10 burst)"
echo "          Requests 11-15 get 429 Too Many Requests"
```

### Test from Different IPs

```bash
# Simulate requests from different IPs (requires Docker/localhost tricks)

# IP 1
curl http://localhost:8080/api/v1/clients -H "X-Forwarded-For: 192.168.1.1"

# IP 2  
curl http://localhost:8080/api/v1/clients -H "X-Forwarded-For: 192.168.1.2"

# IP 1 again (should have separate rate limit from IP 2)
curl http://localhost:8080/api/v1/clients -H "X-Forwarded-For: 192.168.1.1"
```

---

## Error Scenario Testing

### Invalid IP Address
```bash
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "invalid.ip",
    "is_active": true
  }'

# Response (400 Bad Request)
{
  "success": false,
  "message": "Invalid request body",
  "error": "Key: 'CreateClientIPWhitelistRequest.IPAddress' Error:Field validation for 'IPAddress' failed on the 'ip' tag"
}
```

### Missing Required Field
```bash
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "ip_address": "192.168.1.100"
  }'

# Response (400 Bad Request)
{
  "success": false,
  "message": "Invalid request body",
  "error": "Key: 'CreateClientIPWhitelistRequest.ClientID' Error:Field validation for 'ClientID' failed on the 'required' tag"
}
```

### IP Not Found
```bash
curl http://localhost:8080/api/v1/client-ip-whitelist/99999 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response (404 Not Found)
{
  "success": false,
  "message": "IP whitelist tidak ditemukan"
}
```

### Rate Limit Exceeded
```bash
# After 5 requests per second from same IP
curl http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response (429 Too Many Requests)
{
  "success": false,
  "message": "Too many requests, please try again later."
}
```

### Invalid Token
```bash
curl http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer INVALID_TOKEN"

# Response (401 Unauthorized)
{
  "success": false,
  "message": "Invalid token"
}
```

---

## Integration Testing with Postman

### Setup Postman Collection

```json
{
  "info": {
    "name": "Client IP Whitelist API",
    "description": "Test collection for IP whitelist endpoints"
  },
  "item": [
    {
      "name": "Create Whitelist",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{jwt_token}}"
          },
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "url": {
          "raw": "http://localhost:8080/api/v1/client-ip-whitelist",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "client-ip-whitelist"]
        },
        "body": {
          "mode": "raw",
          "raw": "{\n  \"client_id\": 1,\n  \"ip_address\": \"192.168.1.100\",\n  \"description\": \"Office\",\n  \"is_active\": true\n}"
        }
      }
    },
    {
      "name": "Get All IPs",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{jwt_token}}"
          }
        ],
        "url": {
          "raw": "http://localhost:8080/api/v1/client-ip-whitelist/client/1?page=1&limit=10",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "client-ip-whitelist", "client", "1"],
          "query": [
            {"key": "page", "value": "1"},
            {"key": "limit", "value": "10"}
          ]
        }
      }
    }
  ]
}
```

---

## Unit Testing Examples

### Test IP Whitelist Service

```go
// internal/service/client_ip_whitelist_service_test.go

package service

import (
	"testing"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIPWhitelistRepo struct {
	mock.Mock
}

func (m *MockIPWhitelistRepo) Create(ip *model.ClientIPWhitelist) (*model.ClientIPWhitelist, error) {
	args := m.Called(ip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ClientIPWhitelist), args.Error(1)
}

func TestCreateIPWhitelist(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepo)
	service := NewClientIPWhitelistService(mockRepo)

	req := &dto.CreateClientIPWhitelistRequest{
		ClientID:    1,
		IPAddress:   "192.168.1.100",
		Description: "Test",
		IsActive:    true,
	}

	mockRepo.On("Create", mock.MatchedBy(func(ip *model.ClientIPWhitelist) bool {
		return ip.ClientID == 1 && ip.IPAddress == "192.168.1.100"
	})).Return(&model.ClientIPWhitelist{
		ID:        1,
		ClientID:  1,
		IPAddress: "192.168.1.100",
	}, nil)

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestIsIPWhitelisted_Success(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepo)
	service := NewClientIPWhitelistService(mockRepo)

	mockRepo.On("IsIPWhitelisted", int64(1), "192.168.1.100").
		Return(true, nil)

	result, err := service.IsIPWhitelisted(1, "192.168.1.100")

	assert.NoError(t, err)
	assert.True(t, result)
	mockRepo.AssertExpectations(t)
}

func TestIsIPWhitelisted_NotFound(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepo)
	service := NewClientIPWhitelistService(mockRepo)

	mockRepo.On("IsIPWhitelisted", int64(1), "192.168.1.100").
		Return(false, nil)

	result, err := service.IsIPWhitelisted(1, "192.168.1.100")

	assert.NoError(t, err)
	assert.False(t, result)
	mockRepo.AssertExpectations(t)
}
```

### Test Rate Limiter

```go
// internal/middleware/redis_ratelimiter_test.go

package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRateLimiter_AllowedRequest(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rl := NewRedisRateLimiter(client, 5, 10)

	// Should allow first request
	key := "test_key"
	allowed := rl.isAllowed(nil, key)
	assert.True(t, allowed)
}

func TestRedisRateLimiter_RateLimit(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rl := NewRedisRateLimiter(client, 2, 3) // 2 RPS, 3 burst

	key := "rate_limit_test"

	// First 3 requests should succeed (2 + 1 burst)
	for i := 0; i < 3; i++ {
		allowed := rl.isAllowed(nil, key)
		assert.True(t, allowed, "Request %d should be allowed", i+1)
	}

	// Fourth request should fail
	allowed := rl.isAllowed(nil, key)
	assert.False(t, allowed, "Request 4 should be blocked")

	// Wait 1 second for token reset
	time.Sleep(time.Second)

	// Should allow again after reset
	allowed = rl.isAllowed(nil, key)
	assert.True(t, allowed, "Request after reset should be allowed")
}
```

---

## Performance Testing

### Load Test with Apache Bench

```bash
# Test without rate limiting
ab -n 1000 -c 10 http://localhost:8080/api/v1/clients

# Test with rate limiting
ab -n 100 -c 5 http://localhost:8080/api/v1/clients
# Expected: Some 429 responses

# Test with concurrent IPs
ab -n 1000 -c 10 -H "X-Forwarded-For: 192.168.1.1" http://localhost:8080/api/v1/clients
```

### Load Test with k6

```javascript
// load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 10 },   // Ramp up
    { duration: '30s', target: 50 },   // Stress
    { duration: '10s', target: 0 },    // Ramp down
  ],
};

export default function () {
  let res = http.get('http://localhost:8080/api/v1/clients', {
    headers: {
      'Authorization': 'Bearer TOKEN',
    },
  });

  check(res, {
    'status is 200 or 429': (r) => r.status === 200 || r.status === 429,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
```

Run:
```bash
k6 run load_test.js
```

---

## Monitoring & Debugging

### Check Redis Keys
```bash
# Connect to Redis
redis-cli

# List all rate limit keys
KEYS rate_limit:*

# Get details of a specific key
ZRANGE rate_limit:192.168.1.1:window 0 -1 WITHSCORES

# Monitor all Redis commands
MONITOR

# Check memory usage
INFO memory

# Clear rate limit for testing
FLUSHDB  # Careful! Clears all keys
```

### Check Database
```bash
# Connect to PostgreSQL
psql -U postgres -d dimensy_bridge

-- Get all whitelisted IPs for client 1
SELECT * FROM client_ip_whitelists 
WHERE client_id = 1 AND deleted_at IS NULL;

-- Count whitelisted IPs
SELECT COUNT(*) FROM client_ip_whitelists 
WHERE client_id = 1 AND is_active = true AND deleted_at IS NULL;

-- Check for duplicates
SELECT client_id, ip_address, COUNT(*) 
FROM client_ip_whitelists 
WHERE deleted_at IS NULL 
GROUP BY client_id, ip_address 
HAVING COUNT(*) > 1;
```

---

## Test Scenarios

### Scenario 1: Normal Client Usage
```bash
# Setup
1. Create client
2. Add 3 IPs to whitelist
3. Make requests from each IP
4. Verify all requests succeed

# Cleanup
5. Delete one IP
6. Request from deleted IP
7. Verify request blocked (if strict mode)
```

### Scenario 2: Rate Limiting
```bash
# Setup
1. Configure rate limit: 5 RPS, 10 burst

# Test
2. Send 15 rapid requests
3. Verify: First 10 succeed, last 5 get 429
4. Wait 1 second
5. Send 5 requests
6. Verify: All 5 succeed (tokens reset)
```

### Scenario 3: IP Rotation
```bash
# Setup
1. Create client with 2 whitelisted IPs
2. Enable strict whitelist mode

# Test
3. Request from IP 1 → Success
4. Request from IP 2 → Success
5. Request from IP 3 → Blocked (403)
6. Add IP 3 to whitelist
7. Request from IP 3 → Success
```

### Scenario 4: Soft Delete
```bash
# Setup
1. Create IP whitelist entry
2. Record ID

# Test
3. Delete whitelist
4. Query by ID → Null (soft delete)
5. Query by client → Not returned
6. Check deleted_at in DB → Has timestamp
```

