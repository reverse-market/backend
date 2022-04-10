package redis

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
)

type Loggers interface {
	Info() *log.Logger
	Error() *log.Logger
}

func NewRedisPool(c *Config, loggers Loggers) (*redis.Pool, func()) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", c.RedisUrl)
		},
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", c.RedisUrl)
		},
	}
	return pool, func() {
		if err := pool.Close(); err != nil {
			loggers.Error().Printf("Failed closing redis pool: %v", err)
		}
	}
}
