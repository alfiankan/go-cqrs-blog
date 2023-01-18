package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	httpHandlers "github.com/alfiankan/go-cqrs-blog/article/delivery/http/handlers"
	common "github.com/alfiankan/go-cqrs-blog/common/middleware"
	"github.com/alfiankan/go-cqrs-blog/config"
	echoSwagger "github.com/swaggo/echo-swagger"

	articleRepos "github.com/alfiankan/go-cqrs-blog/article/repositories"
	articleUseCases "github.com/alfiankan/go-cqrs-blog/article/usecases"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initInfrastructure(cfg config.ApplicationConfig) (pgConn *sql.DB, esConn *elasticsearch.Client, redisConn *redis.Client) {
	pgConn, err := infrastructure.NewPgConnection(cfg)
	if err != nil {
		panic(err)
	}

	esConn, err = infrastructure.NewElasticSearchClient(cfg)
	if err != nil {
		panic(err)
	}

	redisConn, err = infrastructure.NewRedisConnection(cfg)
	if err != nil {
		panic(err)
	}
	return
}

func initArticleApplication(httpServer *echo.Echo, cfg config.ApplicationConfig) {

	// infrastructure
	pgConn, esConn, redisConn := initInfrastructure(cfg)

	// repositories
	articleWriteRepo := articleRepos.NewArticleWriterPostgree(pgConn)
	articleReadSearchRepo := articleRepos.NewArticleElasticSearch(esConn)
	articleCacheRepo := articleRepos.NewArticleCacheRedis(redisConn, 5*time.Minute)

	// usecases
	articleCommandUseCase := articleUseCases.NewArticleCommand(articleWriteRepo, articleReadSearchRepo)
	articleReadUseCase := articleUseCases.NewArticleQuery(articleReadSearchRepo, articleCacheRepo)

	// handle http request response
	httpHandlers.NewArticleHTTPHandler(articleCommandUseCase, articleReadUseCase).HandleRoute(httpServer)

}

// @title go-cqrs-blog-api
// @version 3.0
// @description Go implemented cqrs.
// @contact.name alfiankan
// @contact.url https://github.com/alfiankan
// @contact.email alfiankan19@gmail.com
// @license.name Apache 2.0
// @BasePath /
func main() {

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(common.MiddlewaresRegistry...)

	cfg := config.Load()
	initArticleApplication(e, cfg)

	// swagger
	url := echoSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger.yaml", cfg.HTTPApiPort))
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	e.Static("/docs", "docs")

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPApiPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 60 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
