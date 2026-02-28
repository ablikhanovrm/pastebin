package routes

import (
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func InitRoutes(h *handler.Handler, jwtManager *jwt.Manager) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ClientInfoMiddleware())
	api := router.Group("/api")

	// публичные
	InitAuthRoutes(api, h)

	// защищённые
	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtManager))

	InitPasteRoutes(auth, h)

	return router
}
