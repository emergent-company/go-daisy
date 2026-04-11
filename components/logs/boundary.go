package logs

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// LogsTableWithBoundary wraps LogsTable with a dev-mode component boundary annotation.
// gallery:token entries
// gallery:hint entries:slice(4)
func LogsTableWithBoundary(entries []LogEntry) templ.Component {
	return devmode.ComponentBoundary("LogsTable", map[string]any{"entryCount": len(entries)}, LogsTable(entries))
}
