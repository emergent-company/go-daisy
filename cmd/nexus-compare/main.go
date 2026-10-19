package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emergent-company/go-daisy/cmd/nexus/compare"
)

func main() {
	base := flag.String("base", "http://localhost:11001", "base URL of the go-daisy nexus server")
	ref := flag.String("ref", "/root/nexus-html/html", "directory of reference nexus-html pages")
	maxDiffPages := flag.Int("max-diff-pages", -1, "diff-page budget: exit 2 when diff count exceeds it (-1 = no budget). Exit 2 takes precedence over exit 1 (passthrough/error).")
	jsonPath := flag.String("json", "/tmp/nexus-compare.json", "path to write JSON results")
	htmlPath := flag.String("html", "/tmp/nexus-compare.html", "path to write HTML report")
	filter := flag.String("filter", "", "substring filter on page name or route")
	flag.Parse()

	compare.SetConfig(compare.Config{BaseURL: *base, RefDir: *ref})

	mappings := compare.Pages()
	if *filter != "" {
		mappings = filterMappings(mappings, *filter)
	}

	results := compare.CompareAll(mappings)

	fmt.Println("=== Nexus go-daisy vs nexus-html Structural Comparison ===")

	compared := len(compare.Compared(mappings))
	ok, diff, passthrough, errCount, skipped := 0, 0, 0, 0, 0
	var passthroughPages, errorPages []string

	for _, r := range results {
		switch r.Status {
		case "skipped":
			skipped++
			fmt.Printf("  %-40s %-11s %s\n", r.Page.Name, "skipped", r.Page.Reason)
			continue
		}
		label := r.Status
		detail := fmt.Sprintf("(%d diffs)", len(r.Diffs))
		switch r.Status {
		case "ok":
			ok++
		case "diff":
			diff++
		case "passthrough":
			passthrough++
			passthroughPages = append(passthroughPages, fmt.Sprintf("%s (%s)", r.Page.Name, r.Page.Route))
			detail = fmt.Sprintf("raw embed %q", r.RawPage)
		case "error":
			errCount++
			errorPages = append(errorPages, fmt.Sprintf("%s (%s): %s", r.Page.Name, r.Page.Route, r.Error))
			detail = r.Error
		}
		fmt.Printf("  %-40s %-11s %s\n", r.Page.Name, label, detail)
	}

	fmt.Printf("\n  compared %d, ok %d, diff %d, passthrough %d, error %d, skipped %d\n\n", compared, ok, diff, passthrough, errCount, skipped)

	if len(passthroughPages) > 0 {
		fmt.Println("Passthrough pages (unported raw-HTML work queue):")
		for _, p := range passthroughPages {
			fmt.Printf("  - %s\n", p)
		}
		fmt.Println()
	}
	if len(errorPages) > 0 {
		fmt.Println("Pages where the comparator errored:")
		for _, p := range errorPages {
			fmt.Printf("  - %s\n", p)
		}
		fmt.Println()
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("marshal results: %v", err)
	}
	if err := os.WriteFile(*jsonPath, jsonData, 0644); err != nil {
		log.Fatalf("write %s: %v", *jsonPath, err)
	}
	fmt.Printf("JSON: %s\n", *jsonPath)

	html := compare.GenerateHTMLReport(results)
	if err := os.WriteFile(*htmlPath, []byte(html), 0644); err != nil {
		log.Fatalf("write %s: %v", *htmlPath, err)
	}
	fmt.Printf("HTML: %s\n", *htmlPath)

	os.Exit(exitCode(diff, passthrough, errCount, *maxDiffPages))
}

func filterMappings(mappings []compare.PageMapping, filter string) []compare.PageMapping {
	var out []compare.PageMapping
	for _, m := range mappings {
		if strings.Contains(m.Name, filter) || strings.Contains(m.Route, filter) || strings.Contains(m.HTMLFile, filter) {
			out = append(out, m)
		}
	}
	return out
}

// exitCode maps the comparison outcome to a process exit code:
//   - 2: diff-page budget breached (only when maxDiffPages >= 0). This takes
//     precedence over passthrough/error so a phase gate can ratchet the budget
//     down and still fail loudly on an out-of-budget run even if raw pages
//     remain.
//   - 1: any passthrough page or comparator error, when no budget is set or the
//     budget is not breached.
//   - 0: clean.
func exitCode(diff, passthrough, errCount, maxDiffPages int) int {
	if maxDiffPages >= 0 && diff > maxDiffPages {
		return 2
	}
	if passthrough > 0 || errCount > 0 {
		return 1
	}
	return 0
}
