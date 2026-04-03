---
name: gallery-install
description: >-
  Install the go-daisy component gallery into any Go module. Use whenever the
  user says "install the gallery", "set up gallery in this project",
  "add go-daisy gallery", "init gallery", "run gallery install",
  "gallery-install", or "how do I get the gallery running". Covers running
  the installer, what files are created (gallery/main.go, gallery/components.go),
  Taskfile patching, port configuration, skill installation, and next steps.
metadata:
  author: emergent
  version: "1.0"
---

# Install the go-daisy Gallery

The gallery is a portable component preview server powered by
`github.com/emergent-company/go-daisy/galleryruntime`. Installing it scaffolds
a `gallery/` package in your module that you run locally during development.

---

## Prerequisites

- A Go module with a `go.mod` file in the current directory
- (Optional) A `Taskfile.yml` — the installer will patch it automatically

---

## Install

Run from the **root of your Go module**:

```bash
go run github.com/emergent-company/go-daisy/cmd/install@latest [flags]
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-dir` | `gallery` | Directory to create the gallery package in |
| `-port` | `11000` | Port for the gallery HTTP server |
| `-title` | module name | Title shown in the gallery UI |

### Example with custom port and title

```bash
go run github.com/emergent-company/go-daisy/cmd/install@latest \
  -port 11001 \
  -title "My App"
```

---

## What the installer creates

### `gallery/main.go`

The gallery entry point. Calls `galleryruntime.Serve` with your registry:

```go
package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/emergent-company/go-daisy/galleryruntime"
)

func main() {
    dbPath := filepath.Join(os.TempDir(), "gallery.db")

    if err := galleryruntime.Serve(galleryruntime.Options{
        Title:      "My App",
        Components: allComponents(),
        Port:       11001,
        StorePath:  dbPath,
    }); err != nil {
        log.Fatal(err)
    }
}
```

### `gallery/components.go`

The component registry stub — **this is where you add your components**:

```go
package main

import "github.com/emergent-company/go-daisy/galleryruntime"

func allComponents() []galleryruntime.GalleryComponent {
    return []galleryruntime.GalleryComponent{
        // Add components here — see gallery-add-component skill
    }
}
```

### Taskfile.yml patch

If `Taskfile.yml` exists, the installer appends:

```yaml
  gallery:
    desc: Run the component gallery on :11001
    cmds:
      - go run ./gallery
```

### `.opencode/skills/`

The installer also copies the gallery skills into `.opencode/skills/`:
- `gallery-add-component` — how to add components to the registry
- `gallery-install` — this skill

---

## Run the gallery

```bash
task gallery
# or directly:
go run ./gallery
```

Open **http://localhost:\<port\>** in your browser.

---

## Change the port after install

Edit two places:

1. `gallery/main.go` — change `Port: 11001`
2. `Taskfile.yml` — update the `desc` line (cosmetic only)

---

## Updating to a newer version

Re-run the installer with `--dir` pointing to the existing gallery dir:

```bash
go run github.com/emergent-company/go-daisy/cmd/install@latest -dir gallery
```

The installer will abort if the directory is non-empty. To update manually:
- Edit `gallery/main.go` to update the import and options
- Run `go mod tidy` to pull the latest `galleryruntime`

---

## Gallery server options

`galleryruntime.Options` full schema:

```go
type Options struct {
    Title      string          // Gallery title in the UI
    Components []GalleryComponent  // Your component registry
    Port       int             // HTTP port (default 11000)
    StorePath  string          // SQLite path for feedback persistence; empty = no persistence
    GitHubCfg  *GitHubConfig   // nil = GitHub export disabled
}
```

To enable GitHub issue export from the feedback panel:

```go
galleryruntime.Serve(galleryruntime.Options{
    // ...
    GitHubCfg: &galleryruntime.GitHubConfig{
        AppID:          123456,
        InstallationID: 789,
        PrivateKeyPath: "/path/to/key.pem",
        Owner:          "my-org",
        Repo:           "my-repo",
    },
})
```
