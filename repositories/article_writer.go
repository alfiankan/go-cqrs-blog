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
func (repo *ArticleWriterPostgree) Save(ctx context.Context, article domains.Article) (id int64, err error) {
	sql := "INSERT INTO articles (title, author, body, created) VALUES ($1, $2, $3, $4) RETURNING id"
	err = repo.db.QueryRowContext(ctx, sql, article.Title, article.Author, article.Body, article.Created).Scan(&id)

	return
}

func (repo *ArticleWriterPostgree) Delete(ctx context.Context, id int64) (err error) {
	sql := "DELETE FROM articles WHERE id = $1"
	_, err = repo.db.ExecContext(ctx, sql, id)

	return
}
