package handler

import (
	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/service"
)

type Handler struct {
	services *service.Services
	cfg      *config.HttpServer
}

func NewHandler(services *service.Services, cfg *config.HttpServer) *Handler {
	return &Handler{services, cfg}
}
