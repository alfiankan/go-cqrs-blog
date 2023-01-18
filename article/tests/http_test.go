package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpHandler "github.com/alfiankan/go-cqrs-blog/article/delivery/http/handlers"
	"github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/article/usecases"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"

	"github.com/jaswdr/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHttpApiCreateArticle(t *testing.T) {
	cfg := config.Load("../../.env")
	pgConn, err := infrastructure.NewPgConnection(cfg)
	assert.NoError(t, err)

	esConn, err := infrastructure.NewElasticSearchClient(cfg)
	assert.NoError(t, err)

	writeRepo := repositories.NewArticleWriterPostgree(pgConn)
	readSearchRepo := repositories.NewArticleElasticSearch(esConn)

	// usecases
	articleCommandUseCase := usecases.NewArticleCommand(writeRepo, readSearchRepo)

	// handle http request response
	handler := httpHandler.NewArticleHTTPHandler(articleCommandUseCase)

	t.Run("http created", func(t *testing.T) {

		fake := faker.New()
		jsonReq, err := json.Marshal(&transport.CreateArticle{
			Title:  fake.Lorem().Sentence(5),
			Author: fake.Person().FirstName(),
			Body:   fake.Lorem().Paragraph(1),
		})

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/articles")

		fmt.Println("response json", rec.Body.String())

		assert.NoError(t, handler.CreateArticle(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

	})

	t.Run("http validation error", func(t *testing.T) {
		fake := faker.New()
		jsonReq, err := json.Marshal(&transport.CreateArticle{
			Title:  "",
			Author: fake.Person().FirstName(),
			Body:   fake.Lorem().Paragraph(1),
		})

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/articles")

		fmt.Println("response json", rec.Body.String())

		assert.NoError(t, handler.CreateArticle(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

}
