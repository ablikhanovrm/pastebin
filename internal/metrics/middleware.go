package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()

		// если роут не сматчился — сразу выходим
		if path == "" {
			return
		}

		// исключаем служебные ручки
		if path == "/api/metrics" {
			return
		}

		// опционально игнорим 404
		if status == "404" {
			return
		}

		duration := time.Since(start).Seconds()

		HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}
