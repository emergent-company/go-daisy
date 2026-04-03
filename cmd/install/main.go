// Command install sets up a gallery in the current Go module.
//
// Usage:
//
//	go run github.com/emergent-company/go-daisy/cmd/install@latest [flags]
//
// Flags:
//
//	-dir  string  Output directory for the gallery binary (default: "gallery")
//	-port int     Port for the gallery server (default: 11000)
//	-title string Gallery title (default: module name)
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("dir", "gallery", "directory to create the gallery binary in")
	port := flag.Int("port", 11000, "port for the gallery server")
	title := flag.String("title", "", "gallery title (defaults to module name)")
	flag.Parse()

	// Detect module name from go.mod in the current directory.
	modName := detectModuleName()
	if *title == "" {
		// Use the last path segment as the title.
		parts := strings.Split(modName, "/")
		*title = parts[len(parts)-1]
	}

	// Ensure the output directory doesn't already exist (or is empty).
	if info, err := os.Stat(*dir); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(*dir)
		if len(entries) > 0 {
			fatalf("directory %q already exists and is not empty — aborting\n", *dir)
		}
	}

	if err := os.MkdirAll(*dir, 0o755); err != nil {
		fatalf("failed to create directory %q: %v\n", *dir, err)
	}

	// Write gallery/main.go.
	mainPath := filepath.Join(*dir, "main.go")
	mainContent := generateMain(modName, *title, *port)
	if err := os.WriteFile(mainPath, []byte(mainContent), 0o644); err != nil {
		fatalf("failed to write %s: %v\n", mainPath, err)
	}
	fmt.Printf("✓ created %s\n", mainPath)

	// Write gallery/components.go (empty stub registry).
	compPath := filepath.Join(*dir, "components.go")
	compContent := generateComponents()
	if err := os.WriteFile(compPath, []byte(compContent), 0o644); err != nil {
		fatalf("failed to write %s: %v\n", compPath, err)
	}
	fmt.Printf("✓ created %s\n", compPath)

	// Patch Taskfile.yml if it exists.
	patchTaskfile(*dir, *port)

	// Run go mod tidy to pull in galleryruntime.
	fmt.Println("\nRunning go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("warning: go mod tidy failed: %v\n", err)
	}

	fmt.Printf(`
✓ Gallery installed!

Next steps:
  1. Add your components to %s/components.go
  2. Run the gallery:

     go run ./%s

  Or if you patched Taskfile.yml:

     task gallery

`, *dir, *dir)
}

// detectModuleName reads go.mod and returns the module path, or "myapp" as fallback.
func detectModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "myapp"
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "myapp"
}

// generateMain returns the content of gallery/main.go.
func generateMain(modName, title string, port int) string {
	return fmt.Sprintf(`package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/emergent-company/go-daisy/galleryruntime"
)

func main() {
	dbPath := filepath.Join(os.TempDir(), "gallery.db")

	if err := galleryruntime.Serve(galleryruntime.Options{
		Title:      %q,
		Components: allComponents(),
		Port:       %d,
		StorePath:  dbPath,
	}); err != nil {
		log.Fatal(err)
	}
}
`, title, port)
}

// generateComponents returns the content of gallery/components.go.
func generateComponents() string {
	return `package main

import "github.com/emergent-company/go-daisy/galleryruntime"

// allComponents returns the component registry for this gallery.
// Add your own components here.
func allComponents() []galleryruntime.GalleryComponent {
	return []galleryruntime.GalleryComponent{
		// Example HTML snippet component:
		// {
		// 	Slug:        "my-button",
		// 	Name:        "My Button",
		// 	Category:    galleryruntime.CategoryBasics,
		// 	Description: "A custom button component.",
		// 	HTML:        ` + "`" + `<button class="btn btn-primary">Click me</button>` + "`" + `,
		// },
	}
}
`
}

// patchTaskfile appends gallery tasks to Taskfile.yml if it exists and doesn't
// already contain a "gallery" task.
func patchTaskfile(dir string, port int) {
	const taskfileName = "Taskfile.yml"
	data, err := os.ReadFile(taskfileName)
	if err != nil {
		// No Taskfile — skip silently.
		return
	}

	content := string(data)
	if strings.Contains(content, "task: gallery") || strings.Contains(content, "\n  gallery:") || strings.Contains(content, "\ngallery:") {
		fmt.Println("✓ Taskfile.yml already has a gallery task — skipping patch")
		return
	}

	patch := fmt.Sprintf(`
  gallery:
    desc: Run the component gallery on :%d
    cmds:
      - go run ./%s
`, port, dir)

	// Find the "tasks:" section and append after it, or just append to end.
	if idx := strings.Index(content, "\ntasks:"); idx >= 0 {
		// Insert after "tasks:\n"
		insertAt := idx + len("\ntasks:\n")
		content = content[:insertAt] + patch + content[insertAt:]
	} else {
		content += "\ntasks:\n" + patch
	}

	if err := os.WriteFile(taskfileName, []byte(content), 0o644); err != nil {
		fmt.Printf("warning: could not patch Taskfile.yml: %v\n", err)
		return
	}
	fmt.Printf("✓ patched Taskfile.yml with gallery task (port %d)\n", port)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "install: "+format, args...)
	os.Exit(1)
}
