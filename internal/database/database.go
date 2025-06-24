package database

import (
        "database/sql"
        "fmt"
        "os"

        "github.com/golang-migrate/migrate/v4"
        "github.com/golang-migrate/migrate/v4/database/sqlite3"
        _ "github.com/golang-migrate/migrate/v4/source/file"
        _ "github.com/mattn/go-sqlite3"
)

func Initialize(databasePath string) (*sql.DB, error) {
        // Create database directory if it doesn't exist
        if err := os.MkdirAll("./data", 0755); err != nil {
                return nil, fmt.Errorf("failed to create data directory: %w", err)
        }
        
        // Use SQLite database file
        if databasePath == "" {
                databasePath = "./data/finderr.db"
        }
        
        db, err := sql.Open("sqlite3", databasePath+"?_foreign_keys=on")
        if err != nil {
                return nil, fmt.Errorf("failed to open database: %w", err)
        }

        if err := db.Ping(); err != nil {
                return nil, fmt.Errorf("failed to ping database: %w", err)
        }

        return db, nil
}

func RunMigrations(databasePath string) error {
        if databasePath == "" {
                databasePath = "./data/finderr.db"
        }
        
        db, err := sql.Open("sqlite3", databasePath+"?_foreign_keys=on")
        if err != nil {
                return fmt.Errorf("failed to open database for migrations: %w", err)
        }
        defer db.Close()

        driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
        if err != nil {
                return fmt.Errorf("failed to create migration driver: %w", err)
        }

        m, err := migrate.NewWithDatabaseInstance(
                "file://migrations",
                "sqlite3", driver)
        if err != nil {
                return fmt.Errorf("failed to create migration instance: %w", err)
        }

        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
                return fmt.Errorf("failed to run migrations: %w", err)
        }

        return nil
}