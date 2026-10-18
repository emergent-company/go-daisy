package table

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func renderSortableTableHeader(t *testing.T, props SortableTableHeaderProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := SortableTableHeader(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func TestSortableTableHeaderRendersActiveAscending(t *testing.T) {
	got := renderSortableTableHeader(t, SortableTableHeaderProps{
		Label:    "CPU",
		Column:   "cpu",
		Sort:     SortContext{Params: SortParams{By: "cpu", Dir: SortAsc}, BasePath: "/hosts", Query: url.Values{"group": {"prod"}}},
		FirstDir: SortAsc,
		Class:    "text-right",
	})
	if !strings.Contains(got, `scope="col"`) {
		t.Fatalf("missing scope: %s", got)
	}
	if !strings.Contains(got, `aria-sort="ascending"`) {
		t.Fatalf("missing aria-sort=ascending: %s", got)
	}
	if !strings.Contains(got, `href="/hosts?group=prod&amp;sort_by=cpu&amp;sort_dir=desc"`) {
		t.Fatalf("unexpected href: %s", got)
	}
	if !strings.Contains(got, `class="text-right"`) {
		t.Fatalf("missing class: %s", got)
	}
	if !strings.Contains(got, `lucide--chevron-up`) {
		t.Fatalf("missing active asc icon: %s", got)
	}
	if !strings.Contains(got, `aria-label="Sort by CPU, currently ascending; activate to sort descending"`) {
		t.Fatalf("missing sort aria-label: %s", got)
	}
}

func TestSortableTableHeaderRendersActiveDescending(t *testing.T) {
	got := renderSortableTableHeader(t, SortableTableHeaderProps{
		Label:    "CPU",
		Column:   "cpu",
		Sort:     SortContext{Params: SortParams{By: "cpu", Dir: SortDesc}, BasePath: "/hosts", Query: url.Values{}},
		FirstDir: SortDesc,
	})
	if !strings.Contains(got, `aria-sort="descending"`) {
		t.Fatalf("missing aria-sort=descending: %s", got)
	}
	if !strings.Contains(got, `sort_dir=asc`) {
		t.Fatalf("expected next direction asc: %s", got)
	}
	if !strings.Contains(got, `lucide--chevron-down`) {
		t.Fatalf("missing active desc icon: %s", got)
	}
}

func TestSortableTableHeaderInactiveOmitsAriaSort(t *testing.T) {
	got := renderSortableTableHeader(t, SortableTableHeaderProps{
		Label:    "Host",
		Column:   "host",
		Sort:     SortContext{BasePath: "/hosts", Query: url.Values{}},
		FirstDir: SortAsc,
	})
	if strings.Contains(got, "aria-sort") {
		t.Fatalf("inactive header must omit aria-sort: %s", got)
	}
	if !strings.Contains(got, `sort_by=host&amp;sort_dir=asc`) {
		t.Fatalf("unexpected first href: %s", got)
	}
	if !strings.Contains(got, `lucide--chevrons-up-down`) {
		t.Fatalf("missing inactive icon: %s", got)
	}
}

func TestSortableTableHeaderThirdClickClearsSort(t *testing.T) {
	got := renderSortableTableHeader(t, SortableTableHeaderProps{
		Label:    "CPU",
		Column:   "cpu",
		Sort:     SortContext{Params: SortParams{By: "cpu", Dir: SortAsc}, BasePath: "/hosts", Query: url.Values{"q": {"web"}}},
		FirstDir: SortDesc,
	})
	if strings.Contains(got, "sort_by") || strings.Contains(got, "sort_dir") {
		t.Fatalf("third click should clear sort params: %s", got)
	}
	if !strings.Contains(got, `href="/hosts?q=web"`) {
		t.Fatalf("expected filter-preserving clear href: %s", got)
	}
}

func TestSortableTableHeaderPlainColumn(t *testing.T) {
	got := renderSortableTableHeader(t, SortableTableHeaderProps{
		Label:  "Actions",
		Column: "",
		Sort:   SortContext{BasePath: "/hosts"},
	})
	if strings.Contains(got, "<a ") {
		t.Fatalf("plain header must not render a link: %s", got)
	}
	if !strings.Contains(got, ">Actions</th>") {
		t.Fatalf("missing label: %s", got)
	}
	if !strings.Contains(got, `scope="col"`) {
		t.Fatalf("missing scope: %s", got)
	}
}

func TestSortableTableHeaderLinkAttrs(t *testing.T) {
	got := renderSortableTableHeader(t, SortableTableHeaderProps{
		Label:    "Name",
		Column:   "name",
		Sort:     SortContext{BasePath: "/items", Query: url.Values{}},
		FirstDir: SortAsc,
		LinkAttrs: templ.Attributes{
			"hx-get":    "/items?sort_by=name&sort_dir=asc",
			"hx-target": "#items-body",
		},
	})
	if !strings.Contains(got, `hx-target="#items-body"`) {
		t.Fatalf("missing hx-target: %s", got)
	}
	if !strings.Contains(got, `hx-get="/items?sort_by=name&amp;sort_dir=asc"`) {
		t.Fatalf("missing escaped hx-get: %s", got)
	}
}
