package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/repository"
	"github.com/ablikhanovrm/pastebin/internal/service"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/routes"
	"github.com/rs/zerolog/log"
)

func Run(cofigPath string) {
	config := config.GetConfig(cofigPath)

	storage, err := repository.NewPostgresStorage(&config.DB)

	if err != nil {
		log.Error().Err(err).Msg("failed to connect database")
	}

	repo := repository.NewRepository(storage.Queries)
	services := service.NewService(repo)
	handler := handler.NewHandler(services)
	router := routes.InitRoutes(handler)
	srv := new(Server)

	go func() {

		if err := srv.NewServer(config, router); err != nil {
			log.Fatal().Err(err).Msg("Failed to run server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Warn().Msg("TodoApp Shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error occured on server shutting down")
	}

	storage.Pool.Close()

}
