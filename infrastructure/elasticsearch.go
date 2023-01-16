package infrastructure

import (
	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearchClient(config config.ApplicationConfig) (es *elasticsearch.Client, err error) {
	cfg := elasticsearch.Config{
		Addresses: config.ElasticSearchAdresses,
		Username:  config.ElasticSearchUsername,
		Password:  config.ElasticSearchPassword,
	}
	es, err = elasticsearch.NewClient(cfg)
	return
}
