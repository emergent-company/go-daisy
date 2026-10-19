#!/usr/bin/env bash
# Compute release versions from the highest reachable semver tag.
#
# Usage:
#   scripts/release-version.sh --current
#   scripts/release-version.sh --bump patch|minor|major
#   scripts/release-version.sh --tag vX.Y.Z
#
# The tag form validates the tag, rejects an existing tag, and rejects a version
# that is not newer than the current one.
set -euo pipefail

usage() {
  echo "usage: $0 --current | --bump patch|minor|major | --tag vX.Y.Z" >&2
  exit 2
}

mode="" value=""
while [ $# -gt 0 ]; do
  case "$1" in
    --current) mode=current; shift ;;
    --bump) [ $# -ge 2 ] || usage; mode=bump; value=$2; shift 2 ;;
    --tag)  [ $# -ge 2 ] || usage; mode=tag;  value=$2; shift 2 ;;
    *) usage ;;
  esac
done
[ -n "$mode" ] || usage

latest="$(git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname \
  | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n1 || true)"
[ -n "$latest" ] || { echo "no semver tag found" >&2; exit 1; }

if [ "$mode" = current ]; then
  printf '%s\n' "$latest"
  exit 0
fi

if [ "$mode" = tag ]; then
  [[ "$value" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] \
    || { echo "invalid tag: $value (expected vX.Y.Z)" >&2; exit 1; }
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

if [ "$(printf '%s\n%s\n' "$latest" "$next" | sort -V | tail -n1)" != "$next" ]; then
  echo "${next} is not newer than ${latest}" >&2
  exit 1
fi

printf '%s\n' "$next"
