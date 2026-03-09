package models

import "time"

type Habit struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	Goal      int       `json:"goal"`
	Increment int       `json:"increment"`
	Color     *string   `json:"color"`
	Unit      *string   `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewHabitRequest struct {
	Name      string  `json:"name"`
	Goal      int     `json:"goal"`
	Increment int     `json:"increment"`
	Color     *string `json:"color"`
	Unit      *string `json:"unit"`
}

type UpdateHabitRequest struct {
	Name      *string `json:"name"`
	Goal      *int    `json:"goal"`
	Increment *int    `json:"increment"`
	Color     *string `json:"color"`
	Unit      *string `json:"unit"`
}

type RecordHabitRequest struct {
	HabitId string `json:"habit_id"`
	Amount  int    `json:"amount"`
	Date    string `json:"date"`
}
