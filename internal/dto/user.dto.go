package dto

type CreateUserDTO struct {
	Id          int     `json:"id" form:"id"`
	Email       string  `json:"email" form:"email"`
	Password    string  `json:"password" form:"password"`
	Role        string  `json:"role" form:"role"`
	FirstName   *string `json:"first_name" form:"first_name"`
	LastName    *string `json:"last_name" form:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	Image       *string `json:"image"`
}
