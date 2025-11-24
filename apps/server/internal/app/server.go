package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ablikhanovrm/pastebin/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) NewServer(cfg *config.Config, handler http.Handler) error {
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
