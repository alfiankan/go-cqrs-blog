package infrastructure

import (
	"context"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/go-redis/redis/v9"
)

// NewRedisConnection create new redis connection
func NewRedisConnection(config config.ApplicationConfig) (client redis.UniversalClient, err error) {

	client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    config.RedisHost,
		Password: config.RedisPass,
	})

	if _, err = client.Ping(context.Background()).Result(); err != nil {
		return
	}

	return
}
