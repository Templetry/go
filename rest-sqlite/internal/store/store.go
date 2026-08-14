// Package store owns the database: connection, schema and repositories.
package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // pure-Go driver: no CGO, no build toolchain
)

// Open connects to SQLite and applies the schema. Use ":memory:" in tests.
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", dsn, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to %s: %w", dsn, err)
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("applying schema: %w", err)
	}
	return db, nil
}

// schema is applied on every start; statements must stay idempotent. Pieces
// append their own tables through the migrations socket below.
const schema = `
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS notes (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT NOT NULL,
	body       TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

// migrations is the piece socket: a piece registers its own idempotent DDL
// from its file's init and Migrate applies it after the base schema.
var migrations []string

// Register adds a migration statement. Call it from an init function.
func Register(ddl string) { migrations = append(migrations, ddl) }

// Migrate applies every registered migration.
func Migrate(db *sql.DB) error {
	for _, ddl := range migrations {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}
	return nil
}
