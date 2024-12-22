package pkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func DB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")

	var connString string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Database connected")
	return conn, nil
}
