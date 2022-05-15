package cache

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

type Cache struct {
	mu sync.RWMutex
	db map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		mu: sync.RWMutex{},
		db: map[string][]byte{},
	}
}

func NewCacheByDb(db map[string][]byte) *Cache {
	return &Cache{
		mu: sync.RWMutex{},
		db: db,
	}
}

func (cache *Cache) Get(key string) ([]byte, error) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.db[key], nil
}

func (cache *Cache) Set(key string, value []byte) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.db[key] = value
	return nil
}

func (cache *Cache) Delete(key string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	delete(cache.db, key)

	return nil
}

func (cache *Cache) Dump() ([]byte, error) {
	byteCache, err := json.Marshal(cache.db)
	if err != nil {
		zap.S().Errorf("can not dump cache.db", err.Error())
		return nil, err
	}

	return byteCache, nil
}
