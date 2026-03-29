package routes

import (
	gometrics "github.com/ablikhanovrm/pastebin/internal/metrics"
	"github.com/ablikhanovrm/pastebin/internal/ratelimit"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func InitRoutes(h *handler.Handler, jwtManager *jwt.Manager, limiter *ratelimit.Limiter, log zerolog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ClientInfoMiddleware())
	router.Use(middleware.RequestLogger(log))
	router.Use(gometrics.GinMiddleware())
	router.Use(middleware.RateLimitMiddleware(limiter))

	api := router.Group("/api")

	// публичные
	InitAuthRoutes(api, h)
	InitMetricRoutes(api)

	// защищённые
	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtManager))

	InitPasteRoutes(auth, h)

	return router
}
