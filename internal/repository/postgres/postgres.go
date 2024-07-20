package postgres

import (
	_ "database/sql"
	"fmt"
	"todo-list/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Client *sqlx.DB

	url string
}

func New(cfg config.DB) (db DB, err error) {
	// postgres://{username}:{password}@{localhost}:{5432}/{dbname}?sslmode=disable
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db.Client, err = sqlx.Connect("postgres", url)
	if err != nil {
		return
	}

	db.url = url

	return
}

func (db *DB) Close() error {
	if db.Client != nil {
		db.Client.Close()
	}
	return nil
}

func (db *DB) Migrate() error {
	if db.url != "" {
		mg, err := migrate.New("file://migrations/postgres", db.url)
		if err != nil {
			return err
		}

		if err = mg.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}
