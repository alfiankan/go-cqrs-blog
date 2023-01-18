package repositories

import (
	"context"
	"encoding/json"
	"time"

	domain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/go-redis/redis"
)

// ArticleCacheRedis implementation from domain.ArticleCacheRepository
// using redis as cache
type ArticleCacheRedis struct {
	cacheClient *redis.Client
	cacheTTL    time.Duration
}

func NewArticleCacheRedis(cacheClient *redis.Client, ttl time.Duration) domain.ArticleCacheRepository {
	return &ArticleCacheRedis{cacheClient, ttl}
}

// Write save (set) search/query cache to redis json format from domain.Article
func (repo *ArticleCacheRedis) Write(ctx context.Context, term string, articles []domain.Article) (err error) {

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		return
	}
	err = repo.cacheClient.WithContext(ctx).Set(term, jsonArticles, repo.cacheTTL).Err()

	return
}

// ReadByQueryTerm read query/search cahce by queryparamterm as index
func (repo *ArticleCacheRedis) ReadByQueryTerm(ctx context.Context, term string) (articles []domain.Article, err error) {

	res := repo.cacheClient.WithContext(ctx).Get(term)
	if res.Err() != nil {
		err = res.Err()
		return
	}

	var cacheResult string
	if err = res.Scan(&cacheResult); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(cacheResult), &articles); err != nil {
		return
	}

	return
}

// InvalidateCache invalidate all cache
func (repo *ArticleCacheRedis) InvalidateCache(ctx context.Context) (err error) {
	err = repo.cacheClient.FlushAll().Err()
	return
}
