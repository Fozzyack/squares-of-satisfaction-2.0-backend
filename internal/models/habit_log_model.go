package models

import "time"

type HabitLog struct {
	Id              string    `json:"id"`
	IncrementAmount int       `json:"increment_amount"`
	HabitId         string    `json:"habit_id"`
	UserId          string    `json:"user_id"`
	Date            string    `json:"date"`
	CreatedAt       time.Time `json:"created_at"`
}
