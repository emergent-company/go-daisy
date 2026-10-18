package table

import (
	"net/url"
	"testing"
)

func TestDataGridSortLinkUsesBaseURLForFragment(t *testing.T) {
	p := DataGridProps{
		ID:      "users",
		BaseURL: "/datagrid?role=admin",
		Sort: SortContext{
			Params:   SortParams{By: "name", Dir: SortAsc},
			BasePath: "/",
			Query:    url.Values{},
		},
	}
	attrs := p.sortLinkAttrs(DataGridColumn{Key: "name", FirstDir: SortAsc})

	if got := attrs["hx-get"]; got != "/datagrid?role=admin&sort_by=name&sort_dir=desc" {
		t.Fatalf("hx-get = %#v, want fragment endpoint", got)
	}
	if got := attrs["hx-push-url"]; got != "/?sort_by=name&sort_dir=desc" {
		t.Fatalf("hx-push-url = %#v, want page URL", got)
	}
	if got := attrs["hx-target"]; got != "#datagrid-users-body" {
		t.Fatalf("hx-target = %#v", got)
	}
}

func TestDataGridSortLinkFallsBackToPageHref(t *testing.T) {
	p := DataGridProps{
		ID: "users",
		Sort: SortContext{
			BasePath: "/datagrid",
			Query:    url.Values{},
		},
	}
	attrs := p.sortLinkAttrs(DataGridColumn{Key: "name", FirstDir: SortAsc})

	if got := attrs["hx-get"]; got != "/datagrid?sort_by=name&sort_dir=asc" {
		t.Fatalf("hx-get = %#v", got)
	}
	if got := attrs["hx-push-url"]; got != "/datagrid?sort_by=name&sort_dir=asc" {
		t.Fatalf("hx-push-url = %#v", got)
	}
}
