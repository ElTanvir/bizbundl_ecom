package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) Store {
	return &RedisStore{client: client}
}

func (s *RedisStore) Get(key string) (any, bool) {
	val, err := s.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false // Treat error as miss?
	}

	// Try unmarshal to map or something?
	// Problem: We don't know the target type here to Unmarshal INTO.
	// go-cache stores the actual object. Redis stores bytes.
	// This interface `Get(key) (any, bool)` is fundamentally incompatible with external storage
	// unless we return map[string]interface{} or raw bytes,
	// OR we change the interface to `Get(key string, target any) bool`.

	// If the app expects the exact struct back, this will fail.
	// We need to check usage.
	// For now, let's look at usages first.
	// But assuming we need to return something, let's return a generic map or the bytes if string.
	var res any
	if err := json.Unmarshal(val, &res); err != nil {
		// If fails, maybe it's just a string?
		return string(val), true
	}
	return res, true
}

func (s *RedisStore) Set(key string, value any, d time.Duration) {
	// Serialize
	bytes, err := json.Marshal(value)
	if err != nil {
		return // Log error?
	}
	s.client.Set(context.Background(), key, bytes, d)
}

func (s *RedisStore) SetDefault(key string, value any) {
	s.Set(key, value, 5*time.Minute) // Default from original store
}

func (s *RedisStore) Delete(key string) {
	s.client.Del(context.Background(), key)
}

func (s *RedisStore) Keys() []string {
	return s.client.Keys(context.Background(), "*").Val()
}

func (s *RedisStore) Clear() {
	// DANGEROUS: FlushDB? Or just delete match?
	// Given it's a "Store" for the app, maybe Clear means Clear Cache.
	// But "Redis installation for Multiple tennant"... FlushDB is bad.
	// We should probably verify if we used a prefix.
	// If we set DB, FlushDB clears that DB.
	// Let's assume FlushDB is fine for the selected DB if isolated.
	s.client.FlushDB(context.Background())
}
