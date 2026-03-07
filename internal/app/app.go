package app

import (
	"database/sql"
	"os"

	"github.com/Fozzyack/habit-tracker/database"
	"github.com/Fozzyack/habit-tracker/internal/env"
	"github.com/Fozzyack/habit-tracker/internal/store"
	"github.com/Fozzyack/habit-tracker/migrations"
	"github.com/rs/zerolog"
)

type Application struct {
	Logger zerolog.Logger
	DB     *sql.DB
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

	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}

	return app, nil

}
