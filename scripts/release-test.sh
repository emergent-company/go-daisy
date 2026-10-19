#!/usr/bin/env bash
# Exercise release-version.sh and release-notes.sh in throwaway repositories.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION="$HERE/release-version.sh"
NOTES="$HERE/release-notes.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

pass=0
fail=0
ok() { pass=$((pass + 1)); }
bad() { echo "FAIL: $*" >&2; fail=$((fail + 1)); }

assert_eq() {
  if [ "$1" = "$2" ]; then ok; else bad "$3 (got '$1', want '$2')"; fi
}
assert_contains() {
  case "$1" in *"$2"*) ok ;; *) bad "$3 (missing '$2')" ;; esac
}
assert_fail() {
  local desc="$1"; shift
  if "$@" >/dev/null 2>&1; then bad "$desc (expected failure)"; else ok; fi
}

mk_repo() {
  local name="$1"
  REPO="$TMP/$name"
  mkdir -p "$REPO" && cd "$REPO" || exit 1
  git init -q -b main
  git config user.email test@example.com
  git config user.name test
  printf 'seed\n' > seed
  git add seed && git commit -qm "chore: seed"
}

# --- version arithmetic and explicit tags -----------------------------------
mk_repo version
git tag v1.2.3
assert_eq "$("$VERSION" --current)" "v1.2.3" "current tag"
assert_eq "$("$VERSION" --bump patch)" "v1.2.4" "patch bump"
assert_eq "$("$VERSION" --bump minor)" "v1.3.0" "minor bump"
assert_eq "$("$VERSION" --bump major)" "v2.0.0" "major bump"
assert_eq "$("$VERSION" --tag v2.1.0)" "v2.1.0" "explicit newer tag"
assert_fail "existing tag" "$VERSION" --tag v1.2.3
assert_fail "older tag" "$VERSION" --tag v1.2.2
assert_fail "non-canonical tag" "$VERSION" --tag v01.2.4
assert_fail "tag without v" "$VERSION" --tag 2.1.0

# Tags on unmerged branches must not move the next version.
git checkout -q -b side
git commit -q --allow-empty -m "feat: side work"
git tag v9.0.0
git checkout -q main
assert_eq "$("$VERSION" --current)" "v1.2.3" "unreachable tag ignored"
assert_eq "$("$VERSION" --bump minor)" "v1.3.0" "bump ignores unreachable tag"

# --- no semver tags ---------------------------------------------------------
mk_repo no-tags
assert_fail "current without tags" "$VERSION" --current
assert_fail "bump without tags" "$VERSION" --bump patch

# --- grouped notes ----------------------------------------------------------
mk_repo notes
git tag v0.1.0
git commit -q --allow-empty -m "feat(ui): add a thing (#1)"
git commit -q --allow-empty -m "fix: correct a thing (#2)"
git commit -q --allow-empty -m "feat(api)!:remove the old call (#3)"
git commit -q --allow-empty -m "chore: tidy docs (#4)"
git commit -q --allow-empty -m "feat(ui): deprecate OldThing (#5)" \
  -m "Deprecated: OldThing, use NewThing."
git commit -q --allow-empty -m "feat:keep me (#6)"
git commit -q --allow-empty -m "feat: prose (#7)" \
  -m "We mention Deprecated: in prose."
notes="$("$NOTES" v0.1.0 HEAD v0.2.0)"

assert_contains "$notes" "### Breaking" "breaking section"
assert_contains "$notes" "- feat(api)!:remove the old call (#3)" "bang without space is breaking"
assert_contains "$notes" "### Features" "features section"
assert_contains "$notes" "- feat:keep me (#6)" "feature without space after colon"
assert_contains "$notes" "- feat(ui): deprecate OldThing (#5)" "deprecating feature stays in Features"
assert_contains "$notes" "### Deprecations" "deprecations section"
assert_contains "$notes" "### Maintenance" "maintenance section"
assert_contains "$notes" "- chore: tidy docs (#4)" "maintenance entry"
assert_contains "$notes" "compare/v0.1.0...v0.2.0" "changelog link"

# An empty range must fail loudly.
assert_fail "empty range" "$NOTES" v0.1.0 v0.1.0 v0.2.0

echo "release script tests: ${pass} passed, ${fail} failed"
[ "$fail" -eq 0 ]
