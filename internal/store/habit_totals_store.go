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
	GetDailyHabitTotal(ctx context.Context, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error)
	GetDailyHabitTotalTx(ctx context.Context, tx *sql.Tx, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error)
	GetHabitYearDailyCounts(ctx context.Context, habitId, userId string) ([]*models.HabitDailyCount, error)
	DeleteDailyHabitTotalsByHabitId(ctx context.Context, userId string, habitId string) error
	DeleteDailyHabitTotalsByHabitIdTx(ctx context.Context, tx *sql.Tx, userId string, habitId string) error
}

func NewHabitTotalStore(db *sql.DB) HabitDailyTotalsStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) deleteDailyHabitTotalsByHabitId(ctx context.Context, e execer, userId string, habitId string) error {
	query := `
	DELETE FROM habit_daily_totals
	WHERE user_id = $1 AND habit_id = $2
	`

	_, err := e.ExecContext(ctx, query, userId, habitId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStore) DeleteDailyHabitTotalsByHabitId(ctx context.Context, userId string, habitId string) error {
	return ps.deleteDailyHabitTotalsByHabitId(ctx, ps.db, userId, habitId)
}

func (ps *PostgresStore) DeleteDailyHabitTotalsByHabitIdTx(ctx context.Context, tx *sql.Tx, userId string, habitId string) error {
	return ps.deleteDailyHabitTotalsByHabitId(ctx, tx, userId, habitId)
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

func (ps *PostgresStore) getDailyHabitTotal(ctx context.Context, q queryRower, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error) {
	query := `
	SELECT id, amount, habit_id, user_id, date, created_at, updated_at
	FROM habit_daily_totals
	WHERE habit_id = $1 AND user_id = $2 AND date = $3
	`

	habitDailyTotal := &models.HabitDailyTotal{}
	err := q.QueryRowContext(ctx, query, habitId, userId, date).Scan(
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

func (ps *PostgresStore) GetDailyHabitTotal(ctx context.Context, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error) {
	return ps.getDailyHabitTotal(ctx, ps.db, habitId, userId, date)
}

func (ps *PostgresStore) GetDailyHabitTotalTx(ctx context.Context, tx *sql.Tx, habitId, userId string, date time.Time) (*models.HabitDailyTotal, error) {
	return ps.getDailyHabitTotal(ctx, tx, habitId, userId, date)
}

func (ps *PostgresStore) GetHabitYearDailyCounts(ctx context.Context, habitId, userId string) ([]*models.HabitDailyCount, error) {
	query := `
	WITH days AS (
		SELECT generate_series(
			CURRENT_DATE - INTERVAL '364 days',
			CURRENT_DATE,
			INTERVAL '1 day'
		)::date AS day
	)
	SELECT d.day, COALESCE(hdt.amount, 0) AS count
	FROM days d
	LEFT JOIN habit_daily_totals hdt
		ON hdt.date = d.day
		AND hdt.habit_id = $1
		AND hdt.user_id = $2
	ORDER BY d.day ASC
	`

	rows, err := ps.db.QueryContext(ctx, query, habitId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dailyCounts := make([]*models.HabitDailyCount, 0)
	for rows.Next() {
		var date time.Time
		dailyCount := &models.HabitDailyCount{}
		err = rows.Scan(&date, &dailyCount.Count)
		if err != nil {
			return nil, err
		}

		dailyCount.Date = date.Format("2006-01-02")
		dailyCounts = append(dailyCounts, dailyCount)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return dailyCounts, nil
}
