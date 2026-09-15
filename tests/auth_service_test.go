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

func newAuthService(db *sql.DB) *services.AuthService {
	txManager := services.NewSQLTxManager(db)
	userStore := store.NewUserStore(db)
	sessionStore := store.NewSessionStore(db)

	return services.NewAuthService(txManager, userStore, sessionStore)
}

func truncateAuthTables(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec("TRUNCATE TABLE sessions, users CASCADE")
	require.NoError(t, err)
}

func TestAuthServiceCreateNewUserTestCases(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	authService := newAuthService(db)

	testCases := []struct {
		name        string
		input       *models.NewUserRequest
		setup       func(t *testing.T)
		expectedErr bool
	}{
		{
			name: "creates user and session",
			input: &models.NewUserRequest{
				Name:     "new-user",
				Email:    "new-user@example.com",
				Password: "password123",
			},
			setup:       func(_ *testing.T) {},
			expectedErr: false,
		},
		{
			name: "duplicate email returns error",
			input: &models.NewUserRequest{
				Name:     "duplicate-user",
				Email:    "duplicate@example.com",
				Password: "password123",
			},
			setup: func(t *testing.T) {
				_, _, err := authService.CreateNewUser(context.Background(), &models.NewUserRequest{
					Name:     "existing-user",
					Email:    "duplicate@example.com",
					Password: "password123",
				})
				require.NoError(t, err)
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			truncateAuthTables(t, db)
			tc.setup(t)

			user, session, err := authService.CreateNewUser(context.Background(), tc.input)

			assert.Equal(t, tc.expectedErr, err != nil)
			if tc.expectedErr {
				assert.Nil(t, user)
				assert.Nil(t, session)
				return
			}

			assert.NotNil(t, user)
			assert.NotNil(t, session)
			assert.Equal(t, tc.input.Email, user.Email)
		})
	}
}

func TestAuthServiceLoginUserTestCases(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	authService := newAuthService(db)

	testCases := []struct {
		name        string
		input       *models.LoginUserRequest
		setup       func(t *testing.T)
		expectedErr error
	}{
		{
			name: "email not found returns invalid credentials",
			input: &models.LoginUserRequest{
				Email:    "missing@example.com",
				Password: "password123",
			},
			setup:       func(_ *testing.T) {},
			expectedErr: services.ErrInvalidCredentials,
		},
		{
			name: "wrong password returns invalid credentials",
			input: &models.LoginUserRequest{
				Email:    "user@example.com",
				Password: "wrong-password",
			},
			setup: func(t *testing.T) {
				_, _, err := authService.CreateNewUser(context.Background(), &models.NewUserRequest{
					Name:     "login-user",
					Email:    "user@example.com",
					Password: "password123",
				})
				require.NoError(t, err)
			},
			expectedErr: services.ErrInvalidCredentials,
		},
		{
			name: "valid login returns session",
			input: &models.LoginUserRequest{
				Email:    "valid-user@example.com",
				Password: "password123",
			},
			setup: func(t *testing.T) {
				_, _, err := authService.CreateNewUser(context.Background(), &models.NewUserRequest{
					Name:     "valid-user",
					Email:    "valid-user@example.com",
					Password: "password123",
				})
				require.NoError(t, err)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			truncateAuthTables(t, db)
			tc.setup(t)

			session, err := authService.LoginUser(context.Background(), tc.input)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, session)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, session)
			assert.Equal(t, tc.input.Email, mustGetUserEmailByID(t, db, session.UserId))
		})
	}
}

func TestAuthServiceUpdateUserTestCases(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	authService := newAuthService(db)

	testCases := []struct {
		name        string
		newName     string
		newPassword string
		currentPass string
		expectedErr error
	}{
		{
			name:        "updates name",
			newName:     "updated-user",
			currentPass: "password123",
		},
		{
			name:        "updates password",
			newPassword: "new-password",
			currentPass: "password123",
		},
		{
			name:        "updates name and password",
			newName:     "updated-user",
			newPassword: "new-password",
			currentPass: "password123",
		},
		{
			name:        "wrong current password returns invalid credentials",
			newName:     "updated-user",
			currentPass: "wrong-password",
			expectedErr: services.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			truncateAuthTables(t, db)
			user, _, err := authService.CreateNewUser(context.Background(), &models.NewUserRequest{
				Name:     "original-user",
				Email:    "update-user@example.com",
				Password: "password123",
			})
			require.NoError(t, err)

			request := &models.UpdateUserRequest{CurrentPassword: tc.currentPass}
			if tc.newName != "" {
				request.Name = &tc.newName
			}
			if tc.newPassword != "" {
				request.NewPassword = &tc.newPassword
			}

			updatedUser, err := authService.UpdateUser(context.Background(), user.Id, request)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, updatedUser)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, updatedUser)
			if tc.newName != "" {
				assert.Equal(t, tc.newName, updatedUser.Name)
			}
			if tc.newPassword != "" {
				_, err = authService.LoginUser(context.Background(), &models.LoginUserRequest{
					Email:    "update-user@example.com",
					Password: tc.newPassword,
				})
				assert.NoError(t, err)
			}
		})
	}
}

func mustGetUserEmailByID(t *testing.T, db *sql.DB, userID string) string {
	t.Helper()

	var email string
	err := db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
	require.NoError(t, err)

	return email
}
