package util

import (
    "fmt"
    "path/filepath"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string, migrationsDir string) error {
    absPath, err := filepath.Abs(migrationsDir)
    if err != nil {
        return fmt.Errorf("failed to get absolute path: %w", err)
    }
    sourceURL := "file://" + absPath

    m, err := migrate.New(sourceURL, dbURL)
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }
    defer m.Close()

    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration failed: %w", err)
    }
    return nil
}