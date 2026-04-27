package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	apperrors "github.com/shownest/pkg/errors"
)

/*
Limit: Maximum number of requests allowed within the window.
Window: Time duration for which the limit applies (e.g., 1 minute).
KeyFunc: A function that generates a unique Redis key for each request based on criteria like IP address or user ID.
*/
type RateLimitConfig struct {
	Limit   int64
	Window  time.Duration
	KeyFunc func(c *gin.Context) string
}

// ByIP generates a rate limit key based on the client's IP address.
func ByIP(c *gin.Context) string {
	return "rl:ip:" + c.ClientIP()
}

// ByUserID generates a rate limit key based on the authenticated user's ID.
func ByUserID(c *gin.Context) string {
	userID, exists := c.Get("userId")
	if !exists {
		return "rl:ip:" + c.ClientIP()
	}
	return fmt.Sprintf("rl:user:%v", userID)
}

// RateLimit limits requests using a fixed window counter in Redis. Returns 429 when the limit is exceeded.
func RateLimit(cache *redis.Client, cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := cfg.KeyFunc(c)
		ctx := context.Background()

		count, err := cache.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			cache.Expire(ctx, key, cfg.Window)
		}

		if count > cfg.Limit {
			ttl, err := cache.TTL(ctx, key).Result()
			if err != nil || ttl <= 0 {
				ttl = cfg.Window
			}
			c.Header("Retry-After", fmt.Sprintf("%.0f", ttl.Seconds()))
			abortWithError(c, apperrors.New(apperrors.CodeResourceExhausted, "too many requests; try again later"))
			return
		}

		c.Next()
	}
}
