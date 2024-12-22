package models

import (
	"RamadanRangkuti/service-tickitz/pkg"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Tes struct {
	Id   int    `json:"id"`
	Name string `json:"name" form:"name"`
}

type ListTes []Tes

func FindALlTes() (ListTes, error) {
	conn, err := pkg.DB()
	if err != nil {
		fmt.Println("connection failed", err)
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), `SELECT id, name from tes`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	tes, err := pgx.CollectRows(rows, pgx.RowToStructByName[Tes])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %v", err)
	}
	return tes, nil
}
