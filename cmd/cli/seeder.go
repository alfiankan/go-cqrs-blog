package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/alfiankan/go-cqrs-blog/repositories"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"
	"github.com/alfiankan/go-cqrs-blog/usecases"
)

// create es indices
func seed() error {

	cfg := config.Load()
	pgConn, _ := infrastructure.NewPgConnection(cfg)
	esConn, _ := infrastructure.NewElasticSearchClient(cfg)
	writeRepo := repositories.NewArticleWriterPostgree(pgConn)
	readRepo := repositories.NewArticleElasticSearch(esConn)

	articleCommandUseCase := usecases.NewArticleCommand(writeRepo, readRepo)

	seedData, err := os.ReadFile("articles_seed.json")
	if err != nil {
		return err
	}

	var articles []map[string]any

	if err := json.Unmarshal(seedData, &articles); err != nil {
		return err
	}

	ctx := context.Background()

	existMap := map[string]int{}

	for _, article := range articles {
		newArticle := transport.CreateArticle{
			Title:  article["title"].(string),
			Author: article["author"].(string),
			Body:   article["text"].(string),
		}
		if existMap[newArticle.Title] == 0 {
			articleCommandUseCase.Create(ctx, newArticle)
			existMap[newArticle.Title] = 1
		}

	}

	fmt.Println("seed success")

	return nil
}
