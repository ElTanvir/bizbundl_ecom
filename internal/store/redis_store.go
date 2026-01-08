package store

import (
	"bizbundl/internal/infra/redis"
	"context"
	"encoding/json"
	"time"

	redisLib "github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redisLib.Client
}

func NewRedisStore(client *redisLib.Client) Store {
	return &RedisStore{client: client}
}

func (s *RedisStore) Get(ctx context.Context, key string) (any, bool) {
	prefixedKey := redis.Key(ctx, key)
	val, err := s.client.Get(ctx, prefixedKey).Bytes()
	if err == redisLib.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false // Treat error as miss?
	}

	var res any
	if err := json.Unmarshal(val, &res); err != nil {
		// If fails, maybe it's just a string?
		return string(val), true
	}
	return res, true
}

func (s *RedisStore) Set(ctx context.Context, key string, value any, d time.Duration) {
	prefixedKey := redis.Key(ctx, key)
	// Serialize
	bytes, err := json.Marshal(value)
	if err != nil {
		return // Log error?
	}
	s.client.Set(ctx, prefixedKey, bytes, d)
}

func (s *RedisStore) SetDefault(ctx context.Context, key string, value any) {
	s.Set(ctx, key, value, 5*time.Minute)
}

func (s *RedisStore) Delete(ctx context.Context, key string) {
	prefixedKey := redis.Key(ctx, key)
	s.client.Del(ctx, prefixedKey)
}

func (s *RedisStore) Keys(ctx context.Context) []string {
	// Keys should probably be prefixed too for search?
	// But KEYS * pattern is dangerous. use "shop_123:*"
	// This is where "Key(ctx, "")" might be useful if it returns prefix

	// Better: Key(ctx, "*") -> "shop_123:*"
	pattern := redis.Key(ctx, "*")
	return s.client.Keys(ctx, pattern).Val()
}

func (s *RedisStore) Clear(ctx context.Context) {
	// Only clear keys for this tenant
	pattern := redis.Key(ctx, "*")
	// Lua script or iterative delete?
	// For MVP, iterative.
	keys := s.client.Keys(ctx, pattern).Val()
	if len(keys) > 0 {
		s.client.Del(ctx, keys...)
	}
}
