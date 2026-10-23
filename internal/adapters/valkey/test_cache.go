package valkey

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/redis/go-redis/v9"
)

var _ ports.Cache[domain.Test] = (*TestCache)(nil)

type TestCache struct {
	rdb *redis.Client
}

func NewTestCache(rdb *redis.Client) *TestCache {
	return &TestCache{rdb: rdb}
}

func (c *TestCache) Create(ctx context.Context, item domain.Test) error {
	if item.ID == 0 {
		return fmt.Errorf("cannot cache empty test")
	}
	cacheKey := fmt.Sprintf("test:%d", item.ID)
	json, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create test error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	return c.rdb.Set(ctx, cacheKey, json, CacheExpTime).Err()
}

func (c *TestCache) Read(ctx context.Context, id int64) (domain.Test, error) {
	cacheKey := fmt.Sprintf("test:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return domain.Test{}, nil
	} else if err != nil {
		return domain.Test{}, fmt.Errorf("cache: read test error: %w", err)
	}
	var test domain.Test
	if err := json.Unmarshal([]byte(data), &test); err != nil {
		return domain.Test{}, err
	}
	return test, nil
}

func (c *TestCache) Delete(ctx context.Context, id int64) error {
	cacheKey := fmt.Sprintf("test:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	if err := c.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete test error: %w", err)
	}
	return nil
}
