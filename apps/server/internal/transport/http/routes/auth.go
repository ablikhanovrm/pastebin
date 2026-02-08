package routes

import (
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}
}
