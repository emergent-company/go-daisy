# lint-deprecated-callers.awk
#
# Flag .templ source that *calls* a deprecated symbol. golangci-lint never
# reads .templ (it is not Go), and the generated *_templ.go that does hold the
# call is excluded from analysis, so a deprecated symbol can keep live callers
# in the render surface indefinitely (#67). This scanner closes that hole.
#
# Input (positional):
#   $1..: first file is the deprecated-symbol list (one `pkg.Symbol` per line,
#        derived by lint-deprecated.sh from live `Deprecated:` markers — the
#        same source of truth as docs/DEPRECATED.md), the rest are .templ files.
#
# Output (to stderr): `file:line: deprecated caller of pkg.Symbol` per hit.
# Exit status: 0 = no callers, 1 = callers found, 2 = scan error.
#
# Detection covers two call shapes:
#   - qualified   `alias.Symbol(`  (render `@` or expression), alias resolved
#     against that file's imports so a same-named symbol in a different package
#     (e.g. `alpine.Toggle` vs deprecated `form.Toggle`) is not a false hit;
#   - unqualified `@Symbol(`       (render call within the symbol's own package
#     only — the file's `package` clause must match).
#
# A mini-lexer masks `"…"` and `` `…` `` string literals plus `//` and `/* */`
# comments before matching, so a deprecated name appearing only in prose, an
# example, or a URL is not flagged. Known limits (documented, not current bugs):
# single-quoted rune literals are not tracked (a `'/'` rune cannot hide a `/*`
# or `//`), and a deprecated call inside a string is intentionally invisible.
# Coverage collapse (empty symbol/file set) is asserted by the caller, which
# fails closed.

BEGIN {
    violations = 0
}

# --- pass 1: deprecated-symbol map (name -> package) ---
NR == FNR {
    n = split($0, parts, ".")
    name = parts[n]
    pkg = ""
    for (i = 1; i < n; i++) pkg = (i == 1 ? parts[i] : pkg "." parts[i])
    sym2pkg[name] = pkg
    next
}

# --- pass 2: .templ files ---
FNR == 1 {
    filePkg = ""
    inRaw = 0
    inStr = 0
    inBlock = 0
    inImport = 0
}

# package clause (first `package X` in the file)
filePkg == "" && $0 ~ /^[[:space:]]*package[[:space:]]+[A-Za-z_]/ {
    p = $0
    sub(/^[[:space:]]*package[[:space:]]+/, "", p)
    sub(/[[:space:]].*/, "", p)
    filePkg = p
    next
}

# import block start/end
/^[[:space:]]*import[[:space:]]*\(/ { inImport = 1; next }
inImport && /^[[:space:]]*\)/ { inImport = 0; next }

# grouped import entry: "path" or alias "path"
inImport {
    if (match($0, /"[^"]+"/)) {
        path = substr($0, RSTART + 1, RLENGTH - 2)
        pre = substr($0, 1, RSTART - 1)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", pre)
        base = path
        sub(/^.*\//, "", base)
        alias = (pre == "" ? base : pre)
        imports[FILENAME, alias] = base
    }
    next
}

# single-line import: import "path" or import alias "path"
/^[[:space:]]*import[[:space:]]+/ {
    if (match($0, /"[^"]+"/)) {
        path = substr($0, RSTART + 1, RLENGTH - 2)
        pre = substr($0, 1, RSTART - 1)
        sub(/^[[:space:]]*import[[:space:]]+/, "", pre)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", pre)
        base = path
        sub(/^.*\//, "", base)
        alias = (pre == "" ? base : pre)
        imports[FILENAME, alias] = base
    }
    next
}

{
    inStr = 0
    masked = mask_line($0)
    scan_calls(masked)
}

function blank(n,    s, i) {
    s = ""
    for (i = 1; i <= n; i++) s = s " "
    return s
}

# mask_line blanks string-literal and comment content so a deprecated name in
# prose, an example, or a URL is not flagged. It tracks backtick raw strings
# and /* */ block comments across lines (inRaw/inBlock persist; `"…"` strings
# cannot span lines, so inStr resets each line).
function mask_line(line,    out, i, len, c, nxt) {
    len = length(line)
    out = ""
    i = 1
    while (i <= len) {
        c = substr(line, i, 1)
        nxt = substr(line, i + 1, 1)
        if (inBlock) {
            if (c == "*" && nxt == "/") { out = out "  "; i = i + 2; inBlock = 0 }
            else { out = out " "; i = i + 1 }
        } else if (inRaw) {
            if (c == "`") { out = out " "; i = i + 1; inRaw = 0 }
            else { out = out " "; i = i + 1 }
        } else if (inStr) {
            if (c == "\\") { out = out "  "; i = i + 2 }
            else if (c == "\"") { out = out " "; i = i + 1; inStr = 0 }
            else { out = out " "; i = i + 1 }
        } else {
            if (c == "`") { out = out " "; i = i + 1; inRaw = 1 }
            else if (c == "\"") { out = out " "; i = i + 1; inStr = 1 }
            else if (c == "/" && nxt == "*") { out = out "  "; i = i + 2; inBlock = 1 }
            else if (c == "/" && nxt == "/") { out = out blank(len - i + 1); i = len + 1 }
            else { out = out c; i = i + 1 }
        }
    }
    return out
}

function report(ln, sym) {
    printf "%s:%d: deprecated caller of %s\n", FILENAME, ln, sym > "/dev/stderr"
    violations++
}

function scan_calls(line,    rest, i, tok, parts2, alias, name, p) {
    # qualified: ident.ident(
    i = 1
    while (1) {
        rest = substr(line, i)
        if (!match(rest, /[A-Za-z_][A-Za-z0-9_]*\.[A-Za-z_][A-Za-z0-9_]*\(/)) break
        tok = substr(rest, RSTART, RLENGTH)
        n2 = split(substr(tok, 1, length(tok) - 1), parts2, ".")
        alias = parts2[1]
        name = parts2[2]
        if (name in sym2pkg) {
            p = imports[FILENAME, alias]
            if (p == "" && alias == filePkg) p = filePkg
            if (p == sym2pkg[name]) report(FNR, sym2pkg[name] "." name)
        }
        i = i + RSTART + RLENGTH - 1
    }

    # unqualified: @ident(
    i = 1
    while (1) {
        rest = substr(line, i)
        if (!match(rest, /@[A-Za-z_][A-Za-z0-9_]*\(/)) break
        tok = substr(rest, RSTART, RLENGTH)
        name = substr(tok, 2, length(tok) - 2)
        if (name in sym2pkg && sym2pkg[name] == filePkg) report(FNR, filePkg "." name)
        i = i + RSTART + RLENGTH - 1
    }
}

END {
    if (violations > 0) exit 1
    exit 0
}
