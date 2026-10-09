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

var _ ports.Cache[domain.Course] = (*CourseCache)(nil)

type CourseCache struct {
	rdb *redis.Client
}

func (c *CourseCache) Create(ctx context.Context, item *domain.Course) error {
	cacheKey := fmt.Sprintf("course:%d", item.ID)
	json, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create course error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	return c.rdb.Set(ctx, cacheKey, json, CourseExpTime).Err()
}

func (c *CourseCache) Read(ctx context.Context, id int64) (*domain.Course, error) {
	cacheKey := fmt.Sprintf("course:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cache: read course error: %w", err)
	}
	var course domain.Course
	if err := json.Unmarshal([]byte(data), &course); err != nil {
		return nil, err
	}
	return &course, nil
}

func (c *CourseCache) Delete(ctx context.Context, id int64) error {
	cacheKey := fmt.Sprintf("course:%d")
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	if err := c.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete course error: %w", err)
	}
	return nil
}
