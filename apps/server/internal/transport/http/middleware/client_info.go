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
