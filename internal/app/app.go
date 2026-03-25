package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/logging"
	"github.com/ablikhanovrm/pastebin/internal/metrics"
	"github.com/ablikhanovrm/pastebin/internal/repository"
	cacheRepo "github.com/ablikhanovrm/pastebin/internal/repository/cache"
	"github.com/ablikhanovrm/pastebin/internal/service"
	"github.com/ablikhanovrm/pastebin/internal/service/storage"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/handler"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/routes"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
)

func Run() {
	newConfig := config.GetConfig()
	metrics.MustRegister()

	logger := logging.New("pastebin")

	db, err := repository.NewPostgresStorage(&newConfig.DB)
	if err != nil {

		logger.Fatal().Err(err).Msg("failed to connect database")
	}

	s3Client, _ := storage.NewS3Client(newConfig.S3)
	s3Storage := storage.NewS3Storage(s3Client, newConfig.S3.Bucket)

	repo := repository.NewRepository(db.Pool)
	jwtManager := jwt.New(newConfig.Server.JwtSecret)

	redisClient := cacheRepo.NewRedis(newConfig.Redis, logger)
	cache := cacheRepo.NewRedisCache(redisClient, logger)

	services := service.NewServices(repo, jwtManager, db.Pool, s3Storage, cache)

	newHandler := handler.NewHandler(services, &newConfig.Server)

	router := routes.InitRoutes(newHandler, jwtManager, logger)

	srv := new(Server)

	go func() {
		if err := srv.NewServer(newConfig, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Failed to run server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()

	logger.Warn().Msg("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Error during shutdown")
	}

	db.Pool.Close()
}
