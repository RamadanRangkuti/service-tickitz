package models

import "time"

type User struct {
	Id          int        `json:"id" form:"id"`
	Email       string     `json:"email" form:"email"`
	Password    string     `json:"password" form:"password"`
	Role        string     `json:"role" form:"role"`
	FirstName   *string    `json:"first_name" form:"first_name"`
	LastName    *string    `json:"last_name" form:"last_name"`
	PhoneNumber *string    `json:"phone_number" form:"phone_number"`
	Image       *string    `json:"image" form:"image"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UserDetails struct {
	Id          int     `json:"id" form:"id"`
	Email       string  `json:"email" form:"email"`
	Password    string  `json:"password" form:"password"`
	Role        string  `json:"role" form:"role"`
	FirstName   *string `json:"first_name" form:"first_name"`
	LastName    *string `json:"last_name" form:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	Image       *string `json:"image"`
}

type Users []User
