package app

import (
	"database/sql"
	"os"

	"github.com/Fozzyack/habit-tracker/internal/database"
	"github.com/Fozzyack/habit-tracker/internal/env"
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

	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}

	return app, nil

}
