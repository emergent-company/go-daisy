package table

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestParseSortPrefixedDefaultsToNone(t *testing.T) {
	got := ParseSortPrefixed(url.Values{}, "", "", SortNone)
	if got.By != "" || got.Dir != SortNone {
		t.Fatalf("expected no sort, got %#v", got)
	}
}

func TestParseSortPrefixedParsesPrefixedParams(t *testing.T) {
	q := url.Values{
		"nodes_sort_by":  {"cpu"},
		"nodes_sort_dir": {"desc"},
		"sort_by":        {"host"},
		"sort_dir":       {"asc"},
	}
	got := ParseSortPrefixed(q, "nodes_", "", SortNone)
	if got.By != "cpu" || got.Dir != SortDesc {
		t.Fatalf("expected nodes cpu desc, got %#v", got)
	}
}

func TestParseSortPrefixedUsesDefaults(t *testing.T) {
	got := ParseSortPrefixed(url.Values{}, "", "host", SortAsc)
	if got.By != "host" || got.Dir != SortAsc {
		t.Fatalf("expected default host asc, got %#v", got)
	}
}

func TestParseSortPrefixedDefaultByWithSortNoneKeepsNaturalOrder(t *testing.T) {
	got := ParseSortPrefixed(url.Values{}, "", "host", SortNone)
	if got.By != "host" || got.Dir != SortNone {
		t.Fatalf("expected default host with no sort, got %#v", got)
	}
}

func TestParseSortPrefixedInvalidDirectionFallsBackToDefault(t *testing.T) {
	q := url.Values{"sort_by": {"cpu"}, "sort_dir": {"sideways"}}
	got := ParseSortPrefixed(q, "", "", SortDesc)
	if got.By != "cpu" || got.Dir != SortDesc {
		t.Fatalf("expected cpu desc, got %#v", got)
	}
}

func TestParseSortPrefixedInvalidDefaultDirBecomesAsc(t *testing.T) {
	q := url.Values{"sort_by": {"cpu"}}
	got := ParseSortPrefixed(q, "", "", SortDir("sideways"))
	if got.By != "cpu" || got.Dir != SortAsc {
		t.Fatalf("expected cpu asc, got %#v", got)
	}
}

func TestParseSortPrefixedByWithoutDirIsAsc(t *testing.T) {
	q := url.Values{"sort_by": {"cpu"}}
	got := ParseSortPrefixed(q, "", "", SortNone)
	if got.By != "cpu" || got.Dir != SortAsc {
		t.Fatalf("expected cpu asc, got %#v", got)
	}
}

func TestParseSortMatchesPrefixed(t *testing.T) {
	q := url.Values{"sort_by": {"cpu"}, "sort_dir": {"desc"}, "q": {"web"}}
	legacy := ParseSort(q, "", SortNone)
	prefixed := ParseSortPrefixed(q, "", "", SortNone)
	if legacy != prefixed {
		t.Fatalf("ParseSort = %#v, ParseSortPrefixed = %#v", legacy, prefixed)
	}
}

func TestParseSortNoParamsMeansNoSort(t *testing.T) {
	got := ParseSort(url.Values{}, "", SortAsc)
	if got.By != "" || got.Dir != SortNone {
		t.Fatalf("expected no sort, got %#v", got)
	}
}

func TestParseSortPrefixedEmptyByForcesNone(t *testing.T) {
	q := url.Values{"sort_dir": {"desc"}}
	got := ParseSortPrefixed(q, "", "", SortAsc)
	if got.By != "" || got.Dir != SortNone {
		t.Fatalf("expected no sort, got %#v", got)
	}
}

func TestParseSortPrefixedDirWithoutByDoesNotActivateDefault(t *testing.T) {
	q := url.Values{"sort_dir": {"desc"}}
	got := ParseSortPrefixed(q, "", "host", SortNone)
	if got.By != "host" || got.Dir != SortNone {
		t.Fatalf("expected default host with no sort, got %#v", got)
	}
}

func TestNextSortDirFromCyclesFirstOppositeNone(t *testing.T) {
	tests := []struct {
		name  string
		sp    SortParams
		first SortDir
		want  SortDir
	}{
		{"new column asc first", SortParams{}, SortAsc, SortAsc},
		{"new column desc first", SortParams{}, SortDesc, SortDesc},
		{"asc first then opposite", SortParams{By: "cpu", Dir: SortAsc}, SortAsc, SortDesc},
		{"desc first then opposite", SortParams{By: "cpu", Dir: SortDesc}, SortDesc, SortAsc},
		{"opposite then none", SortParams{By: "cpu", Dir: SortDesc}, SortAsc, SortNone},
		{"opposite then none (desc first)", SortParams{By: "cpu", Dir: SortAsc}, SortDesc, SortNone},
		{"empty first defaults to asc", SortParams{}, "", SortAsc},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextSortDirFrom(tt.sp, "cpu", tt.first); got != tt.want {
				t.Fatalf("NextSortDirFrom(%#v, cpu, %q) = %q, want %q", tt.sp, tt.first, got, tt.want)
			}
		})
	}
}

