package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type HabitLogStore interface {
	CreateHabitLog(ctx context.Context, tx *sql.Tx, incrementAmount int, userId, habitId string) (*models.HabitLog, error)
	GetHabitLogsByHabitId(habitId, userId string) ([]*models.HabitLog, error)
	GetHabitLogsByUserId(userId string) ([]*models.HabitLog, error)
}

func NewHabitLogStore(db *sql.DB) HabitLogStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) CreateHabitLog(ctx context.Context, tx *sql.Tx, incrementAmount int, userId, habitId string) (*models.HabitLog, error) {
	query := `
	INSERT INTO habit_entries (increment_amount, habit_id, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, increment_amount, habit_id, user_id, created_at
	`

	habitLog := &models.HabitLog{}
	err := tx.QueryRowContext(ctx, query, incrementAmount, habitId, userId).Scan(
		&habitLog.Id,
		&habitLog.IncrementAmount,
		&habitLog.HabitId,
		&habitLog.UserId,
		&habitLog.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return habitLog, nil
}

func (ps *PostgresStore) GetHabitLogsByHabitId(habitId, userId string) ([]*models.HabitLog, error) {
	query := `
	SELECT id, increment_amount, habit_id, user_id, created_at
	FROM habit_entries
	WHERE habit_id = $1 AND user_id = $2
	ORDER BY created_at DESC
	`

	rows, err := ps.db.Query(query, habitId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habitLogs := make([]*models.HabitLog, 0)
	for rows.Next() {
		habitLog := &models.HabitLog{}
		err = rows.Scan(
			&habitLog.Id,
			&habitLog.IncrementAmount,
			&habitLog.HabitId,
			&habitLog.UserId,
			&habitLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		habitLogs = append(habitLogs, habitLog)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return habitLogs, nil
}

func (ps *PostgresStore) GetHabitLogsByUserId(userId string) ([]*models.HabitLog, error) {
	query := `
	SELECT id, increment_amount, habit_id, user_id, created_at
	FROM habit_entries
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := ps.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habitLogs := make([]*models.HabitLog, 0)
	for rows.Next() {
		habitLog := &models.HabitLog{}
		err = rows.Scan(
			&habitLog.Id,
			&habitLog.IncrementAmount,
			&habitLog.HabitId,
			&habitLog.UserId,
			&habitLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		habitLogs = append(habitLogs, habitLog)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return habitLogs, nil
}
