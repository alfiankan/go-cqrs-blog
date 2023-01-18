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
}

func NewArticleCommand(writeRepo domain.ArticleWriterDbRepository, readRepo domain.ArticleReaderDbRepository) domain.ArticleCommand {
	return &ArticleCommand{
		articleWriteRepo:  writeRepo,
		articleReaderRepo: readRepo,
	}
}

func (uc *ArticleCommand) Create(ctx context.Context, article transport.CreateArticle) (err error) {
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
	if err = uc.articleReaderRepo.AddIndex(ctx, newArticle); err != nil {
		// fallback delete article from writedb
		err = uc.articleWriteRepo.Delete(ctx, newArticle.ID)
		return
	}
	return
}
