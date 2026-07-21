# go-daisy — Agent Guide

## What this project is

`go-daisy` is a Go UI component library for building HTMX-driven web interfaces. It provides type-safe, reusable Templ components styled with DaisyUI (Tailwind CSS). The project also ships a live **gallery** app at `cmd/gallery` that showcases every component.

Module path: `github.com/emergent-company/go-daisy`

---

## Repository layout

```
go-daisy/
├── cmd/
│   ├── gallery/          # Gallery showcase app (Echo web server on :11000)
│   │   └── internal/gallery/ # seed.go — component registry for the showcase
│   ├── install/          # Installer: `go run github.com/emergent-company/go-daisy/cmd/install@latest`
│   ├── boundarytoken/    # Codegen: parses gallery:token annotations → tokens_*_gen.go
│   ├── nexus/            # Reference dogfood app using go-daisy components
│   └── nexus-compare/    # Playwright-based visual comparison vs nexus-html
├── components/           # The component library (9 packages)
│   ├── alpine/           # Alpine.js attribute generators (no .templ files)
│   ├── stimulus/         # Stimulus.js attribute generators (no .templ files)
│   ├── form/             # Form inputs, field wrappers, validators, editors (~26 components)
│   ├── layout/           # Page shell, sidebar variants, navbar, topbar (~10 components)
│   ├── logs/             # Log stream display (~1 component)
│   ├── modal/            # Modal dialogs (vanilla, Alpine, Stimulus) (~3 components)
│   ├── nav/              # Page headers, tab menus, breadcrumbs, footer (~16 components)
│   ├── table/            # Data tables, list areas, scroll rows (~6 components)
│   └── ui/               # Primitives: button, badge, card, avatar, toast, alert, chart, etc. (~73 components)
├── pages/                # Pre-built page templates (auth, chat, dashboard, landing, settings)
├── render/               # HTMX-aware rendering helpers
├── stream/               # Turbo-Stream-style HTML fragment broadcasting
├── streamhub/            # SSE/WebSocket pub/sub hub for real-time updates
├── webcomponents/        # Custom Elements for heavy client widgets (calendar, gantt)
├── devmode/              # Zero-overhead dev tooling annotations
├── shared/               # Compose, Ternary, StrComp, RenderInto, ActiveKV
├── staticfs/             # Embedded static assets (CSS, JS, fonts) with cache-busting
├── galleryruntime/       # Reusable gallery package (importable by any project)
│   ├── types.go          # GalleryComponent, GalleryStory, DesignToken, Category constants
│   ├── helpers.go        # ComponentBySlug, SlugifyStoryName, BuildCategoryGroups, TokenGroups
│   ├── serve.go          # Serve(Options) — starts the Echo gallery server
│   ├── handler.go        # HTTP route handlers (package galleryruntime)
│   ├── slugs.go          # ComponentSlugs map (data-component → gallery slug)
│   ├── devoverlay.go     # Hover overlay + component tree iframe bridge JS
│   ├── pages_shell.templ # Gallery shell + sidebar + search modal
│   ├── pages_detail.templ# Component detail page (preview iframe, tokens, component tree)
│   ├── pages_index.templ # Gallery landing/index page
│   ├── pages_group.templ # Category/subcategory group pages
│   └── pages_docs.templ  # SDK documentation reference page
├── assets/               # Tailwind CSS source (app.css)
├── .agents/skills/       # AI agent skills (DaisyUI reference)
├── docs/                 # ROADMAP.md, PATTERNS.md
├── tests/e2e/            # Playwright E2E test suite
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
| `RenderAutoMorph` | `RenderAuto` with DOM morphing enabled |
| `RenderTripleMorph` | `RenderTriple` with DOM morphing enabled |
| `RenderPartialMorph` | Partial render with morphing |
| `RedirectAfterMutation` | HX-Redirect for HTMX, 303 for plain requests |
| `AppendToast` | Write an `hx-swap-oob` toast fragment into the response |
| `SetMorph` | Enable DOM morphing on an HTMX swap |
| `Preserve` | Returns `templ.Attributes` with `hx-preserve="true"` |
| `ForceReload` | Returns a `<meta>` that forces full page reload |
| `CacheFrame(w, maxAge)` | Set private cache headers for frame responses |
| `CacheSharedFrame(w, maxAge)` | Set public cache headers for frame responses |
| `OnAfterRequest(js)` | Returns `hx-on::after-request` attribute |
| `OnAfterSettle(js)` | Returns `hx-on::after-settle` attribute |
| `OnResponseError(js)` | Returns `hx-on::response-error` attribute |
| `OnBeforeRequest(js)` | Returns `hx-on::before-request` attribute |

Detection helpers: `IsHTMX`, `IsPartial`, `IsMainContentTarget`, `IsHistoryRestore`, `IsHistoryRestoreFromContext`, `IsScrollLoad`, `HXTarget`.

### `staticfs`
Embedded static assets with cache-busting content hashes.

| Function | Returns |
|---|---|
| `Stylesheet(path)` | `<link rel="stylesheet" href="/static{path}?v={hash}">` |
| `Script(path)` | `<script src="/static{path}?v={hash}"></script>` |
| `AssetPath(path)` | `/static{path}?v={hash}` (bare URL) |
| `Handler(prefix)` | `http.Handler` serving embedded assets with correct MIME types |
| `Hash()` | 12-char content hash for cache busting |
| `FS()` | Embedded `fs.FS` rooted at `static/` |

**Bundled JS (11 files):** `htmx.js`, `morph.js`, `alpine.js`, `htmx-sse.js`, `htmx-ws.js`, `hx-head.js`, `stimulus.js`, `stimulus-controllers.js`, `go-daisy-calendar.js`, `go-daisy-gantt.js`, `frappe-gantt.js`.

**Bundled CSS:** `app.css` (compiled Tailwind + DaisyUI), `frappe-gantt.css`.

### `components/alpine`
Utility package — no `.templ` files. Exports Go functions returning `templ.Attributes` for Alpine.js directives.

**Attribute generators:** `XData(s State)`, `Show(expr)`, `Transition()`, `TransitionDuration(d)`, `Collapse()`, `On(event, handler)`, `Escape(handler)`, `Outside(handler)`, `Bind(attr, expr)`, `Model(field)`, `Trap(expr)`, `Ref(name)`, `Cloak()`, `Init(expr)`, `Text(expr)`, `HTML(expr)`, `If(expr)`, `For(expr)`, `ClassBinding(expr)`, `BooleanBinding(attr, expr)`, `Effect(expr)`, `Merge(attrs...)`.

**State constructors (data.go):** `Toggle(open bool)`, `TabState(active)`, `DropdownState(open)`, `ThemeState(dark)`, `CounterState(count)`, `AccordionState(openItem)`, `SearchState(items)`, `FormState()`, `ModalState(open)`.

**Lifecycle helpers (lifecycle.go):** `ModalInit()`, `DropdownKeyboardInit()`, `TabsHistoryInit()` — returns init expression strings.

**Component tags:** `Tag()`, `TagDefer()` (Alpine.js script), `MorphTag()` (idiomorph script).

### `components/stimulus`
Utility package — no `.templ` files. Exports Go functions returning `templ.Attributes` for Stimulus.js controllers.

**Attribute generators:** `Controller(name)`, `Action(event, controller, method)`, `Actions(actions...)`, `Target(controller, name)`, `Targets(controller, name)`, `Param(controller, name, value)`, `Value(controller, name, value)`, `Outlet(controller, name, selector)`.

**Pre-built JS controllers** in `controllers/go-daisy-controllers.js`: `ModalController`, `DropdownController`, `TabsController`, `AccordionController`, `ThemeController`, `ClipboardController`.

**Component tag:** `Tag()` (Stimulus.js script).

### `stream`
Turbo-Stream-style HTML fragment broadcasting. Generates `<turbo-stream>` elements with actions: `Append`, `Prepend`, `Replace`, `Update`, `Remove`, `Before`, `After`, `Refresh`, `Morph`. Single file: `stream/stream.go`.

### `streamhub`
SSE/WebSocket pub/sub hub. `streamhub.New()` creates a hub; `hub.Broadcast(ch, msg)` fans out to channel subscribers; `hub.ServeSSE(w, r)` upgrades HTTP to SSE. Files: `hub.go`, `echo.go`, `templ.go`.

### `webcomponents`
Custom Elements for JS-heavy widgets. Self-contained, framework-agnostic. Files: `go-daisy-calendar.js`, `go-daisy-gantt.js`.

### `pages`
Pre-built page templates (auth, chat, dashboard, landing, settings). Follows same conventions as `components/`: `.templ` files, `boundary.go` with `*WithBoundary` wrappers, `types.go` for props structs.

### `components/ui`
Primitive DaisyUI components (~73). All are `templ.Component` values returned by Go functions. Notable ones:

- `Button`, `Badge`, `Card`, `Avatar`
- `Toast`, `Alert`, `Loader`, `EmptyState`
- `Pagination`, `Filter`, `ActionMenu`, `StatCard`
- `Chart` (ApexCharts), `Carousel` (SwiperJS), `Diff`, `Timeline`
- `Accordion` (vanilla/Alpine/Stimulus), `Tabs` (vanilla/Alpine/Stimulus)
- `Dropdown`, `Drawer`, `Hero`, `CommandPalette`
- `ThemeToggle`, `ThemeSwitcher`, `ThemeController`
- `Frame`, `CodePreview`, `Hover3D`, `Swap`
- `LayoutCustomizer`, `NotificationDropdown`, `Dashboard`

### `components/layout`
Full-page shell, sidebar layout, and navbar (~10 components). `AppShell`, `SidebarVariants`, `TopbarVariants`, `LayoutBuilder`, `SidebarDense`, `ViewMenu`.

### `components/nav`
Page headers, tab menus, breadcrumbs, footer (~16 components). `PageHeader`, `TabMenu`, `SimpleTabs`, `TopBar`, `PageTitleVariant`, `SearchModal`, `NotificationPanel`, `ProfileMenuVariant`, `FooterVariant`, `Megamenu`.

### `components/form`
Form inputs, field wrappers, validators, editors (~26 components). `FormField`, `TextInput`, `TextareaInput`, `SelectInput`, `EnhancedSelect`, `SearchInput`, `FilterSelect`, `CheckboxInput`, `Toggle`, `RadioGroup`, `RangeInput`, `FileInput`, `FileUpload`, `Rating`, `OTP`, `PasswordField`, `PasswordMeter`, `DatePicker`, `Clipboard`, `RichTextEditor`, `CalendarWrapper`, `ColorInput`, `Spinner`, `Wizard`.

### `components/table`
`Table`, `DataTable`, `TableCard`, `ListArea` (infinite-scroll container), `ScrollRows` (6 components).

### `components/modal`
Modal dialogs (3 components): `Modal`, `FormModal`, `AlpineFormModal`, `StimulusFormModal`, `ConfirmPopup`, `LoaderModal`, `OpenModalButton`.

### `components/logs`
Log stream display (1 component).

---

## Tech stack

| Layer | Tool |
|---|---|
| HTTP framework | Echo v4 (`github.com/labstack/echo/v4`) |
| Templating | Templ (`github.com/a-h/templ v0.3.1020`) |
| CSS | DaisyUI + Tailwind CSS (via Node CLI) |
| Interactivity | HTMX, Alpine.js, Stimulus.js |
| Client-side JS | idiomorph (DOM morphing), htmx extensions (SSE, WS, head) |
| Real-time | `stream/` (Turbo-Stream-style), `streamhub/` (SSE/WS pub/sub) |
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
2. Add `*WithBoundary` wrapper in `components/<package>/boundary.go` — this is mandatory for every exported component.
3. Run `task build:ui` — this generates `<name>_templ.go` and recompiles CSS.
4. Add a gallery page under `cmd/gallery/internal/handler/` to showcase the component.
5. Register route in `cmd/gallery/internal/handler/handler.go`.
6. Register in `cmd/gallery/internal/gallery/seed.go` via `GalleryComponent{...}` struct with `Preview` and `Source` fields. Slug must match handler route.

### Standard Props struct template

New components MUST use this standard Props struct pattern:

```go
type NewComponentProps struct {
    Class string            // additional CSS classes (appended, not overridden)
    Attrs templ.Attributes  // caller attribute overrides (spread last on root element)
    // ... component-specific fields ...
}
```

```templ
templ NewComponent(props NewComponentProps) {
    <div class={ "base-classes", props.Class }
        { devmode.Attrs(ctx, "pkg/NewComponent")... }
        { props.Attrs... }
    >{ children... }</div>
}
```

Legacy `attrs templ.Attributes` positional param is deprecated — migrate to `props.Attrs` when touching a component. Same for `ExtraClass` → `Class`.

### boundary.go imports reminder

Boundary wrappers for components with children require these imports:

```go
import (
    "context"
    "io"

    "github.com/a-h/templ"
    "github.com/emergent-company/go-daisy/devmode"
)
```

### HTMX method dispatch (form components)

Form-modal components dispatch HTMX method based on a Method field:

```templ
if props.Method == "put" {
    hx-put={ props.Action }
} else if props.Method == "patch" {
    hx-patch={ props.Action }
} else {
    hx-post={ props.Action }
}
```

Used in `FormModal`, `AlpineFormModal`, `StimulusFormModal`, `ActionMenu`.

### Demo vs real component pattern

Components that embed heavy JS libraries (calendars, charts, editors) have two variants:
- `CalendarWrapper` — real composable component (thin wrapper, no default content)
- `CalendarDemo` — pre-configured demo for gallery showcase (hardcoded config, full visuals)

Both need `*WithBoundary` wrappers. The demo variant is only for gallery registration, never dependencies.

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

**Bundled JavaScript** lives in `staticfs/static/js/` — 11 files embedded alongside CSS: `htmx.js`, `morph.js` (idiomorph), `alpine.js`, `htmx-sse.js`, `htmx-ws.js`, `hx-head.js`, `stimulus.js`, `stimulus-controllers.js`, `go-daisy-calendar.js`, `go-daisy-gantt.js`, `frappe-gantt.js`. Served via `staticfs.Script("/js/htmx.js")` etc.

**Cache busting**: `staticfs.Hash()` returns a 12-char content hash updated via `go generate`. `Stylesheet()` and `Script()` helpers append `?v={hash}` query params automatically. Use `data-asset-track="reload"` on tags to trigger full reload on version mismatch.

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

Two patterns exist — choose based on the component's needs:

**Override pattern** (caller's class replaces defaults):
```templ
class={ "base-class", templ.KV(class, class != ""), templ.KV("default-class", class == "") }
```
Use `class string` as last param. When non-empty, replaces defaults; when empty, uses built-in defaults. Used in `SearchInput`, `FilterSelect`.

**Append pattern** (caller's class appended to defaults):
```templ
class={ "base-class", templ.KV("conditional-class", condition), props.Class }
```
Props struct has `Class string` field. Caller-supplied CSS always appended after component's base classes. This is the dominant pattern for props-based components. Use this for new components.

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

### Props struct placement

Define Props structs in the `.templ` file with the component (before the `templ` declaration). This is the dominant pattern — ~85% of components follow it. Only place Props structs in a package `.go` file when shared across multiple `.templ` files within the same package (e.g., `tabs.go` defines `TabsProps` used by vanilla/Alpine/Stimulus variants).

### Alpine / Stimulus utility packages

`components/alpine/` and `components/stimulus/` are utility packages that export Go functions returning `templ.Attributes`. They have no `.templ` files and no `boundary.go`. Components use them by spreading the returned attributes:
```templ
<button { alpine.On("click", "open = !open")... }>Toggle</button>
```
These packages provide attribute generators for common JS framework patterns — never raw template rendering.

### `pages/` package follows component patterns

The top-level `pages/` package (not under `components/`) follows the same conventions: `.templ` files, `boundary.go` with `*WithBoundary` wrappers, `types.go` for props structs. It is architecturally a component package, just organized at the module root for import convenience.

### Sub-component naming

Internal sub-components (helpers rendered within larger components) are exported as `templ` functions so they can be composed, but they do NOT need `*WithBoundary` wrappers or `gallery:token` annotations. Examples: `CardRaw`, `IconSpan`, `FieldError`, `IndicatorBadge`, `AccordionItem`. Sub-components that are standalone and useful on their own SHOULD get boundaries.

### `form.FormField` as generic input

For inputs with placeholder/rows/options/required, use `form.FormField(form.FormFieldProps{...})` instead of `form.TextInput`/`form.TextareaInput`/`form.SelectInput`. The props struct supports all edge cases (placeholder, SelectOption slice, rows, error, attrs, etc.) that the shorthand functions don't cover.

### Dead code removal

Before removing a component function, verify zero callers:
```
rg "FunctionName" --include '*.go' | grep -v '_templ.go' | grep -v '_test.go'
```
`grep -v '_templ.go'` filters generated-file matches (templ generates forwarder funcs). If only `_templ.go` matches remain, the component is dead.

### devmode boundary wrapping

Every exported component MUST have a `*WithBoundary` wrapper in `components/<pkg>/boundary.go`. The wrapper injects `devmode.ComponentBoundary(name, comp, props)` for gallery dev tooling (hover overlay, component tree, source view). In production it's a zero-overhead passthrough.

```go
func ButtonWithBoundary(href string, variant ButtonVariant, ...) templ.Component {
    props := ButtonProps{Href: href, Variant: variant, ...}
    return devmode.ComponentBoundary("Button", Button(props), props)
}
```

For table elements (`<thead>`, `<tbody>`, `<tr>`, `<th>`, `<td>`), use `devmode.ElementBoundary` instead — it injects data attributes directly rather than wrapping with a `<div>`, which would be invalid HTML inside `<table>`.

### gallery:token / gallery:hint annotations

Each `*WithBoundary` function must have Go comments listing controllable design tokens and hint defaults. Parsed by `cmd/boundarytoken` to generate the gallery sidebar.

```go
// gallery:token variant,size,style,typ,shape,icon,loading,block
// gallery:hint href:default(#)
func ButtonWithBoundary(...)
```

Hints support `default()`, `range(min,max,step)`, and `slice(n)` types. Tokens become interactive controls in the gallery.

### Inline JS singleton pattern

Scripts embedded in components use a singleton guard to prevent re-init on HTMX swaps:

```templ
templ myScript() {
    <script>
    if (!window._myInit) {
        window._myInit = true;
        // init code here
    }
    </script>
}
```

Call via `@myScript()` inside the component. This avoids duplicate event listeners and leaked state on partial page swaps.

### Variant dispatch pattern

Multi-variant components use a `switch` on a `Style` string to render different hardcoded layouts. Each case calls a dedicated sub-component.

```templ
templ PageTitleVariant(opts PageTitleVariantOpts) {
    switch opts.Style {
    case "minimal":
        @PageTitleMinimal(...)
    case "editor":
        @PageTitleEditor(...)
    default:
        @PageTitleMinimal(...)
    }
}
```

Used in `PageTitleVariant`, `SearchModal`, `NotificationVariant`, `ProfileMenuVariant`, `FooterVariant`, `SidebarVariants`, `TopbarVariants`.

### children variadic pattern

Optional override slots use Go variadic params. Checked with `len(x) > 0 && x[0] != nil`:

```templ
templ Sidebar(appName string, groups []SidebarGroup, logo ...templ.Component) {
    if len(logo) > 0 && logo[0] != nil {
        @logo[0]
    } else {
        <span>{ appName }</span>
    }
}
```

### Attrs must be spread last

`{ attrs... }` or `{ props.Attrs... }` always goes LAST in the element attribute list. This allows callers to override component defaults via `templ.Attributes`:

```templ
<div class="base-classes" { devmode.Attrs(ctx, "pkg/Component")... } { attrs... }>
```

Pass `nil` when no override needed — it's nil-safe.

### Zero-value = default

Typed string constants use the zero value (empty string) to mean "use default styling". No explicit class emitted for the zero value:

```go
const ButtonMD ButtonSize = ""       // zero value = default size
const SortNone SortDir = ""
const ModalMD ModalSize = ""
```

This keeps generated HTML clean — only non-default values produce CSS classes.

### Shared `ternary` helper

The `ternary` helper lives in `shared/ternary.go`. Each component package provides a thin wrapper so `.templ` files can call `ternary(...)` unqualified:

```go
// shared/ternary.go
package shared
func Ternary(cond bool, a, b string) string { ... }

