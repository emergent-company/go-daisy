// Package behaviors provides a Go helper for including the go-daisy behaviour-init
// runtime. The runtime (staticfs/static/js/go-daisy-init.js) exposes
// window.gdOnce / window.gdBehaviour / window.gdRescan so inline component
// scripts can self-stamp their one-time and swap-safe behaviour guards.
//
// Place @behaviors.Tag() in the page <head> before any component script so the
// gdOnce/gdBehaviour/gdRescan globals are defined by the time inline
// component <script> blocks execute.
package behaviors

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/staticfs"
)

// Tag returns the go-daisy init runtime script tag. It must be placed in the
// <head> before any component script that calls gdOnce or gdBehaviour.
func Tag() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, staticfs.Script("/js/go-daisy-init.js"))
		return err
	})
}

// TagDefer returns the init runtime script tag. The runtime must be present
// before component scripts run (it is loaded synchronously), so this is
// identical to Tag() and is provided for API parity with alpine/stimulus.
func TagDefer() templ.Component {
	return Tag()
}
