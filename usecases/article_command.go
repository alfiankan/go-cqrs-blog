package usecases

import (
	"context"
	"time"

	"github.com/alfiankan/go-cqrs-blog/domains"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"
)

type ArticleCommand struct {
	articleWriteRepo  domains.ArticleWriterDbRepository
	articleReaderRepo domains.ArticleReaderDbRepository
}

func NewArticleCommand(writeRepo domains.ArticleWriterDbRepository, readRepo domains.ArticleReaderDbRepository) domains.ArticleCommand {
	return &ArticleCommand{
		articleWriteRepo:  writeRepo,
		articleReaderRepo: readRepo,
	}
}

func (uc *ArticleCommand) Create(ctx context.Context, article transport.CreateArticle) (err error) {
	// save to write db get insert id

	newArticle := domains.Article{
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
