package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alfiankan/go-cqrs-blog/config"

	httpHandlers "github.com/alfiankan/go-cqrs-blog/delivery/http/handlers"
	middlewares "github.com/alfiankan/go-cqrs-blog/delivery/http/middleware"

	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/alfiankan/go-cqrs-blog/repositories"
	"github.com/alfiankan/go-cqrs-blog/usecases"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
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

func initApplication(httpServer *echo.Echo, cfg config.ApplicationConfig) {

	// infrastructure
	pgConn, esConn, _ := initInfrastructure(cfg)

	// repositories
	writeRepo := repositories.NewArticleWriterPostgree(pgConn)
	readSearchRepo := repositories.NewArticleElasticSearch(esConn)

	// usecases
	articleCommandUseCase := usecases.NewArticleCommand(writeRepo, readSearchRepo)

	// handle http request response
	httpHandlers.NewArticleHTTPHandler(articleCommandUseCase).HandleRoute(httpServer)

}

func main() {

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middlewares.SecureMiddleware())

	cfg := config.Load()
	initApplication(e, cfg)

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
