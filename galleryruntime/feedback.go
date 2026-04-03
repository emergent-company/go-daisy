package galleryruntime

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// FeedbackStatus represents the resolution state of a feedback item.
type FeedbackStatus string

const (
	StatusOpen     FeedbackStatus = "open"
	StatusResolved FeedbackStatus = "resolved"
	StatusWontFix  FeedbackStatus = "wont_fix"
)

// Feedback is a single piece of element-level feedback left on a gallery component.
type Feedback struct {
	ID            int64          `json:"id"`
	ComponentSlug string         `json:"component_slug"`
	Comment       string         `json:"comment"`
	ContextJSON   string         `json:"context_json"`
	Status        FeedbackStatus `json:"status"`
	AgentResponse string         `json:"agent_response"`
	CreatedAt     time.Time      `json:"created_at"`
	ResolvedAt    *time.Time     `json:"resolved_at,omitempty"`
}

// CreateParams are the inputs for creating a new feedback item.
type CreateParams struct {
	ComponentSlug string
	Comment       string
	ContextJSON   string
}

// Create inserts a new feedback item and returns the created row.
func (s *Store) Create(ctx context.Context, p CreateParams) (Feedback, error) {
	if p.ContextJSON == "" {
		p.ContextJSON = "{}"
	}
	const q = `
		INSERT INTO gallery_feedback (component_slug, comment, context_json)
		VALUES (?, ?, ?)
		RETURNING id, component_slug, comment, context_json, status, agent_response, created_at, resolved_at`
	return scanFeedback(s.db.QueryRowContext(ctx, q, p.ComponentSlug, p.Comment, p.ContextJSON))
}

// List returns all feedback for a component slug, newest first.
func (s *Store) List(ctx context.Context, slug string) ([]Feedback, error) {
	const q = `
		SELECT id, component_slug, comment, context_json, status, agent_response, created_at, resolved_at
		FROM gallery_feedback
		WHERE component_slug = ?
		ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, q, slug)
	if err != nil {
		return nil, fmt.Errorf("gallery: list feedback: %w", err)
	}
	defer rows.Close()
	var items []Feedback
	for rows.Next() {
		item, err := scanFeedback(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListOpen returns all open feedback for a component slug, newest first.
func (s *Store) ListOpen(ctx context.Context, slug string) ([]Feedback, error) {
	const q = `
		SELECT id, component_slug, comment, context_json, status, agent_response, created_at, resolved_at
		FROM gallery_feedback
		WHERE component_slug = ? AND status = 'open'
		ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, q, slug)
	if err != nil {
		return nil, fmt.Errorf("gallery: list open feedback: %w", err)
	}
	defer rows.Close()
	var items []Feedback
	for rows.Next() {
		item, err := scanFeedback(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Count returns the number of feedback items for a component slug.
func (s *Store) Count(ctx context.Context, slug string) (int64, error) {
	var n int64
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM gallery_feedback WHERE component_slug = ?`, slug,
	).Scan(&n)
	return n, err
}

// Delete removes a feedback item by ID. Returns an error if not found.
func (s *Store) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM gallery_feedback WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("gallery: delete feedback %d: %w", id, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("gallery: feedback %d not found", id)
	}
	return nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanFeedback(row scanner) (Feedback, error) {
	var f Feedback
	var resolvedAt sql.NullString
	var createdAtStr string
	err := row.Scan(
		&f.ID,
		&f.ComponentSlug,
		&f.Comment,
		&f.ContextJSON,
		&f.Status,
		&f.AgentResponse,
		&createdAtStr,
		&resolvedAt,
	)
	if err != nil {
		return Feedback{}, fmt.Errorf("gallery: scan feedback: %w", err)
	}
	// Parse created_at from SQLite text format.
	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		f.CreatedAt = t
	}
	if resolvedAt.Valid {
		if t, err := time.Parse(time.RFC3339, resolvedAt.String); err == nil {
			f.ResolvedAt = &t
		}
	}
	return f, nil
}
