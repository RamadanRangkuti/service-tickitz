package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/dto"
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/internal/repositories"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllUser godoc
// @Summary Get all users
// @Schemes
// @Description Get all users with pagination, sorting, and search
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(5)
// @Param sortBy query string false "Field to sort by" default(id)
// @Param order query string false "Order of sorting (asc or desc)" default(asc)
// @Param search query string false "Search query"
// @Success 200 {object} pkg.Response{data=[]models.User,meta=pkg.PageInfo}
// @Failure 500 {object} pkg.Response{error=string}
// @Success 200 {object} models.User
// @Security ApiKeyAuth
// @Router /users [get]
func GetAllUser(c *gin.Context) {
	response := pkg.NewResponse(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	sortBy := c.DefaultQuery("sortBy", "id")
	order := c.DefaultQuery("order", "asc")
	search := c.Query("search")
	if order != "asc" {
		order = "desc"
	}
	user, err := repositories.FindALlUser(page, limit, order, sortBy, search)
	fmt.Println(user)
	if err != nil {
		response.InternalServerError("Failed Get All Users", err.Error())
		return
	}
	count := repositories.CountUser(search)
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	pageInfo := &pkg.PageInfo{
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		TotalPage:   totalPage,
		TotalData:   count,
	}
	if page >= totalPage {
		pageInfo.NextPage = 0
	}
	if page <= 1 {
		pageInfo.PrevPage = 0
	}

	response.GetAllSuccess("Success Get All Users", user, pageInfo)
}

// GetDetailUser godoc
// @Summary Get detail user
// @Schemes
// @Description get detail user
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} dto.UserResponse "Successful response"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func GetUserById(c *gin.Context) {
	response := pkg.NewResponse(c)

	userId, exists := c.Get("UserId")
	fmt.Println(userId)
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)

	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	user, err := repositories.FindUserById(id)
	fmt.Println(user)
	if user == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}
	if err != nil {
		response.InternalServerError("Failed Get User", err.Error())
		return
	}

	response.Success("Success get user", user)
}

// CreateUser godoc
// @Summary Creat new user
// @Schemes
// @Description Create a new user with the provided details
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param userInput formData dto.CreateUserDTO true "Create user"
// @Param image formData file false "Profile picture"
// @Success 201 {object} pkg.Response{data=models.User}
// @Failure 400 {object} pkg.Response{error=string}
// @Failure 500 {object} pkg.Response{error=string}
// @Security ApiKeyAuth
// @Router /users [post]
func CreateUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	var userInput dto.CreateUserDTO
	if err := c.ShouldBind(&userInput); err != nil {
		response.BadRequest("Invalid input data", err.Error())
		return
	}
	fmt.Println("data diinput", userInput)
	file, _ := c.FormFile("image")

	if userInput.Email == "" {
		response.BadRequest("Email required", nil)
		return
	}

	if userInput.Password == "" {
		response.BadRequest("Password required", nil)
		return
	}

	if userInput.FirstName == nil {
		response.BadRequest("Firstname required", nil)
		return
	}

	if !IsValidEmail(userInput.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !IsValidPassword(userInput.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	userEmail, _ := repositories.FindUserByEmail(userInput.Email)
	if userEmail != nil {
		response.BadRequest("Email already exist", nil)
		return
	}

	if file != nil {
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		maxSize := int64(2 << 20) // 2 MB
		uploadDir := "public/images"

		imagePath, err := pkg.UploadImage(c, file, allowedExts, maxSize, uploadDir)
		if err != nil {
			response.BadRequest("Failed to upload image", err.Error())
			return
		}

		userInput.Image = &imagePath
	} else {
		imageDefault := ""
		userInput.Image = &imageDefault
	}

	hashedPassword := pkg.GenerateHash(userInput.Password)
	userInput.Password = hashedPassword

	userId, err := repositories.InsertUser((*models.UserDetails)(&userInput))
	if err != nil {
		response.InternalServerError("Failed to create profile", err.Error())
		return
	}

	userInput.Id = userId
	response.Success("Success created user", userInput)
}

// UpdateUser godoc
// @Summary Update user
// @Schemes
// @Description Update an existing user
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param userInput formData dto.UpdateUserDTO true "Update user"
// @Param image formData file false "Profile picture"
// @Success 200 {object} pkg.Response{data=models.User}
// @Failure 400 {object} pkg.Response{error=string}
// @Failure 404 {object} pkg.Response{error=string}
// @Failure 500 {object} pkg.Response{error=string}
// @Security ApiKeyAuth
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	// Ambil UserID dari token
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Ambil data user yang ada
	user, _ := repositories.FindUserById(id)
	if user == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	// Validasi input
	var userInput dto.UpdateUserDTO
	if err := c.ShouldBind(&userInput); err != nil {
		response.BadRequest("Invalid input data", err.Error())
		return
	}

	// Validasi email
	if userInput.Email != nil {
		if !IsValidEmail(*userInput.Email) {
			response.BadRequest("Invalid email format", nil)
			return
		}
		existingUser, _ := repositories.FindUserByEmail(*userInput.Email)
		if existingUser != nil && existingUser.Id != id {
			response.BadRequest("Email already exists", nil)
			return
		}
		user.Email = *userInput.Email
	}

	// Validasi password
	if userInput.Password != nil {
		if !IsValidPassword(*userInput.Password) {
			response.BadRequest("Password must be at least 8 characters long", nil)
			return
		}
		user.Password = pkg.GenerateHash(*userInput.Password)
	}

	// Update data profil jika ada
	if userInput.FirstName != nil {
		user.FirstName = userInput.FirstName
	}
	if userInput.LastName != nil {
		user.LastName = userInput.LastName
	}
	if userInput.PhoneNumber != nil {
		user.PhoneNumber = userInput.PhoneNumber
	}

	// Proses upload gambar
	file, _ := c.FormFile("image")
	if file != nil {
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		maxSize := int64(2 << 20) // 2MB
		uploadDir := "public/images"

		imagePath, err := pkg.UploadImage(c, file, allowedExts, maxSize, uploadDir)
		if err != nil {
			response.BadRequest("Failed to upload image", err.Error())
			return
		}
		user.Image = &imagePath
	}

	// Debugging: Periksa data sebelum update
	fmt.Printf("Updated User Data: %+v\n", user)

	// Update user dan profil di database
	err := repositories.EditUser(id, user)
	if err != nil {
		response.InternalServerError("Failed to update user and profile", err.Error())
		return
	}

	response.Success("User updated successfully", user)
}

// DeleteUser godoc
// @Summary Delete user
// @Schemes
// @Description Delete the logged-in user
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} pkg.Response{data=string}
// @Failure 404 {object} pkg.Response{error=string}
// @Failure 500 {object} pkg.Response{error=string}
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	user, _ := repositories.FindUserById(id)
	if user == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	err := repositories.RemoveUser(id)
	if err != nil {
		response.InternalServerError("Failed to delete user and profile", err.Error())
		return
	}

	response.Success(fmt.Sprintf("User with ID %d deleted successfully", id), nil)
}
