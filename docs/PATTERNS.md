# Component Creation Patterns

> Canonical patterns for go-daisy components. Every new component must follow these.

---

## 1. Props

**Go struct** for 4+ params. **Positional** for 1–3 params.

```go
// Preferred — typed struct
type ButtonProps struct {
    Variant ButtonVariant
    Size    ButtonSize
    Attrs   templ.Attributes
}

// Acceptable — positional for simple components
templ Alert(typ AlertType, icon string, message string, attrs templ.Attributes) { ... }
```

Enum types always use `type Foo string` with `const` block:

```go
type ButtonVariant string
const (
    ButtonPrimary ButtonVariant = "btn-primary"
    ButtonGhost   ButtonVariant = "btn-ghost"
)
```

---

## 2. CSS Classes

Always `templ.KV` for conditionals. Bare strings for always-present classes.

```templ
class={ "btn",
    templ.KV(string(props.Variant), props.Variant != ""),
    templ.KV("btn-block", props.Block),
    props.ExtraClass }
```

**Class override pattern** (when `class` is the last positional param):

```templ
class={ templ.KV(custom, custom != ""), templ.KV("default-class", custom == "") }
```

`ExtraClass` fields are additive — never replace defaults, always appended:

```templ
// Wrong:  class={ props.ExtraClass }           // replaces base "btn" class
// Right:  class={ "btn", props.ExtraClass }    // additive
```

**Borders always come from the line ladder** (`border-line`, `border-line-strong`,
`border-line-translucent`) — never `border-base-300` /
`border-base-content/N`. Surface outlines draw that colour at the theme's border
width with `edge-line` / `edge-line-strong`. See §12.

---

## 3. Icons

**Always `@IconSpan(name, size)`.** Never raw `<span class="iconify ...">`.

```templ
// Right
@IconSpan("lucide--x", "size-5")

// Wrong
<span class="iconify lucide--x size-5"></span>
```

`IconSpan` provides `aria-hidden="true"` automatically — raw spans do not.

Icon size convention: `size-4` for inline content, `size-5` for standalone emphasis, `size-4.5` for label associations.

---

## 4. HTMX Attributes

Only emit when both endpoint and target are non-empty:

```templ
if hxGet != "" && hxTarget != "" {
    hx-get={ hxGet }
    hx-target={ "#" + hxTarget }  // caller omits #, component prepends it
    hx-trigger="change"
    hx-include="closest form"
}
```

For components that don't own the request lifecycle, inject HTMX via `props.Attrs` spread:

```templ
<button { props.Attrs... }>...</button>
```

---

## 5. devmode.Attrs

One annotation per outermost element:

```templ
<div { devmode.Attrs(ctx, "ui/MyComponent")... }>
```

Name format: `"package/ComponentName"` — directory name + exported function name. No double-annotating child elements unless they are independently callable components.

---

## 6. WithBoundary Wrappers

Every gallery-eligible component gets a `*WithBoundary` function in `boundary.go`:

```go
func MyComponentWithBoundary(p MyProps) templ.Component {
    return devmode.ComponentBoundary("MyComponent", MyComponent(p), map[string]any{
        "field1": p.Field1,
        "field2": p.Field2,
    })
}
```

Signature matches the component function exactly. The third map argument drives `data-props` serialization.

---

## 7. Children

`{ children... }` for slotted content. `item.Content` for data-driven lists (Tabs, Accordion).

```templ
// Slotted — component wraps arbitrary content
templ Card(title string) {
    <div class="card">
        <h3>{ title }</h3>
        { children... }
    </div>
}

// Data-driven — items carry their own content
templ Tabs(props TabsProps) {
    for _, item := range props.Items {
        <div>{ @item.Content }</div>
    }
}
```

---

## 8. Default Values

Use Go blocks with sentinel defaults:

```templ
{{ submitText := props.SubmitText }}
{{ if submitText == "" { submitText = "Save" } }}
```

Sentinels: empty string for strings, `0` for numeric (not negative), `<= 0` for rows.

---

## 9. Alpine.js Variants

Conventions for every Alpine component:

```templ
templ AlpineFoo(props FooProps) {
    <div { alpine.XData(alpine.Toggle(false))... }   // outer: x-data, NO display:none
          { alpine.Escape("open = false")... }
          { devmode.Attrs(ctx, "package/AlpineFoo")... }>
        <button { alpine.On("click", "open = !open")... }>Trigger</button>
        <div x-show="open"                            // inner: x-show, x-transition, display:none
             { alpine.Transition()... }
             style="display:none">
            { children... }
        </div>
    </div>
}
```

