package tests

import (
	"context"
	"testing"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/domains"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/alfiankan/go-cqrs-blog/repositories"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestSaveArticleToWriteDb(t *testing.T) {
	t.Run("save valid article to writedb postgree must be success", func(t *testing.T) {

		cfg := config.Load("../.env")
		pgConn, _ := infrastructure.NewPgConnection(cfg)
		repo := repositories.NewArticleWriterPostgree(pgConn)

		faker := faker.New()
		article := domains.Article{
			Title:  faker.Lorem().Sentence(10),
			Author: faker.Person().FirstName(),
			Body:   faker.Lorem().Paragraph(3),
		}

		// save article
		ctx := context.Background()
		err := repo.Save(ctx, article)

		assert.Nil(t, err)

	})
}
