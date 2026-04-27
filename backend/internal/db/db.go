package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func GetDB(dbPath string) (*sql.DB, error) {
	// 1. Ensure the directory exists if a path was provided
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// 2. Prepare Connection String (DSN)
	// Includes WAL mode, busy timeout, and foreign key enforcement
	dsn := fmt.Sprintf("%s?_journal=WAL&_busy_timeout=5000&_pragma=foreign_keys(1)", dbPath)

	// 3. Open the connection
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// 4. Configure pool settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 5. Verify connectivity
	if err := db.Ping(); err != nil {
		db.Close() // Clean up the connection if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func InitSchema(db *sql.DB) error {
	// We can combine multiple CREATE TABLE statements if we use a single string,
	// but executing them separately or in a transaction is clearer for error tracking.
	schema := []string{
		`CREATE TABLE IF NOT EXISTS deployments (
			id TEXT PRIMARY KEY,
			source_type TEXT NOT NULL,
			source TEXT NOT NULL,
			status TEXT NOT NULL,
			image_tag TEXT,
			deploy_url TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			deployment_id TEXT NOT NULL,
			line TEXT NOT NULL,
			sequence INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE
		);`,
	}

	for _, statement := range schema {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return nil
}
