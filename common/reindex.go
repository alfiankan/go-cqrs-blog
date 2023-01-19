package common

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alfiankan/go-cqrs-blog/article"
	articleRepos "github.com/alfiankan/go-cqrs-blog/article/repositories"
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// create es indices
func Reindex(wd string) error {

	cfg := config.Load(fmt.Sprintf("%s/.env", wd))
	pgConn, _ := infrastructure.NewPgConnection(cfg)
	esConn, _ := infrastructure.NewElasticSearchClient(cfg)
	readRepo := articleRepos.NewArticleElasticSearch(esConn)

	// precheck es is alive
	if _, err := esConn.Ping(); err != nil {
		return errors.New("elastic search unreachable, try again wait elasticsearch completly running")
	}

	rows, err := pgConn.Query("SELECT id, title, author, body, created FROM articles")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	ctx := context.Background()

	// delete indices
	req := esapi.IndicesDeleteRequest{Index: []string{"articles"}}

	res, err := req.Do(ctx, esConn)
	if err != nil {
		log.Println("delete: request failed", err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("delete: response: %s \n", res.String())
	}
	for rows.Next() {
		var article article.Article

		err := rows.Scan(&article.ID, &article.Title, &article.Author, &article.Body, &article.Created)
		if err != nil {
			log.Fatal(err)
		}

		if err := readRepo.AddIndex(ctx, article); err != nil {
			log.Println("failed index", article.Title)
		}
	}

	return nil
}
