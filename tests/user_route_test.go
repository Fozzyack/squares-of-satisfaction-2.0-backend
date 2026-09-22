package tests

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fozzyack/habit-tracker/internal/api"
	"github.com/Fozzyack/habit-tracker/internal/app"
	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/routes"
	"github.com/Fozzyack/habit-tracker/internal/services"
	"github.com/Fozzyack/habit-tracker/internal/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCookieName = "test_session"

func newTestApplication(db *sql.DB) *app.Application {
	logger := zerolog.Nop()
	txManager := services.NewSQLTxManager(db)
	userStore := store.NewUserStore(db)
	sessionStore := store.NewSessionStore(db)
	habitStore := store.NewHabitStore(db)
	habitDailyTotalsStore := store.NewHabitTotalStore(db)
	habitLogStore := store.NewHabitLogStore(db)

	authService := services.NewAuthService(txManager, userStore, sessionStore)
	habitService := services.NewHabitService(txManager, habitStore, habitDailyTotalsStore, habitLogStore)
	userService := services.NewUserService(txManager, userStore, sessionStore, habitStore, habitDailyTotalsStore, habitLogStore)

	return &app.Application{
		Logger:       logger,
		DB:           db,
		UserHandler:  api.NewUserHandler(authService, userService, logger),
		HabitHandler: api.NewHabitHandler(habitService, logger),
		HabitService: habitService,
		SessionStore: sessionStore,
	}
}

func countUserRows(t *testing.T, db *sql.DB, table, userID string) int {
	t.Helper()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE user_id = $1", userID).Scan(&count)
	require.NoError(t, err)

	return count
}

func seedHabitActivity(t *testing.T, habitService *services.HabitService, userID string) {
	t.Helper()

	habit, err := habitService.CreateHabit(context.Background(), userID, &models.NewHabitRequest{
		Name:      "Test Habit",
		Goal:      10,
		Increment: 1,
	})
	require.NoError(t, err)

	_, err = habitService.RecordHabit(context.Background(), userID, &models.RecordHabitRequest{
		HabitId: habit.Id,
		Amount:  2,
		Date:    "2026-09-16",
	})
	require.NoError(t, err)
}

func TestDeleteUserRoute(t *testing.T) {
	t.Setenv("COOKIE_NAME", testCookieName)

	db, err := SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE habit_entries, habit_daily_totals, habits, sessions, users CASCADE")
	require.NoError(t, err)

	application := newTestApplication(db)
	authService := services.NewAuthService(
		services.NewSQLTxManager(db),
		store.NewUserStore(db),
		store.NewSessionStore(db),
	)

	deletedUser, deletedSession, err := authService.CreateNewUser(context.Background(), &models.NewUserRequest{
		Name:     "delete-me",
		Email:    "delete-me@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	otherUser, _, err := authService.CreateNewUser(context.Background(), &models.NewUserRequest{
		Name:     "keep-me",
		Email:    "keep-me@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	seedHabitActivity(t, application.HabitService, deletedUser.Id)
	seedHabitActivity(t, application.HabitService, otherUser.Id)

	router := routes.SetupRoutes(application)

	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: deletedSession.Token})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	deleteCookies := rec.Result().Cookies()
	require.Len(t, deleteCookies, 1)
	assert.Equal(t, testCookieName, deleteCookies[0].Name)
	assert.Less(t, deleteCookies[0].MaxAge, 0)

	for _, table := range []string{"habit_entries", "habit_daily_totals", "habits", "sessions"} {
		assert.Zero(t, countUserRows(t, db, table, deletedUser.Id), table)
	}

	var remainingUsers int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", deletedUser.Id).Scan(&remainingUsers)
	require.NoError(t, err)
	assert.Zero(t, remainingUsers)

	assert.Equal(t, 1, countUserRows(t, db, "habits", otherUser.Id))
	assert.Equal(t, 1, countUserRows(t, db, "habit_entries", otherUser.Id))
	assert.Equal(t, 1, countUserRows(t, db, "habit_daily_totals", otherUser.Id))
	assert.Equal(t, 1, countUserRows(t, db, "sessions", otherUser.Id))

	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: deletedSession.Token})
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
