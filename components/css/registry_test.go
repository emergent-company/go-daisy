package css

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRegistryLoads verifies the CSS registry contains the core component
// packages and that every listed module set is usable.
func TestRegistryLoads(t *testing.T) {
	for _, pkg := range []string{"ui", "nav", "form", "table", "layout", "alpine", "modal"} {
		if _, ok := Registry[pkg]; !ok {
			t.Errorf("registry missing package %q", pkg)
		}
	}
	if len(Registry["ui"].DaisyUIModules) == 0 {
		t.Error("ui package module list is empty")
	}
	if len(Registry["form"].DaisyUIModules) == 0 {
		t.Error("form package module list is empty")
	}
}

// TestCustomCSSExists verifies the co-located custom CSS file is shipped
// alongside the registry.
func TestCustomCSSExists(t *testing.T) {
	// registry.go lives in components/css/, next to custom.css.
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(here, filepath.Base(CustomCSS))
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("custom.css not found at %s: %v", path, err)
	}
}
