package models

type Auth struct {
	Id       int    `json:"-" form:"id"`
	Email    string `json:"email" form:"email"`
	Password string `json:"-" form:"password"`
	RoleId   int    `json:"role_id" form:"role_id"`
}
