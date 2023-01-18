package common

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// migration migrate to database
func Migration(wd string) error {
	cfg := config.Load(fmt.Sprintf("%s/.env", wd))
	pgConn, err := infrastructure.NewPgConnection(cfg)
	driver, err := postgres.WithInstance(pgConn, &postgres.Config{})
	if err != nil {
		return errors.New("CANNOT CONNECT TO DATABASE")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", wd),
		"postgres",
		driver,
	)

	if len(os.Args) > 2 {
		if os.Args[2] == "down" {
			if err := m.Down(); err != nil {
				return err
			}
			log.Println("migration down success")

			return nil
		}
	}

	if err := m.Up(); err != nil {
		return err
	}
	log.Println("migration up success")

	return nil
}
