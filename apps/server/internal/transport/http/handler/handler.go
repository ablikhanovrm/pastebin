package handler

import (
	"github.com/ablikhanovrm/pastebin/internal/service"
	"github.com/rs/zerolog"
)

type Handler struct {
	services *service.Services
	log      zerolog.Logger
}

func NewHandler(services *service.Services, log zerolog.Logger) *Handler {
	return &Handler{services, log}
}
