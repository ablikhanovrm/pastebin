package routes

import (
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func InitPasteRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	paste := rg.Group("/paste")
	{
		paste.GET("", h.GetPastes)
		paste.GET("/my", h.GetMyPastes)
		paste.GET("/:id", h.GetPaste)
		paste.POST("/", h.CreatePaste)
		paste.PUT("/:id", h.UpdatePaste)
		paste.DELETE("/:id", h.DeletePaste)
	}

}
