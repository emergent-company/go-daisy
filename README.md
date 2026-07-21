# go-daisy

A Go UI component library for building HTMX-driven web interfaces. Provides type-safe, reusable [Templ](https://templ.guide) components styled with [DaisyUI](https://daisyui.com) (Tailwind CSS).

Module path: `github.com/emergent-company/go-daisy`

---

## Installation

```bash
go get github.com/emergent-company/go-daisy
```

---

## Components

Components live under `components/` and are organized by package:

| Package | Contents |
|---|---|
| `components/ui` | Primitives: Button, Badge, Card, Avatar, Toast, Alert, Pagination, Chart, ProgressBar, Frame, Dashboard, Aura, etc. |
| `components/form` | Form field wrappers, inputs, selects, wizards, OTP, ColorInput, DatalistInput |
| `components/layout` | Full-page shell, sidebar, navbar, LayoutCustomizer |
| `components/nav` | PageHeader, TabMenu, TopBar, Megamenu, breadcrumbs, menus |
| `components/table` | Table, ListArea (infinite-scroll), ScrollRows |
| `components/modal` | Modal dialogs (CSS, Alpine.js, Stimulus.js variants) |
| `components/logs` | Log stream display |
| `components/alpine` | Alpine.js attribute helpers + pre-built data states |
| `components/stimulus` | Stimulus.js attribute helpers + 6 pre-built controllers |

All component functions return `templ.Component`. Props are passed as plain Go structs or positional arguments — no global state.

---

## HTMX rendering

Use helpers from the `render` package in every HTTP handler instead of calling `templ.Component.Render` directly:

| Function | When to use |
|---|---|
| `render.RenderPage` | Always render the full HTML shell |
| `render.RenderPartial` | Always render a content fragment only |
| `render.RenderAuto` | Choose page vs. partial based on HTMX headers |
| `render.RenderTriple` | Full shell / sidebar nav swap / tab swap |
| `render.RedirectAfterMutation` | HX-Redirect for HTMX, 303 for plain requests |
| `render.AppendToast` | Out-of-band toast fragment into `#toast-container` |
| `render.SetMorph` | Enable DOM morphing for this swap (HTMX + idiomorph) |
| `render.RenderAutoMorph` | RenderAuto with DOM morphing |
| `render.Preserve` | Mark element as hx-preserve (survives any swap) |
| `render.OnAfterRequest` | hx-on::after-request handler |
| `render.OnAfterSettle` | hx-on::after-settle handler |
| `render.CacheFrame` | Set cache headers for frame response |
| `render.ForceReload` | Meta tag forcing full page reload |

---

## Micro-Interactions (Alpine.js / Stimulus.js)

Every interactive component has three render modes: CSS-only (default), Alpine.js, and Stimulus.js. Pick one per page.

### Alpine.js (recommended)

```go
import "github.com/emergent-company/go-daisy/components/alpine"

// Use pre-built state generators + attribute helpers
<div { alpine.XData(alpine.ModalState(false))... }>
    <div x-show="open" { alpine.Transition()... } { alpine.Trap("open")... }>
        ...
    </div>
</div>
```

Or use the ready-made Alpine variants:

```go
@modal.AlpineModal("Edit Profile", modal.ModalMD) {
    <p>Content here</p>
}
@ui.AlpineTabs(ui.TabsProps{...})
@ui.AlpineDropdown(ui.DropdownBottom) { ... }
@ui.AlpineAccordionItem("Section", false) { ... }
@ui.AlpineThemeToggle()
```

### Stimulus.js (explicit controller pattern)

```go
import "github.com/emergent-company/go-daisy/components/stimulus"

<div { stimulus.Controller("modal")... }>
    <button { stimulus.Action("click", "modal", "open")... }>Open</button>
</div>
```

Pre-built controllers: `ModalController`, `DropdownController`, `TabsController`, `AccordionController`, `ThemeController`, `ClipboardController`.

### Real-Time Updates

```go
import "github.com/emergent-company/go-daisy/stream"
import "github.com/emergent-company/go-daisy/streamhub"

hub := streamhub.New()
e.GET("/events", streamhub.EchoHandler(hub))
hub.BroadcastAppend("#messages", MessageRow(msg))
hub.BroadcastRefresh()
```

```go
stream.Append("#list", ItemRow(item))
stream.Morph("#detail", DetailView(data))
```

### DOM Morphing

```go
render.RenderAutoMorph(w, r, page, partial)
render.SetMorph(w, r)
```

Requires `layout.PageFull(layout.PageProps{Morph: true, ...})`.

### View Transitions

```go
layout.PageFull(layout.PageProps{ViewTransitions: true, ...})
```

Browser-native crossfade + slide animations between page navigations. Zero JS.

### Charts

```go
ui.Chart(ui.ChartProps{
    ID: "revenue", Type: ui.ChartBar, Stacked: true,
    Series: []ui.ChartSeries{
        {Name: "Revenue", Data: []float64{30, 40, 45, 50}},
        {Name: "Cost",    Data: []float64{20, 25, 30, 35}},
    },
    Categories: []string{"Q1", "Q2", "Q3", "Q4"},
    FillType: "gradient", FillOpacity: 0.9,
    ForecastCount: 2, GoalValues: []float64{55},
})
```

ApexCharts 4.5.0 (lazy-loaded from CDN). Supports area, bar, column, line, pie, donut, radial, rangeBar, heatmap. Stacking, gradients, annotations, forecasts, synced groups.

---

## Dev mode: `data-component` attribute annotations

Every component in the library emits a `data-component="package/ComponentName"` attribute on its outermost HTML element when **dev mode is active**. In production the attribute is never emitted — zero overhead.

### What it looks like

In the browser inspector (dev mode on):

```html
<button data-component="ui/Button" class="btn btn-primary">Save</button>
<div data-component="layout/Sidebar" id="layout-sidebar" …>…</div>
<a data-component="layout/sidebarNavItem" href="/dashboard" …>Dashboard</a>
```

This makes it trivial to find the Go source for any element on the page.

### Enabling dev mode

Call `devmode.WithDevMode(ctx)` once per request in your middleware or handler, then pass the derived context when rendering:

```go
import "github.com/emergent-company/go-daisy/devmode"

func myMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        if isDev {
            c.SetRequest(c.Request().WithContext(
                devmode.WithDevMode(c.Request().Context()),
            ))
        }
        return next(c)
    }
}
```

### `devmode` package API

```go
// Activate dev mode in a context (call once per request in dev).
ctx = devmode.WithDevMode(ctx)

// Check whether dev mode is active.
devmode.IsDevMode(ctx) // bool

// Return templ.Attributes{"data-component": name} in dev mode,
// or a shared empty map (no allocation) in production.
// Spread into any templ element:
//
//   <div { devmode.Attrs(ctx, "mypkg/MyComp")... }>
devmode.Attrs(ctx, "package/ComponentName") // templ.Attributes

// Wrap a component in a display:contents <div> with data-component
// and optionally data-props attributes (gallery hover overlay / component tree).
// props is optional — omit it to suppress data-props entirely.
devmode.ComponentBoundary("Button", ui.Button(props))           // no props
devmode.ComponentBoundary("Button", ui.Button(props), props)    // with props

// Inject data-component and optionally data-props directly onto the first
// opening tag (use for table structural elements: thead, tbody, tr, td, th).
devmode.ElementBoundary("TableRow", table.TableRow(id, hover))                          // no props
devmode.ElementBoundary("TableRow", table.TableRow(id, hover), map[string]any{...})    // with props
```

`ComponentBoundary` and `ElementBoundary` are used by the gallery infrastructure. `Attrs` is used by every component internally and is the right primitive when adding annotations to your own components.

### Value format

`"package/ComponentName"` — the package directory name and the exported (or unexported) function name exactly as written in Go source. Examples:

- `"ui/Button"`
- `"form/Field"`
- `"table/Table"`
- `"layout/Sidebar"`
- `"layout/sidebarNavItem"`

### Skipped elements

`data-component` is **not** added to `<table>`, `<thead>`, `<tbody>`, `<tr>`, `<th>`, or `<td>` outermost elements (browser table-parsing rules would strip injected wrapper `<div>` elements). Those components use `devmode.ElementBoundary` from gallery infrastructure instead.

The `<html>` element in `layout.Page` is also skipped intentionally.

---

## AI agent skills

go-daisy ships OpenCode skills that teach AI coding agents how to work with the library. The installer copies them into `.opencode/skills/` in your project automatically.

### Skills included

| Skill | Trigger phrases | What it covers |
|---|---|---|
| `gallery-install` | "install the gallery", "set up gallery", "init gallery" | Running the installer, generated files, Taskfile patching, port config |
| `gallery-add-component` | "add a component to the gallery", "add a story for X", "register this component" | `GalleryComponent` schema, HTML snippets, Templ components, variants, design tokens |
| `use-go-daisy` | "add a button/badge/table/modal", "create a page with sidebar", "render HTMX partial", "use go-daisy" | Full component API reference: layout, render helpers, ui primitives, nav, table, form, modal, logs |

### Installing the skills manually

The skills are installed automatically when you run the gallery installer:

```bash
go run github.com/emergent-company/go-daisy/cmd/install@latest
```

To install them without the gallery, copy them from `cmd/install/skills/` into `.opencode/skills/` in your project, or run the installer in a directory that already has a gallery:

```bash
# Re-run the installer to refresh skills (will skip non-empty gallery dir gracefully)
go run github.com/emergent-company/go-daisy/cmd/install@latest
```

Once installed, any OpenCode session in the project can load a skill on demand, and the agent will automatically activate the right skill based on what you ask for.

---

## Gallery

A live showcase app is included at `cmd/gallery`. Run it with:

```bash
task gallery
# → http://localhost:11000
```

The gallery is also importable as a package (`galleryruntime`) for embedding a component showcase in any project.

---

## Building

```bash
# Generate Templ files + compile Tailwind CSS + build Go binaries
task build

# Generate Templ + CSS only (no go build)
task build:ui

# Run the gallery
task gallery

# Watch CSS during development (separate terminal)
task dev:ui
```

> Always run `task build:ui` after editing any `.templ` file or CSS. The generated `*_templ.go` files must be committed alongside their `.templ` sources.

---

## Tech stack

| Layer | Tool |
|---|---|
| HTTP framework | Echo v4 |
| Templating | Templ |
| CSS | DaisyUI 5 + Tailwind CSS 4 |
| Interactivity | HTMX v4 |
| Micro-interactions | Alpine.js 3 or Stimulus 3 (optional, pick one) |
| DOM morphing | idiomorph (optional) |
| Real-time | HTMX SSE / WebSocket extensions (optional) |
| Charts | ApexCharts 4 (lazy-loaded CDN) |
| Static assets | Go `embed` (`staticfs/`) |
| Build | go-task (`Taskfile.yml`) |
