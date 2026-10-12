// Package valkey handles cache connection and respective cache methods
package valkey

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CacheQueryTimeout = 500 * time.Millisecond

	CacheExpTime = time.Hour
)

type HealthChecker struct {
	rdb *redis.Client
}

func NewHealthChecker(rdb *redis.Client) *HealthChecker {
	return &HealthChecker{rdb: rdb}
}

func NewConn(address, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		Protocol: 2,
	})
	return rdb
}

func (v *HealthChecker) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return v.rdb.Ping(ctx).Err()
}
