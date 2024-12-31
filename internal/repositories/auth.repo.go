package repositories

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"fmt"
)

func RegisterUser(user *models.Auth) (*models.Auth, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())

	var newUser models.Auth

	fmt.Println("Query Input:", user.Email, user.Password, user.RoleId)
	err = conn.QueryRow(context.Background(), `
		INSERT INTO users(email, password, role_id)
		VALUES ($1, $2, $3)
		RETURNING id, email, password, role_id
		`, user.Email, user.Password, user.RoleId).Scan(&newUser.Id, &newUser.Email, &newUser.Password, &newUser.RoleId)

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &newUser, nil
}

func FindUserByEmail(email string) (*models.Auth, error) {
	var user models.Auth
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), `
		SELECT id, email, password, role_id
		FROM users WHERE email = $1`, email).Scan(&user.Id, &user.Email, &user.Password, &user.RoleId)
	if err != nil {
		if err.Error() == "no rows in result set" { // Periksa jika tidak ada hasil
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return &user, nil
}
