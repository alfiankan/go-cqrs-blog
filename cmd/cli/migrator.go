package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alfiankan/go-cqrs-blog/config"
	"github.com/alfiankan/go-cqrs-blog/infrastructure"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// migration migrate to database
func migration() error {
	config := config.Load()
	pgConn, err := infrastructure.NewPgConnection(config)
	driver, err := postgres.WithInstance(pgConn, &postgres.Config{})
	if err != nil {
		return errors.New("CANNOT CONNECT TO DATABASE")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)

	if len(os.Args) > 2 {
		if os.Args[2] == "down" {
			if err := m.Down(); err != nil {
				return err
			}
			fmt.Println("migration down success")

			return nil
		}
	}

	if err := m.Up(); err != nil {
		return err
	}
	fmt.Println("migration up success")

	return nil
}
