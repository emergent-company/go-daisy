#!/usr/bin/env bash
# Guard: every `Deprecated:` marker in source must be registered in
# docs/DEPRECATED.md, and every inventory entry must have a live marker.
#
# Fails when:
#   - a `Deprecated:` marker exists in source but is missing from the inventory
#     (a symbol was deprecated without being registered), or
#   - an inventory entry has no matching marker (a stale entry left after the
#     symbol was deleted).
#
# The inventory is the single source of truth; keep the two in sync in the
# same change. See docs/DEPRECATED.md for the full policy.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

inventory_file="docs/DEPRECATED.md"

if [ ! -f "$inventory_file" ]; then
  echo "error: $inventory_file not found" >&2
  exit 1
fi

# --- 1. Extract `pkg.Symbol` for every live `Deprecated:` marker ---
# Scan .templ and .go sources, skipping generated files (*_templ.go, *_gen.go).
# For each line containing "Deprecated:", read forward a short window for the
# next top-level declaration (templ/func/type/const) and record its name.
# Uses awk (mawk-compatible: no 3-arg match, no GNU-only features).
mapfile -t src_symbols < <(
  find . \
    -path ./node_modules -prune -o \
    -path ./.git -prune -o \
    -type f \( -name '*.templ' -o -name '*.go' \) \
    ! -name '*_templ.go' \
    ! -name '*_gen.go' \
    -print0 |
    xargs -0 -r awk '
      FNR == 1 { pkg = ""; pending = 0; depth = 0 }

      # Capture the package name from the first package clause in each file.
      pkg == "" && /^[[:space:]]*package[[:space:]]+[A-Za-z_]/ {
        p = $0
        sub(/^[[:space:]]*package[[:space:]]+/, "", p)
        sub(/[[:space:]].*/, "", p)
        pkg = p
      }

      /Deprecated:/ { pending = 1; depth = 0; next }

      pending {
        if ($0 ~ /^[[:space:]]*(templ|func|type|const)[[:space:]]+[A-Za-z_(]/) {
          s = $0
          sub(/^[[:space:]]*(templ|func|type|const)[[:space:]]+/, "", s)
          if (s ~ /^\(/) { sub(/^\([^)]*\)[[:space:]]*/, "", s) }
          split(s, t, /[^A-Za-z0-9_]/)
          if (t[1] ~ /^[A-Za-z_][A-Za-z0-9_]*$/ && pkg != "") {
            print pkg "." t[1]
          }
          pending = 0
        } else {
          depth++
          if (depth > 30) {
            print "WARN: no declaration after Deprecated marker: " FILENAME ":" FNR > "/dev/stderr"
            pending = 0
          }
        }
      }
    '
)

# --- 2. Extract registered entries from the inventory block ---
mapfile -t inv_symbols < <(
  awk '
    /<!-- DEPRECATED-INVENTORY-START -->/ { inblock = 1; next }
    /<!-- DEPRECATED-INVENTORY-END -->/ { inblock = 0 }
    inblock && $0 ~ /^- `[^`]+`/ {
      s = $0
      sub(/^- `/, "", s)
      sub(/`.*/, "", s)
      print s
    }
  ' "$inventory_file"
)

# --- 3. Diff the two sorted, de-duplicated sets ---
src_sorted="$(printf '%s\n' "${src_symbols[@]}" | sort -u)"
inv_sorted="$(printf '%s\n' "${inv_symbols[@]}" | sort -u)"

only_src="$(comm -23 <(printf '%s\n' "$src_sorted") <(printf '%s\n' "$inv_sorted"))"
only_inv="$(comm -13 <(printf '%s\n' "$src_sorted") <(printf '%s\n' "$inv_sorted"))"

status=0

if [ -n "$only_src" ]; then
  echo "error: Deprecated: markers missing from $inventory_file:" >&2
  printf '  - `%s`\n' $only_src >&2
  echo "Add each symbol to the inventory (see docs/DEPRECATED.md)." >&2
  status=1
fi

if [ -n "$only_inv" ]; then
  echo "error: stale inventory entries in $inventory_file (marker removed or never existed):" >&2
  printf '  - `%s`\n' $only_inv >&2
  echo "Remove the entry or restore the Deprecated: marker." >&2
  status=1
fi

if [ "$status" -ne 0 ]; then
  exit "$status"
fi

echo "ok: deprecation inventory in sync ($(printf '%s\n' "$src_sorted" | grep -c . || true) symbols)"
