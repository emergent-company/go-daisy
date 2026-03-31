// Package logs provides a generalized log table Templ component.
package logs

import "time"

// LogEntry is a generalized workflow/event log entry.
// Type maps to the status dot colour (e.g. "success", "error", "info", "warn").
// Message is the human-readable event description.
type LogEntry struct {
	Type      string
	Message   string
	CreatedAt time.Time
}
