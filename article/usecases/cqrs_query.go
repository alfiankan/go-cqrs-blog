package usecases

import (
	"context"
	"fmt"

	domain "github.com/alfiankan/go-cqrs-blog/article"
)

type ArticleQuery struct {
	articleReaderRepo domain.ArticleReaderDbRepository
	articleCacheRepo  domain.ArticleCacheRepository
}

func NewArticleQuery(articleReaderRepo domain.ArticleReaderDbRepository, articleCacheRepo domain.ArticleCacheRepository) domain.ArticleQuery {

	return &ArticleQuery{articleReaderRepo, articleCacheRepo}

}

func (uc *ArticleQuery) Get(ctx context.Context, keyword, author string, page uint64) (articles []domain.Article, err error) {

	queryTerm := fmt.Sprintf("keyword=%s&author=%s&page=%d", keyword, author, page)

	// get from cache first
	articles, err = uc.articleCacheRepo.ReadByQueryTerm(ctx, queryTerm)
	if err != nil {

		// get from readdb elasticsearch
		articles, err = uc.articleReaderRepo.Find(ctx, keyword, author, page)
		if err != nil {
			return
		}

		// set cache
		if len(articles) > 0 {
			if err = uc.articleCacheRepo.Write(ctx, queryTerm, articles); err != nil {
				return
			}
		}

	}

	return
}
