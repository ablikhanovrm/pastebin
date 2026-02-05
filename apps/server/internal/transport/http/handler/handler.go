package handler

import (
	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/service"
	"github.com/rs/zerolog"
)

type Handler struct {
	services *service.Services
	cfg      *config.HttpServer
	log      zerolog.Logger
}

func NewHandler(services *service.Services, cfg *config.HttpServer, log zerolog.Logger) *Handler {
	return &Handler{services, cfg, log}
}
