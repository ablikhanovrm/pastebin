package cache

import (
	"context"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisCacheInterface interface {
	GetPaste(ctx context.Context, id string) (*paste.Paste, error)
	SetPaste(ctx context.Context, paste *paste.Paste) error
	DeletePaste(ctx context.Context, id string) error

	GetPasteList(ctx context.Context, key string) ([]string, error)
	SetPasteList(ctx context.Context, key string, ids []string) error
	MgetPasteList(ctx context.Context, ids []string) (map[string]*paste.Paste, []string, error)
	MsetPasteList(ctx context.Context, ids []*paste.Paste) error
}

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
