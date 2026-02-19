package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/redis/go-redis/v9"
)

func (c *RedisCache) SetPaste(ctx context.Context, paste *paste.Paste) error {
	key := pasteKey(paste.Uuid.String())

	data, err := json.Marshal(&paste)

	if err != nil {
		return err
	}

	err = c.rdb.SetEx(ctx, key, data, time.Minute*5).Err()

	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) GetPaste(ctx context.Context, id string) (*paste.Paste, error) {
	key := pasteKey(id)

	data, err := c.rdb.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	cachePaste := &paste.Paste{}

	err = json.Unmarshal(data, cachePaste)
	if err != nil {
		return nil, err
	}
	return cachePaste, nil
}

func (c *RedisCache) DeletePaste(ctx context.Context, id string) error {
	key := pasteKey(id)
	err := c.rdb.Del(ctx, key).Err()

	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) SetPasteList(ctx context.Context, limit int32, cursor *time.Time, ids []string) error {
	key := pasteListKey(limit, cursor)

	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, time.Minute*5).Err()
}

func (c *RedisCache) GetPasteList(ctx context.Context, limit int32, cursor *time.Time) ([]string, error) {
	key := pasteListKey(limit, cursor)

	val, err := c.rdb.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		c.log.Error().Err(err).Msg("redis get paste list failed")
		return nil, err
	}

	ids := make([]string, 0, limit)
	if err := json.Unmarshal(val, &ids); err != nil {
		c.log.Error().Err(err).Msg("redis get paste list failed")
		return nil, err
	}

	return ids, nil
}

func (c *RedisCache) InvalidatePasteLists(ctx context.Context) error {
	iter := c.rdb.Scan(ctx, 0, "paste_list:*", 0).Iterator()

	for iter.Next(ctx) {
		if err := c.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

func (c *RedisCache) SetPasteContent(ctx context.Context, id string, data []byte) error {
	key := pasteContentKey(id)

	err := c.rdb.Set(ctx, key, data, time.Minute*5).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) GetPasteContent(ctx context.Context, id string) ([]byte, error) {
	key := pasteContentKey(id)

	val, err := c.rdb.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		c.log.Error().Err(err).Msg("redis get paste content failed")
		return nil, err
	}

	return val, nil
}
