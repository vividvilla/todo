package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func DataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "todo")
}

func DBPath() string {
	return filepath.Join(DataDir(), "todo.db")
}

func Init() error {
	dir := DataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite", DBPath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	return migrate()
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		description TEXT DEFAULT '',
		priority    TEXT DEFAULT 'medium',
		status      TEXT DEFAULT 'pending',
		due_at      DATETIME,
		notified    BOOLEAN DEFAULT 0,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS config (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`
	_, err := DB.Exec(schema)
	return err
}
