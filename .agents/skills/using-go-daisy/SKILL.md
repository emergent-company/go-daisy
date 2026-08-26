# Using go-daisy — AI Agent Skill

Skill for AI coding agents to install, wire, and use go-daisy components in consumer Go applications. Does not cover maintaining go-daisy itself, editing component internals, or release processes.

---

## 1. Install go-daisy

```bash
go get github.com/emergent-company/go-daisy@latest
templ generate  # only for YOUR .templ files, not go-daisy's (pre-generated)
```

go-daisy ships pre-generated `*_templ.go` files. Consumer agents do **not** run `templ generate` inside the go-daisy module.

---

## 2. Serve Bundled Assets

go-daisy bundles CSS + JS (HTMX, Alpine.js, morph, Stimulus, calendar, gantt). Mount them once:

```go
import "github.com/emergent-company/go-daisy/staticfs"

// In your router setup:
e := echo.New()
e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static", staticfs.Handler("/"))))
```

In your page `<head>`, use the helpers with automatic cache-busting:

```go
import "github.com/emergent-company/go-daisy/staticfs"
```

```templ
<head>
    { staticfs.Stylesheet("/css/app.css") }
    { staticfs.Script("/js/htmx.js") }
    { staticfs.Script("/js/alpine.js") }
    { staticfs.Script("/js/morph.js") }
</head>
```

**Common mistake:** Double `StripPrefix`. Handler already handles `/static/*` — don't wrap it again.

---

## 3. Render Components

Every component is a `templ.Component`. Import the package and call the function:

```templ
import "github.com/emergent-company/go-daisy/components/ui"

templ MyPage() {
    @ui.Button(ui.ButtonProps{
        Variant: ui.ButtonPrimary,
        Size:    ui.ButtonLG,
    }) {
        Click Me
    }
    @ui.Alert(ui.AlertProps{
        Type:    ui.AlertInfo,
        Message: "Welcome to go-daisy",
    })
}
```

**Components return `templ.Component` — not HTML strings.** Treat them like any Templ component. Use `@` syntax or call `.Render(ctx, w)`.

---

## 4. HTMX Rendering Rules

Use the `render` package helpers. Never call `component.Render(ctx, w)` directly in handlers:

```go
import "github.com/emergent-company/go-daisy/render"

// For pages inside a sidebar layout (3-tier: shell / content-area / tab):
render.RenderTriple(w, r, page, content, partial)

// For simpler pages inside a full shell:
render.RenderAuto(w, r, page, partial)

// For always-full-page or always-partial:
render.RenderPage(w, r, content)
render.RenderPartial(w, r, content)

// For morph-enabled variants:
render.RenderAutoMorph(w, r, page, partial)
render.RenderTripleMorph(w, r, page, content, partial)

// After a mutation (POST/PUT/DELETE), redirect the correct way:
render.RedirectAfterMutation(w, r, "/items/42")

// Append a toast OOB swap:
render.AppendToast(w, toast.Success, "Item saved.")
```

**Never** set `HX-Redirect` header manually. Use `RedirectAfterMutation`.

---

## 5. Layout Components

Use `layout.Page()` or `layout.PageFull()` as the outer HTML shell:

```go
import "github.com/emergent-company/go-daisy/components/layout"

layout.Page(layout.PageProps{
    Title:      "My App",
    Theme:      "light",
    Sidebar:    sidebarContent,
    MainContent: mainContent,
})
```

For custom shells, compose manually:

```go
@layout.AppShell(appName, sidebarGroups)
@layout.SidebarVariant(sidebarVariantProps)
```

---

## 6. Forms

Use `FormField` as the generic input — it handles labels, errors, placeholders, selects, textareas:

```go
import "github.com/emergent-company/go-daisy/components/form"

@form.FormField(form.FormFieldProps{
    Label:       "Email",
    Name:        "email",
    Type:        form.InputEmail,
    Placeholder: "you@example.com",
    Required:    true,
    Error:       validationError,
    Attrs: templ.Attributes{
        "hx-post": "/validate/email",
    },
})
```

For modals with form content:

```go
@form.FormModal(form.FormModalProps{
    ID:     "create-item",
    Action: "/items",
    Method: "post",
    Title:  "New Item",
})
```

**HTMX method dispatch:** Form components auto-dispatch `hx-post`/`hx-put`/`hx-patch` based on the `Method` field.

---

## 7. Alpine.js Components

For client-side micro-interactions (toggles, tabs, dropdowns, modals), use Alpine variants:

```go
import "github.com/emergent-company/go-daisy/components/alpine"
```

