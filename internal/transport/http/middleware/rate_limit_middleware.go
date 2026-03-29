package middleware

import (
	"github.com/ablikhanovrm/pastebin/internal/ratelimit"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiter *ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)

		allowed, err := limiter.Allow(c.Request.Context(), ip)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal error"})
			return
		}

		if !allowed {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
