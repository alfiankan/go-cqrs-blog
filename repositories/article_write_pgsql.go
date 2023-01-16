package repositories

import "github.com/alfiankan/go-cqrs-blog/domains"

type ArticleWriterPostgree struct {
}

func NewArticleWriterPostgree() domains.ArticleWriterDbRepository {
	return &ArticleWriterPostgree{}
}

func (repo *ArticleWriterPostgree) Save(article domains.Article) (err error) {
	return
}
