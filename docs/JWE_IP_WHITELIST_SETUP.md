# JWE-Based IP Whitelist Setup

## Overview

Setup ketat untuk IP whitelist di PSRE routes. Setiap request harus:
1. Memiliki valid JWE token
2. IP address dari request harus terdaftar di database untuk client_id yang ada di token
3. Status is_active harus true

Jika ada yang tidak sesuai → **403 Forbidden**

---

## Architecture

```
Request dengan JWE Token
        ↓
AuthJWE() Middleware
    ├─ Validate & decode token
    ├─ Extract client_id
    └─ Set ke context: c.Set("client_id", clientID)
        ↓
JWEIPWhitelistMiddleware
    ├─ Get client_id dari context
    ├─ Get IP dari request
    ├─ Query: SELECT * FROM client_ip_whitelists
    │         WHERE client_id = ? AND ip_address = ? AND is_active = true
    ├─ Tidak ada? → Return 403 Forbidden
    └─ Ada? → Continue ke handler
        ↓
Handler (Proses Request)
```

---

## Implementation Details

### 1. Middleware: JWEIPWhitelistMiddleware

**Location**: `internal/middleware/ip_whitelist_middleware.go`

```go
func JWEIPWhitelistMiddleware(ipWhitelistRepo repository.ClientIPWhitelistRepository) gin.HandlerFunc {
    // Extract client_id dari JWE context
    clientID, exists := c.Get("client_id")
    
    // Get IP dari request
    ip := c.ClientIP()
    
    // Check database: is IP registered & active for this client?
    isWhitelisted, err := ipWhitelistRepo.IsIPWhitelisted(clientID, ip)
    
    // Jika tidak whitelisted → 403 Forbidden
    if !isWhitelisted {
        return 403 error
    }
    
    // Set ke context untuk handler
    c.Set("client_ip", ip)
}
```

### 2. Routes Integration

**Location**: `routes/routes_psre.go`

Applied ke semua resource groups:

```go
// Company routes
company := psre.Group("/company")
company.Use(middleware.AuthJWE())
company.Use(middleware.JWEIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))

// User routes
user := psre.Group("/user")
user.Use(middleware.AuthJWE())
user.Use(middleware.JWEIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))

// Certificate routes
certificate := psre.Group("/certificate")
certificate.Use(middleware.AuthJWE())
certificate.Use(middleware.JWEIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))

// Document routes
document := psre.Group("/document")
document.Use(middleware.AuthJWE())
document.Use(middleware.JWEIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))

// Client routes
client := psre.Group("/client")
client.POST("/login", ...) // Login tanpa IP check
client.Use(middleware.AuthJWE())
client.Use(middleware.JWEIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))
// Endpoints setelahnya semua protected

// Backend routes
backend := psre.Group("/backend")
backend.POST("/login", ...) // Login tanpa IP check
backend.Use(middleware.AuthJWE())
backend.Use(middleware.JWEIPWhitelistMiddleware(deps.ClientIPWhitelistRepo))
// Endpoints setelahnya semua protected
```

---

## API Flow

### 1. Register IP (Admin)

```bash
POST /api/v1/client-ip-whitelist
Content-Type: application/json
Authorization: Bearer {admin_token}

{
  "client_id": 1,
  "ip_address": "192.168.1.100",
  "description": "Office Server",
  "is_active": true
}

Response (201 Created):
{
  "success": true,
  "message": "IP whitelist created successfully",
  "data": {
    "id": 1,
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Office Server",
    "is_active": true
  }
}
```

### 2. Client Login (PSRE)

```bash
POST /api/v1/psre/client/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}

Response (200 OK):
{
  "success": true,
  "data": {
    "token": "eyJhbGc...", // JWE token
    "client_id": 1,
    "client_name": "PT ABC"
  }
}
```

### 3. Client Request Resource (Protected)

