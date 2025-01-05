package models

import "time"

type GetCinema struct {
	Id   int       `db:"id"`
	Name string    `db:"name"`
	City string    `db:"city"`
	Date time.Time `db:"date"`
	Time time.Time `db:"time"`
}

type CinemaSeat struct {
	SeatNumber  string  `json:"seat_number"`
	IsAvailable bool    `json:"is_available"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
}

type Transaction struct {
	ID         int     `json:"id"`
	UserId     int     `json:"user_id"`
	MovieId    int     `json:"movie_id"`
	Date       string  `json:"date"`
	Location   string  `json:"location"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type PaymentDetail struct {
	MovieName  string   `json:"movie_name"`
	Date       string   `json:"date"`
	Time       string   `json:"time"`
	Quantity   int      `json:"quantity"`
	Seats      []string `json:"seats"`
	TotalPrice float64  `json:"total_price"`
}

type Ticket struct {
	MovieTitle string   `json:"movie_title"`
	Date       string   `json:"date"` // Format: 07 Jul 2025
	Time       string   `json:"time"` // Format: 2:00 PM
	Quantity   int      `json:"quantity"`
	Seats      []string `json:"seats"` // Example: ["A1", "A2", "A3"]
	TotalPrice float64  `json:"total_price"`
}

type GetCinemas []GetCinema
