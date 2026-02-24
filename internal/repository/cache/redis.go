package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRedis(cfg config.RedisConfig, log zerolog.Logger) *redis.Client {
	rdbAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         rdbAddr,
		Password:     cfg.Password,
		DB:           0,
		PoolSize:     30,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Err(err).Msg("failed to connect to redis")
	}

	log.Info().Msg("redis connected")

	return rdb
}
