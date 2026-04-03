package galleryruntime

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // register "sqlite" driver
)

const schema = `
CREATE TABLE IF NOT EXISTS gallery_feedback (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  component_slug TEXT    NOT NULL,
  comment        TEXT    NOT NULL,
  context_json   TEXT    NOT NULL DEFAULT '{}',
  status         TEXT    NOT NULL DEFAULT 'open',
  agent_response TEXT    NOT NULL DEFAULT '',
  created_at     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),
  resolved_at    TEXT
);

CREATE INDEX IF NOT EXISTS gallery_feedback_slug_idx
  ON gallery_feedback (component_slug);
`

// Store wraps a SQLite database connection for gallery feedback operations.
type Store struct {
	db *sql.DB
}

// Open opens (or creates) the SQLite database at the given path and applies the
// schema. The caller is responsible for calling Close when done.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("gallery: open sqlite %s: %w", path, err)
	}
	// Single writer — no WAL needed for this dev tool.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("gallery: apply schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the underlying SQLite connection.
func (s *Store) Close() error { return s.db.Close() }
