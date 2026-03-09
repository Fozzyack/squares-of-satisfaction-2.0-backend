package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type SessionStore interface {
	CreateSession(ctx context.Context, tx *sql.Tx, userId, token string, expiresAt time.Time) (*models.Session, error)
	GetSessionById(ctx context.Context, id string) (*models.Session, error)
	GetSessionByIdTx(ctx context.Context, tx *sql.Tx, id string) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
	GetSessionByTokenTx(ctx context.Context, tx *sql.Tx, token string) (*models.Session, error)
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

func (ps *PostgresStore) getSessionById(ctx context.Context, q queryRower, id string) (*models.Session, error) {
	query := `
	SELECT id, user_id, token, expires_at, created_at
	FROM sessions
	WHERE id = $1
	`

	session := &models.Session{}
	err := q.QueryRowContext(ctx, query, id).Scan(
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

func (ps *PostgresStore) GetSessionById(ctx context.Context, id string) (*models.Session, error) {
	return ps.getSessionById(ctx, ps.db, id)
}

func (ps *PostgresStore) GetSessionByIdTx(ctx context.Context, tx *sql.Tx, id string) (*models.Session, error) {
	return ps.getSessionById(ctx, tx, id)
}

func (ps *PostgresStore) getSessionByToken(ctx context.Context, q queryRower, token string) (*models.Session, error) {
	query := `
	SELECT id, user_id, token, expires_at, created_at
	FROM sessions
	WHERE token = $1
	`

	session := &models.Session{}
	err := q.QueryRowContext(ctx, query, token).Scan(
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

func (ps *PostgresStore) GetSessionByToken(ctx context.Context, token string) (*models.Session, error) {
	return ps.getSessionByToken(ctx, ps.db, token)
}

func (ps *PostgresStore) GetSessionByTokenTx(ctx context.Context, tx *sql.Tx, token string) (*models.Session, error) {
	return ps.getSessionByToken(ctx, tx, token)
}
