package database

import (
	"database/sql"

	"github.com/Fozzyack/habit-tracker/internal/env"
	_"github.com/jackc/pgx/v5/stdlib"
)

func Open() (*sql.DB, error) {
	connStr, err := env.GetDbConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if  err != nil {
		return nil, err
	}

	return db, nil

}
