package tests

import (
	"testing"

	"github.com/alfiankan/go-cqrs-blog/domains"
	"github.com/jaswdr/faker"
)

func TestMain(m *testing.M) {
	// setup repository and database
	// if postgree and elasticsearch not running
}

func TestSaveArticleToWriterDb(t *testing.T) {
	t.Run("save valid article to writedb postgree must be success", func(t *testing.T) {
		faker := faker.New()
		article := domains.Article{
			Title:  faker.Lorem().Text(10),
			Author: faker.Person().FirstName(),
			Body:   faker.Lorem().Paragraph(3),
		}

		// save article

	})
}
