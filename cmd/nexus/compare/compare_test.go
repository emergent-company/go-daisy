package compare

import (
	"reflect"
	"testing"
)

func TestCanonicalClassesMapsToneTokens(t *testing.T) {
	got := canonicalClasses([]string{"text-base-content/60", "text-sm", "text-muted"})
	want := []string{"text-muted", "text-sm"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("canonicalClasses() = %v, want %v", got, want)
	}
}

func TestCanonicalClassesKeepsUnknown(t *testing.T) {
	got := canonicalClasses([]string{"card", "bg-base-100", "text-base-content/20"})
	want := []string{"card", "bg-base-100", "text-base-content/20"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("canonicalClasses() = %v, want %v", got, want)
	}
}

func TestClassAliasesCoverEveryToneToken(t *testing.T) {
	// Every opacity the token layer replaces must have an alias, otherwise a
	// converted page reports a phantom class diff.
	for _, raw := range []string{"80", "70", "60", "50", "40"} {
		key := "text-base-content/" + raw
		if _, ok := classAliases[key]; !ok {
			t.Errorf("missing alias for %s", key)
		}
	}
}
