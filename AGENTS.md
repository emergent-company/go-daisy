# go-daisy — Agent Guide

## What this project is

`go-daisy` is a Go UI component library for building HTMX-driven web interfaces. It provides type-safe, reusable Templ components styled with DaisyUI (Tailwind CSS). The project also ships a live **gallery** app at `cmd/gallery` that showcases every component.

Module path: `github.com/emergent-company/go-daisy`

---

## Repository layout

```
go-daisy/
├── cmd/gallery/          # go-daisy's own gallery showcase app (Echo web server on :11000)
│   └── internal/gallery/ # seed.go — component registry for the showcase
├── cmd/install/          # Installer: `go run github.com/emergent-company/go-daisy/cmd/install@latest`
├── galleryruntime/       # Reusable gallery package (importable by any project)
│   ├── types.go          # GalleryComponent, GalleryStory, DesignToken, Category constants
│   ├── helpers.go        # ComponentBySlug, SlugifyStoryName, BuildCategoryGroups, TokenGroups
│   ├── store.go          # SQLite feedback persistence (Open, Create, List, Delete, Count)
│   ├── feedback.go       # Feedback types and CRUD helpers
│   ├── github.go         # GitHub App client for issue export
│   ├── serve.go          # Serve(Options) — starts the Echo gallery server
│   ├── handler.go        # HTTP route handlers (package galleryruntime)
│   ├── pages_shell.templ # Gallery shell + sidebar + search modal
│   ├── pages_detail.templ# Component detail page (preview iframe, tokens, feedback)
│   └── pages_index.templ # Gallery landing/index page
├── components/           # The component library
│   ├── form/             # Form inputs and field wrappers
│   ├── layout/           # Page shell, sidebar, navbar
│   ├── logs/             # Log display components
│   ├── modal/            # Modal dialogs
│   ├── nav/              # Page headers, tab menus, top bar
│   ├── table/            # Data tables, list areas, scroll rows
│   └── ui/               # Primitives: button, badge, card, avatar, toast, alert, pagination, etc.
├── render/               # HTMX-aware rendering helpers
├── assets/               # Tailwind CSS source (app.css)
├── staticfs/             # Embedded static assets served via Go embed
├── Taskfile.yml          # Build tasks
├── go.mod                # Module definition
├── package.json          # Node deps: DaisyUI, Tailwind CLI
└── tailwind.config.js    # Tailwind + DaisyUI config
```

---

## Key packages

### `render`
HTMX-aware rendering helpers. Use these in every HTTP handler instead of calling `templ.Component.Render` directly.

| Function | When to use |
|---|---|
| `RenderPage` | Always render the full HTML shell |
| `RenderPartial` | Always render a content fragment only |
| `RenderAuto` | Choose page vs. partial based on HTMX headers |
| `RenderTriple` | Three-tier: full shell / sidebar nav swap / tab swap |
| `RedirectAfterMutation` | HX-Redirect for HTMX, 303 for plain requests |
| `AppendToast` | Write an `hx-swap-oob` toast fragment into the response |

### `components/ui`
Primitive DaisyUI components. All are `templ.Component` values returned by Go functions. Notable ones:

- `Button`, `Badge`, `Card`, `Avatar`
- `Toast`, `Alert`, `Loader`, `EmptyState`
- `Pagination`, `Filter`, `ActionMenu`, `StatCard`

### `components/layout`
Full-page shell, sidebar layout, and navbar.

### `components/nav`
`PageHeader`, `TabMenu`, `SimpleTabs`, `TopBar`.

### `components/form`
Form field wrappers and input primitives.

### `components/table`
`Table`, `ListArea` (infinite-scroll container), `ScrollRows`.

### `components/modal`
Modal dialog components.

### `components/logs`
Log stream display.

---

## Tech stack

| Layer | Tool |
|---|---|
| HTTP framework | Echo v4 (`github.com/labstack/echo/v4`) |
| Templating | Templ (`github.com/a-h/templ v0.3.1001`) |
| CSS | DaisyUI + Tailwind CSS (via Node CLI) |
| Interactivity | HTMX |
| Static assets | Go `embed` package (`staticfs/`) |
| Build orchestration | go-task (`Taskfile.yml`) |

---

## Build & run

```bash
# Generate Templ files + compile Tailwind, then build Go binaries
task build

# Generate Templ + CSS only (no go build)
task build:ui

# Run the gallery at http://localhost:11000
task gallery

# Watch CSS (separate terminal during development)
task dev:ui
```

> **Important:** Always run `task build:ui` (or `task build`) after editing any `.templ` file or CSS. The generated `*_templ.go` files must be committed alongside their `.templ` sources.

---

## Adding a new component

1. Create `components/<package>/<name>.templ` with your Templ component(s).
2. Expose a thin Go wrapper in `components/<package>/<package>.go` if props structs or helper types are needed.
3. Run `task build:ui` — this generates `<name>_templ.go` and recompiles CSS.
4. Add a gallery page under `cmd/gallery/internal/handler/` to showcase the component.
5. Register the new route in `cmd/gallery/internal/handler/handler.go`.

---

## HTMX rendering conventions

- Detect HTMX context with helpers in `render/render.go` (`IsHTMX`, `IsPartial`, `IsMainContentTarget`).
- Use `RenderTriple` for pages that live inside the sidebar layout (most pages).
- Use `RenderAuto` for simpler two-tier pages.
- Never set `HX-Redirect` manually — use `RedirectAfterMutation`.
- Toast notifications are appended via `AppendToast` (out-of-band swap into `#toast-container`).

---

## Static assets

CSS is compiled from `assets/app.css` → `staticfs/static/css/app.css` and served at `/static/css/app.css`. The `staticfs` package embeds this file into the binary. Do not edit the output file directly; edit `assets/app.css` and re-run `task build:ui`.

---

## Conventions

- All component functions return `templ.Component`.
- Props are passed as plain Go structs or positional arguments — no global state.
- Tailwind classes are written directly in `.templ` files; DaisyUI component classes (`btn`, `badge`, `card`, etc.) are preferred over raw utility classes.
- File naming: `<component-name>.templ` (kebab-case), package name matches directory name.
- Generated files (`*_templ.go`) are checked in to source control.
