package compare

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// DiffType categorizes a structural difference.
type DiffType string

const (
	DiffMissingElement DiffType = "missing_element"
	DiffExtraElement   DiffType = "extra_element"
	DiffWrongClass     DiffType = "wrong_class"
	DiffMissingClass   DiffType = "missing_class"
	DiffExtraClass     DiffType = "extra_class"
	DiffWrongTag       DiffType = "wrong_tag"
	DiffTextMismatch   DiffType = "text_mismatch"
	DiffChildCount     DiffType = "child_count"
	DiffOK             DiffType = "ok"
)

// DiffEntry describes a single structural difference.
type DiffEntry struct {
	Type     DiffType `json:"type"`
	Path     string   `json:"path"`     // CSS selector path to element
	NexusVal string   `json:"nexus_val"` // value from nexus-html reference
	GoDaisy  string   `json:"godaisy"`   // value from go-daisy output
}

// PageResult holds the comparison result for one page.
type PageResult struct {
	Page        PageMapping `json:"page"`
	Status      string      `json:"status"` // "ok", "error", "diff"
	Error       string      `json:"error,omitempty"`
	NexusChars  int         `json:"nexus_chars"`
	GoDaisyChars int        `json:"godaisy_chars"`
	Diffs       []DiffEntry `json:"diffs"`
}

// CompareAll fetches all pages and compares them.
func CompareAll(mappings []PageMapping) []PageResult {
	var results []PageResult
	for _, m := range mappings {
		r := comparePage(m)
		results = append(results, r)
	}
	return results
}

func comparePage(m PageMapping) PageResult {
	r := PageResult{Page: m, Status: "ok", Diffs: []DiffEntry{}}

	nexusHTML, err := os.ReadFile(refPath(m.HTMLFile))
	if err != nil {
		r.Status = "error"
		r.Error = fmt.Sprintf("read nexus-html: %v", err)
		return r
	}
	r.NexusChars = len(nexusHTML)

	resp, err := http.Get(goDaisyURL(m.Route))
	if err != nil {
		r.Status = "error"
		r.Error = fmt.Sprintf("fetch go-daisy: %v", err)
		return r
	}
	defer resp.Body.Close()
	goDaisyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Status = "error"
		r.Error = fmt.Sprintf("read go-daisy: %v", err)
		return r
	}
	r.GoDaisyChars = len(goDaisyBytes)

	nexusContent := extractContentNodes(nexusHTML)
	goDaisyContent := extractContentNodes(goDaisyBytes)

	r.Diffs = compareNodeLists(nexusContent, goDaisyContent, "")
	if len(r.Diffs) > 0 {
		r.Status = "diff"
	}
	return r
}

// extractContentNodes extracts all child nodes from inside #layout-content.
func extractContentNodes(raw []byte) []*html.Node {
	doc, err := html.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return nil
	}
	contentDiv := findElementByID(doc, "layout-content")
	if contentDiv == nil {
		return nil
	}
	var children []*html.Node
	for c := contentDiv.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode || (c.Type == html.TextNode && strings.TrimSpace(c.Data) != "") {
			children = append(children, c)
		}
	}
	return children
}

func findElementByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findElementByID(c, id); found != nil {
			return found
		}
	}
	return nil
}

// nodePath builds a CSS-selector-like path for an element node.
func nodePath(n *html.Node, prefix string) string {
	if n == nil {
		return prefix
	}
	tag := n.Data
	if n.Type == html.ElementNode {
		id := getAttr(n, "id")
		if id != "" {
			tag += "#" + id
		}
		cls := getAttr(n, "class")
		if cls != "" {
			parts := strings.Fields(cls)
			if len(parts) > 0 && len(parts) <= 2 {
				tag += "." + strings.Join(parts, ".")
			}
		}
	}
	if prefix == "" {
		return tag
	}
	return prefix + " > " + tag
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// collectClasses returns sorted class names from a node.
func collectClasses(n *html.Node) []string {
	cls := getAttr(n, "class")
	if cls == "" {
		return nil
	}
	return strings.Fields(cls)
}

// collectAllText collects all text content recursively from a node.
func collectAllText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(collectAllText(c))
	}
	return b.String()
}

