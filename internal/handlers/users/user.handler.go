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
	user, err := repository.FindALlUser()
	if err != nil {
		response.InternalServerError("Failed Get All Users", err.Error())
		return
	}
	response.Success("Success Get All Users", user)
}

func GetUserById(c *gin.Context) {
	response := pkg.NewResponse(c)
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}
	user, err := repository.FindUserById(id)
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

	file, err := c.FormFile("image")

	if err != nil {
		response.BadRequest("Image file is required", err.Error())
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

	// Step 3: Simpan file gambar
	allowedExts := []string{".jpg", ".jpeg", ".png"}
	maxSize := int64(5 << 20) // 500 KB
	uploadDir := "public/images"

	imagePath, err := pkg.UploadImage(c, file, allowedExts, maxSize, uploadDir)
	if err != nil {
		response.BadRequest("Failed to upload image", err.Error())
		return
	}

	// Step 4: Assign path file ke struct
	userInput.Image = &imagePath

	// Hash password sebelum menyimpan ke database
	hashedPassword := pkg.GenerateHash(userInput.Password)
	userInput.Password = hashedPassword

	// Panggil fungsi repository untuk menyimpan data
	userId, err := repository.InsertUser(&userInput)
	if err != nil {
		response.InternalServerError("Failed to create user and profile", err.Error())
		return
	}

	// Tambahkan ID yang baru saja dibuat ke dalam respons
	userInput.Id = userId
	response.Success("User and profile created successfully", userInput)
}

// func CreateUser(c *gin.Context) {
// 	response := pkg.NewResponse(c)

// 	var userInput models.UserDetails
// 	if err := c.ShouldBind(&userInput); err != nil {
// 		fmt.Printf("Failed to bind JSON: %v\n", err)
// 		response.BadRequest("Invalid input data", err.Error())
// 		return
// 	}

// 	if !handlers.IsValidEmail(userInput.Email) {
// 		response.BadRequest("Invalid email format", nil)
// 		return
// 	}

// 	if !handlers.IsValidPassword(userInput.Password) {
// 		response.BadRequest("Password must be at least 8 characters long", nil)
// 		return
// 	}

// 	userEmail, _ := repository.FindUserByEmail(userInput.Email)
// 	if userEmail != nil {
// 		response.BadRequest("Email already exist", nil)
// 		return
// 	}

// 	// Hash password sebelum menyimpan ke database
// 	hashedPassword := pkg.GenerateHash(userInput.Password)
// 	userInput.Password = hashedPassword

// 	// Panggil fungsi repository untuk menyimpan data
// 	userId, err := repository.InsertUser(&userInput)
// 	if err != nil {
// 		response.InternalServerError("Failed to create user and profile", err.Error())
// 		return
// 	}

// 	// Tambahkan ID yang baru saja dibuat ke dalam respons
// 	userInput.Id = userId
// 	response.Success("User and profile created successfully", userInput)
// }

func UpdateUser(c *gin.Context) {
	response := pkg.NewResponse(c)

	// Parse user ID dari parameter URL
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

	// Bind data dari request body
	var userInput models.UserDetails
	if err := c.ShouldBind(&userInput); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		response.BadRequest("Invalid input data", err.Error())
		return
	}

	if !handlers.IsValidEmail(userInput.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	userEmail, _ := repository.FindUserByEmail(userInput.Email)
	if userEmail != nil {
		response.BadRequest("Email already exist", nil)
		return
	}

	if !handlers.IsValidPassword(userInput.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	fmt.Printf("Request Body: %+v\n", userInput)

	if userInput.Email != "" {
		user.Email = userInput.Email
	}
	if userInput.Password != "" {
		hashded := pkg.GenerateHash(userInput.Password)
		user.Password = hashded
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
	if userInput.Image != nil {
		user.Image = userInput.Image
	}

	// Panggil fungsi repository untuk update data
	err = repository.EditUser(id, user)
	if err != nil {
		response.InternalServerError("Failed to update user and profile", err.Error())
		return
	}

	response.Success("Success updated user", user)
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
