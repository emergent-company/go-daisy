package staticfs

import "testing"

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

// TestComputeHashIsOrderIndependent guards the reason the generator sorts: a
// hash that depends on filesystem walk order gives the same content two
// different values across checkouts.
func TestComputeHashIsOrderIndependent(t *testing.T) {
	first, err := computeHash(files)
	if err != nil {
		t.Fatalf("computeHash: %v", err)
	}
	for i := 0; i < 5; i++ {
		again, err := computeHash(files)
		if err != nil {
			t.Fatalf("computeHash: %v", err)
		}
		if again != first {
			t.Fatalf("computeHash is not stable: %s != %s", again, first)
		}
	}
}
