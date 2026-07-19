package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type sDistributedCache struct {
	client *redis.Client
	locker *redislock.Client
}

func NewRedisCacheImpl(client *redis.Client) *sDistributedCache {
	return &sDistributedCache{
		client: client,
		locker: redislock.New(client),
	}
}

func (s *sDistributedCache) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	fmt.Println("val:", val)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return val, nil
		}
		return val, fmt.Errorf("redis get error: %w", err)
	}

	return val, nil
}

func (s *sDistributedCache) Set(ctx context.Context, key string, value interface{}, expirationSeconds int) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}
	if err := s.client.Set(ctx, key, b, time.Duration(expirationSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

func (s *sDistributedCache) Del(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del error: %w", err)
	}

	return nil
}

func (s *sDistributedCache) Incr(ctx context.Context, key string) (int64, error) {
	val, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis incr error: %w", err)
	}

	return val, nil
}

func (s *sDistributedCache) Decr(ctx context.Context, key string) (int64, error) {
	val, err := s.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis decr error: %w", err)
	}

	return val, nil
}

func (s *sDistributedCache) Exists(ctx context.Context, key string) (bool, error) {
	val, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}

	return val == 1, nil
}

func (s *sDistributedCache) WithDistributedLock(ctx context.Context, key string, ttlSeconds int, fn func(ctx context.Context) error) error {
	lockTTL := time.Duration(ttlSeconds) * time.Second
	lock, err := s.locker.Obtain(ctx, key, lockTTL, nil)
	if err != redislock.ErrNotObtained {
		return fmt.Errorf("could not obtain lock for key: %s", key)
	} else if err != nil {
		return fmt.Errorf("failed to obtain lock: %w", err)
	}
	defer lock.Release(ctx)

	return fn(ctx)
}