```bash
GET /api/v1/psre/user/
Authorization: Bearer eyJhbGc...

Request comes from: 192.168.1.100
Client ID in token: 1

Middleware checks:
  1. Decode & validate token ✓
  2. Extract client_id = 1 ✓
  3. Get request IP = 192.168.1.100 ✓
  4. Query DB: SELECT * FROM client_ip_whitelists 
              WHERE client_id = 1 
              AND ip_address = '192.168.1.100' 
              AND is_active = true
  5. Found? ✓ → Continue to handler
  
Response (200 OK):
[... user list ...]
```

### 4. Client Request from Unregistered IP

```bash
GET /api/v1/psre/user/
Authorization: Bearer eyJhbGc...

Request comes from: 10.0.0.50
Client ID in token: 1

Middleware checks:
  1. Decode & validate token ✓
  2. Extract client_id = 1 ✓
  3. Get request IP = 10.0.0.50 ✓
  4. Query DB: SELECT * FROM client_ip_whitelists 
              WHERE client_id = 1 
              AND ip_address = '10.0.0.50' 
              AND is_active = true
  5. NOT Found? ✗ → Return 403 Forbidden
  
Response (403 Forbidden):
{
  "success": false,
  "message": "Your IP address is not whitelisted. Please register your IP in the system.",
  "ip": "10.0.0.50"
}
```

---

## Database Query

### Check IP Whitelisted

```go
// Repository method
func (r *clientIPWhitelistRepository) IsIPWhitelisted(clientID int64, ipAddress string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.ClientIPWhitelist{}).
		Where("client_id = ? AND ip_address = ? AND is_active = ?", clientID, ipAddress, true).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
```

**SQL Generated**:
```sql
SELECT COUNT(*) FROM client_ip_whitelists
WHERE client_id = $1 
  AND ip_address = $2 
  AND is_active = true
  AND deleted_at IS NULL
```

### Get All Active IPs for Client

```go
func (r *clientIPWhitelistRepository) GetActiveIPsByClientID(clientID int64) ([]string, error) {
	var ips []string
	if err := r.db.Model(&model.ClientIPWhitelist{}).
		Where("client_id = ? AND is_active = ?", clientID, true).
		Pluck("ip_address", &ips).Error; err != nil {
		return nil, err
	}
	return ips, nil
}
```

**SQL Generated**:
```sql
SELECT ip_address FROM client_ip_whitelists
WHERE client_id = $1 
  AND is_active = true
  AND deleted_at IS NULL
```

---

## Testing

### Setup Test Data

```bash
# 1. Create client (already exists)
# Client ID: 1

# 2. Register IP for client 1
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "192.168.1.100",
    "description": "Test Office",
    "is_active": true
  }'

# 3. Login client (get JWE token)
curl -X POST http://localhost:8080/api/v1/psre/client/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@test.com",
    "password": "password123"
  }'
# Response: token = "eyJhbGc..."

# 4. Request resource dengan registered IP (SUKSES)
curl http://localhost:8080/api/v1/psre/user/ \
  -H "Authorization: Bearer eyJhbGc..."
# Response: 200 OK (user list)

# 5. Request dari unregistered IP (GAGAL)
curl http://localhost:8080/api/v1/psre/user/ \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "X-Forwarded-For: 10.0.0.50"
# Response: 403 Forbidden
```

### cURL Test Script

```bash
#!/bin/bash

# Variables
ADMIN_TOKEN="your_admin_token"
CLIENT_ID=1
REGISTERED_IP="192.168.1.100"
UNREGISTERED_IP="10.0.0.50"

echo "=== Test 1: Register IP ==="
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"client_id\": $CLIENT_ID,
    \"ip_address\": \"$REGISTERED_IP\",
    \"is_active\": true
  }" | jq .

echo -e "\n=== Test 2: Login Client ==="
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/psre/client/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@test.com",
    "password": "password"
  }' | jq -r '.data.token')

echo "Token: $TOKEN"

echo -e "\n=== Test 3: Access with Registered IP ==="
curl -v http://localhost:8080/api/v1/psre/user/ \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n=== Test 4: Access with Unregistered IP ==="
curl -v http://localhost:8080/api/v1/psre/user/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Forwarded-For: $UNREGISTERED_IP"
```

