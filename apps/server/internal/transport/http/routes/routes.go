package routes

import (
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, h *handler.Handler) {
	api := r.Group("/api")

	InitAuthRoutes(api, h)
	InitPasteRoutes(api, h)
}
