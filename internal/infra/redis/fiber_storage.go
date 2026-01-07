package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// FiberStorage implements fiber.Storage interface
type FiberStorage struct {
	client *redis.Client
}

func NewFiberStorage(client *redis.Client) *FiberStorage {
	return &FiberStorage{client: client}
}

func (s *FiberStorage) Get(key string) ([]byte, error) {
	val, err := s.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil // Fiber expects nil, nil for miss
	}
	return val, err
}

func (s *FiberStorage) Set(key string, val []byte, exp time.Duration) error {
	return s.client.Set(context.Background(), key, val, exp).Err()
}

func (s *FiberStorage) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}

func (s *FiberStorage) Reset() error {
	return s.client.FlushDB(context.Background()).Err()
}

func (s *FiberStorage) Close() error {
	return s.client.Close()
}
