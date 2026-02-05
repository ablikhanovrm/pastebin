package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	clientIPKey  = "client_ip"
	userAgentKey = "user_agent"
)

func ClientInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ua := c.Request.UserAgent()

		c.Set(clientIPKey, ip)
		c.Set(userAgentKey, ua)

		c.Next()
	}
}

func GetClientIP(c *gin.Context) string {
	if v, ok := c.Get(clientIPKey); ok {
		if ip, ok := v.(string); ok {
			return ip
		}
	}
	return ""
}

func GetUserAgent(c *gin.Context) string {
	if v, ok := c.Get(userAgentKey); ok {
		if ua, ok := v.(string); ok {
			return ua
		}
	}
	return ""
}
