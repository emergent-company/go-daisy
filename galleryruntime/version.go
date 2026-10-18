package galleryruntime

// Version is the gallery build version, reported by the /version endpoint.
//
// It is injected at build time via:
//
//	-ldflags "-X github.com/emergent-company/go-daisy/galleryruntime.Version=$VERSION"
//
// When no value is injected it defaults to "dev".
var Version = "dev"

// Commit is the git commit SHA the gallery was built from, reported by the
// /version endpoint alongside Version.
//
// It is injected at build time via:
//
//	-ldflags "-X github.com/emergent-company/go-daisy/galleryruntime.Commit=$COMMIT"
//
// When no value is injected it defaults to the empty string.
var Commit = ""
