package routes

import (
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoutes(h *handler.Handler) *gin.Engine {
	router := gin.New()
	router.Use(middleware.ClientInfoMiddleware())
	api := router.Group("/api")

	InitAuthRoutes(api, h)
	InitPasteRoutes(api, h)

	return router
}
