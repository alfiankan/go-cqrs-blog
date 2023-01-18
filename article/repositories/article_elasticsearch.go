package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	articleDomain "github.com/alfiankan/go-cqrs-blog/article"

	"github.com/alfiankan/go-cqrs-blog/common"
	transport "github.com/alfiankan/go-cqrs-blog/transport/response"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ArticleElasticSearch struct {
	es *elasticsearch.Client
}

func NewArticleElasticSearch(es *elasticsearch.Client) articleDomain.ArticleReaderDbRepository {
	return &ArticleElasticSearch{es}
}

// AddIndex request to index data to elastic search
func (repo *ArticleElasticSearch) AddIndex(ctx context.Context, article articleDomain.Article) (err error) {

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
func (repo *ArticleElasticSearch) Find(ctx context.Context, keyword, author string) (articles []articleDomain.Article, err error) {

	// setup query
	esBoolQuery := esquery.Bool()

	if keyword != common.EmptyString {
		esBoolQuery.Boost(2.0)
		esBoolQuery.MinimumShouldMatch(1)
		esBoolQuery.Should(esquery.MatchPhrase("title", keyword), esquery.MatchPhrase("body", keyword))
	}

	if author != common.EmptyString {
		esBoolQuery.Filter(esquery.MatchPhrase("author", strings.ToLower(author)))
	}

	res, err := esquery.Search().Query(esBoolQuery).From(0).Size(1000).Run(
		repo.es,
		repo.es.Search.WithIndex("articles"),
		repo.es.Search.WithContext(ctx),
	)

	defer res.Body.Close()

	// mapping to article domain
	hits := transport.EsHits[*articleDomain.Article]{}
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
