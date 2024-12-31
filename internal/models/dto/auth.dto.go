package dto

type RegisterDTO struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
	RoleID   int    `json:"role_id" form:"role_id" binding:"required"`
}

type LoginUserDTO struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserResponseDTO struct {
	Id     int    `json:"id"`
	Email  string `json:"email"`
	RoleId int    `json:"role_id"`
}
