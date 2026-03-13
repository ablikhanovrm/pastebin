package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestLogger(base zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		log := base.With().
			Str("request_id", reqID).
			Logger()

		ctx := log.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()
	}
}
