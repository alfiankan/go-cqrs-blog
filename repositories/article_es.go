package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/alfiankan/go-cqrs-blog/domains"
	transport "github.com/alfiankan/go-cqrs-blog/transport/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ArticleElasticSearch struct {
	es *elasticsearch.Client
}

func NewArticleElasticSearch(es *elasticsearch.Client) domains.ArticleReaderDbRepository {
	return &ArticleElasticSearch{es}
}

// AddIndex request to index data to elastic search
func (repo *ArticleElasticSearch) AddIndex(ctx context.Context, article domains.Article) (err error) {

	data, err := json.Marshal(article)
	if err != nil {
		return
	}

	// create es index request
	req := esapi.IndexRequest{
		Index:      "articles",
		DocumentID: strconv.Itoa(int(article.ID)),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	_, err = req.Do(ctx, repo.es)

	return
}

// FindAll get all articles from elastic search
func (repo *ArticleElasticSearch) FindAll(ctx context.Context) (articles []domains.Article, err error) {

	// setup query
	var buf bytes.Buffer
	query := map[string]any{
		"query": map[string]any{
			"match_all": map[string]any{},
		},
		"sort": []any{
			map[string]any{"created": map[string]any{
				"order":         "desc",
				"missing":       "_first",
				"unmapped_type": "date",
			}},
		},
	}
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		return
	}

	// call es api
	res, err := repo.es.Search(
		repo.es.Search.WithContext(ctx),
		repo.es.Search.WithIndex("articles"),
		repo.es.Search.WithBody(&buf),
	)
	defer res.Body.Close()

	// mapping to article domain
	hits := transport.EsHits[*domains.Article]{}
	if err = json.NewDecoder(res.Body).Decode(&hits); err != nil {
		return
	}

	if len(hits.Hits.Hits) == 0 {
		return
	}

	for _, source := range hits.Hits.Hits {
		articles = append(articles, *source.Source)
	}

	return
}

func (repo *ArticleElasticSearch) Find(ctx context.Context, keyword, author string) (articles []domains.Article, err error) {
	return
}
