package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisCache struct {
	rdb *redis.Client
	log zerolog.Logger
}

func NewRedisCache(db *redis.Client, log zerolog.Logger) *RedisCache {
	return &RedisCache{
		rdb: db,
		log: log,
	}
}
