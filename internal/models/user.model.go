package models

import "time"

type User struct {
	Id        int        `json:"id" form:"id"`
	Email     string     `json:"email" form:"email"`
	Password  string     `json:"password" form:"password"`
	RoleId    int        `json:"role_id" form:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `jsin:"updated_at"`
}

type UserProfile struct {
	Id          int        `json:"id" form:"id"`
	UserId      int        `json:"user_id" form:"user_id"`
	FirstName   string     `json:"first_name" form:"first_name"`
	LastName    string     `json:"last_name" form:"last_name"`
	PhoneNumber string     `json:"phone_number" form:"phone_number"`
	Image       *string    `json:"image" form:"image"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `jsin:"updated_at"`
}

type UserRoles struct {
	Id        int        `json:"id" form:"id"`
	Role      string     `json:"role" form:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `jsin:"updated_at"`
}

type UserDetails struct {
	Id          int     `json:"id" form:"id"`
	Email       string  `json:"email" form:"email"`
	Password    string  `json:"-" form:"password"`
	FirstName   *string `json:"first_name" form:"first_name"`
	LastName    *string `json:"last_name" form:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	Image       *string `json:"image"`
}

type Users []User
