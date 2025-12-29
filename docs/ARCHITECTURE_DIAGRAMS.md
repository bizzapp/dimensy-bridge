# Architecture & Flow Diagrams

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Request                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────┐
        │  Rate Limiter Middleware       │
        │  (Redis Token Bucket)          │
        │                                │
        │  Check: rate_limit:{ip}:window │
        └────────────────┬───────────────┘
                         │
              ┌──────────▼──────────┐
              │  Allowed?           │
              └──┬───────────────┬──┘
                 │ NO            │ YES
                 │               │
            429  │               │
          (Abort)│               ▼
                 │    ┌─────────────────────────────┐
                 │    │ IP Whitelist Middleware     │
                 │    │ (Optional)                  │
                 │    │                             │
                 │    │ Query: SELECT * FROM        │
                 │    │ client_ip_whitelists        │
                 │    └────────────┬────────────────┘
                 │                 │
                 │      ┌──────────▼──────────┐
                 │      │  Whitelisted?       │
                 │      └──┬───────────────┬──┘
                 │         │ (if exists)   │
                 │         │ NO            │ YES
                 │         │               │
                 │    403  │               │
                 │  (Abort)│               ▼
                 │         │    ┌──────────────────────┐
                 │         │    │ JWT Auth Middleware  │
                 │         │    └──────────┬───────────┘
                 │         │               │
                 │         │    ┌──────────▼──────────┐
                 │         │    │  JWT Valid?         │
                 │         │    └──┬───────────────┬──┘
                 │         │       │ NO            │ YES
                 │         │       │               │
                 │         │  401  │               ▼
                 │         │(Abort)│    ┌──────────────────┐
                 │         │       │    │  Handler Layer   │
                 │         │       │    │                  │
                 │         │       │    │  Process Request │
                 │         │       │    └────────┬─────────┘
                 │         │       │             │
                 │         │       │             ▼
                 ▼         ▼       ▼    ┌────────────────┐
              ┌──────────────────────┐  │   Response     │
              │   Error Response     │  │   (200/201)    │
              │   (429/403/401/500)  │  └────────────────┘
              └──────────────────────┘
```

## Request Flow - IP Whitelist Validation

```
GET /api/v1/clients?client_id=1
     │
     ▼
IPWhitelistMiddleware
     │
     ├─→ Extract client_id from query/params
     │
     ├─→ Get client_id from JWT token context (optional)
     │
     ├─→ Query DB: SELECT ip_address FROM client_ip_whitelists
     │            WHERE client_id = 1 AND is_active = true AND deleted_at IS NULL
     │
     ├─→ Check if request IP in whitelist
     │
     ├─→ YES: Next()
     │
     └─→ NO: Return 403 Forbidden
```

## Rate Limiter - Token Bucket Algorithm

```
Time: T0
Request IP: 192.168.1.1

Redis Entry: rate_limit:192.168.1.1:window
ZSET:
  Member 1: T0 - 0.001 (Score: T0 - 0.001)
  Member 2: T0 - 0.002 (Score: T0 - 0.002)
  Member 3: T0 - 0.003 (Score: T0 - 0.003)
  Member 4: T0 - 0.004 (Score: T0 - 0.004)
  Member 5: T0 - 0.005 (Score: T0 - 0.005)
  
  Current Window: T0 to T0 + 1 second
  Count: 5 (exactly at RPS limit)
  
  Next Request: ZRANGE gets count, if count < RPS, allow
              : if count >= RPS, reject (429)
              : ZADD new entry
              : ZREM entries older than 1 second
              : EXPIRE set to 2 seconds
```

## Database - Client IP Whitelist

```
                ┌─────────────────────┐
                │  clients table      │
                │  ──────────────     │
                │  id (PK)            │
                │  company_name       │
                │  pic_name           │
                │  user_id            │
                └──────────┬──────────┘
                           │ (1:N)
                           │ FK: client_id
                           │
                ┌──────────▼────────────────────┐
                │client_ip_whitelists table     │
                │────────────────────────────   │
                │id (PK)                        │
                │client_id (FK, indexed)        │
                │ip_address (indexed)           │
                │description                    │
                │is_active (boolean)            │
                │created_at                     │
                │updated_at                     │
                │deleted_at (soft delete idx)   │
                └───────────────────────────────┘
```

## Data Model Relationships

```
┌─────────────────────────────────────────────────────┐
│                    User                             │
│                                                     │
│  id (PK)                                            │
│  email                                              │
│  password                                           │
└────────────┬────────────────────────────────────────┘
             │ (1:N)
             │
┌────────────▼────────────────────────────────────────┐
│                    Client                           │
│                                                     │
│  id (PK)                                            │
│  company_name                                       │
│  user_id (FK)                                       │
└────────────┬────────────────────────────────────────┘
             │ (1:N)
             │
