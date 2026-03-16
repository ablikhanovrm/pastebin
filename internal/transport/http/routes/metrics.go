package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetricRoutes(rg *gin.RouterGroup) {
	rg.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
