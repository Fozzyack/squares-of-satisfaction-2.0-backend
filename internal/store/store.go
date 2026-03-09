package store

import (
	"context"
	"database/sql"
)


type PostgresStore struct {
	db *sql.DB
}

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
