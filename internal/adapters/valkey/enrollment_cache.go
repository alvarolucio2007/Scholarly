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

var _ ports.Cache[domain.Enrollment] = (*EnrollmentCache)(nil)

type EnrollmentCache struct {
	rdb *redis.Client
}

func (c *EnrollmentCache) Create(ctx context.Context, item *domain.Enrollment) error {
	cacheKey := fmt.Sprintf("enrollment:%d", item.ID)
	json, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create enrollment error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	return c.rdb.Set(ctx, cacheKey, json, CacheExpTime).Err()
}

func (c *EnrollmentCache) Read(ctx context.Context, id int64) (*domain.Enrollment, error) {
	cacheKey := fmt.Sprintf("enrollment:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cache: read enrollment error: %w", err)
	}
	var enrollment domain.Enrollment
	if err := json.Unmarshal([]byte(data), &enrollment); err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (c *EnrollmentCache) Delete(ctx context.Context, id int64) error {
	cacheKey := fmt.Sprintf("enrollment:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	if err := c.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete enrollment error: %w", err)
	}
	return nil
}
