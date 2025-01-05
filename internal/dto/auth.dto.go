package dto

type RegisterDTO struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginUserDTO struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserResponseDTO struct {
	Id     int    `json:"id"`
	Email  string `json:"email"`
	RoleId int    `json:"role_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
