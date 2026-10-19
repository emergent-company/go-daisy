#!/usr/bin/env bash
# Pick the newest published release that is reachable from the tag being
# backfilled. Candidates arrive in publication order (newest first).
#
# Usage:
#   scripts/release-base.sh vX.Y.Z <candidate-tag>...
#
# The workflow supplies candidates in publication order (newest first):
#   gh release list --limit 100 --exclude-drafts --json tagName \
#     --jq '.[] | .tagName'
set -euo pipefail

version="${1:?usage: release-base.sh vX.Y.Z <candidate-tag>...}"
shift

for candidate in "$@"; do
  [ -n "$candidate" ] || continue
  [ "$candidate" = "$version" ] && continue
  git rev-parse -q --verify "refs/tags/${candidate}^{commit}" >/dev/null 2>&1 || continue
  if git merge-base --is-ancestor "$candidate" "$version"; then
    printf '%s\n' "$candidate"
    exit 0
  fi
done

exit 1
