package repository

import (
	"RamadanRangkuti/service-tickitz/internal/models"
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func FindALlUser() (models.Users, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())
	query := `
	SELECT u.id, u.email, u.password, ur.role, p.first_name, p.last_name, 
	p.phone_number, p.image, u.created_at, u.updated_at 
	FROM users u LEFT JOIN user_roles ur ON u.role_id = ur.id 
	LEFT JOIN profiles p ON p.user_id = u.id ORDER BY id ASC;
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	user, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %v", err)
	}
	return user, nil
}

func FindUserById(id int) (*models.UserDetails, error) {
	var user models.UserDetails
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())

	query := `
	SELECT u.id, u.email, u.password, ur.role, p.first_name, p.last_name, 
	p.phone_number, p.image
	FROM users u LEFT JOIN user_roles ur ON u.role_id = ur.id 
	LEFT JOIN profiles p ON p.user_id = u.id WHERE u.id = $1;
	`

	err = conn.QueryRow(context.Background(),
		query, id).
		Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %v", err)
	}

	fmt.Println(user)

	return &user, nil
}

func InsertUser(userDetails *models.UserDetails) (int, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("data diterima ===> : ", userDetails)
	// Mulai transaksi
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	var userId int
	// Update tabel `users`
	userQuery := `
		INSERT INTO users (email, password, role_id)
		VALUES($1, $2, $3) returning id
	`
	err = tx.QueryRow(context.Background(), userQuery, userDetails.Email, userDetails.Password, 1).
		Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into users table: %v", err)
	}

	// Insert ke tabel `profiles`
	profileQuery := `
		INSERT INTO profiles (user_id, first_name, last_name, phone_number, image) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(context.Background(), profileQuery,
		userId, userDetails.FirstName, userDetails.LastName, userDetails.PhoneNumber, userDetails.Image)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into profiles table: %v", err)
	}

	// Commit transaksi
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return userId, nil
}

func EditUser(userId int, userDetails *models.UserDetails) error {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())
	// Mulai transaksi
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()
	fmt.Println("user diterima dari handler :", userDetails)
	// Dereference nilai pointer
	// firstName := ""
	// if userDetails.FirstName != nil {
	// 	firstName = *userDetails.FirstName
	// }
	// lastName := ""
	// if userDetails.LastName != nil {
	// 	lastName = *userDetails.LastName
	// }
	// phoneNumber := ""
	// if userDetails.PhoneNumber != nil {
	// 	phoneNumber = *userDetails.PhoneNumber
	// }
	// image := ""
	// if userDetails.Image != nil {
	// 	image = *userDetails.Image
	// }

	// Update tabel `users`
	userQuery := `
		UPDATE users
		SET email = $1, password = $2, updated_at = current_timestamp
		WHERE id = $3
	`
	_, err = tx.Exec(context.Background(), userQuery, userDetails.Email, userDetails.Password, userId)
	if err != nil {
		return fmt.Errorf("failed to update users table: %v", err)
	}

	// Update tabel `profiles`
	profileQuery := `
		UPDATE profiles
		SET first_name = $1, last_name = $2, phone_number = $3, image = $4, updated_at = current_timestamp
		WHERE user_id = $5
	`
	_, err = tx.Exec(context.Background(), profileQuery, userDetails.FirstName, userDetails.LastName, userDetails.PhoneNumber, userDetails.Image, userId)
	if err != nil {
		return fmt.Errorf("failed to update profiles table: %v", err)
	}

	// Commit transaksi
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func RemoveUser(userId int) error {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
		return err
	}
	defer conn.Close(context.Background())

	// Mulai transaksi
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Hapus data dari tabel `profiles`
	profileQuery := `DELETE FROM profiles WHERE user_id = $1`
	_, err = tx.Exec(context.Background(), profileQuery, userId)
	if err != nil {
		return fmt.Errorf("failed to delete from profiles table: %v", err)
	}

	// Hapus data dari tabel `users`
	userQuery := `DELETE FROM users WHERE id = $1`
	_, err = tx.Exec(context.Background(), userQuery, userId)
	if err != nil {
		return fmt.Errorf("failed to delete from users table: %v", err)
	}

	// Commit transaksi
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
