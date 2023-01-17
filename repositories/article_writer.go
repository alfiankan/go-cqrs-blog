package repositories

import (
	"context"
	"database/sql"

	"github.com/alfiankan/go-cqrs-blog/domains"
)

type ArticleWriterPostgree struct {
	db *sql.DB
}

func NewArticleWriterPostgree(db *sql.DB) domains.ArticleWriterDbRepository {
	return &ArticleWriterPostgree{db}
}

// Save save article to write database (postgree)
func (repo *ArticleWriterPostgree) Save(ctx context.Context, article domains.Article) (err error) {
	sql := "INSERT INTO articles (title, author, body) VALUES ($1, $2, $3)"
	_, err = repo.db.ExecContext(ctx, sql, article.Title, article.Author, article.Body)
	return
}