```templ
// Alpine tabs — instant client-side switching
<div { alpine.XData(alpine.TabState("tab1")) }>
    <button { alpine.On("click", "tab = 'tab1'") } :class="tab === 'tab1' ? 'tab-active' : ''">
        Tab 1
    </button>
    <div x-show="tab === 'tab1'" { alpine.Transition() }>Content 1</div>
</div>

// Alpine dropdown — keyboard nav + outside click dismiss
<div { alpine.XData(alpine.DropdownState(false)) }
     { alpine.Escape("open = false") }>
    <button { alpine.On("click", "open = !open") }>Menu</button>
    <div x-show="open" { alpine.Transition() } @click.outside="open = false">
        Items...
    </div>
</div>
```

**Alpine component conventions:**
- Outer div: `x-data` + `Escape` handler. Never `display:none`.
- Inner toggled element: `x-show` + `Transition()` + `style="display:none"`.
- `x-trap` for modals with focus trapping.
- `@click.outside` on the content panel (not outer wrapper).

Pre-built Alpine data components exist for: `Toggle`, `TabState`, `DropdownState`, `ThemeState`, `CounterState`, `AccordionState`, `SearchState`, `FormState`, `ModalState`.

---

## 8. Stimulus.js Components

For teams that prefer explicit controller-based JS, use Stimulus variants:

```go
import "github.com/emergent-company/go-daisy/components/stimulus"
```

```templ
<div { stimulus.Controller("dropdown") }
     { stimulus.Action("keydown.escape", "dropdown", "close") }>
    <button { stimulus.Action("click", "dropdown", "toggle") }
            data-dropdown-target="trigger">Menu</button>
    <div data-dropdown-target="menu" style="display:none">
        Items...
    </div>
</div>
```

Pre-built controllers: `ModalController`, `DropdownController`, `TabsController`, `AccordionController`, `ThemeController`, `ClipboardController`.

---

## 9. Real-Time (SSE / WebSocket)

Use `streamhub` for pub/sub and `stream` for HTML fragment broadcasting:

```go
import "github.com/emergent-company/go-daisy/streamhub"

hub := streamhub.New()

// Broadcast to channel subscribers:
hub.Broadcast("orders", stream.Append("#items", itemRow))
hub.Broadcast("notifications", stream.Replace("#badge", badge))

// SSE handler:
e.GET("/events", func(c echo.Context) error {
    return hub.ServeSSE(c.Response(), c.Request())
})
```

---

## 10. Web Components

For JS-heavy widgets, use Custom Elements:

```templ
<go-daisy-calendar data-value="2025-01-15" data-locale="en-US"></go-daisy-calendar>
<go-daisy-gantt data-tasks='[{"id":"1","name":"Design","start":"2025-01-01","end":"2025-01-05"}]'></go-daisy-gantt>
```

Include their JS: `staticfs.Script("/js/go-daisy-calendar.js")`.

---

## 11. Common Integration Mistakes

| Mistake | Fix |
|---|---|
| Calling `templ generate` inside go-daisy | Don't — components are pre-generated |
| Double `StripPrefix` on asset handler | `staticfs.Handler("/")` already handles path |
| Setting `HX-Redirect` manually | Use `render.RedirectAfterMutation(w, r, path)` |
| Using `render.RenderPage` for sidebar pages | Use `RenderTriple` (3-tier: shell / content / tab) |
| Raw `<span class="iconify">` for icons | Use `@ui.IconSpan(name, size)` |
| `ExtraClass` or `class string` positional | Use `Class string` in Props struct (new components) |
| Missing `staticfs.Script("/js/htmx.js")` | HTMX is not auto-loaded — you must include it |
| Missing `staticfs.Script("/js/alpine.js")` | Alpine variants require the JS tag in `<head>` |
| Forgetting `style="display:none"` on Alpine `x-show` | Required for initial hidden state |
| Server returns 200 OK on form POST | Must return 303 redirect (use `RedirectAfterMutation`) |
| Using `innerHTML` swaps without morphing | Enable morph with `RenderAutoMorph` / `RenderTripleMorph` |

---

## 12. Page Templates

go-daisy ships pre-built page templates:

```go
import "github.com/emergent-company/go-daisy/pages"

@pages.AuthPageWithBoundary(pages.AuthPageProps{...})
@pages.DashboardPageWithBoundary(pages.DashboardProps{...})
@pages.ChatLayoutWithBoundary(pages.ChatLayoutProps{...})
@pages.LandingPageWithBoundary(pages.LandingPageProps{...})
@pages.SettingsPageWithBoundary(pages.SettingsPageProps{...})
```

---

## Skill Metadata

- **Module:** `github.com/emergent-company/go-daisy`
- **Scope:** Consumer application integration only
- **Excludes:** Library maintenance, component internals, release processes, gallery development
- **Prerequisites:** Go 1.21+, Templ (`github.com/a-h/templ v0.3.1020+`), Echo v4
