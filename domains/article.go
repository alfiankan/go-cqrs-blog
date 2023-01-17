package domains

import "context"

// Article base domain
type Article struct {
	ID      int64  `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

// ArticleCommand is a usecase interface for (C) Command from CQRS
// Save write data to persistence db and write to search db
type ArticleCommand interface {
	Create(article Article) (err error)
}

// ArticleQuery is a usecase interface for (Q) Query from CQRS
// Get will find and get data from cache first if the data doesnt exist will continue to use search db
type ArticleQuery interface {
	Get(ctx context.Context, keyword, author string) (articles []Article, err error)
}

// ArticleWriterDbRepository is a repository interface for writing data to db
// Save to any database implemetation need
type ArticleWriterDbRepository interface {
	Save(ctx context.Context, article Article) (err error)
}

// ArticleReaderDbRepository is a repository interface for read and search
// Read from search database
type ArticleReaderDbRepository interface {
	AddIndex(ctx context.Context, article Article) (err error)
	FindAll() (ctx context.Context, articles []Article, err error)
	Find(ctx context.Context, keyword, author string) (articles []Article, err error)
}

// ArticleCacheRepository is a interface for rw cache
// ReadByQueryTerm accept term parameter, term parameter notated by query param combination
// cache.ReadByQueryTerm("keyword=lorem&author=john")
type ArticleCacheRepository interface {
	Write(ctx context.Context, article Article) (err error)
	ReadAll(ctx context.Context) (articles []Article, err error)
	ReadByQueryTerm(ctx context.Context, term string) (article []Article, err error)
}
