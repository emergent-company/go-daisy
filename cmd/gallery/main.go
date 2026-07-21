package main

import (
	"log"
	"os"
	"strconv"

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
	if err := galleryruntime.Serve(galleryruntime.Options{
		Title:      "go-daisy",
		Components: gallery.AllComponents(),
		Port:       port,
		DevMode:    true,
	}); err != nil {
		log.Fatal(err)
	}
}
