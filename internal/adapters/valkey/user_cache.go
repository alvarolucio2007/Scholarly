package valkey

import (
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/redis/go-redis/v9"
)

var _ ports.UserCache = (*UserCache)(nil)

type UserCache struct {
	rdb *redis.Client
}

func (c *UserCache) Test() {
}
