#!/usr/bin/env bash
# Render grouped release notes from commits between two refs.
#
# Usage: scripts/release-notes.sh <previous-tag> <head> [next-version]
#
# Sections are emitted only when non-empty, in this order:
# Breaking, Features, Fixes, Improvements, Deprecations, Maintenance.
#
# A commit body line that starts with `Deprecated:` adds the commit to
# Deprecations in addition to its type section, so a feature does not disappear
# from Features just because it deprecates something.
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
  local subject="$1" body="$2" pr="" clean="" entry="" type="" bang=""
  if [[ "$subject" =~ \(#([0-9]+)\)$ ]]; then
    pr=" (#${BASH_REMATCH[1]})"
  fi
  clean="$(printf '%s' "$subject" | sed -E 's/ \(#[0-9]+\)$//')"
  entry="- ${clean}${pr}"

  if [[ "$clean" =~ ^([A-Za-z]+)(\([^\)]*\))?(!)?: ]]; then
    type="$(printf '%s' "${BASH_REMATCH[1]}" | tr '[:upper:]' '[:lower:]')"
    bang="${BASH_REMATCH[3]}"
  fi

  if [ -n "$bang" ] || printf '%s' "$body" | grep -qE '^BREAKING CHANGE:'; then
    breaking+=("$entry")
  elif [ "$type" = feat ]; then
    features+=("$entry")
  elif [ "$type" = fix ]; then
    fixes+=("$entry")
  elif [ "$type" = perf ] || [ "$type" = refactor ]; then
    improvements+=("$entry")
  else
    maintenance+=("$entry")
  fi

  if printf '%s' "$body" | grep -qE '^Deprecated:'; then
    deprecations+=("$entry")
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

[ "$count" -gt 0 ] || { echo "no commits between ${prev} and ${head}; nothing to release" >&2; exit 1; }

section() {
  local title="$1"; shift
  [ "$#" -gt 0 ] || return 0
  printf '\n### %s\n\n' "$title"
  printf '%s\n' "$@"
}

# ${arr[@]+...} keeps empty arrays safe under `set -u` on bash 3.2 (macOS).
section "Breaking" ${breaking[@]+"${breaking[@]}"}
section "Features" ${features[@]+"${features[@]}"}
section "Fixes" ${fixes[@]+"${fixes[@]}"}
section "Improvements" ${improvements[@]+"${improvements[@]}"}
section "Deprecations" ${deprecations[@]+"${deprecations[@]}"}
section "Maintenance" ${maintenance[@]+"${maintenance[@]}"}

printf '\n**Full changelog:** %s/%s/compare/%s...%s\n' "$server" "$repo" "$prev" "$next"
