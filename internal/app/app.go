package app

import (
	"database/sql"
	"os"

	"github.com/Fozzyack/habit-tracker/database"
	"github.com/Fozzyack/habit-tracker/internal/api"
	"github.com/Fozzyack/habit-tracker/internal/env"
	"github.com/Fozzyack/habit-tracker/internal/services"
	"github.com/Fozzyack/habit-tracker/internal/store"
	"github.com/Fozzyack/habit-tracker/migrations"
	"github.com/rs/zerolog"
)

type Application struct {
	Logger       zerolog.Logger
	DB           *sql.DB
	UserHandler  *api.UserHandler
	HabitService *services.HabitService
	SessionStore store.SessionStore
}

func NewApplication() (*Application, error) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if !env.GetProduction() {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	pgDB, err := database.Open()
	if err != nil {
		return nil, err
	}

	err = database.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	// store init
	userStore := store.NewUserStore(pgDB)
	sessionStore := store.NewSessionStore(pgDB)
	habitStore := store.NewHabitStore(pgDB)

	// services init
	txManager := services.NewSQLTxManager(pgDB)
	authService := services.NewAuthService(txManager, userStore, sessionStore)
	habitService := services.NewHabitService(txManager, habitStore)

	// handler init
	userHandler := api.NewUserHandler(authService, logger)

	app := &Application{
		Logger:       logger,
		DB:           pgDB,
		UserHandler:  userHandler,
		HabitService: habitService,
		SessionStore: sessionStore,
	}

	return app, nil

}
