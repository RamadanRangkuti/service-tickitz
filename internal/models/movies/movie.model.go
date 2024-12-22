package models

import "time"

type Movies struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Synopsis    string     `json:"synopsis"`
	DirectorID  int        `json:"director_id"`
	Duration    int        `json:"duration"`
	ReleaseDate time.Time  `json:"release_date"`
	Image       string     `json:"image"`
	Banner      string     `json:"banner"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
