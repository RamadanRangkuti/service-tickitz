package models

import "time"

type MovieGenre struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	MovieID   int        `json:"movie_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
