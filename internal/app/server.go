package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/rs/zerolog/log"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) NewServer(cfg *config.Config, handler http.Handler) error {
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Info().
		Dur("read", cfg.Server.ReadTimeout).
		Dur("write", cfg.Server.WriteTimeout).
		Dur("idle", cfg.Server.IdleTimeout).
		Str("Host", cfg.Server.Host).
		Str("Port", cfg.Server.Port).
		Msg("CONFIG READY")

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