┌────────────▼────────────────────────────────────────┐
│              ClientIPWhitelist (NEW)                │
│                                                     │
│  id (PK)                                            │
│  client_id (FK, indexed) ← Query by client         │
│  ip_address (indexed) ← Query by IP                │
│  description                                       │
│  is_active                                          │
│  created_at, updated_at, deleted_at                │
└─────────────────────────────────────────────────────┘
```

## Rate Limiter - Redis Storage Structure

```
Redis Memory:
│
├─ rate_limit:192.168.1.1:window (ZSET)
│  │ └─ Members: [req1_ts, req2_ts, req3_ts, req4_ts, req5_ts]
│  │ └ Expire: 2 seconds
│  │
├─ rate_limit:10.0.0.5:window (ZSET)
│  │ └─ Members: [req1_ts, req2_ts, ...]
│  │ └─ Expire: 2 seconds
│  │
└─ rate_limit:203.0.113.42:window (ZSET)
   └─ Members: [req1_ts, req2_ts, ...]
   └─ Expire: 2 seconds
```

## API Endpoint Flow

```
CREATE /api/v1/client-ip-whitelist
    ↓
ClientIPWhitelistHandler.Create()
    ↓
Validate Request DTO
    ├─ client_id required
    ├─ ip_address required & valid
    ├─ is_active required
    └─ description optional
    ↓
ClientIPWhitelistService.Create(req)
    ↓
Create Model & Insert to DB
    ├─ Save to: client_ip_whitelists table
    ├─ Trigger: created_at timestamp
    ├─ Return: model with ID
    └─ Handle: foreign key constraint
    ↓
Convert to DTO
    ↓
Return 201 Created + Response


GET /api/v1/client-ip-whitelist/client/:client_id
    ↓
ClientIPWhitelistHandler.GetByClientID()
    ↓
Parse Params
    ├─ client_id (required)
    ├─ page (optional, default: 1)
    └─ limit (optional, default: 10)
    ↓
ClientIPWhitelistService.GetByClientID()
    ↓
Query Database
    ├─ WHERE client_id = ? AND deleted_at IS NULL
    ├─ COUNT total rows
    ├─ OFFSET & LIMIT for pagination
    └─ ORDER BY created_at DESC
    ↓
Convert to DTO Array
    ↓
Return 200 OK + Paginated Response
```

## Middleware Chain

```
Request Comes In
    │
    ▼
┌──────────────────────────┐
│ CORS Middleware          │ (Allow cross-origin)
└─────────────┬────────────┘
              │
              ▼
┌──────────────────────────────────┐
│ Rate Limiter Middleware          │ (Redis check)
│ (Global or per-group)            │
└──────────────┬───────────────────┘
               │
               ▼ (Optional)
┌──────────────────────────────────┐
│ IP Whitelist Middleware          │ (DB check)
│ (Only on specific routes)        │
└──────────────┬───────────────────┘
               │
               ▼
┌──────────────────────────────────┐
│ JWT Auth Middleware              │ (Token validate)
│ (On protected routes)            │
└──────────────┬───────────────────┘
               │
               ▼
┌──────────────────────────────────┐
│ Handler Function                 │ (Process request)
└──────────────┬───────────────────┘
               │
               ▼
         Response
```

## Redis Commands Executed

```
Request 1 (192.168.1.1):
  ├─ ZREMRANGEBYSCORE rate_limit:192.168.1.1:window -inf T0-1
  │   (Remove old entries outside 1-second window)
  │
  ├─ ZCOUNT rate_limit:192.168.1.1:window T0-1 +inf
  │   (Count requests in current window)
  │   → Result: 0 (first request)
  │
  ├─ ZADD rate_limit:192.168.1.1:window T0 "req-1"
  │   (Add current request)
  │
  └─ EXPIRE rate_limit:192.168.1.1:window 2
      (Auto-cleanup after 2 seconds)

Request 5 (192.168.1.1):
  ├─ ZREMRANGEBYSCORE rate_limit:192.168.1.1:window -inf T0-1
  │   (Remove old entries)
  │
  ├─ ZCOUNT rate_limit:192.168.1.1:window T0-1 +inf
  │   (Count requests)
  │   → Result: 4 (previous requests)
  │
  ├─ ZADD rate_limit:192.168.1.1:window T0 "req-5"
  │
  └─ EXPIRE rate_limit:192.168.1.1:window 2

Request 6 (192.168.1.1):
  ├─ ZCOUNT rate_limit:192.168.1.1:window T0-1 +inf
  │   → Result: 5 (at RPS limit)
  │
  └─ RETURN 429 Too Many Requests (no ZADD)
```

## Database Query Flow

```
IsIPWhitelisted(clientID: 1, ipAddress: "192.168.1.1")
    │
    ▼
SELECT COUNT(*) FROM client_ip_whitelists
WHERE client_id = 1
  AND ip_address = '192.168.1.1'
  AND is_active = true
  AND deleted_at IS NULL
    │
    ▼
    ├─ Index: (client_id, ip_address) → Fast lookup
    │
    └─ Return: count > 0 → true/false


GetByClientID(clientID: 1, page: 1, limit: 10)
    │
    ▼
-- Count total
SELECT COUNT(*) FROM client_ip_whitelists
WHERE client_id = 1
  AND deleted_at IS NULL
    │
    ├─ Index: client_id → Fast count
    │
    └─ Return: 5 (example)
    
-- Get paginated results
SELECT * FROM client_ip_whitelists
WHERE client_id = 1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 10 OFFSET 0
    │
    └─ Return: 5 rows
```

