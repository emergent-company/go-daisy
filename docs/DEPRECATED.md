# Deprecation Policy

This document is the single source of truth for the deprecation lifecycle of
public symbols in go-daisy. Every symbol marked `Deprecated:` in source must be
listed in the inventory below, and every inventory entry must correspond to a
live `Deprecated:` marker. `task lint` enforces this via `scripts/lint-deprecated.sh`.

## Policy

1. **Mark first, remove later.** When a public symbol is superseded, annotate
   its doc comment with a `Deprecated:` line that names the replacement
   (including a usable call signature or props example), and add the symbol to
   the inventory below in the same change.
2. **Never deprecate without a path forward.** A `Deprecated:` marker must name
   the replacement. If no replacement exists, do not mark the symbol
   deprecated — propose the replacement first.
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
