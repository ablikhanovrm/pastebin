package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/logging"
	"github.com/ablikhanovrm/pastebin/internal/repository"
	"github.com/ablikhanovrm/pastebin/internal/service"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/routes"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/rs/zerolog/log"
)

func Run(configPath string) {
	newConfig := config.GetConfig(configPath)

	logger := logging.New("pastebin")

	storage, err := repository.NewPostgresStorage(&newConfig.DB)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect database")
	}

	repo := repository.NewRepository(storage.Pool, logger.With().Str("layer", "repository").Logger())
	manager := jwt.New(newConfig.Server.JwtSecret)
	services := service.NewServices(repo, manager, storage.Pool, logger.With().Str("layer", "service").Logger())
	handlerLogger := logger.With().Str("layer", "handler").Logger()
	newHandler := handler.NewHandler(services, handlerLogger)
	router := routes.InitRoutes(newHandler)

	srv := new(Server)

	go func() {
		if err := srv.NewServer(newConfig, router); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Warn().Msg("TodoApp Shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error occurred on server shutting down")
	}

	storage.Pool.Close()
}
