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

---

## 10. Stimulus.js Variants

Identical HTML structure to Alpine variants. Only data attributes differ:

```templ
templ StimulusFoo(props FooProps) {
    <div { stimulus.Controller("foo")... }
         { stimulus.Action("keydown.escape", "foo", "close")... }
         { devmode.Attrs(ctx, "package/StimulusFoo")... }>
        <button { stimulus.Action("click", "foo", "toggle")... }
                data-foo-target="trigger">Trigger</button>
        <div data-foo-target="menu" style="display:none">
            { children... }
        </div>
    </div>
}
```

- Controller on outer wrapper.
- Actions on interactive elements.
- Targets for element references.
- Keyboard events: always add `keydown.escape`, `keydown.arrow-down`, `keydown.arrow-up`, `keydown.enter`.
- `aria-expanded="false"` hardcoded initially — controller manages the toggle.

---

## 11. HTMX Method Dispatch

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

## 12. Boolean Attributes

Use `checked?={ boolVal }` or `disabled?={ boolVal }` — never raw `if` blocks or string comparisons:

```templ
// Right
<input type="checkbox" checked?={ props.Checked } disabled?={ props.Disabled }/>

// Wrong
if props.Checked { checked }
<input type="checkbox" checked?={ props.Value == "true" }/>
```
