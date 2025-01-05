package models

type Auth struct {
	Id       int `json:"-"`
	Email    string
	Password string `json:"-"`
	RoleId   int    `json:"-"`
}
