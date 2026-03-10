package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/Fozzyack/habit-tracker/database"
	"github.com/Fozzyack/habit-tracker/internal/auth"
	"github.com/Fozzyack/habit-tracker/migrations"
	"github.com/joho/godotenv"
)

const (
	seedUserName     = "John Doe"
	seedUserEmail    = "john@example.com"
	seedUserPassword = "password123"
	seedHabitName    = "Drink Water"
	seedHabitGoal    = 10
	seedHabitIncr    = 1
	seedHabitColor   = "#3f7e9e"
	seedHabitUnit    = "cups"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not parse .env file: %v", err)
	}

	db, err := database.Open()
	if err != nil {
		log.Fatalf("failed opening database: %v", err)
	}
	defer db.Close()

	err = database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	err = seed(context.Background(), db)
	if err != nil {
		log.Fatalf("failed seeding database: %v", err)
	}

	log.Println("seed complete: 1 user, 1 habit")
}

func seed(ctx context.Context, db *sql.DB) error {
	passwordHash, err := auth.HashPassword(seedUserPassword)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		TRUNCATE TABLE users, sessions, habits, habit_entries, habit_daily_totals
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		return err
	}

	var userID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`, seedUserName, seedUserEmail, passwordHash).Scan(&userID)
	if err != nil {
		return err
	}

	var habitID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO habits (user_id, name, goal, increment, color, unit)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, userID, seedHabitName, seedHabitGoal, seedHabitIncr, seedHabitColor, seedHabitUnit).Scan(&habitID)
	if err != nil {
		return err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	today := time.Now().UTC().Truncate(24 * time.Hour)
	start := today.AddDate(0, 0, -364)

	for day := start; !day.After(today); day = day.AddDate(0, 0, 1) {
		dailyTotal := 0
		logsForDay := rng.Intn(5)

		for i := 0; i < logsForDay; i++ {
			incrementAmount := rng.Intn(3) + 1
			dailyTotal += incrementAmount

			logCreatedAt := day.Add(time.Duration(rng.Intn(24))*time.Hour + time.Duration(rng.Intn(60))*time.Minute)
			_, err = tx.ExecContext(ctx, `
				INSERT INTO habit_entries (increment_amount, habit_id, user_id, created_at)
				VALUES ($1, $2, $3, $4)
			`, incrementAmount, habitID, userID, logCreatedAt)
			if err != nil {
				return err
			}
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO habit_daily_totals (amount, habit_id, user_id, date)
			VALUES ($1, $2, $3, $4)
		`, dailyTotal, habitID, userID, day)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