func TestSortContextHrefPreservesFiltersAndResetsPagination(t *testing.T) {
	q := url.Values{"q": {"web"}, "group": {"prod"}, "page": {"3"}}
	sc := SortContext{
		Params:   SortParams{By: "name", Dir: SortAsc},
		BasePath: "/hosts",
		Query:    q,
	}
	got := sc.Href("name", SortAsc)
	want := "/hosts?group=prod&q=web&sort_by=name&sort_dir=desc"
	if got != want {
		t.Fatalf("Href = %q, want %q", got, want)
	}
	if q.Get("sort_by") != "" || q.Get("sort_dir") != "" || q.Get("page") != "3" {
		t.Fatalf("Href mutated the input query: %#v", q)
	}
}

func TestSortContextHrefClearsSort(t *testing.T) {
	sc := SortContext{
		Params:   SortParams{By: "name", Dir: SortDesc},
		BasePath: "/hosts",
		Query:    url.Values{"q": {"web"}, "sort_by": {"name"}, "sort_dir": {"desc"}},
	}
	got := sc.Href("name", SortAsc)
	want := "/hosts?q=web"
	if got != want {
		t.Fatalf("Href = %q, want %q", got, want)
	}
}

func TestSortContextHrefPrefixed(t *testing.T) {
	sc := SortContext{
		Prefix:   "nodes_",
		Params:   SortParams{By: "cpu", Dir: SortAsc},
		BasePath: "/proxmox",
		Query:    url.Values{"nodes_sort_by": {"cpu"}, "nodes_sort_dir": {"asc"}},
	}
	// Next direction after asc is desc; the existing params are replaced, not duplicated.
	got := sc.Href("cpu", SortAsc)
	want := "/proxmox?nodes_sort_by=cpu&nodes_sort_dir=desc"
	if got != want {
		t.Fatalf("Href = %q, want %q", got, want)
	}
	if sc.ByParam() != "nodes_sort_by" || sc.DirParam() != "nodes_sort_dir" {
		t.Fatalf("unexpected param names: %q %q", sc.ByParam(), sc.DirParam())
	}
}

func TestSortContextHrefWithoutBasePath(t *testing.T) {
	sc := SortContext{Query: url.Values{}}
	if got := sc.Href("name", SortAsc); got != "" {
		t.Fatalf("expected empty href without base path, got %q", got)
	}
}

func TestSortContextAriaSort(t *testing.T) {
	sc := SortContext{Params: SortParams{By: "cpu", Dir: SortDesc}}
	if got := sc.AriaSort("cpu"); got != "descending" {
		t.Fatalf("AriaSort(cpu) = %q, want descending", got)
	}
	if got := sc.AriaSort("host"); got != "" {
		t.Fatalf("AriaSort(host) = %q, want empty", got)
	}
	sc.Params.Dir = SortNone
	if got := sc.AriaSort("cpu"); got != "" {
		t.Fatalf("AriaSort with SortNone = %q, want empty", got)
	}
}

func TestNewSortContextUsesEscapedPath(t *testing.T) {
	r := httptest.NewRequest("GET", "/reports/2024%2F05?sort_by=cpu&sort_dir=asc", nil)
	sc := NewSortContext(r, "", "", SortNone)
	if sc.BasePath != "/reports/2024%2F05" {
		t.Fatalf("BasePath = %q", sc.BasePath)
	}
	if got := sc.Href("cpu", SortAsc); got != "/reports/2024%2F05?sort_by=cpu&sort_dir=desc" {
		t.Fatalf("Href = %q", got)
	}
}

func TestNewSortContextReadsRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/proxmox?nodes_sort_by=cpu&nodes_sort_dir=desc&tab=guests", nil)
	sc := NewSortContext(r, "nodes_", "", SortNone)
	if sc.Params.By != "cpu" || sc.Params.Dir != SortDesc {
		t.Fatalf("unexpected params: %#v", sc.Params)
	}
	if sc.BasePath != "/proxmox" {
		t.Fatalf("BasePath = %q", sc.BasePath)
	}
	if !sc.Active("cpu") {
		t.Fatal("cpu should be active")
	}
	if got := sc.Href("cpu", SortDesc); !strings.Contains(got, "tab=guests") {
		t.Fatalf("Href should preserve unrelated params, got %q", got)
	}
}
