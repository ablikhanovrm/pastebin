package routes

import (
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func InitRoutes(h *handler.Handler) *gin.Engine {
	r := gin.New()
	api := r.Group("/api")

	InitAuthRoutes(api, h)
	InitPasteRoutes(api, h)

	return r
}
