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
	common "github.com/alfiankan/go-cqrs-blog/common"
	middlewares "github.com/alfiankan/go-cqrs-blog/common/middleware"

	"github.com/alfiankan/go-cqrs-blog/config"
	echoSwagger "github.com/swaggo/echo-swagger"

	articleRepos "github.com/alfiankan/go-cqrs-blog/article/repositories"
	articleUseCases "github.com/alfiankan/go-cqrs-blog/article/usecases"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

// initInfrastructure init all infrastructure needs to run this application
func initInfrastructure(cfg config.ApplicationConfig) (pgConn *sql.DB, esConn *elasticsearch.Client, redisConn redis.UniversalClient) {

	pgConn, err := infrastructure.NewPgConnection(cfg)
	common.LogExit(err, common.LOG_LEVEL_ERROR)

	esConn, err = infrastructure.NewElasticSearchClient(cfg)
	common.LogExit(err, common.LOG_LEVEL_ERROR)

	redisConn, err = infrastructure.NewRedisConnection(cfg)
	common.LogExit(err, common.LOG_LEVEL_ERROR)

	return
}

// initArticleApplication init app by injecting deps
func initArticleApplication(httpServer *echo.Echo, cfg config.ApplicationConfig) {

	pgConn, esConn, redisConn := initInfrastructure(cfg)

	// repositories
	articleWriteRepo := articleRepos.NewArticleWriterPostgree(pgConn)
	articleReadSearchRepo := articleRepos.NewArticleElasticSearch(esConn)
	articleCacheRepo := articleRepos.NewArticleCacheRedis(redisConn, 5*time.Minute)

	// usecases
	articleCommandUseCase := articleUseCases.NewArticleCommand(articleWriteRepo, articleReadSearchRepo, articleCacheRepo)
	articleReadUseCase := articleUseCases.NewArticleQuery(articleReadSearchRepo, articleCacheRepo)

	// handle http request response
	httpHandlers.NewArticleHTTPHandler(articleCommandUseCase, articleReadUseCase).HandleRoute(httpServer)

}

// @title go-cqrs-blog-api
// @version 1.0
// @description Go implemented cqrs.
// @contact.name alfiankan
// @contact.url https://github.com/alfiankan
// @contact.email alfiankan19@gmail.com
// @license.name Apache 2.0
// @BasePath /
func main() {

	cfg := config.Load()
	e := echo.New()
	e.Use(middlewares.MiddlewaresRegistry...)

	initArticleApplication(e, cfg)

	// swagger api docs
	url := echoSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger.yaml", cfg.HTTPApiPort))
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	e.Static("/docs", "docs")

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPApiPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err.Error())
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
