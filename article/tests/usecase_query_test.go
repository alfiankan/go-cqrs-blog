package tests

import (
	"context"
	"testing"
	"time"

	"github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/article/usecases"
	"github.com/alfiankan/go-cqrs-blog/common"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestQueryArticleAll(t *testing.T) {

	cfg := config.Load("../../.env")
	redisConn, _ := infrastructure.NewRedisConnection(cfg)
	esConn, _ := infrastructure.NewElasticSearchClient(cfg)

	cacheRepo := repositories.NewArticleCacheRedis(redisConn, 5*time.Minute)
	readRepo := repositories.NewArticleElasticSearch(esConn)

	articleQueryUseCase := usecases.NewArticleQuery(readRepo, cacheRepo)
	ctx := context.Background()

	res, err := articleQueryUseCase.Get(ctx, "", "Adam Geitgey", 1)

	assert.True(t, len(res) > 0)
	assert.NoError(t, err)

}

func TestQueryArticleSearchFilter(t *testing.T) {

	cfg := config.Load("../../.env")
	redisConn, _ := infrastructure.NewRedisConnection(cfg)
	esConn, _ := infrastructure.NewElasticSearchClient(cfg)

	cacheRepo := repositories.NewArticleCacheRedis(redisConn, 5*time.Minute)
	readRepo := repositories.NewArticleElasticSearch(esConn)

	articleQueryUseCase := usecases.NewArticleQuery(readRepo, cacheRepo)
	ctx := context.Background()

	res, err := articleQueryUseCase.Get(ctx, common.EmptyString, common.EmptyString, 1)
	assert.True(t, len(res) > 0)
	assert.NoError(t, err)

}
