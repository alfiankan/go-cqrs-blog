package infrastructure

import (
	"fmt"
	"log"
	"time"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPgConnection create db connection with pool
// return *sqlx.DB
func NewPgConnection(config config.ApplicationConfig) (db *sqlx.DB) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
			config.PostgreeHost,
			config.PostgreeUser,
			config.PostgreePass,
			config.PostgreeDb,
			config.PostgreePort,
			config.PostgreeSsl,
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	// setup db pool
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(25)
	return
}