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

## Agent rules

- **Build only via `task`**: Never run `go build`, `templ generate`, or Tailwind directly. Use `task build` (full build) or `task build:ui` (templ + CSS only).
- **Gallery server only via `task`**: Never run `go run ./cmd/gallery` or binary directly.
  - `task gallery` — interactive run (stdout)
  - `task gallery:start` — background daemon
  - `task gallery:stop` — stop daemon
  - `task gallery:restart` — stop + start (preferred after code changes)
  - `task gallery:log` — tail daemon log
- Only one gallery instance on port 11000. Use `task gallery:restart` rather than starting a second.
- After editing `.templ` files, always run `task build:ui` to regenerate `*_templ.go` before committing.

## Conventions

### int-to-string conversion

Use `strconv.Itoa(n)` instead of custom helper functions. For `colspan`/`rowspan` where 0 is invalid, use `strconv.Itoa(max(n, 1))`. Go's built-in `max`/`min` (1.21+) are preferred over manual guards.

### CSS class override pattern

```templ
class={ "base-class", templ.KV(class, class != ""), templ.KV("default-class", class == "") }
```

Use `class string` as last param. When non-empty, replaces defaults; when empty, uses built-in defaults. This keeps backward compat — existing callers pass `""`.

### Conditional HTMX pattern

Only emit HTMX attrs when both `hxGet` and `hxTarget` are non-empty:

```templ
if hxGet != "" && hxTarget != "" {
    hx-get={ hxGet }
    hx-target={ "#" + hxTarget }  // caller omits # prefix
    hx-trigger="change"
    hx-include="closest form"
}
```

Prevents empty `hx-get="" hx-target=""` in rendered DOM.

### Do NOT import `github.com/a-h/templ` in `.templ` files

Templ auto-injects `import "github.com/a-h/templ"` in the generated `*_templ.go`. Adding it in the `.templ` file's import block causes a duplicate import compile error. Use `templ.Component`, `templ.Attributes`, etc. in `.templ` files without importing the package explicitly.

### `form.FormField` as generic input

For inputs with placeholder/rows/options/required, use `form.FormField(form.FormFieldProps{...})` instead of `form.TextInput`/`form.TextareaInput`/`form.SelectInput`. The props struct supports all edge cases (placeholder, SelectOption slice, rows, error, attrs, etc.) that the shorthand functions don't cover.

### Dead code removal

Before removing a component function, verify zero callers:
```
rg "FunctionName" --include '*.go' | grep -v '_templ.go' | grep -v '_test.go'
```
`grep -v '_templ.go'` filters generated-file matches (templ generates forwarder funcs). If only `_templ.go` matches remain, the component is dead.

### Refactoring milestones

| Phase | Status | Description |
|-------|--------|-------------|
| 1 — intStr cleanup | ✅ Done | Replaced all `intStr/sprintCount/promptIntStr/spinnerIntStr/intToStr` with `strconv.Itoa`; removed 6 custom functions |
| 2 — IconSpan composition | ✅ Done | Replaced ~120 manual `<span class="iconify ...">` with `@ui.IconSpan(name, size)` across nav/layout/form; added `aria-hidden` to IconSpan |
| 3 — Badge composition | ✅ Done | Replaced 9 manual `<div class="badge ...">` with `@ui.Badge(props)`; added `BadgeSizeXS` constant |
| 4 — Button composition | ✅ Done | Refactored `ui.Button` to `ButtonProps` struct with `Style`, `Block`, `ExtraClass` fields; replaced ~31 manual `<button class="btn ...">` with `@ui.Button(props)` across nav/layout/form/modal; updated seed.go callers |
| 5 — Card composition | ✅ Done | Made `ui.StatCardMinimal` (non-icon path) use `@CardRaw` instead of raw `<div class="card">` |

No cross-package import issues: all packages that need `ui` already import it or can safely do so (no circular deps). `form/range.templ` was the reference for `strconv.Itoa` inside templ expressions.

- All component functions return `templ.Component`.
- Props are passed as plain Go structs or positional arguments — no global state.
- Tailwind classes are written directly in `.templ` files; DaisyUI component classes (`btn`, `badge`, `card`, etc.) are preferred over raw utility classes.
- File naming: `<component-name>.templ` (kebab-case), package name matches directory name.
- Generated files (`*_templ.go`) are checked in to source control.
