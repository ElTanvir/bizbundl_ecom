package store

import (
	"errors"
	"sync"
)

type Store interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
	Keys() []string
	Clear()
}

type memStore struct {
	mu   sync.RWMutex
	data map[string]any
}

func (m *memStore) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *memStore) Set(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *memStore) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *memStore) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *memStore) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]any)
}

var (
	inst     Store
	initOnce sync.Once
	mu       sync.Mutex
)

func initDefault() {
	inst = &memStore{data: make(map[string]any)}
}

func Get() Store {
	initOnce.Do(initDefault)
	return inst
}
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
