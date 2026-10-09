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

var _ ports.Cache[domain.Teacher] = (*TeacherCache)(nil)

type TeacherCache struct {
	rdb *redis.Client
}

func (c *TeacherCache) Create(ctx context.Context, item *domain.Teacher) error {
	cacheKey := fmt.Sprintf("teacher:%d", item.UserID)
	json, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache: create user error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	return c.rdb.Set(ctx, cacheKey, json, UserExpTime).Err()
}

func (c *TeacherCache) Read(ctx context.Context, id int64) (*domain.Teacher, error) {
	cacheKey := fmt.Sprintf("teacher:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cache: read teacher error: %w", err)
	}
	var teacher domain.Teacher
	if err := json.Unmarshal([]byte(data), &teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (c *TeacherCache) Delete(ctx context.Context, id int64) error {
	cacheKey := fmt.Sprintf("teacher:%d", id)
	ctx, cancel := context.WithTimeout(ctx, CacheQueryTimeout)
	defer cancel()
	if err := c.rdb.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("cache: delete teacher error: %w", err)
	}
	return nil
}
