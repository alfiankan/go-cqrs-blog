package repositories

import (
	"context"
	"encoding/json"
	"time"

	articleDomain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/go-redis/redis"
)

type ArticleCacheRedis struct {
	cacheClient *redis.Client
	cacheTTL    time.Duration
}

func NewArticleCacheRedis(cacheClient *redis.Client, ttl time.Duration) articleDomain.ArticleCacheRepository {
	return &ArticleCacheRedis{cacheClient, ttl}
}

func (repo *ArticleCacheRedis) Write(ctx context.Context, term string, articles []articleDomain.Article) (err error) {

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		return
	}
	err = repo.cacheClient.WithContext(ctx).Set(term, jsonArticles, repo.cacheTTL).Err()

	return
}

func (repo *ArticleCacheRedis) ReadByQueryTerm(ctx context.Context, term string) (articles []articleDomain.Article, err error) {

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
func (repo *ArticleCacheRedis) InvalidateCache(ctx context.Context) (err error) {

	err = repo.cacheClient.FlushAll().Err()

	return
}
