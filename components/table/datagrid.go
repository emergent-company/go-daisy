package table

import "github.com/a-h/templ"

// sortContext returns the sort state with the page URL as BasePath. On a
// fragment render the request path is the fragment endpoint, so PagePath is
// what keeps href and hx-push-url on the full page.
func (p DataGridProps) sortContext() SortContext {
	sc := p.Sort
	if p.PagePath != "" {
		sc.BasePath = p.PagePath
	}
	return sc
}

// sortPageHref is the page URL for a header, used for the non-JS href and for
// hx-push-url.
func (p DataGridProps) sortPageHref(col DataGridColumn) string {
	base := p.PagePath
	if base == "" {
		base = p.Sort.BasePath
	}
	if base == "" {
		return ""
	}
	return p.Sort.HrefFor(base, col.Key, col.FirstDir)
}

// sortHref is the fragment endpoint for a header click. DataGrid's search and
// pagination already use BaseURL for this; sorting must do the same or the
// response will be a whole page instead of the table fragment. An empty
// BaseURL means the grid has no fragment endpoint, so the caller gets a plain
// page navigation through href instead of an HTMX swap.
func (p DataGridProps) sortHref(col DataGridColumn) string {
	if p.BaseURL == "" {
		return ""
	}
	return p.Sort.HrefFor(p.BaseURL, col.Key, col.FirstDir)
}

// sortLinkAttrs builds the HTMX attributes for a sortable header. hx-get points
// at the fragment endpoint while hx-push-url keeps the browser on the page URL.
func (p DataGridProps) sortLinkAttrs(col DataGridColumn) templ.Attributes {
	if p.BaseURL == "" {
		return nil
	}
	attrs := templ.Attributes{
		"hx-get":    p.sortHref(col),
		"hx-target": "#datagrid-" + p.ID + "-body",
		"hx-swap":   "innerHTML",
	}
	if page := p.sortPageHref(col); page != "" {
		attrs["hx-push-url"] = page
	}
	return attrs
}
