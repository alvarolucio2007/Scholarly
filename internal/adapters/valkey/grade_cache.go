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

var _ ports.Cache[domain.Grade] = (*GradeCache)(nil)

type GradeCache struct {
	rdb *redis.Client
}

func (c *GradeCache) Create(ctx context.Context, item *domain.Grade) error {
	cacheKey := fmt.Sprintf("grade:%d", item.ID)
	json, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create grade error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	return c.rdb.Set(ctx, cacheKey, json, CacheExpTime).Err()
}

func (c *GradeCache) Read(ctx context.Context, id int64) (*domain.Grade, error) {
	cacheKey := fmt.Sprintf("grade:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cache: read grade error: %w", err)
	}
	var grade domain.Grade
	if err := json.Unmarshal([]byte(data), &grade); err != nil {
		return nil, err
	}
	return &grade, nil
}

func (c *GradeCache) Delete(ctx context.Context, id int64) error {
	cacheKey := fmt.Sprintf("grade:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	if err := c.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete grade error: %w", err)
	}
	return nil
}
