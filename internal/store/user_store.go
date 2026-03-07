package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type UserStore interface {
	CreateUser(ctx context.Context, tx *sql.Tx, passwordHash string, userReq *models.NewUserRequest) (*models.User, error)
	GetUserById(id string) (*models.User, error)
}

func NewUserStore(db *sql.DB) UserStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) CreateUser(ctx context.Context, tx *sql.Tx, passwordHash string, userReq *models.NewUserRequest) (*models.User, error) {
	query := `
	INSERT into users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, email, password_hash, name, created_at, updated_at
	`

	newUser := &models.User{}
	err := tx.QueryRowContext(ctx, query, userReq.Name, userReq.Email, passwordHash).Scan(
		&newUser.Id,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.Name,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (ps *PostgresStore) GetUserById(id string) (*models.User, error) {
	query := `
	SELECT id, email, password_hash, name, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	user := &models.User{}
	err := ps.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
