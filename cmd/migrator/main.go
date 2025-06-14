package main // Note: changed from 'migrator' to 'main' for the executable

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/lavish-gambhir/dashbeam/shared/config"
)

const (
	_migrationsPath = "file://shared/database/migrations"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application configuration: %v", err)
	}
	m, err := migrate.New(
		_migrationsPath,
		cfg.Database.Address(),
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer func() {
		_, dbErr := m.Close()
		if dbErr != nil {
			log.Printf("Warning: error closing migrate database connection: %v", dbErr)
		}
	}()

	if len(os.Args) < 2 {
		log.Fatal("Usage: migrator <command> [args]\nCommands: up, down, create")
	}

	command := os.Args[1]

	switch command {
	case "up":
		log.Println("Applying all pending migrations...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			log.Println("Migrations applied successfully!")
		}

	case "down":
		log.Println("Rolling back the last migration...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to roll back migration: %v", err)
		}
		if err == migrate.ErrNoChange {
			log.Println("No migrations to roll back.")
		} else {
			log.Println("Last migration rolled back successfully!")
		}

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrator create <name_of_migration>")
		}
		name := os.Args[2]
		// Create a new migration file with a timestamp prefix
		timestamp := time.Now().Format("20060102150405")                               // YYYYMMDDHHmmss
		upFile := fmt.Sprintf("%s/%s_%s.up.sql", _migrationsPath[7:], timestamp, name) // Remove "file://" prefix
		downFile := fmt.Sprintf("%s/%s_%s.down.sql", _migrationsPath[7:], timestamp, name)

		if err := os.WriteFile(upFile, []byte(""), 0644); err != nil {
			log.Fatalf("Failed to create up migration file: %v", err)
		}
		if err := os.WriteFile(downFile, []byte(""), 0644); err != nil {
			log.Fatalf("Failed to create down migration file: %v", err)
		}
		log.Printf("Created migration files:\n%s\n%s", upFile, downFile)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
