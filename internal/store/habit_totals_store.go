package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type HabitDailyTotalsStore interface {
	CreateDailyHabitTotal(ctx context.Context, tx *sql.Tx, amount int, userId string, habitId string, date time.Time) (*models.HabitDailyTotal, error)
	UpdateDailyHabitTotal(ctx context.Context, tx *sql.Tx, habitDailyTotal *models.HabitDailyTotal) (*models.HabitDailyTotal, error)
	GetDailyHabitTotal(habitId, userId string, date time.Time) (*models.HabitDailyTotal, error)
	GetDailyHabitTotalInTx(ctx context.Context, tx *sql.Tx, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error)
}

func NewHabitTotalStore(db *sql.DB) HabitDailyTotalsStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) CreateDailyHabitTotal(ctx context.Context, tx *sql.Tx, amount int, userId string, habitId string, date time.Time) (*models.HabitDailyTotal, error) {
	query := `
	INSERT into habit_daily_totals ( amount, habit_id, user_id, date)
	VALUES ($1, $2, $3, $4)
	RETURNING id, amount, habit_id, user_id, date, created_at, updated_at
	`

	dbHabitDailyTotal := &models.HabitDailyTotal{}
	err := tx.QueryRowContext(ctx, query, amount, habitId, userId, date).Scan(
		&dbHabitDailyTotal.Id,
		&dbHabitDailyTotal.Amount,
		&dbHabitDailyTotal.HabitId,
		&dbHabitDailyTotal.UserId,
		&dbHabitDailyTotal.Date,
		&dbHabitDailyTotal.CreatedAt,
		&dbHabitDailyTotal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return dbHabitDailyTotal, nil
}

func (ps *PostgresStore) UpdateDailyHabitTotal(ctx context.Context, tx *sql.Tx, habitDailyTotal *models.HabitDailyTotal) (*models.HabitDailyTotal, error) {
	query := `
	UPDATE habit_daily_totals
	SET
		amount = $4,
		updated_at = NOW()
	WHERE habit_id = $1 AND user_id = $2 AND date = $3
	RETURNING id, amount, habit_id, user_id, date, created_at, updated_at
	`

	updatedHabitDailyTotal := &models.HabitDailyTotal{}
	err := tx.QueryRowContext(ctx, query, habitDailyTotal.HabitId, habitDailyTotal.UserId, habitDailyTotal.Date, habitDailyTotal.Amount).Scan(
		&updatedHabitDailyTotal.Id,
		&updatedHabitDailyTotal.Amount,
		&updatedHabitDailyTotal.HabitId,
		&updatedHabitDailyTotal.UserId,
		&updatedHabitDailyTotal.Date,
		&updatedHabitDailyTotal.CreatedAt,
		&updatedHabitDailyTotal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return updatedHabitDailyTotal, nil
}

func (ps *PostgresStore) GetDailyHabitTotal(habitId, userId string, date time.Time) (*models.HabitDailyTotal, error) {
	query := `
	SELECT id, amount, habit_id, user_id, date, created_at, updated_at
	FROM habit_daily_totals
	WHERE habit_id = $1 AND user_id = $2 AND date = $3
	`

	habitDailyTotal := &models.HabitDailyTotal{}
	err := ps.db.QueryRow(query, habitId, userId, date).Scan(
		&habitDailyTotal.Id,
		&habitDailyTotal.Amount,
		&habitDailyTotal.HabitId,
		&habitDailyTotal.UserId,
		&habitDailyTotal.Date,
		&habitDailyTotal.CreatedAt,
		&habitDailyTotal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return habitDailyTotal, nil
}

func (ps *PostgresStore) GetDailyHabitTotalInTx(ctx context.Context, tx *sql.Tx, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error) {
	query := `
	SELECT id, amount, habit_id, user_id, date, created_at, updated_at
	FROM habit_daily_totals
	WHERE habit_id = $1 AND user_id = $2 AND date = $3
	`

	habitDailyTotal := &models.HabitDailyTotal{}
	err := tx.QueryRowContext(ctx, query, habitId, userId, date).Scan(
		&habitDailyTotal.Id,
		&habitDailyTotal.Amount,
		&habitDailyTotal.HabitId,
		&habitDailyTotal.UserId,
		&habitDailyTotal.Date,
		&habitDailyTotal.CreatedAt,
		&habitDailyTotal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return habitDailyTotal, nil
}
