package models

import "time"

type Habit struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	Goal      int       `json:"goal"`
	Increment int       `json:"increment"`
	Color     *string   `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewHabitRequest struct {
	Name      string  `json:"name"`
	Goal      int     `json:"goal"`
	Increment int     `json:"increment"`
	Color     *string `json:"color"`
}
