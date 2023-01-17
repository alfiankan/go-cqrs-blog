package tests

import (
	"fmt"
	"testing"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestPostgreeConnection(t *testing.T) {
	t.Run("test connect to postgreesql with config", func(t *testing.T) {

		config := config.Load("../.env")
		_, err := infrastructure.NewPgConnection(config)

		assert.Nil(t, err)
	})
}

func TestESConnection(t *testing.T) {
	t.Run("test connection to elasticsearch with config", func(t *testing.T) {

		config := config.Load("../.env")

		esConn, err := infrastructure.NewElasticSearchClient(config)

		fmt.Println(esConn.Info())

		assert.Nil(t, err)

	})
}

func TestRedisConnection(t *testing.T) {
	t.Run("test redis initialize connection", func(t *testing.T) {
		config := config.Load("../.env")

		_, err := infrastructure.NewRedisConnection(config)

		assert.Nil(t, err)
	})
}
