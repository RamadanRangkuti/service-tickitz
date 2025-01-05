package dto

import "time"

type CreateUserDTO struct {
	Id          int     `json:"id" form:"id" swaggerignore:"true"`
	Email       string  `json:"email" form:"email"`
	Password    string  `json:"password" form:"password"`
	Role        string  `json:"role" form:"role"`
	FirstName   *string `json:"first_name" form:"first_name"`
	LastName    *string `json:"last_name" form:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	Image       *string `json:"image" swaggerignore:"true"`
}

type UpdateUserDTO struct {
	Email       *string `json:"email" form:"email"`
	Password    *string `json:"password" form:"password"`
	FirstName   *string `json:"first_name" form:"first_name"`
	LastName    *string `json:"last_name" form:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	Image       *string `json:"image" swaggerignore:"true"`
}

type UserResponse struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
