package tests

import (
	"context"
	"testing"
	"time"

	"github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"

	"github.com/alfiankan/go-cqrs-blog/article/usecases"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewArticle(t *testing.T) {

	t.Run("create new article write to writedb and elastic search for read db must be no error", func(t *testing.T) {
		cfg := config.Load("../../.env")
		pgConn, _ := infrastructure.NewPgConnection(cfg)
		esConn, _ := infrastructure.NewElasticSearchClient(cfg)
		redisConn, _ := infrastructure.NewRedisConnection(cfg)

		writeRepo := repositories.NewArticleWriterPostgree(pgConn)
		readRepo := repositories.NewArticleElasticSearch(esConn)
		cacheRepo := repositories.NewArticleCacheRedis(redisConn, 5*time.Minute)

		articleCommandUseCase := usecases.NewArticleCommand(writeRepo, readRepo, cacheRepo)

		faker := faker.New()
		article := transport.CreateArticle{
			Title:  faker.Lorem().Sentence(10),
			Author: faker.Person().FirstName(),
			Body:   faker.Lorem().Paragraph(3),
		}

		ctx := context.Background()
		err := articleCommandUseCase.Create(ctx, article)
		assert.NoError(t, err)
	})

}
