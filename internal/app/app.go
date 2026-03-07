package app

import (
	"os"

	envutils "github.com/Fozzyack/habit-tracker/internal/utils/env_utils"
	"github.com/rs/zerolog"
)

type Application struct {
	Logger zerolog.Logger
}

func NewApplication() (*Application, error) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if !envutils.GetProduction() {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	app := &Application{
		Logger: logger,
	}

	return app, nil

}
