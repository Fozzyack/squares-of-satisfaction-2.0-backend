package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/Fozzyack/habit-tracker/internal/env"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	connStr, err := env.GetDbConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil

}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)

}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Goose Dialect set: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("Goose up: %w", err)
	}
	return nil

}
