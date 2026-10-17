# Deprecation Policy

This document is the single source of truth for the deprecation lifecycle of
public symbols in go-daisy. Every symbol marked `Deprecated:` in source must be
listed in the inventory below, and every inventory entry must correspond to a
live `Deprecated:` marker. `task lint` enforces this via
`scripts/lint-deprecated.sh`, which also fails when a `.templ` source *calls* a
deprecated symbol (golangci-lint never reads `.templ`, so those callers are
otherwise invisible — see #67).

## Policy

1. **Mark first, remove later.** When a public symbol is superseded, annotate
   its doc comment with a `Deprecated:` line that names the replacement
   (including a usable call signature or props example), and add the symbol to
   the inventory below in the same change.
2. **Never deprecate without a path forward.** A `Deprecated:` marker must name
   the replacement. If no replacement exists, do not mark the symbol
   deprecated — propose the replacement first. The named replacement must also
   be *faithful*: see "Choosing the named replacement" below.
3. **Do not silently delete.** A deprecated symbol may be removed only after a
   release cycle that contained the marker, and only once the inventory shows
   the removal is intentional. Deletion and inventory removal happen together.
4. **The inventory is authoritative.** `scripts/lint-deprecated.sh` fails CI
   when:
   - a `Deprecated:` marker in source is missing from the inventory (forgotten
     registration), or
   - an inventory entry no longer has a matching marker (stale entry left after
     deletion).
5. **Deprecation of a component with a boundary wrapper.** The *canonical*
   `*WithBoundary` wrapper stays **undeprecated** even when the inner component
   is deprecated: only the inner component is marked and inventoried. For
   example, `ui.Avatar`, `ui.CardCompact`, `layout.Page`, and `nav.SimpleTabs`
   each carry `Deprecated:`, while their canonical wrappers
   (`AvatarWithBoundary`, `CardCompactWithBoundary`, `PageWithBoundary`,
   `SimpleTabsWithBoundary`) carry no marker and are absent from the inventory.
   Only *alias* wrappers — deprecated shims that delegate to a different
   component or wrapper (e.g. `ui.PaginationCircleWithBoundary`,
   `ui.IconSpanColoredWithBoundary`, `ui.PersonChipWithBoundary`, and the
   `form.*WithBoundary` aliases) — carry their own `Deprecated:` marker and are
   inventoried alongside the component they alias.

### Choosing the named replacement

The shorthand components (`form.Toggle`, `form.Checkbox`) and their `Form*`
successors (`form.FormToggle`, `form.FormCheckbox`) are not interchangeable:

| | shorthand (`form.Toggle`) | `Form*` (`form.FormToggle`) | building block (`form.ToggleInput`) |
|---|---|---|---|
| colour | forces `toggle-primary` | neutral `toggle` | caller's `Class`, none forced |
| DOM | `<label class="label">` + input | `<div class="form-control">` + `<label>` + input | bare `<input>` |
| label text | rendered by the component | rendered by the component | rendered by the caller |
| error slot | no | yes (`FieldError`) | no |

The emphasis drop from `toggle-primary` to neutral is a deliberate part of the
`Form*` design (neutral fields, emphasis added by the page), so it is **not** a
faithful migration for a site whose switch is visibly primary today — and the
`form-control` wrapper (a class with no rules in the shipped DaisyUI 5
stylesheet) is noise inside a `flex items-center justify-between` row whose
title and description the caller already renders.

When migrating a deprecated shorthand, choose by layout:

- **Does the component render its own visible label, stacked above/next to the
  field?** Then the `Form*` replacement is right; accept the neutral colour or
  restyle deliberately.
- **Does the caller render the visible title and description, with the control
  as a right-hand flex item?** Then use the bare building block
  (`ToggleInput(ToggleInputProps{...})`) and pass the colour explicitly
  (`Class: "toggle toggle-primary"`) to keep the page's emphasis unchanged.
  `pages/settings.templ` is the reference migration.

`form.Checkbox` has no bare counterpart yet (`ToggleInput` covers toggles only),
so its `FormCheckbox` path is still the only forward option; a row-layout
checkbox currently means accepting the neutral colour or hand-rolling the
input.

## Inventory

<!-- DEPRECATED-INVENTORY-START -->
- `form.Checkbox`
- `form.CheckboxInput`
- `form.CheckboxInputWithBoundary`
- `form.CheckboxWithBoundary`
- `form.FormActionsWithLoading`
- `form.SelectInput`
- `form.SelectInputWithBoundary`
- `form.TextInput`
- `form.TextInputWithBoundary`
- `form.TextareaInput`
- `form.TextareaInputWithBoundary`
- `form.Toggle`
- `form.ToggleWithBoundary`
- `galleryruntime.CategoryBasics`
- `layout.Page`
- `nav.SimpleTabs`
- `ui.Avatar`
- `ui.CardCompact`
- `ui.IconSpanColored`
- `ui.IconSpanColoredWithBoundary`
- `ui.PaginationCircle`
- `ui.PaginationCircleWithBoundary`
- `ui.PersonChipWithBoundary`
<!-- DEPRECATED-INVENTORY-END -->

## Removed symbols

Symbols removed in issue #34 (PR #62), recorded so a downstream `go get`
upgrade surfaces the breakage explicitly instead of failing silently. 27 net
deletions: 12 `templ` components, 2 types, 13 `*WithBoundary` wrappers.

### `ui` (10 templ components, 1 type, 10 wrappers)

- Components: `AlertStyled`, `InlineAlert`, `InlineAlertWithIcon`,
  `AlpineThemeToggle`, `CompactFilterBar`, `Hover3DCardLayered`, `Pagination`
  (square alias), `PersonChip`, `PlaceholderAvatar`, `StatCardIconCorner`
- Types: `StatCardIconCornerItem`
- Wrappers: `AlertWithIconBoundary`, `AlertStyledWithBoundary`,
  `AlpineThemeToggleWithBoundary`, `AlpineToastWithBoundary`,
  `AnimatedGradientTextWithBoundary`, `ButtonGlowWithBoundary`,
  `Hover3DCardLayeredWithBoundary`, `SimpleButtonGlowWithBoundary`,
  `StatCardIconCornerWithBoundary`, `ThemeControllerBtnWithBoundary`

### `form` (2 templ components, 1 type, 1 wrapper)

- Components: `Label`, `LabelAlt`
- Types: `LabelProps`
- Wrappers: `LabelWithBoundary`

### `nav` (2 wrappers)

- Wrappers: `PageTitleMinimalWithBoundary`, `PageTitleEditorWithBoundary`

## Guard limitations

`scripts/lint-deprecated.sh` is a line-based heuristic over source text, not a
type-checker. The following are known robustness gaps, not current bugs — none
of these shapes exist in the tree today (all 23 live markers are
`// Deprecated:` on their own line, directly above a top-level
`templ`/`func`/`type`/`const` declaration). They are documented rather than
fixed here. The separate *fail-open* class — an environmental failure silently
passing the guard — was closed in #60: the guard now fails closed when the scan
cannot run, and self-checks that coverage has not collapsed.

- **Trailing marker on a declaration line** (`func Foo() {} // Deprecated: ...`)
  is silently missed — the marker must be on its own line.
- **Marker on a struct field or inside a `const (...)` block** is mis-attributed
  to the next top-level declaration, or skipped entirely.
- **Marker/declaration gap over 30 lines** warns and continues, so the symbol is
  missed rather than failing the build.

The `.templ` caller scan (step 4) is a line-based mini-lexer over `.templ`
sources, not a type-checker. It detects `alias.Symbol(` (render or expression
calls, alias resolved against the file's imports) and `@Symbol(` (render calls
within the symbol's own package). It masks `"…"` and `` `…` `` string literals
and `//` / `/* */` comments first, so a deprecated name appearing only in prose,
a doc example, or a URL is not flagged. Known gaps: single-quoted rune literals
are not tracked, and a deprecated call written inside a string literal is
invisible by design. Coverage collapse (an empty symbol set or `.templ` file
set, or a missing known-good sentinel) fails the guard with a distinct status,
so a broken scan never reports clean.
