package main

import (
	"calometer/internal/db"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [args...]\n", os.Args[0])
	}

	command := os.Args[1]
	args := os.Args[2:]

	dbURL, err := db.GetDBUrl()
	if err != nil {
		log.Fatalf("Failed to load DB URL: %v", err)
	}

	switch command {
	case "create":
		var cmd *exec.Cmd

		if len(args) < 1 {
			log.Fatalf("Usage: %s create <migration_name>\n", os.Args[0])
		}
		migrationName := args[0]
		cmd = exec.Command("migrate", "create", "-ext", "sql", "-dir", "internal/db/migrations", "-seq", "-digits", "1", migrationName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to execute command: %v\n", err)
		}

		fmt.Println("Command executed successfully.")
	case "up":
		migrationsDir := "internal/db/migrations" // Adjust this path if needed
		sourceURL := "file://" + migrationsDir

		m, err := migrate.New(
			sourceURL,
			dbURL,
		)
		if err != nil {
			log.Fatalf("Failed to create migration instance: %v", err)
		}

		// Attempt to fix dirty state by forcefully marking migration as complete
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Printf("Migration error: %v\n", err)
			// Handling dirty state
			if err := m.Force(-1); err != nil {
				log.Fatalf("Failed to force migration state: %v", err)
			}
			log.Println("Dirty state fixed. Retrying migration.")
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to reapply migrations: %v", err)
			}
			log.Println("Migrations applied successfully")
		} else {
			log.Println("Migrations applied successfully")
		}
	default:
		log.Fatalf("Unknown command: %s\n", command)
	}
}