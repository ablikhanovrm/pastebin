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
	if l == nil {
		return true, nil
	}

	if l.rdb == nil {
		return true, nil
	}

	redisKey := "rate_limit:" + key

	count, err := l.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		_ = l.rdb.Expire(ctx, redisKey, l.ttl).Err()
	}

	if count > int64(l.limit) {
		return false, nil
	}

	return true, nil
}