---

## Whitelist Management

### Get All IPs for Client

```bash
curl http://localhost:8080/api/v1/client-ip-whitelist/client/1 \
  -H "Authorization: Bearer $TOKEN"

Response:
{
  "success": true,
  "data": {
    "data": [
      {
        "id": 1,
        "client_id": 1,
        "ip_address": "192.168.1.100",
        "description": "Office",
        "is_active": true
      },
      {
        "id": 2,
        "client_id": 1,
        "ip_address": "10.0.0.5",
        "description": "Dev",
        "is_active": true
      }
    ],
    "total": 2,
    "page": 1,
    "limit": 10
  }
}
```

### Add Another IP

```bash
curl -X POST http://localhost:8080/api/v1/client-ip-whitelist \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "ip_address": "203.0.113.42",
    "description": "Remote Office",
    "is_active": true
  }'
```

### Disable IP Temporarily

```bash
curl -X PUT http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

### Delete IP Permanently

```bash
curl -X DELETE http://localhost:8080/api/v1/client-ip-whitelist/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## Context Values

Middleware menambahkan ke context:

```go
c.Set("client_id", clientID)    // int64 - from JWE token
c.Set("client_ip", ip)           // string - dari request
```

Handler bisa access:

```go
func (h *Handler) SomeAction(c *gin.Context) {
    // Get dari context
    clientID, _ := c.Get("client_id")     // int64
    clientIP, _ := c.Get("client_ip")     // string
    
    // Use untuk logging, audit trail, dll
    log.Printf("Client %d accessed from IP %s", clientID, clientIP)
}
```

---

## Security Notes

1. **Strict Mode**: Tidak ada fallback jika IP tidak registered
2. **Token-Based**: client_id extracted dari JWE token (tidak dari user input)
3. **Case Sensitive**: IP address validation case-sensitive
4. **IPv4 & IPv6**: Support keduanya (validated by GORM binding:"ip")
5. **Soft Delete**: Deleted IPs tidak terlihat di query
6. **Audit Trail**: created_at, updated_at, deleted_at tersimpan

---

## Endpoints Protected

### Sebelum Perubahan (No IP Check)

```
POST   /api/v1/psre/client/login        ← Login only
POST   /api/v1/psre/backend/login       ← Login only
```

### Sesudah Perubahan (With IP Check)

```
PROTECTED (IP Whitelist Required):
├─ /api/v1/psre/company/*               ← All company endpoints
├─ /api/v1/psre/user/*                  ← All user endpoints
├─ /api/v1/psre/certificate/*           ← All certificate endpoints
├─ /api/v1/psre/document/*              ← All document endpoints
├─ /api/v1/psre/client/*                ← Except /login
└─ /api/v1/psre/backend/*               ← Except /login

UNPROTECTED (No IP Check):
├─ POST /api/v1/psre/client/login
└─ POST /api/v1/psre/backend/login
```

---

## Error Responses

### 403 Forbidden - IP Not Whitelisted

```json
{
  "success": false,
  "message": "Your IP address is not whitelisted. Please register your IP in the system.",
  "ip": "10.0.0.50"
}
```

### 401 Unauthorized - Invalid Token

```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

### 401 Unauthorized - No Token

```json
{
  "success": false,
  "message": "Authorization header required"
}
```

### 500 Internal Server Error - DB Error

```json
{
  "success": false,
  "message": "Failed to verify IP whitelist",
  "error": "database connection error"
}
```

---

## Summary

✅ **Implementasi Lengkap:**
- JWE token parsing untuk extract client_id
- Strict IP whitelist check (403 jika IP tidak registered)
- Applied ke semua PSRE endpoints (kecuali login)
- Database query optimized dengan index
- Full audit trail (created_at, updated_at, deleted_at)
- Context management untuk handler use

🔒 **Security:**
- Token-based client identification
- Strict matching (no IP can access tanpa registered)
- Soft delete untuk maintain history
- Error messages helpful but secure

