#!/usr/bin/env bash
# Compute release versions from the highest semver tag reachable from HEAD.
#
# Usage:
#   scripts/release-version.sh --current
#   scripts/release-version.sh --before vX.Y.Z
#   scripts/release-version.sh --bump patch|minor|major
#   scripts/release-version.sh --tag vX.Y.Z
#
# The tag form validates the tag, rejects an existing tag, and rejects a version
# that is not newer than the current one. Tags on branches that are not merged
# into HEAD are ignored.
set -euo pipefail

usage() {
  echo "usage: $0 --current | --before vX.Y.Z | --bump patch|minor|major | --tag vX.Y.Z" >&2
  exit 2
}

# Canonical segments only: no leading zeros, so v01.2.3 cannot coexist with v1.2.3.
semver_re='^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$'

mode="" value=""
while [ $# -gt 0 ]; do
  case "$1" in
    --current) mode=current; shift ;;
    --before) [ $# -ge 2 ] || usage; mode=before; value=$2; shift 2 ;;
    --bump) [ $# -ge 2 ] || usage; mode=bump; value=$2; shift 2 ;;
    --tag)  [ $# -ge 2 ] || usage; mode=tag;  value=$2; shift 2 ;;
    *) usage ;;
  esac
done
[ -n "$mode" ] || usage

# ver_gt A B: true when A is strictly newer than B. Same numeric comparison for
# every caller, so tag selection and the "newer than" guard cannot disagree.
ver_gt() {
  local a="${1#v}" b="${2#v}" am aj ap bm bj bp
  IFS=. read -r am aj ap <<< "$a"
  IFS=. read -r bm bj bp <<< "$b"
  (( am > bm )) && return 0
  (( am < bm )) && return 1
  (( aj > bj )) && return 0
  (( aj < bj )) && return 1
  (( ap > bp ))
}

# Highest semver tag reachable from the given ref.
latest_from() {
  git tag --merged "$1" --list 'v*' --sort=-v:refname \
    | grep -E "$semver_re" | head -n1 || true
}

# Highest semver tag strictly older than the target and reachable from it.
# Version order, not topology: this repo's v0.16.5 tag sits on an older commit
# than v0.16.4, which is why the workflow prefers release order for backfill.
latest_before() {
  local ref="$1" target="$2" tag
  while IFS= read -r tag; do
    [ -n "$tag" ] || continue
    [ "$tag" = "$target" ] && continue
    if ver_gt "$target" "$tag"; then
      printf '%s\n' "$tag"
      return 0
    fi
  done < <(git tag --merged "$ref" --list 'v*' --sort=-v:refname | grep -E "$semver_re" || true)
  return 1
}

if [ "$mode" = before ]; then
  [[ "$value" =~ $semver_re ]] \
    || { echo "invalid tag: $value (expected canonical vX.Y.Z)" >&2; exit 1; }
  git rev-parse -q --verify "${value}^{commit}" >/dev/null \
    || { echo "unknown tag: $value" >&2; exit 1; }
  if ! previous="$(latest_before "${value}^{commit}" "$value")"; then
    echo "no semver tag older than ${value} reachable from it" >&2
    exit 1
  fi
  printf '%s\n' "$previous"
  exit 0
fi

latest="$(latest_from HEAD)"
[ -n "$latest" ] || { echo "no semver tag reachable from HEAD" >&2; exit 1; }

if [ "$mode" = current ]; then
  printf '%s\n' "$latest"
  exit 0
fi

if [ "$mode" = tag ]; then
  [[ "$value" =~ $semver_re ]] \
    || { echo "invalid tag: $value (expected canonical vX.Y.Z)" >&2; exit 1; }
  next="$value"
else
  case "$value" in patch|minor|major) ;; *) usage ;; esac
  IFS=. read -r major minor patch <<< "${latest#v}"
  case "$value" in
    major) major=$((major + 1)); minor=0; patch=0 ;;
    minor) minor=$((minor + 1)); patch=0 ;;
    patch) patch=$((patch + 1)) ;;
  esac
  next="v${major}.${minor}.${patch}"
fi

if git rev-parse -q --verify "refs/tags/${next}" >/dev/null; then
  echo "tag already exists: ${next}" >&2
  exit 1
fi

if ! ver_gt "$next" "$latest"; then
  echo "${next} is not newer than ${latest}" >&2
  exit 1
fi

printf '%s\n' "$next"
