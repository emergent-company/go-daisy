#!/usr/bin/env bash
# Guard: every committed Go source file must be gofmt-clean.
#
# Fails when `gofmt -l` reports any file (i.e. its formatting differs from
# gofmt's canonical output). This covers hand-written sources and generated
# files alike — `*_templ.go` and `*_gen.go` are committed and must be formatted
# at their source of generation, not hand-fixed after the fact (see #68).
#
# Fail-closed design (see #60, #68): the gate must never report a clean tree
# having scanned nothing. It therefore:
#   - fails with exit 127 and a distinct message when `gofmt` is missing from
#     PATH (the tool itself is unavailable);
#   - builds the file list first and asserts it is non-empty (a broken `find`
#     prunes everything → fail);
#   - self-checks against a known-good file that must be in the scan set
#     (over-pruning the list → fail);
#   - treats any non-zero `gofmt` exit (parse error, I/O error, unreadable
#     path) as a scan failure rather than a clean pass.
# Only when `gofmt` runs to completion and reports zero files does the gate
# print `ok` and exit 0.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

# --- 0. Tool availability: fail closed if gofmt is unreachable ---
if ! command -v gofmt >/dev/null 2>&1; then
  echo "error: lint:fmt cannot run: gofmt not found on PATH" >&2
  exit 127
fi

# --- 1. Build the scan list (all .go files, pruned of node_modules/.git) ---
# Written to a temp file first so find's exit status is observable (a command
# substitution would strip NUL bytes and mask a scan error).
list_file="$(mktemp)"
trap 'rm -f "$list_file"' EXIT

if ! find . \
      -path ./node_modules -prune -o \
      -path ./.git -prune -o \
      -type f -name '*.go' \
      -print0 > "$list_file"; then
  echo "error: lint:fmt scan failed: find could not enumerate Go files" >&2
  exit 2
fi

mapfile -d '' files < "$list_file"
if [ "${#files[@]}" -eq 0 ]; then
  echo "error: lint:fmt scanned 0 files; refusing to report a clean tree" >&2
  exit 2
fi

# Self-check: a known-good source file must be present in the scan set, so a
# broken or over-pruned `find` fails instead of passing vacuously.
sentinel='./components/ui/ui.go'
found=0
for f in "${files[@]}"; do
  [ "$f" = "$sentinel" ] && { found=1; break; }
done
if [ "$found" -ne 1 ]; then
  echo "error: lint:fmt self-check failed: known-good file '$sentinel' not in scan set" >&2
  exit 2
fi

# --- 2. Run gofmt -l over the whole list in one invocation ---
# `gofmt -l` prints the name of each file whose formatting differs and exits 0
# on success; a non-zero exit means the scan itself errored (parse error, I/O,
# unreadable path), which must fail closed rather than look clean.
unformatted="$(gofmt -l -- "${files[@]}")" || {
  echo "error: lint:fmt scan failed: gofmt exited $? over the file list" >&2
  exit 2
}

if [ -n "$unformatted" ]; then
  printf '%s\n' "$unformatted"
  echo "error: the above files are not gofmt-formatted; run 'gofmt -w' on them (or fix the generator that emits them)" >&2
  exit 1
fi

echo "ok: all $(printf '%s\n' "${#files[@]}") Go files are gofmt-formatted"
