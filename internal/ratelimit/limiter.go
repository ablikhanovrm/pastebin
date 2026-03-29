package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb   *redis.Client
	limit int
	ttl   time.Duration
}

func NewLimiter(rdb *redis.Client, limit int, ttl time.Duration) *Limiter {
	return &Limiter{
		rdb:   rdb,
		limit: limit,
		ttl:   ttl,
	}
}

func (l *Limiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := "rate_limit:" + key

	count, err := l.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		l.rdb.Expire(ctx, redisKey, l.ttl)
	}

	if count > int64(l.limit) {
		return false, nil
	}

	return true, nil
}
