package staticfs

import (
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
