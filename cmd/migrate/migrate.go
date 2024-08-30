package main

import (
	"calometer/internal/db"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Required for PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Required for file-based migrations
)

func main() {
	dbURL, err := db.GetDBUrl()
	if err != nil {
		log.Fatalf("Failed to load DB URL: %v", err)
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	sourceURL := "file://" + migrationsDir

	m, err := migrate.New(
		sourceURL,
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully")
}
