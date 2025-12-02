package handler

import (
	"github.com/ablikhanovrm/pastebin/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services}
}
