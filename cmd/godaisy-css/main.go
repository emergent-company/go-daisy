// Command godaisy-css prints the CSS a set of go-daisy component packages needs,
// from the components/css registry: the union of daisyUI modules (for a
// consumer's `@plugin "daisyui"` include/exclude list) and the co-located
// custom CSS files (for @import) — so consumers compile exactly what they use
// without hand-maintaining the list.
//
// Usage: pass component package directory names or full import paths:
//
//	go run ./cmd/godaisy-css ui nav form table layout alpine
//	go run ./cmd/godaisy-css github.com/emergent-company/go-daisy/components/ui
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/emergent-company/go-daisy/components/css"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: godaisy-css <package-dir-or-import>...")
		os.Exit(2)
	}
	pkgs := os.Args[1:]

	modules := css.ModulesFor(pkgs...)
	files := css.CSSFilesFor(pkgs...)

	fmt.Printf("daisyui modules (union of %d packages):\n", len(pkgs))
	fmt.Printf("  %s\n", strings.Join(modules, ", "))
	if len(files) > 0 {
		fmt.Println("co-located custom CSS to @import (relative to the components/ root):")
		for _, f := range files {
			fmt.Printf("  %s\n", f)
		}
	}
}