// components/ui/ui.go
import "github.com/emergent-company/go-daisy/shared"
func ternary(cond bool, a, b string) string { return shared.Ternary(cond, a, b) }
```

Do NOT define `ternary` with an inline body in your package — always delegate to the shared implementation.

### CSS class prop naming

Props structs use `Class string` for caller-injected CSS. New components MUST use `Class` (not `ExtraClass`, not `class` as positional param):

```go
type NewComponentProps struct {
    Class string // additional CSS classes
}
```

Legacy components using `ExtraClass struct field` or `class string` positional param still work but are deprecated for new code. The semantics are identical — the value is appended to the root element's class list, or replaces defaults when non-empty.

### Props struct over positional args

Prefer Props structs over flat positional parameters for all new components. Props structs are extensible, self-documenting, and enable the gallery's design token system. The migration path:

| Component | Current | Preferred |
|---|---|---|
| Alert | positional: `typ, icon, message, attrs` | `AlertProps{...}` |
| Pagination | positional: `currentPage, totalPages, baseURL, targetID` | `PaginationProps{...}` |
| Toast | positional: `typ, message` | `ToastProps{...}` |
| Avatar | positional: `name, src, icon, size, attrs` | `AvatarProps{...}` |

When adding a new feature to a component, prefer extending its Props struct (or creating one if none exists) over adding positional params.

### Refactoring milestones

| Phase | Status | Description |
|-------|--------|-------------|
| 1 — intStr cleanup | ✅ Done | Replaced all `intStr/sprintCount/promptIntStr/spinnerIntStr/intToStr` with `strconv.Itoa`; removed 6 custom functions |
| 2 — IconSpan composition | ✅ Done | Replaced ~120 manual `<span class="iconify ...">` with `@ui.IconSpan(name, size)` across nav/layout/form; added `aria-hidden` to IconSpan |
| 3 — Badge composition | ✅ Done | Replaced 9 manual `<div class="badge ...">` with `@ui.Badge(props)`; added `BadgeSizeXS` constant |
| 4 — Button composition | ✅ Done | Refactored `ui.Button` to `ButtonProps` struct with `Style`, `Block`, `ExtraClass` fields; replaced ~31 manual `<button class="btn ...">` with `@ui.Button(props)` across nav/layout/form/modal; updated seed.go callers |
| 5 — Card composition | ✅ Done | Made `ui.StatCardMinimal` (non-icon path) use `@CardRaw` instead of raw `<div class="card">` |
| 6 — ternary dedup | ✅ Done | Extracted duplicated `ternary` into `shared/ternary.go`; each package has thin wrapper |
| 7 — AlertType/ToastType merge | ✅ Done | Merged identical types into `shared/` and `ui/ui.go`; `ToastType = AlertType` alias |

No cross-package import issues: all packages that need `ui` already import it or can safely do so (no circular deps). `form/range.templ` was the reference for `strconv.Itoa` inside templ expressions.

- All component functions return `templ.Component`.
- Props are passed as plain Go structs or positional arguments — no global state.
- Tailwind classes are written directly in `.templ` files; DaisyUI component classes (`btn`, `badge`, `card`, etc.) are preferred over raw utility classes.
- File naming: `<component-name>.templ` (kebab-case), package name matches directory name.
- Generated files (`*_templ.go`) are checked in to source control.

### `shared.Compose` — combining multiple Components

When you need to concatenate multiple `templ.Component` values into a single renderable, use `shared.Compose`:

```go
import "github.com/emergent-company/go-daisy/shared"

