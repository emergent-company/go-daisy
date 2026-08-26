package compare

import (
	"fmt"
	"strings"
)

func GenerateHTMLReport(results []PageResult) string {
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Nexus go-daisy vs nexus-html Comparison</title>
<style>
* { margin:0; padding:0; box-sizing:border-box; }
body { font-family: system-ui, sans-serif; background: #f5f5f5; color: #333; padding: 20px; }
h1 { margin-bottom: 10px; }
.summary { display: flex; gap: 20px; margin-bottom: 30px; }
.stat { padding: 10px 20px; border-radius: 8px; color: #fff; font-weight: bold; }
.stat-ok { background: #22c55e; }
.stat-diff { background: #f59e0b; }
.stat-error { background: #ef4444; }
.category { margin-bottom: 30px; }
.category h2 { border-bottom: 2px solid #d1d5db; padding-bottom: 5px; margin-bottom: 15px; }
.page { background: #fff; border-radius: 8px; padding: 15px; margin-bottom: 10px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.page-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.page-name { font-weight: 600; }
.status-badge { padding: 2px 10px; border-radius: 12px; font-size: 12px; font-weight: 600; }
.badge-ok { background: #dcfce7; color: #16a34a; }
.badge-diff { background: #fef3c7; color: #d97706; }
.badge-error { background: #fee2e2; color: #dc2626; }
.diff-list { font-size: 13px; max-height: 400px; overflow-y: auto; }
.diff-item { padding: 4px 8px; border-left: 3px solid #e5e7eb; margin: 4px 0; font-family: monospace; font-size: 12px; }
.diff-missing_element { border-left-color: #ef4444; background: #fef2f2; }
.diff-extra_element { border-left-color: #f59e0b; background: #fffbeb; }
.diff-wrong_class, .diff-missing_class, .diff-extra_class { border-left-color: #8b5cf6; background: #f5f3ff; }
.diff-text_mismatch { border-left-color: #06b6d4; background: #ecfeff; }
.diff-wrong_tag { border-left-color: #ec4899; background: #fdf2f8; }
.diff-type { font-weight: 600; color: #374151; }
.diff-path { color: #6b7280; }
.collapsible { cursor: pointer; user-select: none; }
.no-diffs { color: #9ca3af; font-style: italic; }
</style>
</head>
<body>
<h1>Nexus Structural Comparison</h1>
<p>go-daisy (port 11001) vs nexus-html (/root/nexus-html/html/)</p>
`)

	// Summary
	ok, diff, errCount := 0, 0, 0
	for _, r := range results {
		switch r.Status {
		case "ok":
			ok++
		case "diff":
			diff++
		case "error":
			errCount++
		}
	}
	b.WriteString( fmt.Sprintf(`<div class="summary">
<div class="stat stat-ok">%d OK</div>
<div class="stat stat-diff">%d DIFF</div>
<div class="stat stat-error">%d ERROR</div>
</div>`, ok, diff, errCount))

	// By category
	cats := CategoryGroups(results)
	for _, cat := range cats {
		pages := FilterByCategory(results, cat)
		b.WriteString( fmt.Sprintf(`<div class="category"><h2>%s</h2>`, cat))
		for _, r := range pages {
			badgeClass := "badge-ok"
			if r.Status == "diff" {
				badgeClass = "badge-diff"
			} else if r.Status == "error" {
				badgeClass = "badge-error"
			}
			b.WriteString( fmt.Sprintf(`<div class="page">
<div class="page-header collapsible" onclick="this.nextElementSibling.classList.toggle('no-diffs')">
<span class="page-name">%s</span>
<span class="status-badge %s">%s</span>
<span style="color:#6b7280;font-size:12px">(%d diffs)</span>
<span style="color:#9ca3af;font-size:11px">nexus=%d chars / godaisy=%d chars</span>
</div>
<div class="diff-list">`, r.Page.Name, badgeClass, r.Status, len(r.Diffs), r.NexusChars, r.GoDaisyChars))
			if len(r.Diffs) == 0 {
				b.WriteString( `<div class="no-diffs">No structural differences</div>`)
			}
			for _, d := range r.Diffs {
				b.WriteString( fmt.Sprintf(`<div class="diff-item diff-%s">
<span class="diff-type">%s</span>
<span class="diff-path">%s</span>
<div>nexus: %s</div>
<div>godaisy: %s</div>
</div>`, d.Type, d.Type, d.Path, d.NexusVal, d.GoDaisy))
			}
			b.WriteString( `</div></div>`)
		}
		b.WriteString( `</div>`)
	}
	b.WriteString( `</body></html>`)
	return b.String()
}
