package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func setupHandler() *httpHandler.ArticleHTTPHandler {
	cfg := config.Load("../../.env")
	pgConn, _ := infrastructure.NewPgConnection(cfg)
	esConn, _ := infrastructure.NewElasticSearchClient(cfg)
	redisConn, _ := infrastructure.NewRedisConnection(cfg)

	writeRepo := repositories.NewArticleWriterPostgree(pgConn)
	readSearchRepo := repositories.NewArticleElasticSearch(esConn)
	cacheRepo := repositories.NewArticleCacheRedis(redisConn, time.Minute*5)

	// usecases
	articleCommandUseCase := usecases.NewArticleCommand(writeRepo, readSearchRepo)
	articleQueryUseCase := usecases.NewArticleQuery(readSearchRepo, cacheRepo)

	// handle http request response
	return httpHandler.NewArticleHTTPHandler(articleCommandUseCase, articleQueryUseCase)
}

func TestHttpApiCreateArticle(t *testing.T) {
	handler := setupHandler()
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

		assert.NoError(t, handler.CreateArticle(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

}

func TestHttpFindArticles(t *testing.T) {

	handler := setupHandler()

	t.Run("http success", func(t *testing.T) {

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/articles?keyword=part 2&author=Adam Geitgey", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/articles")

		assert.NoError(t, handler.FindArticle(c))
		assert.Equal(t, http.StatusOK, rec.Code)

	})

}
