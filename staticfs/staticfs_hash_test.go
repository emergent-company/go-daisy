package staticfs

import (
	"strings"
	"testing"
	"testing/fstest"
)

// TestStaticFilesHashMatchesCommitted fails when staticfs_hash.go is stale —
// i.e. a file under static/ changed without re-running `go generate
// ./staticfs/`. A stale constant freezes every cache-busting ?v= at its old
// value, so a changed asset is served under a URL clients have already cached.
//
// Run `go generate ./staticfs/` to update the constant.
func TestStaticFilesHashMatchesCommitted(t *testing.T) {
	recomputed, err := computeHash(files)
	if err != nil {
		t.Fatalf("computeHash: %v", err)
	}
	if recomputed != staticFilesHash {
		t.Fatalf("staticfs_hash.go is stale:\n  committed:  %s\n  recomputed: %s\nrun `go generate ./staticfs/`",
			staticFilesHash, recomputed)
	}
}

// TestComputeHashGolden pins the exact algorithm shape the //go:generate shell
// pipeline mirrors: sha256 over "<sha256(content)>  <path>\n" lines in sorted
// path order. computeHash cannot be order-dependent (fs.WalkDir sorts), so a
// stability loop would assert nothing; a golden value instead fails loudly if
// the separator, the path form, or the sort ever changes without the shell
// directive being updated to match.
func TestComputeHashGolden(t *testing.T) {
	fake := fstest.MapFS{
		// Inserted out of order on purpose: the result must not depend on it.
		"static/b.js": &fstest.MapFile{Data: []byte("b\n")},
		"static/a.js": &fstest.MapFile{Data: []byte("a\n")},
	}
	const want = "d35f5e00dc694d55bac392908a05adfbcefeaa1c7a249b94c420ebc8335d1f56"
	got, err := computeHash(fake)
	if err != nil {
		t.Fatalf("computeHash: %v", err)
	}
	if got != want {
		t.Fatalf("computeHash format changed:\n  got:  %s\n  want: %s\nif this is intentional, keep the //go:generate pipeline in staticfs.go in sync",
			got, want)
	}
}

// TestScriptNonce locks the CSP tag contract: without a nonce the output must
// stay byte-identical to the pre-CSP Script(), and with one the attribute is
// present and escaped.
func TestScriptNonce(t *testing.T) {
	plain := Script("/js/htmx.js")
	if strings.Contains(plain, "nonce=") {
		t.Errorf("empty nonce should render no attribute: %s", plain)
	}
	if got := ScriptNonce("/js/htmx.js", ""); got != plain {
		t.Errorf("ScriptNonce with empty nonce diverged from Script:\n%s\n%s", got, plain)
	}

	nonced := ScriptNonce("/js/htmx.js", "abc123")
	if !strings.Contains(nonced, `nonce="abc123"`) {
		t.Errorf("nonce attribute missing: %s", nonced)
	}
	if !strings.Contains(nonced, "/static/js/htmx.js?v=") {
		t.Errorf("resolved asset path lost: %s", nonced)
	}

	escaped := ScriptNonce("/js/htmx.js", `a"b<c`)
	if !strings.Contains(escaped, `nonce="a&#34;b&lt;c"`) {
		t.Errorf("nonce was not escaped: %s", escaped)
	}
}

// TestInlineScriptNonce covers the helper used for hand-built inline blocks.
func TestInlineScriptNonce(t *testing.T) {
	if got := InlineScript("", "var x = 1;"); got != `<script>var x = 1;</script>` {
		t.Errorf("no-nonce inline script = %s", got)
	}
	if got := InlineScript("n1", "var x = 1;"); got != `<script nonce="n1">var x = 1;</script>` {
		t.Errorf("nonced inline script = %s", got)
	}
}
