package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/emergent-company/go-daisy/cmd/gallery/internal/gallery"
	"github.com/emergent-company/go-daisy/galleryruntime"
)

func main() {
	port := 11000
	if p := os.Getenv("GALLERY_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	title := "go-daisy"
	if t := os.Getenv("GALLERY_TITLE"); t != "" {
		title = t
	}

	devMode := parseBool("GALLERY_DEV_MODE", true)

	log.Printf("go-daisy gallery starting: version=%s commit=%s port=%d devMode=%t",
		galleryruntime.Version, galleryruntime.Commit, port, devMode)

	if err := galleryruntime.Serve(galleryruntime.Options{
		Title:      title,
		Components: gallery.AllComponents(),
		Port:       port,
		DevMode:    devMode,
	}); err != nil {
		log.Fatal(err)
	}
}

// parseBool reads a boolean environment variable. It returns def when the
// variable is unset or empty, and falls back to def when the value is not a
// recognised boolean (strconv.ParseBool syntax: 1/t/true/TRUE/... and
// 0/f/false/FALSE/...).
func parseBool(name string, def bool) bool {
	v, ok := os.LookupEnv(name)
	if !ok || strings.TrimSpace(v) == "" {
		return def
	}
	b, err := strconv.ParseBool(strings.TrimSpace(v))
	if err != nil {
		return def
	}
	return b
}
