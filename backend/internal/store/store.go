package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store provides access to the DockPulse SQLite database.
type Store struct {
	db *sql.DB
}

// New opens the SQLite database at dbPath, enables WAL mode, and runs
// schema migrations. Use ":memory:" for an in-memory database in tests.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err = db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if _, err = db.ExecContext(context.Background(), "PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}

	s := &Store{db: db}
	if err = s.migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return s, nil
}

// Close releases the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS image_checks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			container_name TEXT NOT NULL,
			container_id TEXT NOT NULL,
			image_ref TEXT NOT NULL,
			local_digest TEXT NOT NULL DEFAULT '',
			remote_digest TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'unknown',
			checked_at TEXT NOT NULL,
			registry TEXT NOT NULL DEFAULT 'dockerhub'
		)`,
		`CREATE INDEX IF NOT EXISTS idx_image_checks_container_id ON image_checks(container_id)`,
		`CREATE INDEX IF NOT EXISTS idx_image_checks_checked_at ON image_checks(checked_at)`,
		`CREATE TABLE IF NOT EXISTS preferences (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
	}

	for _, ddl := range migrations {
		if _, err := s.db.ExecContext(context.Background(), ddl); err != nil {
			return fmt.Errorf("exec migration: %w", err)
		}
	}

	return nil
}
