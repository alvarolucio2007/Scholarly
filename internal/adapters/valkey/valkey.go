// Package valkey handles cache connection and respective cache methods
package valkey

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CacheQueryTimeout = 500 * time.Millisecond

	CacheExpTime = time.Hour
)

func NewConn(address, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		Protocol: 2,
	})
	return rdb
}
