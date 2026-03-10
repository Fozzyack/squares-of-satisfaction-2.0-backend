package tests

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Fozzyack/habit-tracker/database"
	"github.com/Fozzyack/habit-tracker/migrations"
)

const defaultTestDatabaseURL = "postgresql://postgres:postgres@localhost:5501/squares"

func OpenTestDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_TEST_URL")
	if connStr == "" {
		connStr = defaultTestDatabaseURL
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupTestDB() (*sql.DB, error) {
	db, err := OpenTestDB()
	if err != nil {
		return nil, err
	}

	err = database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate test db: %w", err)
	}

	return db, nil
}
