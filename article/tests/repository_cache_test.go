package tests

import (
	"context"
	"testing"
	"time"

	articleDomain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestAddArticleQueryCache(t *testing.T) {
	cfg := config.Load("../../.env")
	redisConn, _ := infrastructure.NewRedisConnection(cfg)
	ttl := time.Minute * 5
	repo := repositories.NewArticleCacheRedis(redisConn, ttl)

	fake := faker.New()
	articles := []articleDomain.Article{}
	articles = append(articles, articleDomain.Article{
		Title:  fake.Lorem().Sentence(5),
		Author: fake.Person().FirstName(),
		Body:   fake.Lorem().Paragraph(1),
	})

	searchTerm := "keyword=lorem&author=ipsum"
	ctx := context.Background()
	t.Run("test to add cache from query (cqrs) must be no error", func(t *testing.T) {

		err := repo.Write(ctx, searchTerm, articles)

		assert.NoError(t, err)
	})

	t.Run("test to get cache from query (cqrs) must be no error", func(t *testing.T) {

		articles, err := repo.ReadByQueryTerm(ctx, searchTerm)

		assert.True(t, len(articles) > 0)
		assert.NoError(t, err)

	})

	t.Run("test invalidate cache must be no error", func(t *testing.T) {

		err := repo.InvalidateCache(ctx)

		assert.NoError(t, err)

	})

}
