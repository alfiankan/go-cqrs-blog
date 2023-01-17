package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alfiankan/go-cqrs-blog/config"
	_ "github.com/lib/pq"
)

// NewPgConnection create db connection with pool
// return *sql.DB
func NewPgConnection(config config.ApplicationConfig) (db *sql.DB, err error) {

	dbURI := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		config.PostgreeHost,
		config.PostgreeUser,
		config.PostgreePass,
		config.PostgreeDb,
		config.PostgreePort,
		config.PostgreeSsl,
	)
	db, err = sql.Open("postgres", dbURI)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(5)

	return
}
