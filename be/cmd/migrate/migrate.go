package main

import (
	"calometer/internal/db"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

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

		// Run migrations only if there are unapplied migrations
		err = m.Up()
		if err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No new migrations to apply.")
			} else {
				log.Printf("Migration error: %v\n", err)

				// Handle dirty state and retry
				version, err := db.GetLatestMigrationVersion()
				if err != nil {
					log.Fatalf("Failed to get latest migration version: %v", err)
				}

				if dirtyErr := m.Force(*version); dirtyErr != nil {
					log.Fatalf("Failed to force migration state: %v", dirtyErr)
				}

				log.Println("Dirty state fixed. Retrying migration.")
				if err := m.Up(); err != nil && err != migrate.ErrNoChange {
					log.Fatalf("Failed to reapply migrations: %v", err)
				} else {
					log.Println("Migrations applied successfully after fixing dirty state.")
				}
			}
		} else {
			log.Println("Migrations applied successfully.")
		}
	case "down":
		if len(args) < 1 {
			log.Fatalf("Usage: %s down <number_of_steps>\n", os.Args[0])
		}

		stepsStr := args[0]
		steps, _ := strconv.Atoi(stepsStr)

		migrationsDir := "internal/db/migrations"
		sourceURL := "file://" + migrationsDir

		m, err := migrate.New(
			sourceURL,
			dbURL,
		)
		if err != nil {
			log.Fatalf("Failed to create migration instance: %v", err)
		}

		// Rollback by the specified number of steps
		err = m.Steps(-1 * steps)
		if err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}

		log.Println("Rollback action executed successfully!")
	default:
		log.Fatalf("Unknown command: %s\n", command)
	}
}
