package main

import (
	"fmt"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
)

// create es indices
func createESIndex() error {

	cfg := config.Load()
	es, err := infrastructure.NewElasticSearchClient(cfg)

	if err != nil {
		return err
	}

	res, err := es.Indices.Create("articles")

	fmt.Println(res)

	return nil
}
