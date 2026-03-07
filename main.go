package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/app"
	"github.com/Fozzyack/habit-tracker/internal/env"
	"github.com/Fozzyack/habit-tracker/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not parse .env file: %v", err)
	}

	var port int
	flag.IntVar(&port, "port", 8800, "Starting port of the app")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Could not start app: %v", err)
	}
	app.Logger.Info().Msg("App Initialized")
	environment := "DEVELOPMENT"
	if env.GetProduction() {
		environment = "PRODUCTION"
	}
	app.Logger.Info().Str("Env", fmt.Sprintf(environment)).Msg("Env")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes.SetupRoutes(app),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute * 5,
	}

	app.Logger.Info().Int("port", port).Msg("Server Starting")
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Error Listening and Serving")
	}

}
