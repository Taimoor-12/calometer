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
	dbURL, err := GetDBUrl()
	if err != nil {
		return nil, err
	}

	// Create a connection pool
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func GetDBUrl() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	return os.Getenv("DB_URL"), nil
}

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

func GetPool() *pgxpool.Pool {
	return pool
}
