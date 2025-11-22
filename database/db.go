package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dsn := "postgres://postgres:todo_pass@localhost:5432/todo?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	DB = pool
	fmt.Println("connected")

	createTable()
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS Tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		is_completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		completed_at TIMESTAMP
	);
	`

	if _, err := DB.Exec(context.Background(), query); err != nil {
		log.Fatal("failed to create table:", err)
	}

	fmt.Println("table 'todos' is ready")
}
