package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Package-level variable to hold the connection pool
var pool *pgxpool.Pool

// init function to initialize the connection pool
func Init() (*pgxpool.Pool, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// Retrieve the database details from environment variables
	dbURL := os.Getenv("DB_URL")

	// Create a connection pool
	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

func GetPool() *pgxpool.Pool {
	return pool
}
