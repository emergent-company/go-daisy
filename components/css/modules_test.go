package css

import (
	"reflect"
	"sort"
	"testing"
)

func TestModulesFor(t *testing.T) {
	// single package
	got := ModulesFor("form")
	if !reflect.DeepEqual(got, Registry["form"].DaisyUIModules) {
		t.Errorf("form modules = %v", got)
	}

	// union is sorted + deduped, full import paths accepted
	got = ModulesFor("github.com/emergent-company/go-daisy/components/form", "ui", "form")
	want := Registry["form"].DaisyUIModules
	seen := map[string]bool{}
	for _, m := range want {
		seen[m] = true
	}
	for _, m := range Registry["ui"].DaisyUIModules {
		seen[m] = true
	}
	if !reflect.DeepEqual(got, sortedKeys(seen)) {
		t.Errorf("union = %v", got)
	}

	// unknown package ignored
	if got := ModulesFor("does-not-exist"); len(got) != 0 {
		t.Errorf("unknown pkg modules = %v, want empty", got)
	}
}

func TestCSSFilesFor(t *testing.T) {
	got := CSSFilesFor("ui", "layout", "nope")
	if !reflect.DeepEqual(got, []string{CustomCSS}) {
		t.Errorf("CSSFilesFor = %v, want [%s]", got, CustomCSS)
	}
	if got := CSSFilesFor("table"); len(got) != 0 {
		t.Errorf("table css = %v, want empty", got)
	}
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
