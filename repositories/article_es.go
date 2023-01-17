package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/alfiankan/go-cqrs-blog/domains"
	transport "github.com/alfiankan/go-cqrs-blog/transport/elasticsearch"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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
func (repo *ArticleElasticSearch) Find(ctx context.Context, keyword, author string) (articles []domains.Article, err error) {

	// setup query
	esBoolQuery := esquery.Bool()

	if keyword != "" {
		esBoolQuery.Boost(1.0)
		esBoolQuery.MinimumShouldMatch(1)
		esBoolQuery.Should(esquery.Term("title", strings.ToLower(keyword)), esquery.Term("body", strings.ToLower(keyword)))
	}

	if author != "" {
		esBoolQuery.Filter(esquery.Term("author", strings.ToLower(author)))
	}

	res, err := esquery.Search().Query(esBoolQuery).Run(
		repo.es,
		repo.es.Search.WithIndex("articles"),
		repo.es.Search.WithContext(ctx),
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
