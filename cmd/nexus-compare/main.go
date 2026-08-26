package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/emergent-company/go-daisy/cmd/nexus/compare"
)

func main() {
	fmt.Println("=== Nexus go-daisy vs nexus-html Structural Comparison ===")

	mappings := compare.Pages()
	results := compare.CompareAll(mappings)

	ok, diff, errCount := 0, 0, 0
	for _, r := range results {
		status := r.Summary()
		fmt.Printf("  %-35s %s\n", r.Page.Name, status)
		switch r.Status {
		case "ok":
			ok++
		case "diff":
			diff++
		case "error":
			errCount++
		}
	}

	total := ok + diff + errCount
	fmt.Printf("\n  %d tested: %d OK  %d DIFF  %d ERROR\n\n", total, ok, diff, errCount)

	jsonData, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("/tmp/nexus-compare.json", jsonData, 0644)
	fmt.Println("JSON: /tmp/nexus-compare.json")

	html := compare.GenerateHTMLReport(results)
	os.WriteFile("/tmp/nexus-compare.html", []byte(html), 0644)
	fmt.Println("HTML: /tmp/nexus-compare.html")
}
