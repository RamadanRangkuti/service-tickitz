package pkg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func DB() (*pgx.Conn, error) {
	godotenv.Load()
	config, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}
	// var connString string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
	conn, err := pgx.Connect(context.Background(), config.ConnString())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Database connected")
	return conn, nil
}