combined := shared.Compose(header, body, footer)
```

The `layout` package provides a local wrapper for convenience:
```go
// components/layout/sidebar-elements.go
func children(components ...templ.Component) templ.Component {
    return shared.Compose(components...)
}
```

Do NOT re-implement this loop pattern inline — use `shared.Compose`.

### Boundary wrapper for components with children

Components that accept `templ.Component` children (or render strings as children) need a two-tier `templ.ComponentFunc` wrapper in `boundary.go`:

**String-as-child pattern** (e.g., `Divider`, `Kbd`, `Link`):
```go
func DividerWithBoundary(color DividerColor, vertical bool, label string) templ.Component {
    child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        _, err := io.WriteString(w, label)
        return err
    })
    inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return Divider(color, vertical).Render(templ.WithChildren(ctx, child), w)
    })
    return devmode.ComponentBoundary("Divider", inner, map[string]any{
        "color": string(color), "vertical": vertical, "label": label,
    })
}
```

**Component-as-child pattern** (e.g., `Mask`, `Fieldset`, `Stack`):
```go
func MaskWithBoundary(shape MaskShape, content templ.Component) templ.Component {
    outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return Mask(shape).Render(templ.WithChildren(ctx, content), w)
    })
    return devmode.ComponentBoundary("Mask", outer, map[string]any{"shape": string(shape)})
}
```

**Item-list pattern** (e.g., `Carousel`, `Timeline`, `List`, `MockupCode`):
```go
func TimelineWithBoundary(items []TimelineItemProps) templ.Component {
    inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        for i, item := range items {
            if err := TimelineItem(item, i == 0, i == len(items)-1).Render(ctx, w); err != nil {
                return err
            }
        }
        return nil
    })
    outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return Timeline().Render(templ.WithChildren(ctx, inner), w)
    })
    return devmode.ComponentBoundary("Timeline", outer, map[string]any{"itemCount": len(items)})
}
```

For simple item iteration that only needs concatenation (no per-item logic), use `shared.Compose` inside the outer `ComponentFunc`.

**Mockup-with-placeholder pattern** (e.g., `MockupBrowser`, `MockupPhone`, `MockupWindow`):
```go
func MockupBrowserWithBoundary(url string) templ.Component {
    inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return MockupBrowser(url).Render(templ.WithChildren(ctx, MockupBrowserPlaceholder()), w)
    })
    return devmode.ComponentBoundary("MockupBrowser", inner, map[string]any{"url": url})
}
```

**Variadic-children pattern** (e.g., `Join`, `Stack`):
```go
func JoinWithBoundary(vertical bool, children ...templ.Component) templ.Component {
    content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        for _, c := range children {
            if err := c.Render(ctx, w); err != nil {
                return err
            }
        }
        return nil
    })
    outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return Join(vertical).Render(templ.WithChildren(ctx, content), w)
    })
    return devmode.ComponentBoundary("Join", outer, map[string]any{"vertical": vertical})
}
```

**Nested item-list with sub-boundaries** (e.g., `Carousel`):
Each sub-item gets its own `devmode.ComponentBoundary` so they appear individually in DevTools:
```go
func CarouselWithBoundary(snap CarouselSnap, vertical bool, width string, items []CarouselItemProps) templ.Component {
    children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        for _, item := range items {
            it := item
            inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
                return CarouselItem(it.ID, it.ItemWidth).Render(templ.WithChildren(ctx, it.Content), w)
            })
            itemBoundary := devmode.ComponentBoundary("CarouselItem", inner, map[string]any{"id": it.ID})
            if err := itemBoundary.Render(ctx, w); err != nil {
                return err
            }
        }
        return nil
    })
    outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return Carousel(snap, vertical, width).Render(templ.WithChildren(ctx, children), w)
    })
    return devmode.ComponentBoundary("Carousel", outer, map[string]any{"snap": string(snap), "itemCount": len(items)})
}
```

**Multi-sub-component pattern** (e.g., `Diff`, `Hero`):
Multiple sub-components composed in sequence, wrapped in a single outer container:
```go
func DiffWithBoundary(before string, after string) templ.Component {
    inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        if err := DiffItem1().Render(templ.WithChildren(ctx, DiffItemContent(before, false)), w); err != nil {
            return err
        }
        if err := DiffItem2().Render(templ.WithChildren(ctx, DiffItemContent(after, true)), w); err != nil {
            return err
        }
        return DiffResizer().Render(ctx, w)
    })
    outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        return DiffContainer().Render(templ.WithChildren(ctx, inner), w)
    })
    return devmode.ComponentBoundary("Diff", outer, map[string]any{"before": before, "after": after})
}
```

**No-props / minimal boundary** (for components with no configurable parameters):
```go
func NoPermissionsWithBoundary() templ.Component {
    return devmode.ComponentBoundary("NoPermissions", NoPermissions())
}
```

### Package `.go` file naming

| File | Purpose |
|---|---|
| `<pkg>.go` (e.g., `ui.go`, `form.go`) | Ternary wrapper, shared type aliases, helper constructors |
| `types.go` | Props/type-only declarations (used by `pages/`; not required by component sub-packages) |
| `boundary.go` | All `*WithBoundary` wrappers + `gallery:token` annotations |
| `<component>-variants.go` | Variant-default helpers (e.g., `sidebar-variants.go`, `topbar-variants.go`) |
| `<component>-elements.go` | Sub-item types and helpers (e.g., `sidebar-elements.go`) |
| `helpers.go` | Package-level defaults and utilities (e.g., `modal/helpers.go`) |

Do NOT scatter boundary wrappers across multiple files — all go in `boundary.go`.

### `devmode.Attrs()` placement in `.templ`

Every component's root element calls `devmode.Attrs(ctx, "pkg/Component")` to annotate with `data-component` for dev tooling. Place it after base CSS classes but before `{ attrs... }`:

```templ
templ Button(props ButtonProps) {
    <button class={ "btn", ..., props.ExtraClass }
        { devmode.Attrs(ctx, "ui/Button")... }
        { props.Attrs... }
    >{ children... }</button>
}
```

Do NOT skip this — gallery dev tools depend on it for hover overlays and component tree navigation.

### Gallery registration pattern

Components are registered in `cmd/gallery/internal/gallery/seed.go` via `GalleryComponent` structs:

```go
{Name: "Button", Slug: "button", Package: "ui", Category: "Primitives",
    Preview: ui.ButtonWithBoundary("#", ui.ButtonPrimary, ...),
    Source:  ui.ButtonWithBoundary as templ.Component,
}
```

Each component gets a `Preview` (shown in gallery list/cards) and a `Source` (full render in detail view). The `Slug` must match the route registered in `cmd/gallery/internal/handler/handler.go`.
