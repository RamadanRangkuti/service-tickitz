package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/internal/repository"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	response := pkg.NewResponse(c)
	var newUser models.Auth

	if err := c.ShouldBind(&newUser); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if !IsValidEmail(newUser.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !IsValidPassword(newUser.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	existingUser, err := repository.FindUserByEmail(newUser.Email)
	if err != nil {
		response.InternalServerError("Failed to query user", err.Error())
		return
	}

	if existingUser != nil {
		response.BadRequest(fmt.Sprintf("Email %s is already registered", newUser.Email), nil)
		return
	}

	hashed := pkg.GenerateHash(newUser.Password)
	newUser.Password = hashed

	fmt.Println(newUser)
	createdUser, err := repository.RegisterUser(&newUser)

	if err != nil {
		response.InternalServerError("Failed to create user", err.Error())
		return
	}

	response.Created("Success register user", createdUser)
}

func Login(c *gin.Context) {
	response := pkg.NewResponse(c)
	var loginData models.Auth
	if err := c.ShouldBind(&loginData); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	user, _ := repository.FindUserByEmail(loginData.Email)
	if user == nil || !pkg.VerifyHash(loginData.Password, user.Password) {
		response.BadRequest("Invalid email or password", nil)
		return
	}
	token, err := pkg.GenerateToken(user.Id)
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