- Outer div: `x-data` + `Escape` handler + `devmode.Attrs`. **Never** `style="display:none"` on outer.
- Inner toggled element: `x-show="open"` + `alpine.Transition()` + `style="display:none"`.
- `x-trap="open"` for modals with focus trapping.
- `@click.outside="open = false"` on the content panel (not the outer wrapper).
- **`x-data` must declare every identifier the markup reads or assigns.** Alpine evaluates `x-show`, `x-model` and `x-bind:*` during init, so a property missing from the state object is a hard `ReferenceError` — the binding silently never works. `alpine.Toggle` only provides `open`; a component that tracks something else (`picked`, `show`, `tab`, …) needs a state constructor that declares it (`alpine.PaletteState`, `alpine.TabState`, …).
- **`x-bind:*` values are JavaScript expressions, never literal data.** To emit a static attribute (`name`, `id`, `href`) use the plain Templ attribute — `x-bind:name={ props.ID }` compiles `demo-palette` as `demo - palette` and throws `demo is not defined`.

---

## 10. HTMX Method Dispatch

When a component supports multiple HTTP methods, use this pattern:

```templ
if props.Method == "put" {
    hx-put={ props.Action }
} else if props.Method == "patch" {
    hx-patch={ props.Action }
} else {
    hx-post={ props.Action }
}
```

---

## 11. Boolean Attributes

Use `checked?={ boolVal }` or `disabled?={ boolVal }` — never raw `if` blocks or string comparisons:

```templ
// Right
<input type="checkbox" checked?={ props.Checked } disabled?={ props.Disabled }/>

// Wrong
if props.Checked { checked }
<input type="checkbox" checked?={ props.Value == "true" }/>
```

---

## 12. Borders — the line ladder

Every border, divider and rule comes from the **line ladder** in
`components/css/custom.css` (the only go-daisy stylesheet consumers load, which is
why it is declared there — `assets/app.css` imports that file, so the library
registers it too). Three named steps replace the three ad-hoc tokens
(`--color-base-200`, `--color-base-300`, `--color-base-content` at 10%) that
used to compete for the same visual role.

| Token | Value | Use it for |
|---|---|---|
| `border-line` | `--color-base-200` | **Default.** The outline of a surface that sits on the page (card, panel, list, table wrapper), a region divider *inside* a surface, and a floating surface's edge (popover, menu, modal) — there the shadow carries the elevation. |
| `border-line-strong` | `--color-base-300` | The emphasised step. Control outlines (input, outline button, chip, swatch, editor wrapper), separator bars and keyboard chips, rules that must assert over scrolling content (sticky bar), and edges drawn on an opaque `bg-base-200` fill. |
| `border-line-translucent` | `--color-base-content` at 10% | A translucent hairline for a surface the surface ramp cannot serve — a `bg-base-200/40` fill, where `--color-base-200` itself nearly vanishes (soft box, disclosure, repo-doc callouts). |

```templ
// Right — card outline, divider, popover edge
class="card bg-base-100 edge-line shadow-sm"
class="px-4 py-3 border-b border-line"
class="dropdown-content bg-base-100 edge-line rounded-box shadow-lg"

// Right — control outline, separator bar, hairline on a tinted fill
class="input input-sm w-full border-line-strong"
class="border-line-strong h-7 w-px border-e"
class="rounded-box border border-line-translucent bg-base-200/40"

// Wrong — raw base tokens; no ladder to match against
class="card bg-base-100 border border-base-300"
class="rounded-box border border-base-content/10"
```

### Width comes from the theme

A ladder step names a *colour* and never a width, so it is always paired with
one. On a **surface outline** — card, banner, menu panel, accordion item — that
width is the theme's `--border` token, which daisyUI raises to `2px` under
`prefers-contrast: more`. Use `edge-line` / `edge-line-strong` there: one class
that sets width, style and colour together.

```templ
// Right — 1px normally, 2px for a high-contrast user
class="card bg-base-100 edge-line shadow-sm"
class="card bg-base-100 edge-line-strong"

// Wrong — Tailwind's `border` utility pins border-width to 1px, so the outline
// never thickens. This is exactly what `.card-border` used to get right.
class="card bg-base-100 border border-line"
```

`edge-line` is declared in the `components` layer of `components/css/custom.css`,
below `utilities`, so a caller can still override it (`border-0`, `border-2`,
`border-line-strong`) through `props.Class` or attrs — the same override path the
old `border border-line` pair offered. It falls back to `1px` for a theme that
does not define `--border`, and daisyUI's `forced-colors: active` rule on `.card`
still wins over it (unlayered, so forced-colors keeps `1px solid ButtonText`).

It is hand-written component CSS rather than a Tailwind `@utility` on purpose: a
generated utility only exists where a scanner sees the class name, and go-daisy's
class names live in Go strings and `.templ` files a consumer's own Tailwind build
does not scan.

