package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/dto"
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/internal/repositories"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register a new user
// @Schemes
// @Description Register a new user with email, password, and role
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password (min 8 characters)"
// @Success 201 {object} dto.UserResponse "User successfully registered"
// @Failure 400 {object} dto.ErrorResponse "Invalid input or email already registered"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	response := pkg.NewResponse(c)
	var input dto.RegisterDTO

	if err := c.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if !IsValidEmail(input.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !IsValidPassword(input.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	existingUser, err := repositories.FindUserByEmail(input.Email)
	if err != nil {
		response.InternalServerError("Failed to query user", err.Error())
		return
	}

	if existingUser != nil {
		response.BadRequest(fmt.Sprintf("Email %s is already registered", input.Email), nil)
		return
	}

	hashed := pkg.GenerateHash(input.Password)
	newUser := models.Auth{
		Email:    input.Email,
		Password: hashed,
		RoleId:   1,
	}

	createdUser, err := repositories.RegisterUser(&newUser)
	if err != nil {
		response.InternalServerError("Failed to create user", err.Error())
		return
	}

	response.Created("Success register user", createdUser)
}

// LoginUser godoc
// @Summary Login user
// @Schemes
// @Description Authenticate user and return a JWT token
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 200 {object} dto.LoginResponse "Login successful with token"
// @Failure 400 {object} dto.ErrorResponse "Invalid email or password"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	response := pkg.NewResponse(c)
	var loginDTO dto.LoginUserDTO

	if err := c.ShouldBind(&loginDTO); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	user, _ := repositories.FindUserByEmail(loginDTO.Email)
	if user == nil || !pkg.VerifyHash(loginDTO.Password, user.Password) {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	token, err := pkg.GenerateToken(user.Id, user.RoleId)
	if err != nil {
		response.InternalServerError("Failed to generate token", err.Error())
		return
	}

	response.Success("Login success", map[string]interface{}{
		"token": token,
	})
}

func IsValidEmail(email string) bool {
	if len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}
