package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/services"
	"github.com/Fozzyack/habit-tracker/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHabitServiceUndoHabit(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE habit_entries, habit_daily_totals, habits, users CASCADE")
	require.NoError(t, err)

	var userID, habitID string
	err = db.QueryRow(`
		INSERT INTO users (name, email, password_hash)
		VALUES ('undo-user', 'undo-user@example.com', 'test-hash')
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)

	err = db.QueryRow(`
		INSERT INTO habits (user_id, name, goal, increment, unit)
		VALUES ($1, 'Undo Habit', 10, 1, 'time')
		RETURNING id
	`, userID).Scan(&habitID)
	require.NoError(t, err)

	habitService := services.NewHabitService(
		services.NewSQLTxManager(db),
		store.NewHabitStore(db),
		store.NewHabitTotalStore(db),
		store.NewHabitLogStore(db),
	)
	date := "2026-09-16"
	for _, amount := range []int{2, 3} {
		_, err = habitService.RecordHabit(context.Background(), userID, &models.RecordHabitRequest{
			HabitId: habitID,
			Amount:  amount,
			Date:    date,
		})
		require.NoError(t, err)
	}

	total, err := habitService.UndoHabit(context.Background(), userID, habitID, date)
	require.NoError(t, err)
	assert.Equal(t, 2, total.Amount)

	var entries, amount int
	err = db.QueryRow(`
		SELECT COUNT(*), COALESCE((SELECT amount FROM habit_daily_totals WHERE habit_id = $1), 0)
		FROM habit_entries
		WHERE habit_id = $1
	`, habitID).Scan(&entries, &amount)
	require.NoError(t, err)
	assert.Equal(t, 1, entries)
	assert.Equal(t, 2, amount)

	_, err = habitService.UndoHabit(context.Background(), userID, habitID, date)
	require.NoError(t, err)
	_, err = habitService.UndoHabit(context.Background(), userID, habitID, date)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
