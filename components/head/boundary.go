package head

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// DependenciesWithBoundary returns the Dependencies component wrapped for gallery dev tooling.
func DependenciesWithBoundary(props DepsProps) templ.Component {
	return devmode.ComponentBoundary("Dependencies", Dependencies(props), map[string]any{
		"alpine":   props.Alpine,
		"morph":    props.Morph,
		"stimulus": props.Stimulus,
		"sse":      props.SSE,
		"ws":       props.WS,
	})
}
