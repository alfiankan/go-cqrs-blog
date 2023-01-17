package main

import (
	"fmt"
	"strings"

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

	mapping := `
{
  "settings": {
    "number_of_shards": 1
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      },
      "title": {
        "type": "keyword"
      },
      "author": {
        "type": "keyword"
      },
      "body": {
        "type": "keyword"
      },
      "created": {
        "type": "date"
      }
    }
  }
}`

	res, err := es.Indices.Create("articles", es.Indices.Create.WithBody(strings.NewReader(mapping)))

	fmt.Println(res)

	return nil
}
