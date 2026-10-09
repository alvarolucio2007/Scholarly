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

var _ ports.Cache[domain.Student] = (*StudentCache)(nil)

type StudentCache struct {
	rdb *redis.Client
}

func (c *StudentCache) Create(ctx context.Context, item *domain.Student) error {
	cacheKey := fmt.Sprintf("student:%d", item.UserID)
	json, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create user error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	return c.rdb.Set(ctx, cacheKey, json, CacheExpTime).Err()
}

func (c *StudentCache) Read(ctx context.Context, id int64) (*domain.Student, error) {
	cacheKey := fmt.Sprintf("student:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cache: read student error: %w", err)
	}
	var student domain.Student
	if err := json.Unmarshal([]byte(data), &student); err != nil {
		return nil, err
	}
	return &student, nil
}

func (c *StudentCache) Delete(ctx context.Context, id int64) error {
	cacheKey := fmt.Sprintf("student:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	if err := c.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete student error: %w", err)
	}
	return nil
}
