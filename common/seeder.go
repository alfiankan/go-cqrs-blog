package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	articleRepos "github.com/alfiankan/go-cqrs-blog/article/repositories"
	articleUseCases "github.com/alfiankan/go-cqrs-blog/article/usecases"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"
)

// create es indices
func Seed(wd string) error {

	cfg := config.Load(fmt.Sprintf("%s/.env", wd))
	pgConn, _ := infrastructure.NewPgConnection(cfg)
	esConn, _ := infrastructure.NewElasticSearchClient(cfg)
	writeRepo := articleRepos.NewArticleWriterPostgree(pgConn)
	readRepo := articleRepos.NewArticleElasticSearch(esConn)

	articleCommandUseCase := articleUseCases.NewArticleCommand(writeRepo, readRepo)

	seedData, err := os.ReadFile(fmt.Sprintf("%s/articles_seed.json", wd))
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
			err := articleCommandUseCase.Create(ctx, newArticle)
			fmt.Println("inserted", newArticle.Title, err)
			existMap[newArticle.Title] = 1
		}

	}

	fmt.Println("seed success")

	return nil
}
