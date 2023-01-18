package usecases

import (
	"context"
	"time"

	domain "github.com/alfiankan/go-cqrs-blog/article"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"
)

// ArticleCommand implementation from domain.ArticleCommand
// hold usecase/busineeslogic for command related operation (CQRS)
type ArticleCommand struct {
	articleWriteRepo  domain.ArticleWriterDbRepository
	articleReaderRepo domain.ArticleReaderDbRepository
	articleCacheRepo  domain.ArticleCacheRepository
}

func NewArticleCommand(
	writeRepo domain.ArticleWriterDbRepository,
	readRepo domain.ArticleReaderDbRepository,
	cacheRepo domain.ArticleCacheRepository,
) domain.ArticleCommand {
	return &ArticleCommand{
		articleWriteRepo:  writeRepo,
		articleReaderRepo: readRepo,
		articleCacheRepo:  cacheRepo,
	}
}

func (uc *ArticleCommand) Create(ctx context.Context, article transport.CreateArticle) (err error) {

	// invalidate cache
	err = uc.articleCacheRepo.InvalidateCache(ctx)
	if err != nil {
		return
	}
	// save to write db get insert id
	newArticle := domain.Article{
		Title:   article.Title,
		Author:  article.Author,
		Body:    article.Body,
		Created: time.Now(),
	}

	articleID, err := uc.articleWriteRepo.Save(ctx, newArticle)
	if err != nil {
		return
	}

	newArticle.ID = articleID
	// save index to readdb elasticsearch
	if esErr := uc.articleReaderRepo.AddIndex(ctx, newArticle); esErr != nil {
		// fallback delete article from writedb
		if err = uc.articleWriteRepo.Delete(ctx, newArticle.ID); err != nil {
			return
		}
		err = esErr
		return
	}
	return
}
