package handlers

import (
	"RamadanRangkuti/service-tickitz/internal/handlers"
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/internal/repository"
	"RamadanRangkuti/service-tickitz/pkg"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	user, err := repository.FindALlUser(page, limit, order, sortBy, search)
	if err != nil {
		response.InternalServerError("Failed Get All Users", err.Error())
		return
	}
	response.Success("Success Get All Users", user)
}

func GetUserById(c *gin.Context) {
	response := pkg.NewResponse(c)
	// id, err := strconv.Atoi(c.Param("id"))
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

	user, err := repository.FindUserById(id)
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

func CreateUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	var userInput models.UserDetails
	if err := c.ShouldBind(&userInput); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		response.BadRequest("Invalid input data", err.Error())
		return
	}
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

	if !handlers.IsValidEmail(userInput.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !handlers.IsValidPassword(userInput.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	userEmail, _ := repository.FindUserByEmail(userInput.Email)
	if userEmail != nil {
		response.BadRequest("Email already exist", nil)
		return
	}

	if file != nil {
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		maxSize := int64(5 << 20) // 5 MB
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

	userId, err := repository.InsertUser(&userInput)
	if err != nil {
		response.InternalServerError("Failed to create profile", err.Error())
		return
	}

	userInput.Id = userId
	response.Success("Success created user", userInput)
}

func UpdateUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	// Ambil UserID dari token
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	fmt.Println("ID dari token:", id)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Cari data user dari database
	user, _ := repository.FindUserById(id)
	fmt.Println("EXISTING USER ", user)
	if user == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	// Bind data input dari request
	var userInput models.UserDetails
	if err := c.ShouldBind(&userInput); err != nil {
		response.BadRequest("Invalid input data", err.Error())
		return
	}

	// Validasi Email
	if userInput.Email != "" {
		if !handlers.IsValidEmail(userInput.Email) {
			response.BadRequest("Invalid email format", nil)
			return
		}
		// Cek jika email sudah digunakan user lain
		userEmail, _ := repository.FindUserByEmail(userInput.Email)
		if userEmail != nil && userEmail.Id != id {
			response.BadRequest("Email already exists", nil)
			return
		}
	}

	// Validasi Password
	if userInput.Password != "" {
		if !handlers.IsValidPassword(userInput.Password) {
			response.BadRequest("Password must be at least 8 characters long", nil)
			return
		}
		hashed := pkg.GenerateHash(userInput.Password)
		user.Password = hashed
	}

	if userInput.FirstName != nil {
		user.FirstName = userInput.FirstName
	}
	if userInput.LastName != nil {
		user.LastName = userInput.LastName
	}
	if userInput.PhoneNumber != nil {
		user.PhoneNumber = userInput.PhoneNumber
	}

	// Proses upload gambar jika ada
	file, _ := c.FormFile("image")
	if file != nil {
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		maxSize := int64(5 << 20) // 5MB
		uploadDir := "public/images"

		imagePath, err := pkg.UploadImage(c, file, allowedExts, maxSize, uploadDir)
		if err != nil {
			response.BadRequest("Failed to upload image", err.Error())
			return
		}
		user.Image = &imagePath
	}

	// Debugging: Periksa data userInput sebelum mengupdate
	fmt.Printf("USER YANG DIKIRIM : %+v\n", user)

	// Update data user
	err := repository.EditUser(id, user)
	if err != nil {
		response.InternalServerError("Failed to update user and profile", err.Error())
		return
	}

	response.Success("User updated successfully", user)
}

func DeleteUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest("Invalid user ID", err.Error())
		return
	}

	user, _ := repository.FindUserById(id)
	if user == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	err = repository.RemoveUser(id)
	if err != nil {
		response.InternalServerError("Failed to delete user and profile", err.Error())
		return
	}

	response.Success(fmt.Sprintf("User with ID %d deleted successfully", id), nil)
}
