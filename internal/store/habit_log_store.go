package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type HabitLogStore interface {
	CreateHabitLog(ctx context.Context, tx *sql.Tx, incrementAmount int, userId, habitId string, date time.Time) (*models.HabitLog, error)
	GetHabitLogsByDate(ctx context.Context, userId string, date time.Time) ([]*models.HabitLog, error)
	GetHabitLogsByHabitId(ctx context.Context, habitId, userId string) ([]*models.HabitLog, error)
	GetHabitLogsByHabitIdTx(ctx context.Context, tx *sql.Tx, habitId, userId string) ([]*models.HabitLog, error)
	GetHabitLogsByUserId(ctx context.Context, userId string) ([]*models.HabitLog, error)
	GetHabitLogsByUserIdTx(ctx context.Context, tx *sql.Tx, userId string) ([]*models.HabitLog, error)
	DeleteHabitLogByHabitId(ctx context.Context, userId string, habitId string) error
	DeleteHabitLogByHabitIdTx(ctx context.Context, tx *sql.Tx, userId string, habitId string) error
}

func NewHabitLogStore(db *sql.DB) HabitLogStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) deleteHabitById(ctx context.Context, e execer, userId string, habitId string) error {
	query := `
	DELETE FROM habit_entries
	WHERE user_id = $1 AND habit_id = $2
	`

	_, err := e.ExecContext(ctx, query, userId, habitId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStore) DeleteHabitLogByHabitId(ctx context.Context, userId string, habitId string) error {
	return ps.deleteHabitById(ctx, ps.db, userId, habitId)
}

func (ps *PostgresStore) DeleteHabitLogByHabitIdTx(ctx context.Context, tx *sql.Tx, userId string, habitId string) error {
	return ps.deleteHabitById(ctx, tx, userId, habitId)
}

func (ps *PostgresStore) GetHabitLogsByDate(ctx context.Context, userId string, date time.Time) ([]*models.HabitLog, error) {
	query := `
		SELECT id, increment_amount, habit_id, user_id, date, created_at
		FROM habit_entries
		WHERE user_id = $1 AND date = $2
		ORDER BY created_at DESC
	`
	habitLogs := make([]*models.HabitLog, 0)
	rows, err := ps.db.QueryContext(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var logDate time.Time
		habitLog := &models.HabitLog{}
		err := rows.Scan(
			&habitLog.Id,
			&habitLog.IncrementAmount,
			&habitLog.HabitId,
			&habitLog.UserId,
			&logDate,
			&habitLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		habitLog.Date = logDate.Format("2006-01-02")
		habitLogs = append(habitLogs, habitLog)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return habitLogs, nil
}

func (ps *PostgresStore) CreateHabitLog(ctx context.Context, tx *sql.Tx, incrementAmount int, userId, habitId string, date time.Time) (*models.HabitLog, error) {
	query := `
	INSERT INTO habit_entries (increment_amount, habit_id, user_id, date)
	VALUES ($1, $2, $3, $4)
	RETURNING id, increment_amount, habit_id, user_id, date, created_at
	`

	habitLog := &models.HabitLog{}
	var logDate time.Time
	err := tx.QueryRowContext(ctx, query, incrementAmount, habitId, userId, date).Scan(
		&habitLog.Id,
		&habitLog.IncrementAmount,
		&habitLog.HabitId,
		&habitLog.UserId,
		&logDate,
		&habitLog.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	habitLog.Date = logDate.Format("2006-01-02")

	return habitLog, nil
}

func (ps *PostgresStore) getHabitLogsByHabitId(ctx context.Context, q queryer, habitId, userId string) ([]*models.HabitLog, error) {
	query := `
	SELECT id, increment_amount, habit_id, user_id, date, created_at
	FROM habit_entries
	WHERE habit_id = $1 AND user_id = $2
	ORDER BY created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, habitId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habitLogs := make([]*models.HabitLog, 0)
	for rows.Next() {
		var logDate time.Time
		habitLog := &models.HabitLog{}
		err = rows.Scan(
			&habitLog.Id,
			&habitLog.IncrementAmount,
			&habitLog.HabitId,
			&habitLog.UserId,
			&logDate,
			&habitLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		habitLog.Date = logDate.Format("2006-01-02")

		habitLogs = append(habitLogs, habitLog)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return habitLogs, nil
}

func (ps *PostgresStore) GetHabitLogsByHabitId(ctx context.Context, habitId, userId string) ([]*models.HabitLog, error) {
	return ps.getHabitLogsByHabitId(ctx, ps.db, habitId, userId)
}

func (ps *PostgresStore) GetHabitLogsByHabitIdTx(ctx context.Context, tx *sql.Tx, habitId, userId string) ([]*models.HabitLog, error) {
	return ps.getHabitLogsByHabitId(ctx, tx, habitId, userId)
}

func (ps *PostgresStore) getHabitLogsByUserId(ctx context.Context, q queryer, userId string) ([]*models.HabitLog, error) {
	query := `
	SELECT id, increment_amount, habit_id, user_id, date, created_at
	FROM habit_entries
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habitLogs := make([]*models.HabitLog, 0)
	for rows.Next() {
		var logDate time.Time
		habitLog := &models.HabitLog{}
		err = rows.Scan(
			&habitLog.Id,
			&habitLog.IncrementAmount,
			&habitLog.HabitId,
			&habitLog.UserId,
			&logDate,
			&habitLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		habitLog.Date = logDate.Format("2006-01-02")

		habitLogs = append(habitLogs, habitLog)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return habitLogs, nil
}

func (ps *PostgresStore) GetHabitLogsByUserId(ctx context.Context, userId string) ([]*models.HabitLog, error) {
	return ps.getHabitLogsByUserId(ctx, ps.db, userId)
}

func (ps *PostgresStore) GetHabitLogsByUserIdTx(ctx context.Context, tx *sql.Tx, userId string) ([]*models.HabitLog, error) {
	return ps.getHabitLogsByUserId(ctx, tx, userId)
}
