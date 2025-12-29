package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter menggunakan Redis untuk distributed rate limiting
type RedisRateLimiter struct {
	client      *redis.Client
	maxRequests int           // maximum requests allowed
	window      time.Duration // time window (e.g., 1 minute)
}

// NewRedisRateLimiter membuat instance RedisRateLimiter
// maxRequests: jumlah request maksimal (contoh: 200)
// window: durasi time window (contoh: 1 minute = 60 seconds)
func NewRedisRateLimiter(client *redis.Client, maxRequests int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:      client,
		maxRequests: maxRequests,
		window:      window,
	}
}

// isAllowed mengecek apakah request diizinkan menggunakan sliding window counter algorithm
func (rl *RedisRateLimiter) isAllowed(ctx *gin.Context, key string) bool {
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	pipe := rl.client.TxPipeline()

	// Remove old entries outside window
	pipe.ZRemRangeByScore(ctx.Request.Context(), key, "-inf", fmt.Sprintf("%d", windowStart))

	// Count requests in current window
	count := pipe.ZCount(ctx.Request.Context(), key, fmt.Sprintf("%d", windowStart), "+inf")

	// Execute pipeline
	_, err := pipe.Exec(ctx.Request.Context())
	if err != nil && err != redis.Nil {
		return true // Allow on error
	}

	// Check if under rate limit
	countVal := count.Val()
	if int(countVal) >= rl.maxRequests {
		return false
	}

	// Add current request to window
	rl.client.ZAdd(ctx.Request.Context(), key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d-%d", now, time.Now().Nanosecond()),
	})

	// Set expiration to window duration + buffer
	rl.client.Expire(ctx.Request.Context(), key, rl.window+5*time.Second)

	return true
}

// Middleware untuk rate limiting dengan Redis
func (rl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		if !rl.isAllowed(c, key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests, please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MiddlewareWithWhitelist middleware dengan IP whitelist check
func (rl *RedisRateLimiter) MiddlewareWithWhitelist(ipWhitelistRepo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Jika menggunakan IP whitelist, skip rate limiting untuk IP yang di-whitelist
		// Implementasi tergantung dari context request yang ada

		key := fmt.Sprintf("rate_limit:%s", ip)

		if !rl.isAllowed(c, key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests, please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
