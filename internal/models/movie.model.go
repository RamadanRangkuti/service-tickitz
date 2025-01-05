package models

import "time"

type MovieDetails struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Synopsis    string     `json:"synopsis"`
	Duration    int        `json:"duration"`
	ReleaseDate time.Time  `json:"release_date"`
	Image       string     `json:"image"`
	Banner      string     `json:"banner"`
	Director    string     `json:"director_name"`
	Casts       []*string  `json:"casts"`
	Genres      []*string  `json:"genres"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ListMovies []MovieDetails
