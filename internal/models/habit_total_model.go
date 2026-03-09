package models

import "time"

type HabitDailyTotal struct {
	Id        string    `json:"id"`
	Amount    int       `json:"amount"`
	HabitId   string    `json:"habit_id"`
	UserId    string    `json:"user_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HabitDailyCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
