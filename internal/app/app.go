package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/logging"
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

	logger := logging.New("pastebin")

	db, err := repository.NewPostgresStorage(&newConfig.DB)
	if err != nil {

		logger.Fatal().Err(err).Msg("failed to connect database")
	}

	s3Client, _ := storage.NewS3Client(newConfig.S3)
	s3Storage := storage.NewS3Storage(s3Client, newConfig.S3.Bucket, logger.With().Str("layer", "storage").Logger())

	repo := repository.NewRepository(db.Pool, logger.With().Str("layer", "repository").Logger())
	jwtManager := jwt.New(newConfig.Server.JwtSecret)

	redisClient := cacheRepo.NewRedis(newConfig.Redis, logger.With().Str("layer", "redis_client").Logger())
	cache := cacheRepo.NewRedisCache(redisClient, logger.With().Str("layer", "cache").Logger())

	services := service.NewServices(repo, jwtManager, db.Pool, s3Storage, cache, logger.With().Str("layer", "service").Logger())

	handlerLogger := logger.With().Str("layer", "handler").Logger()
	newHandler := handler.NewHandler(services, &newConfig.Server, handlerLogger)

	router := routes.InitRoutes(newHandler, jwtManager)

	srv := new(Server)

	go func() {
		if err := srv.NewServer(newConfig, router); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Warn().Msg("Pastebin app Shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error().Err(err).Msg("Error occurred on server shutting down")
	}

	db.Pool.Close()
}
