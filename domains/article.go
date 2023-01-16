package domains

// Article base domain
type Article struct {
	ID      int64  `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

// ArticleCommand is interface for (C) Command from CQRS
// Save write data to persistence db and write to search db
type ArticleCommand interface {
	Save(article Article) (err error)
}

// ArticleQuery is interface for (Q) Query from CQRS
// Get will find and get data from cache first if the data doesnt exist will continue to use search db
type ArticleQuery interface {
	Get(keyword, author string) (articles []Article, err error)
}
