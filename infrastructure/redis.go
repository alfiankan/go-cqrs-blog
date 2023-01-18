package infrastructure

import (
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/go-redis/redis"
)

// NewRedisConnection create new redis connection
func NewRedisConnection(config config.ApplicationConfig) (client *redis.Client, err error) {

	client = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPass,
		DB:       0,
	})
	if _, err = client.Ping().Result(); err != nil {
		return
	}

	return
}
