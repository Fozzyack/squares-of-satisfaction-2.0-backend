package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/habit-tracker/internal/models"
)

type HabitStore interface {
	CreateHabit(ctx context.Context, tx *sql.Tx, userId string, habitReq *models.NewHabitRequest) (*models.Habit, error)
	GetHabitById(id, userId string) (*models.Habit, error)
	GetHabitsByUserId(userId string) ([]*models.Habit, error)
	UpdateHabit(ctx context.Context, tx *sql.Tx, id, userId string, habitReq *models.UpdateHabitRequest) (*models.Habit, error)
}

func NewHabitStore(db *sql.DB) HabitStore {
	return &PostgresStore{db: db}
}

func (ps *PostgresStore) CreateHabit(ctx context.Context, tx *sql.Tx, userId string, habitReq *models.NewHabitRequest) (*models.Habit, error) {
	query := `
	INSERT INTO habits (user_id, name, goal, increment, color, unit)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, user_id, name, goal, increment, color, unit, created_at, updated_at
	`

	newHabit := &models.Habit{}
	err := tx.QueryRowContext(ctx, query, userId, habitReq.Name, habitReq.Goal, habitReq.Increment, habitReq.Color, habitReq.Unit).Scan(
		&newHabit.Id,
		&newHabit.UserId,
		&newHabit.Name,
		&newHabit.Goal,
		&newHabit.Increment,
		&newHabit.Color,
		&newHabit.Unit,
		&newHabit.CreatedAt,
		&newHabit.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return newHabit, nil
}

func (ps *PostgresStore) GetHabitById(id, userId string) (*models.Habit, error) {
	query := `
	SELECT id, user_id, name, goal, increment, color, unit, created_at, updated_at
	FROM habits
	WHERE id = $1 AND user_id = $2
	`

	habit := &models.Habit{}
	err := ps.db.QueryRow(query, id, userId).Scan(
		&habit.Id,
		&habit.UserId,
		&habit.Name,
		&habit.Goal,
		&habit.Increment,
		&habit.Color,
		&habit.Unit,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (ps *PostgresStore) GetHabitsByUserId(userId string) ([]*models.Habit, error) {
	query := `
	SELECT id, user_id, name, goal, increment, color, unit, created_at, updated_at
	FROM habits
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := ps.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := make([]*models.Habit, 0)
	for rows.Next() {
		habit := &models.Habit{}
		err = rows.Scan(
			&habit.Id,
			&habit.UserId,
			&habit.Name,
			&habit.Goal,
			&habit.Increment,
			&habit.Color,
			&habit.Unit,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		habits = append(habits, habit)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (ps *PostgresStore) UpdateHabit(ctx context.Context, tx *sql.Tx, id, userId string, habitReq *models.UpdateHabitRequest) (*models.Habit, error) {
	query := `
	UPDATE habits
	SET
		name = COALESCE($3, name),
		goal = COALESCE($4, goal),
		increment = COALESCE($5, increment),
		color = COALESCE($6, color),
		unit = COALESCE($7, unit),
		updated_at = NOW()
	WHERE id = $1 AND user_id = $2
	RETURNING id, user_id, name, goal, increment, color, unit, created_at, updated_at
	`

	habit := &models.Habit{}
	err := tx.QueryRowContext(ctx, query, id, userId, habitReq.Name, habitReq.Goal, habitReq.Increment, habitReq.Color, habitReq.Unit).Scan(
		&habit.Id,
		&habit.UserId,
		&habit.Name,
		&habit.Goal,
		&habit.Increment,
		&habit.Color,
		&habit.Unit,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return habit, nil
}
