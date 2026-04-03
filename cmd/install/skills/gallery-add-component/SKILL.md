---
name: gallery-add-component
description: >-
  Add a component to the project's go-daisy gallery. Use whenever the user says
  "add a component to the gallery", "register this component in the gallery",
  "show this in the gallery", "add a gallery entry for X", "add a story for X",
  "add variants for X", "preview this component in the gallery", or
  "gallery-add-component". Covers editing gallery/components.go, GalleryComponent
  fields, Category constants, HTML snippets, Templ components, named Variants/stories,
  DesignToken controls, and running the gallery to verify.
metadata:
  author: emergent
  version: "1.0"
---

# Add a Component to the Gallery

The gallery registry lives in **`gallery/components.go`** at the root of the project.
All components are returned by `allComponents()` as a `[]galleryruntime.GalleryComponent` slice.

---

## GalleryComponent schema

```go
galleryruntime.GalleryComponent{
    // Required
    Slug        string              // URL-safe unique identifier, e.g. "my-button"
    Name        string              // Display name, e.g. "My Button"
    Category    galleryruntime.Category  // see Category constants below
    Description string              // One-line description shown in the sidebar

    // Pick ONE of these for a single-variant component:
    HTML        string              // Raw HTML snippet (simplest — no Go imports needed)
    Templ       templ.Component     // Pre-instantiated templ.Component

    // Optional
    Subcategory string              // Groups components under a heading within a category
    FrameHeight string              // iframe height, e.g. "300px", "100vh" (default: "400px")
    Tokens      []galleryruntime.DesignToken  // Design token controls (see below)

    // Multiple named variants (Storybook-style):
    Variants    []galleryruntime.GalleryStory
}
```

### Category constants

```go
galleryruntime.CategoryFoundation   // "Foundation"  — colors, typography, spacing tokens
galleryruntime.CategoryBasics       // "Basics"       — buttons, badges, avatars
galleryruntime.CategoryDataDisplay  // "Data Display" — tables, lists, stat cards
galleryruntime.CategoryFeedback     // "Feedback"     — alerts, toasts, loaders
galleryruntime.CategoryOverlays     // "Overlays"     — modals, drawers, dropdowns
galleryruntime.CategoryNavigation   // "Navigation"   — tabs, breadcrumbs, sidebar
galleryruntime.CategoryForms        // "Forms"        — inputs, selects, checkboxes
galleryruntime.CategoryLayout       // "Layout"       — page shell, grid, cards
```

---

## Workflow

### 1. Open gallery/components.go

The file is at `gallery/components.go` in the project root. The `allComponents()` function
returns the full registry. Add your entry to the slice.

### 2. HTML snippet (simplest)

Use `HTML` for static DaisyUI/Tailwind markup — no imports, no build step required:

```go
{
    Slug:        "status-badge",
    Name:        "Status Badge",
    Category:    galleryruntime.CategoryBasics,
    Subcategory: "Badges",
    Description: "Pill badge showing live/pending/archived states.",
    HTML: `<div class="flex gap-2 p-6 justify-center">
  <span class="badge badge-success">Live</span>
  <span class="badge badge-warning">Pending</span>
  <span class="badge badge-ghost">Archived</span>
</div>`,
},
```

### 3. Templ component

Use `Templ` when the component is a Go `templ.Component`. Import the package at the top of
`gallery/components.go`:

```go
import (
    "github.com/emergent-company/go-daisy/galleryruntime"
    "github.com/your-module/internal/ui/components/mycomp"
)

// inside allComponents():
{
    Slug:        "case-card",
    Name:        "Case Card",
    Category:    galleryruntime.CategoryDataDisplay,
    Description: "Card summarising a legal case.",
    Templ:       mycomp.CaseCard(mycomp.CaseCardProps{Title: "Smith v Jones", Status: "Open"}),
    FrameHeight: "220px",
},
```

### 4. Named variants (Storybook-style)

Use `Variants` to show multiple states side-by-side with tabs:

```go
{
    Slug:        "alert",
    Name:        "Alert",
    Category:    galleryruntime.CategoryFeedback,
    Description: "Success, warning, and error alert styles.",
    Variants: []galleryruntime.GalleryStory{
        {
            Name:        "Success",
            Description: "Operation completed successfully.",
            HTML:        `<div class="alert alert-success m-6">Saved!</div>`,
        },
        {
            Name:        "Warning",
            HTML:        `<div class="alert alert-warning m-6">Check your input.</div>`,
        },
        {
            Name:        "Error",
            HTML:        `<div class="alert alert-error m-6">Something went wrong.</div>`,
        },
    },
},
```

### 5. Design token controls (optional)

Add interactive CSS variable sliders/pickers to the token panel:

```go
Tokens: []galleryruntime.DesignToken{
    {
        CSSVar:   "--rounded-btn",
        Selector: ":root",
        Label:    "Border radius",
        Group:    "Shape",
        Type:     galleryruntime.TokenTypeRange,
        Default:  "0.5rem",
        Unit:     "rem",
        Min:      0, Max: 2, Step: 0.125,
    },
    {
        CSSVar:   "--btn-color",
        Selector: ".btn-primary",
        Label:    "Button color",
        Group:    "Color",
        Type:     galleryruntime.TokenTypeColor,
        Default:  "#6d62d4",
    },
},
```

### 6. Run the gallery

```bash
task gallery
# or: go run ./gallery
```

Open http://localhost:<port> and find your component in the sidebar. The `Slug` becomes the
URL: `/gallery/<slug>`.

---

## Rules

- `Slug` must be **unique** across all components — use kebab-case.
- Do **not** set both `HTML` and `Templ` on the same entry (or the same story) — pick one.
- `Variants` and root-level `HTML`/`Templ` are mutually exclusive: if `Variants` is non-empty,
  the root fields are ignored.
- `FrameHeight` defaults to `"400px"` — set it explicitly for tall or short components.
- Keep `Description` to one sentence; it appears in the sidebar tooltip and detail header.
