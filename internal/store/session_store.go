package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type SessionStore interface {
	CreateSession(ctx context.Context, tx *sql.Tx, userId, token string, expiresAt time.Time) (*models.Session, error)
	GetSessionById(id string) (*models.Session, error)
	GetSessionByToken(token string) (*models.Session, error)
}

func NewSessionStore(db *sql.DB) SessionStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) CreateSession(ctx context.Context, tx *sql.Tx, userId, token string, expiresAt time.Time) (*models.Session, error) {
	query := `
	INSERT INTO sessions (user_id, token, expires_at)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, token, expires_at, created_at
	`

	newSession := &models.Session{}
	err := tx.QueryRowContext(ctx, query, userId, token, expiresAt).Scan(
		&newSession.Id,
		&newSession.UserId,
		&newSession.Token,
		&newSession.ExpiresAt,
		&newSession.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (ps *PostgresStore) GetSessionById(id string) (*models.Session, error) {
	query := `
	SELECT id, user_id, token, expires_at, created_at
	FROM sessions
	WHERE id = $1
	`

	session := &models.Session{}
	err := ps.db.QueryRow(query, id).Scan(
		&session.Id,
		&session.UserId,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (ps *PostgresStore) GetSessionByToken(token string) (*models.Session, error) {
	query := `
	SELECT id, user_id, token, expires_at, created_at
	FROM sessions
	WHERE token = $1
	`

	session := &models.Session{}
	err := ps.db.QueryRow(query, token).Scan(
		&session.Id,
		&session.UserId,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}
