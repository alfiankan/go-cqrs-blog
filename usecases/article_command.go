package usecases

import (
	"context"
	"time"

	"github.com/alfiankan/go-cqrs-blog/domains"
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

func (uc *ArticleCommand) Create(ctx context.Context, article domains.Article) (err error) {
	// save to write db get insert id
	article.Created = time.Now()

	articleID, err := uc.articleWriteRepo.Save(ctx, article)
	if err != nil {
		return
	}

	article.ID = articleID
	// save index to readdb elasticsearch
	if err = uc.articleReaderRepo.AddIndex(ctx, article); err != nil {
		// fallback delete article from writedb
		err = uc.articleWriteRepo.Delete(ctx, article.ID)
		return
	}
	return
}
