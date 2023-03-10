package tests

import (
	"testing"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestPostgreeConnection(t *testing.T) {
	t.Run("test connect to postgreesql with config", func(t *testing.T) {

		cfg := config.Load("../.env")
		_, err := infrastructure.NewPgConnection(cfg)

		assert.NoError(t, err)
	})
}

func TestESConnection(t *testing.T) {
	t.Run("test connection to elasticsearch with config", func(t *testing.T) {

		cfg := config.Load("../.env")

		_, err := infrastructure.NewElasticSearchClient(cfg)

		assert.NoError(t, err)

	})
}

func TestRedisConnection(t *testing.T) {
	t.Run("test redis initialize connection", func(t *testing.T) {
		cfg := config.Load("../.env")

		_, err := infrastructure.NewRedisConnection(cfg)

		assert.NoError(t, err)
	})
}
