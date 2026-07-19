package impl

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

type sRistrettoCacheCacheImpl struct {
	cache *ristretto.Cache
}

func NewRistrettoCacheImpl() (*sRistrettoCacheCacheImpl, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, errors.New("failed to create ristretto cache")
	}

	return &sRistrettoCacheCacheImpl{cache}, nil
}

func (s *sRistrettoCacheCacheImpl) Get(ctx context.Context, key string) (interface{}, bool) {
	return s.cache.Get(key)
}

func (s *sRistrettoCacheCacheImpl) Set(ctx context.Context, key string, value interface{}) bool {
	return s.cache.Set(key, value, 1) // Cost mặc định = 1
}

func (s *sRistrettoCacheCacheImpl) SetWithTTL(ctx context.Context, key string, value interface{}) bool {
	dataJson, _ := json.Marshal(value)
	return s.cache.SetWithTTL(key, string(dataJson), 1, 5*time.Minute) // Cost mặc dinh = 1
}

func (s *sRistrettoCacheCacheImpl) Del(ctx context.Context, key string) {
	s.cache.Del(key)
}
