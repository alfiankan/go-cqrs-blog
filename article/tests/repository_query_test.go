package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	domain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/alfiankan/go-cqrs-blog/config"

	"github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticleIndex(t *testing.T) {
	t.Run("create index article to elasticsearch must be success", func(t *testing.T) {
		faker := faker.New()
		article := domain.Article{
			ID:      time.Now().Unix(),
			Title:   faker.Lorem().Sentence(10),
			Author:  faker.Person().FirstName(),
			Body:    faker.Lorem().Paragraph(3),
			Created: time.Now(),
		}

		// save article
		cfg := config.Load("../../.env")
		esClient, _ := infrastructure.NewElasticSearchClient(cfg)
		repo := repositories.NewArticleElasticSearch(esClient)

		ctx := context.Background()
		err := repo.AddIndex(ctx, article)

		assert.NoError(t, err)

	})
}

func TestGetAllFromES(t *testing.T) {
	t.Run("get all articles from elastic search must be no error", func(t *testing.T) {
		// save article
		cfg := config.Load("../../.env")
		esClient, _ := infrastructure.NewElasticSearchClient(cfg)
		repo := repositories.NewArticleElasticSearch(esClient)

		ctx := context.Background()
		articels, err := repo.Find(ctx, "", "Adam Geitgey")

		for _, article := range articels {
			fmt.Println(article.ID, article.Title, article.Author)
		}

		assert.NoError(t, err)

	})
}
