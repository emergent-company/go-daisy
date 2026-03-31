// Package ui provides reusable primitive Templ components (button, badge, avatar, etc.).
package ui

import (
	"io/fs"
	"net/http"
)

// StaticHandlerFS returns an http.Handler that serves static assets from fsys,
// stripping the given URL prefix. Pass staticfs.FS() from the staticfs package.
func StaticHandlerFS(prefix string, fsys fs.FS) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.FS(fsys)))
}
