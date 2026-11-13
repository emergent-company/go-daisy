#!/usr/bin/env bash
#
# Guard the one invariant a consumer of the brand theme cannot see.
#
# `theme-memory.css` and `tokens.css` both declare `--color-muted`; the theme's
# 50% step must win, which only holds because consumers import the theme LAST.
# The Go tests can assert the source says 50%, but nothing there compiles CSS, so
# a regression in the documented order (or in Tailwind's handling of `@plugin`
# inside an imported file) would only surface as wrong colours in the two apps.
#
# This compiles the exact entry documented in theme-memory.css and asserts the
# resolved values. Runs in the CSS job of pr-checks.yml, which has node_modules.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

if [[ ! -x ./node_modules/.bin/tailwindcss ]]; then
  echo "error: tailwind not installed; run npm ci first" >&2
  exit 2
fi

tmp=".theme-css-check"
rm -rf "$tmp"
mkdir -p "$tmp"
trap 'rm -rf "$tmp"' EXIT

cat >"$tmp/app.css" <<'CSS'
@import "tailwindcss" source(none);
@source inline("text-muted");
@source inline("rounded-lg");
@plugin "daisyui" { themes: dark --default; }
@import "../components/css/tokens.css";
@import "../components/css/custom.css";
@import "../components/css/theme-memory.css";
CSS

if ! ./node_modules/.bin/tailwindcss -i "$tmp/app.css" -o "$tmp/out.css" >"$tmp/build.log" 2>&1; then
  echo "error: the documented import order (tokens -> custom -> theme-memory) does not compile" >&2
  cat "$tmp/build.log" >&2
  exit 1
fi

failures=0
expect() {
  local what="$1" needle="$2"
  if ! grep -qF -- "$needle" "$tmp/out.css"; then
    echo "error: $what — expected the compiled CSS to contain: $needle" >&2
    failures=$((failures + 1))
  fi
}

# The brand palette is registered from an *imported* file: a `@plugin
# "daisyui/theme"` that Tailwind stopped honouring outside the entry stylesheet
# would leave daisyUI's stock dark theme in place, with no error.
expect "brand palette not registered from the imported theme" "oklch(0.74 0.135 78)"

# The radius bridge: `rounded-lg` must resolve through the semantic token, not
# the literal Tailwind scale, or every go-daisy radius ignores the theme. Read
# the utility's own block — daisyUI's rules emit `var(--radius-box)` in several
# places, so a bare grep here would pass even with the bridge deleted.
radius_block="$(grep -A3 -F '.rounded-lg {' "$tmp/out.css" || true)"
if [[ -z "$radius_block" ]]; then
  echo "error: .rounded-lg was not generated; the check cannot verify the radius bridge" >&2
  failures=$((failures + 1))
elif [[ "$radius_block" != *"var(--radius-box)"* ]]; then
  echo "error: the radius bridge is missing — .rounded-lg does not resolve through the theme:" >&2
  printf '%s\n' "$radius_block" >&2
  failures=$((failures + 1))
fi

# Direct `var()` consumers (the gateway's own rules) read the theme's steps from
# :root, which `@theme inline` does not emit on its own.
expect "muted scale missing from :root" "--color-muted-strong:"

# Import order, read off the utility rather than the :root block: whichever
# `@theme` block comes last owns the generated value. tokens.css would make it
# 60%; the theme must make it 50%.
muted_block="$(grep -A4 -F '.text-muted {' "$tmp/out.css" || true)"
if [[ -z "$muted_block" ]]; then
  echo "error: .text-muted was not generated; the check cannot verify import order" >&2
  failures=$((failures + 1))
elif [[ "$muted_block" != *"50%"* ]] || [[ "$muted_block" == *"60%"* ]]; then
  echo "error: theme-memory.css must be imported after tokens.css — .text-muted did not resolve to the theme's 50% step:" >&2
  printf '%s\n' "$muted_block" >&2
  failures=$((failures + 1))
fi

if [[ "$failures" -ne 0 ]]; then
  echo "error: $failures theme-memory assertion(s) failed" >&2
  exit 1
fi

echo "ok: theme-memory.css compiles in the documented order and wins on --color-muted"