func normalizeText(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

func classSet(classes []string) map[string]bool {
	m := make(map[string]bool)
	for _, c := range classes {
		m[c] = true
	}
	return m
}

func compareNodeLists(nexus, goDaisy []*html.Node, parentPath string) []DiffEntry {
	var diffs []DiffEntry

	maxLen := len(nexus)
	if len(goDaisy) > maxLen {
		maxLen = len(goDaisy)
	}

	for i := range maxLen {
		path := fmt.Sprintf("%s > [%d]", parentPath, i)
		if i >= len(nexus) {
			diffs = append(diffs, DiffEntry{
				Type:    DiffExtraElement,
				Path:    path,
				GoDaisy: nodeSummary(goDaisy[i]),
			})
			continue
		}
		if i >= len(goDaisy) {
			diffs = append(diffs, DiffEntry{
				Type:     DiffMissingElement,
				Path:     path,
				NexusVal: nodeSummary(nexus[i]),
			})
			continue
		}
		diffs = append(diffs, compareTwoNodes(nexus[i], goDaisy[i], path)...)
	}
	return diffs
}

func compareTwoNodes(nexus, goDaisy *html.Node, path string) []DiffEntry {
	var diffs []DiffEntry

	// Compare tag names
	nTag := nexus.Data
	gTag := goDaisy.Data
	if nTag != gTag && nexus.Type == html.ElementNode && goDaisy.Type == html.ElementNode {
		diffs = append(diffs, DiffEntry{
			Type:     DiffWrongTag,
			Path:     path,
			NexusVal: nTag,
			GoDaisy:  gTag,
		})
		return diffs
	}

	// Compare classes
	if nexus.Type == html.ElementNode && goDaisy.Type == html.ElementNode {
		nClasses := collectClasses(nexus)
		gClasses := collectClasses(goDaisy)
		nSet := classSet(nClasses)
		gSet := classSet(gClasses)

		for _, c := range nClasses {
			if !gSet[c] {
				diffs = append(diffs, DiffEntry{
					Type:     DiffMissingClass,
					Path:     path,
					NexusVal: c,
				})
			}
		}
		for _, c := range gClasses {
			if !nSet[c] {
				diffs = append(diffs, DiffEntry{
					Type:    DiffExtraClass,
					Path:    path,
					GoDaisy: c,
				})
			}
		}
	}

	// Compare text content for leaf nodes
	if nexus.Type == html.TextNode || goDaisy.Type == html.TextNode {
		nText := normalizeText(collectAllText(nexus))
		gText := normalizeText(collectAllText(goDaisy))
		if nText != gText && nText != "" && gText != "" {
			diffs = append(diffs, DiffEntry{
				Type:     DiffTextMismatch,
				Path:     path,
				NexusVal: truncate(nText, 80),
				GoDaisy:  truncate(gText, 80),
			})
		}
		return diffs
	}

	// Compare children recursively
	nChildren := elementChildren(nexus)
	gChildren := elementChildren(goDaisy)
	diffs = append(diffs, compareNodeLists(nChildren, gChildren, nodePath(nexus, path))...)

	return diffs
}

func elementChildren(n *html.Node) []*html.Node {
	var children []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode || (c.Type == html.TextNode && strings.TrimSpace(c.Data) != "") {
			children = append(children, c)
		}
	}
	return children
}

func nodeSummary(n *html.Node) string {
	if n == nil {
		return "nil"
	}
	if n.Type == html.TextNode {
		return truncate(strings.TrimSpace(n.Data), 60)
	}
	tag := n.Data
	id := getAttr(n, "id")
	cls := getAttr(n, "class")
	parts := []string{"<" + tag}
	if id != "" {
		parts = append(parts, "id="+id)
	}
	if cls != "" {
		parts = append(parts, "class=\""+truncate(cls, 60)+"\"")
	}
	return strings.Join(parts, " ") + ">"
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// Summary returns counts of each diff type.
func (r *PageResult) Summary() string {
	if r.Status == "error" {
		return fmt.Sprintf("ERROR: %s", r.Error)
	}
	if r.Status == "ok" {
		return "OK — no structural differences"
	}
	byType := make(map[DiffType]int)
	for _, d := range r.Diffs {
		byType[d.Type]++
	}
	var parts []string
	keys := []DiffType{DiffMissingElement, DiffExtraElement, DiffWrongTag, DiffMissingClass, DiffExtraClass, DiffTextMismatch}
	for _, k := range keys {
		if n := byType[k]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", k, n))
		}
	}
	return "DIFF: " + strings.Join(parts, " ")
}

// CategoryGroups returns results grouped by category for the report.
func CategoryGroups(results []PageResult) []string {
	seen := make(map[string]bool)
	var cats []string
	for _, r := range results {
		cat := pageCategory(r.Page.Route)
		if !seen[cat] {
			seen[cat] = true
			cats = append(cats, cat)
		}
	}
	sort.Strings(cats)
	return cats
}

func pageCategory(route string) string {
	switch {
	case strings.Contains(route, "/dashboards"):
		return "Dashboards"
	case strings.Contains(route, "/ecommerce"):
		return "Ecommerce CRUD"
	case strings.Contains(route, "/gen-ai"):
		return "Gen AI"
	case strings.Contains(route, "/chat"):
		return "Chat"
	case strings.Contains(route, "/file-manager"):
		return "File Manager"
	case strings.Contains(route, "/pages"):
		return "Pages"
	case strings.Contains(route, "/auth"):
		return "Auth"
	case strings.Contains(route, "/landing"):
		return "Landing"
	default:
		return "Other"
	}
}

// FilterByCategory returns results for a specific category.
func FilterByCategory(results []PageResult, cat string) []PageResult {
	var out []PageResult
	for _, r := range results {
		if pageCategory(r.Page.Route) == cat {
			out = append(out, r)
		}
	}
	return out
}
