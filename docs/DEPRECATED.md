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
5. **Deprecation of a boundary pair.** When a component has a
   `*WithBoundary` gallery wrapper, both the component and its wrapper carry the
   marker and appear in the inventory independently.

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
