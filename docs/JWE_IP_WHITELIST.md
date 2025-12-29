# JWE-based IP Whitelist Implementation

## Overview

Implementasi IP whitelist yang menggunakan JWE token untuk ekstrak client information dan validasi IP address terhadap database `client_ip_whitelists`.

## Flow

```
Request dengan JWE Token
        ↓
Extract Authorization header
        ↓
Decrypt JWE token
        ↓
Extract data.id (external_id) dari payload
        ↓
Query client_psre table untuk dapatkan client_id
        ↓
Get active IPs dari client_ip_whitelists
        ↓
Validasi request IP ada di whitelist
        ↓
Allow/Block (403 Forbidden jika tidak)
```

## JWE Token Payload Format

```json
{
  "data": {
    "id": "7c111ffa-e965-48c1-b803-ebb70878dacb",  // external_id di client_psre
    "name": "R17 Group"
  },
  "exp": 1767076711,
  "iat": 1766990311,
  "iss": "Dimensy"
}
```

## Middleware Implementation

### Function: `JWEIPWhitelistWithClientPsreMiddleware`

```go
middleware.JWEIPWhitelistWithClientPsreMiddleware(
    deps.ClientIPWhitelistRepo,
    deps.ClientPsreRepo,
)
```

**Fitur:**
- Extract & decrypt JWE token dari Authorization header
- Parse `data.id` sebagai external_id
- Query `client_psre` table untuk dapatkan `client_id`
- Cek IP whitelist di `client_ip_whitelists` table
- STRICT MODE: Block jika:
  - Tidak ada IP yang terdaftar untuk client
  - Request IP tidak ada di whitelist

**Error Responses:**

| Status | Condition |
|--------|-----------|
| 401 | Authorization header missing |
| 401 | Invalid/expired JWE token |
| 401 | Data.id not found in token |
| 401 | Client not found in database |
| 403 | No IP addresses registered for account |
| 403 | Request IP not whitelisted |
| 500 | Database query error |

## Integration dengan PSRE Routes

Middleware sudah diterapkan di `/api/v1/psre` group:

```go
psre := api.Group("/psre")
psre.Use(rl.Middleware())
psre.Use(middleware.JWEIPWhitelistWithClientPsreMiddleware(
    deps.ClientIPWhitelistRepo,
    deps.ClientPsreRepo,
))
```

Semua PSRE routes akan:
1. Rate limited (5 RPS, 10 burst)
2. Require valid JWE token
3. Require registered IP address

## Context Values Set

Middleware set nilai-nilai ke gin.Context:

```go
c.Set("client_id", int64)           // ID dari clients table
c.Set("client_ip", string)          // Request IP address
c.Set("external_id", string)        // External ID dari client_psre
c.Set("client_name", interface{})   // Name dari token
```

Dapat diakses di handlers:

```go
func (h *Handler) GetData(c *gin.Context) {
    clientID := c.GetInt64("client_id")
    ip := c.GetString("client_ip")
    externalID := c.GetString("external_id")
    
    // Use values in handler logic
}
```

## Database Tables Involved

### 1. client_psre
```
Columns: id, client_id, external_id, ...
Query: SELECT client_id FROM client_psre WHERE external_id = ?
```

### 2. client_ip_whitelists
```
Columns: id, client_id, ip_address, is_active, ...
Query: SELECT ip_address FROM client_ip_whitelists 
       WHERE client_id = ? AND is_active = true AND deleted_at IS NULL
```

### 3. clients (via relationship)
```
Used to get client information through client_psre
```

## Testing

### Success Case (IP Whitelisted)

```bash
# Assume JWE token with data.id = "7c111ffa-e965-48c1-b803-ebb70878dacb"
# And IP 192.168.1.100 is whitelisted for this client

curl -X GET http://localhost:8080/api/v1/psre/user \
  -H "Authorization: Bearer eyJ..." \
  -H "Content-Type: application/json"

# Response: 200 OK with user data
```

### Error Case 1 - No IP Registered

```bash
# Client tidak memiliki IP yang terdaftar

curl -X GET http://localhost:8080/api/v1/psre/user \
  -H "Authorization: Bearer eyJ..." 

# Response: 403 Forbidden
{
  "success": false,
  "message": "No IP addresses registered for this account. Please register your IP first.",
  "ip": "192.168.1.100"
}
```

### Error Case 2 - IP Not Whitelisted

```bash
# Request dari IP yang tidak terdaftar

curl -X GET http://localhost:8080/api/v1/psre/user \
  -H "Authorization: Bearer eyJ..."

# Response: 403 Forbidden
{
  "success": false,
  "message": "Your IP address is not whitelisted",
  "ip": "192.168.1.100"
}
```

### Error Case 3 - Invalid Token

```bash
curl -X GET http://localhost:8080/api/v1/psre/user \
  -H "Authorization: Bearer invalid_token"

# Response: 401 Unauthorized
{
  "success": false,
  "message": "Invalid or expired token",
  "error": "failed to decrypt: ..."
}
```

## Configuration

No additional configuration needed. The middleware:
- Uses existing `ClientIPWhitelistRepository`
- Uses existing `ClientPsreRepository`
- Uses `VerifyJWE` from `pkg/utils/jwe-utils.go`
- Integrates with existing rate limiter

## Key Points

✅ STRICT MODE - Blocks jika tidak ada IP registered
✅ JWE-based - Secure token handling dengan encryption
✅ Database lookup - Verifies client exists before checking IP
✅ Context tracking - Sets client info untuk use di handlers
✅ Rate limited - Combined dengan existing rate limiter
✅ Error handling - Clear error messages untuk debugging

## Files Modified

- `internal/middleware/ip_whitelist_middleware.go` - Added new middleware
- `routes/routes_psre.go` - Integrated middleware to /psre routes