Rules that follow:

- **Nothing quieter than `border-line` is a token.** Use Tailwind's opacity
  modifier (`border-line/50`) instead — a sub-visible step would be a phantom.
- **Directional forms exist**: `border-t-line`, `border-b-line`, `border-e-line`,
  `border-s-line`, and `divide-line` / `divide-line-strong` for row separators.
- **Surface outline vs internal rule.** `edge-line*` for the outline of a surface
  that sits on the page; `border border-line` for rules *inside* one. Only the
  former follows `--border` — daisyUI gave the token width to card outlines, not
  to every hairline in a layout.
- **A surface root, not a border anywhere.** `edge-line*` belongs on the element
  that owns one rounded outline for a card, panel, list, table wrapper,
  page-title bar or *floating* surface (dropdown content, popover, notification
  panel). Everything else stays on the fixed utility even when its colour is a
  ladder step: control outlines (input, trigger, swatch, editor shell, draggable
  handle, keyboard chip), hairlines (`border-t`, `divide-*`, a popover's arrow),
  decorative frames (mockup browser/window, dashed skeleton placeholders) and the
  translucent inset panels above. An edge that carries a state variant
  (`hover:border-*`, `focus-within:border-primary`, an `aria-pressed` colour rule)
  keeps its utilities: `edge-line` is plain component CSS and has no variants. A
  filled surface whose outline is `currentColor` rather than a ladder step — the
  highlighted `pricingCard` on the landing page — keeps `border` for the same
  reason: there is no `edge-line-current`.
- **Let the theme own the hue.** The steps derive from `--color-base-200`,
  `--color-base-300` and `--color-base-content`, so every theme (nord, dracula,
  flow) keeps its own colour — never hard-code an `oklch()` border.
- **`border-line-translucent` is not a positional slot.** In the consumer's
  `--color-muted-*` ladder `-bright` is the *most emphatic* step (content at 80%);
  here the third step is a translucent hairline that reads *weakest* and exists
  only for fills the opaque ramp cannot serve, so it is named for what it is. It
  is the most present edge of the three in `nord` (sRGB Δ32 vs Δ28 for
  `line-strong`) and `dracula` (Δ35 vs Δ18), but marginally *below* `line-strong`
  in `flow` (Δ39 vs Δ46). It is translucent, so it composites over whatever sits
  behind it: its contrast depends on the fill and it flips direction per theme
  (darkest rule on a light theme, lightest on a dark one). Reach for it only when
  the surface ramp can't do the job — otherwise stay on `border-line`.
- **Positional adjectives only where the ordering is real.** `line` <
  `line-strong` is a genuine strength ramp; `line-translucent` sits outside it
  (opaque vs translucent, not weaker vs stronger), which is why it carries a
  mechanism name instead of a third adjective.
- The ladder shares the consumer's foreground ladder
  (`--color-muted-subtle | -faint | muted | -strong | -bright`) *shape*: one
  `--color-*` scale in a plain `@theme` block, steps derived from the theme's own
  variables. It does **not** borrow the `-bright` slot adjective, precisely so the
  two scales can't be misread against each other.

## 13. Server-side table sorting

Use `table.SortContext` and `table.SortableTableHeader` for column sorting that
is driven by the request URL and rendered by the server. Do not hand-build sort
query strings or render a clickable `<th>` without `aria-sort`.

```go
sc := table.NewSortContext(r, "", "", table.SortNone)
// later, after collecting rows, apply the parsed params with your own
// stable comparator map (for example sort.SliceStable):
rows = sortRows(rows, sc.Params)
```

```templ
@table.SortableTableHeader(table.SortableTableHeaderProps{
    Label:    "CPU",
    Column:   "cpu",
    Sort:     sc,
    FirstDir: table.SortDesc,
    Class:    "text-right",
})
```

- `SortContext` clones the request query, so filters survive a sort click.
- `FirstDir` chooses the useful first click per column (`desc` for metrics,
  `asc` for names).
- The direction cycles first -> opposite -> none. `none` removes the sort
  parameters and returns to the handler's default row order.
- The component emits `aria-sort` only on the active column and uses a muted
  double-chevron on inactive columns so sorting is discoverable.
- Multi-table pages use a unique prefix per table, e.g.
  `table.NewSortContext(r, "guests_", "", table.SortNone)`.
- The handler is responsible for applying `SortParams` to the rows. Keep
  unknown metrics, empty text and zero times last in both directions.
- Sorting clears the `page` query param. `DataGrid` sort links target
  `#datagrid-<id>-body`; the handler must return that body-only fragment, the
  same contract as search and pagination.
