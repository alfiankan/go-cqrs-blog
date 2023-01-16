package tests

import (
	"testing"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestPostgreeConnection(t *testing.T) {
	t.Run("test connect to postgreesql", func(t *testing.T) {

		config := config.Load("../.env")
		dbConn := infrastructure.NewPgConnection(config)

		// test ping db
		err := dbConn.Ping()

		assert.Nil(t, err)
	})
}
