package main

import (
	"log"

	"github.com/emergent-company/go-daisy/cmd/gallery/internal/gallery"
	"github.com/emergent-company/go-daisy/galleryruntime"
)

func main() {
	if err := galleryruntime.Serve(galleryruntime.Options{
		Title:      "go-daisy",
		Components: gallery.AllComponents(),
		Port:       11000,
		DevMode:    true,
	}); err != nil {
		log.Fatal(err)
	}
}
