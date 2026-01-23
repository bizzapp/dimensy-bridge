package middleware

import (
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// isIPEqual mengecek apakah dua IP address sama
// Khusus mendukung localhost (IPv4: 127.0.0.1 dan IPv6: ::1 dianggap sama)
func isIPEqual(ip1, ip2 string) bool {
	// Exact match
	if ip1 == ip2 {
		return true
	}

	// Parse sebagai net.IP untuk normalization
	parsedIP1 := net.ParseIP(ip1)
	parsedIP2 := net.ParseIP(ip2)

	// Jika keduanya valid IP
	if parsedIP1 != nil && parsedIP2 != nil {
		// Check jika keduanya localhost (IPv4 atau IPv6)
		if parsedIP1.IsLoopback() && parsedIP2.IsLoopback() {
			return true
		}

		// Check direct equality setelah parsing
		return parsedIP1.Equal(parsedIP2)
	}

	return false
}

// IPWhitelistMiddleware mengecek apakah client IP ada dalam whitelist
func IPWhitelistMiddleware(ipWhitelistRepo repository.ClientIPWhitelistRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil client_id dari request (bisa dari URL param, query, atau token)
		clientIDStr := c.Query("client_id")
		if clientIDStr == "" {
			clientIDStr = c.Param("client_id")
		}

		if clientIDStr == "" {
			ip := c.ClientIP()
			log.Printf("[IP_WHITELIST_BLOCKED] client_id is required. IP: %s, Method: %s, Path: %s", ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "client_id is required",
			})
			c.Abort()
			return
		}

		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			ip := c.ClientIP()
			log.Printf("[IP_WHITELIST_BLOCKED] Invalid client_id format: %s. IP: %s, Method: %s, Path: %s", clientIDStr, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid client_id",
			})
			c.Abort()
			return
		}

		ip := c.ClientIP()

		// Check apakah IP di-whitelist
		isWhitelisted, err := ipWhitelistRepo.IsIPWhitelisted(clientID, ip)
		if err != nil {
			log.Printf("[IP_WHITELIST_ERROR] Failed to verify IP whitelist. ClientID: %d, IP: %s, Error: %v", clientID, ip, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify IP whitelist",
			})
			c.Abort()
			return
		}

		if !isWhitelisted {
			log.Printf("[IP_WHITELIST_BLOCKED] IP not whitelisted. ClientID: %d, IP: %s, Method: %s, Path: %s", clientID, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Your IP address is not whitelisted",
				"ip":      ip,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalIPWhitelistMiddleware mengecek whitelist jika ada, tapi tidak memblokir jika tidak ada whitelist entry
func OptionalIPWhitelistMiddleware(ipWhitelistRepo repository.ClientIPWhitelistRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIDStr := c.Query("client_id")
		if clientIDStr == "" {
			clientIDStr = c.Param("client_id")
		}

		// Skip jika tidak ada client_id
		if clientIDStr == "" {
			c.Next()
			return
		}

		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			c.Next()
			return
		}

		// Ambil semua active IPs untuk client
		activeIPs, err := ipWhitelistRepo.GetActiveIPsByClientID(clientID)
		if err != nil {
			log.Printf("[IP_WHITELIST_ERROR] Failed to get active IPs for client. ClientID: %d, Error: %v", clientID, err)
			// Jika error, skip validation
			c.Next()
			return
		}

		// Jika tidak ada IP whitelist, skip validation
		if len(activeIPs) == 0 {
			c.Next()
			return
		}

		// Check apakah current IP ada dalam whitelist
		ip := c.ClientIP()
		isWhitelisted := false

		for _, whitelistedIP := range activeIPs {
			if isIPEqual(whitelistedIP, ip) {
				isWhitelisted = true
				break
			}
		}

		if !isWhitelisted {
			log.Printf("[IP_WHITELIST_BLOCKED] IP not whitelisted (optional mode). ClientID: %d, IP: %s, Method: %s, Path: %s", clientID, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Your IP address is not whitelisted",
				"ip":      ip,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWEIPWhitelistMiddleware mengecek IP berdasarkan client_id dari JWE token
// STRICT MODE: Tidak boleh akses jika IP tidak registered atau tidak ada data

// JWEIPWhitelistWithClientPsreMiddleware mengecek IP berdasarkan JWE token + client_psre lookup
// Flow:
// 1. Extract JWE token dari Authorization header
// 2. Decrypt dan extract client_psre.external_id dari payload data.id
// 3. Query client_psre table untuk dapatkan client_id
// 4. Validasi IP terhadap client_ip_whitelists
// 5. STRICT MODE: Block jika IP tidak registered
func JWEIPWhitelistWithClientPsreMiddleware(
	ipWhitelistRepo repository.ClientIPWhitelistRepository,
	clientPsreRepo repository.ClientPsreRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract JWE token dari Authorization header
		authHeader := c.GetHeader("Authorization")

		// OPTIONAL MODE: Jika tidak ada Authorization header (untuk login endpoints), skip validation
		// Login endpoints akan melewati middleware ini tanpa validation
		if authHeader == "" {
			c.Next()
			return
		}

		// Get token (format: "Bearer <token>")
		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}

		// 2. Decrypt JWE token
		payload, err := utils.VerifyJWE(token)
		if err != nil {
			ip := c.ClientIP()
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] Invalid or expired token. IP: %s, Method: %s, Path: %s, Error: %v", ip, c.Request.Method, c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 3. Extract data.id dari payload (client_psre external_id)
		dataInterface, exists := payload["data"]
		if !exists {
			ip := c.ClientIP()
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] Data field not found in token. IP: %s, Method: %s, Path: %s", ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Data field not found in token",
				"payload": payload,
			})
			c.Abort()
			return
		}

		dataMap, ok := dataInterface.(map[string]interface{})
		if !ok {
			// Coba handle jika data adalah map[interface{}]interface{}
			if dataMapAlt, okAlt := dataInterface.(map[interface{}]interface{}); okAlt {
				dataMap = make(map[string]interface{})
				for k, v := range dataMapAlt {
					if kStr, ok := k.(string); ok {
						dataMap[kStr] = v
					}
				}
			} else {
				ip := c.ClientIP()
				log.Printf("[JWE_IP_WHITELIST_BLOCKED] Invalid data format in token (type: %T). IP: %s, Method: %s, Path: %s", dataInterface, ip, c.Request.Method, c.Request.URL.Path)
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Invalid data format in token",
					"type":    fmt.Sprintf("%T", dataInterface),
				})
				c.Abort()
				return
			}
		}

		externalID, exists := dataMap["id"]
		if !exists {
			ip := c.ClientIP()
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] Client ID not found in token. IP: %s, Method: %s, Path: %s", ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Client ID not found in token",
				"dataKeys": func() []string {
					var keys []string
					for k := range dataMap {
						keys = append(keys, k)
					}
					return keys
				}(),
			})
			c.Abort()
			return
		}

		externalIDStr, ok := externalID.(string)
		if !ok {
			ip := c.ClientIP()
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] Invalid client ID format in token (type: %T). IP: %s, Method: %s, Path: %s", externalID, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid client ID format in token",
			})
			c.Abort()
			return
		}

		// 4. Query client_psre untuk dapatkan client_id
		client, err := clientPsreRepo.FindByExternalID(externalIDStr)
		ip := c.ClientIP()
		if err != nil {
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] Client not found or invalid token. ExternalID: %s, IP: %s, Method: %s, Path: %s, Error: %v", externalIDStr, ip, c.Request.Method, c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Client not found or invalid token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if client == nil {
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] Client information not found. ExternalID: %s, IP: %s, Method: %s, Path: %s", externalIDStr, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Client information not found",
			})
			c.Abort()
			return
		}

		clientID := client.ID

		// 6. Check apakah ada IP whitelist untuk client ini
		activeIPs, err := ipWhitelistRepo.GetActiveIPsByClientID(clientID)
		if err != nil {
			log.Printf("[JWE_IP_WHITELIST_ERROR] Failed to verify IP whitelist. ClientID: %d, ExternalID: %s, IP: %s, Error: %v", clientID, externalIDStr, ip, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify IP whitelist",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 7. STRICT MODE: Jika tidak ada IP yang registered, block akses
		if len(activeIPs) == 0 {
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] No IP addresses registered for this account. ClientID: %d, ExternalID: %s, IP: %s, Method: %s, Path: %s", clientID, externalIDStr, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "No IP addresses registered for this account. Please register your IP first.",
				"ip":      ip,
			})
			c.Abort()
			return
		}

		// 8. Check apakah current IP ada dalam whitelist
		isWhitelisted := false
		for _, whitelistedIP := range activeIPs {
			if isIPEqual(whitelistedIP, ip) {
				isWhitelisted = true
				break
			}
		}

		// 9. Jika tidak whitelisted, block akses
		if !isWhitelisted {
			log.Printf("[JWE_IP_WHITELIST_BLOCKED] IP not whitelisted. ClientID: %d, ExternalID: %s, IP: %s, Method: %s, Path: %s", clientID, externalIDStr, ip, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Your IP address is not whitelisted",
				"ip":      ip,
			})
			c.Abort()
			return
		}

		// 10. Set ke context
		c.Set("client_id", clientID)
		c.Set("client_ip", ip)
		c.Set("external_id", externalIDStr)
		c.Set("client_name", dataMap["name"])

		c.Next()
	}
}
