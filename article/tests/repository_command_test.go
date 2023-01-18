package tests

import (
	"context"
	"testing"

	domain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestSaveArticleToWriteDb(t *testing.T) {
	t.Run("save valid article to writedb postgree must be success", func(t *testing.T) {

		cfg := config.Load("../../.env")
		pgConn, _ := infrastructure.NewPgConnection(cfg)
		repo := repositories.NewArticleWriterPostgree(pgConn)

		faker := faker.New()
		article := domain.Article{
			Title:  faker.Lorem().Sentence(10),
			Author: faker.Person().FirstName(),
			Body:   faker.Lorem().Paragraph(3),
		}

		// save article
		ctx := context.Background()
		_, err := repo.Save(ctx, article)

		assert.NoError(t, err)

	})
}
