package shell

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

func NexusPageWithBoundary(title string, content templ.Component) templ.Component {
	return devmode.ComponentBoundary("NexusPage", NexusShell(content, ""), map[string]any{"title": title})
}
