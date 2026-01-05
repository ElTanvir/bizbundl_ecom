package store

import (
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type Store interface {
	Get(key string) (any, bool)
	Set(key string, value any, d time.Duration)
	Delete(key string)
	// Keys is not natively supported efficiently by go-cache without locking, but we can iterate items
	Keys() []string
	Clear()
	// Helper for default expiration
	SetDefault(key string, value any)
}

type goCacheStore struct {
	c *cache.Cache
}

func (s *goCacheStore) Get(key string) (any, bool) {
	return s.c.Get(key)
}

func (s *goCacheStore) Set(key string, value any, d time.Duration) {
	s.c.Set(key, value, d)
}

func (s *goCacheStore) SetDefault(key string, value any) {
	s.c.Set(key, value, cache.DefaultExpiration)
}

func (s *goCacheStore) Delete(key string) {
	s.c.Delete(key)
}

func (s *goCacheStore) Keys() []string {
	// This is expensive as it returns all items, use carefully
	items := s.c.Items()
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	return keys
}

func (s *goCacheStore) Clear() {
	s.c.Flush()
}

var (
	inst     Store
	initOnce sync.Once
	mu       sync.Mutex
)

func initDefault() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	inst = &goCacheStore{
		c: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func Get() Store {
	initOnce.Do(initDefault)
	return inst
}

// Init allows injecting a specific store (or re-configuring)
// Note: go-cache is robust enough we usually don't need to replace it unless testing.
func Init(s Store) error {
	if s == nil {
		return errors.New("nil store")
	}
	mu.Lock()
	defer mu.Unlock()
	if inst != nil {
		return errors.New("store already initialized")
	}
	inst = s
	initOnce.Do(func() {})
	return nil
}
