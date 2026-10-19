#!/usr/bin/env bash
# Render grouped release notes from commits between two refs.
#
# Usage: scripts/release-notes.sh <previous-tag> <head> [next-version]
#
# Sections are emitted only when non-empty, in this order:
# Breaking, Features, Fixes, Improvements, Deprecations, Maintenance.
set -euo pipefail

prev="${1:?usage: release-notes.sh <previous-tag> <head> [next-version]}"
head="${2:?usage: release-notes.sh <previous-tag> <head> [next-version]}"
next="${3:-$head}"

git rev-parse -q --verify "${prev}^{commit}" >/dev/null \
  || { echo "unknown previous ref: ${prev}" >&2; exit 1; }
git rev-parse -q --verify "${head}^{commit}" >/dev/null \
  || { echo "unknown head ref: ${head}" >&2; exit 1; }

server="${GITHUB_SERVER_URL:-https://github.com}"
repo="${GITHUB_REPOSITORY:-emergent-company/go-daisy}"

breaking=(); features=(); fixes=(); improvements=(); deprecations=(); maintenance=()

emit_entry() {
  # $1 = subject, $2 = body
  local subject="$1" body="$2" pr="" clean=""
  if [[ "$subject" =~ \(#([0-9]+)\)$ ]]; then
    pr=" (#${BASH_REMATCH[1]})"
  fi
  clean="$(printf '%s' "$subject" | sed -E 's/ \(#[0-9]+\)$//')"

  if [[ "$clean" =~ ^[A-Za-z]+(\([^\)]*\))?!:\  ]] \
    || printf '%s' "$body" | grep -q 'BREAKING CHANGE'; then
    breaking+=("- ${clean}${pr}")
  elif printf '%s' "$body" | grep -q 'Deprecated:'; then
    deprecations+=("- ${clean}${pr}")
  elif [[ "$clean" =~ ^feat(\([^\)]*\))?:\  ]]; then
    features+=("- ${clean}${pr}")
  elif [[ "$clean" =~ ^fix(\([^\)]*\))?:\  ]]; then
    fixes+=("- ${clean}${pr}")
  elif [[ "$clean" =~ ^(perf|refactor)(\([^\)]*\))?:\  ]]; then
    improvements+=("- ${clean}${pr}")
  else
    maintenance+=("- ${clean}${pr}")
  fi
}

count=0
while IFS= read -r sha; do
  [ -n "$sha" ] || continue
  subject="$(git log -1 --pretty=%s "$sha")"
  body="$(git log -1 --pretty=%b "$sha")"
  emit_entry "$subject" "$body"
  count=$((count + 1))
done < <(git log --no-merges --reverse --pretty=tformat:'%H' "${prev}..${head}")

[ "$count" -gt 0 ] || { echo "no commits between ${prev} and ${head}" >&2; exit 1; }

section() {
  local title="$1"; shift
  [ "$#" -gt 0 ] || return 0
  printf '\n### %s\n\n' "$title"
  printf '%s\n' "$@"
}

section "Breaking" "${breaking[@]}"
section "Features" "${features[@]}"
section "Fixes" "${fixes[@]}"
section "Improvements" "${improvements[@]}"
section "Deprecations" "${deprecations[@]}"
section "Maintenance" "${maintenance[@]}"

printf '\n**Full changelog:** %s/%s/compare/%s...%s\n' "$server" "$repo" "$prev" "$next"
