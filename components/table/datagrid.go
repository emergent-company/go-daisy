package table

import "github.com/a-h/templ"

// sortPageHref is the page URL for a header, used for the non-JS href and for
// hx-push-url.
func (p DataGridProps) sortPageHref(col DataGridColumn) string {
	return p.Sort.Href(col.Key, col.FirstDir)
}

// sortHref is the fragment endpoint for a header click. DataGrid's search and
// pagination already use BaseURL for this; sorting must do the same or the
// response will be a whole page instead of the table fragment.
func (p DataGridProps) sortHref(col DataGridColumn) string {
	if p.BaseURL != "" {
		if href := p.Sort.HrefFor(p.BaseURL, col.Key, col.FirstDir); href != "" {
			return href
		}
	}
	return p.sortPageHref(col)
}

// sortLinkAttrs builds the HTMX attributes for a sortable header. hx-get points
// at the fragment endpoint while hx-push-url keeps the browser on the page URL.
func (p DataGridProps) sortLinkAttrs(col DataGridColumn) templ.Attributes {
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
