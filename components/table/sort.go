package table

import (
	"net/http"
	"net/url"
	"strings"
)

// SortContext bundles the parsed sort state with the request path and query
// needed to build sort links. One context is used per table; multi-table pages
// use a distinct Prefix for each table (for example "nodes_").
type SortContext struct {
	Params   SortParams
	BasePath string     // escaped request path without a query string
	Query    url.Values // full request query, including filters
	Prefix   string     // optional query-parameter prefix, e.g. "nodes_"
}

// NewSortContext parses sort state from r and captures the path/query used to
// build sort links. defaultBy/defaultDir are applied when the request has no
// sort parameters. Pass an empty defaultBy and SortNone to preserve the
// handler's natural row order until a column is chosen.
func NewSortContext(r *http.Request, prefix string, defaultBy string, defaultDir SortDir) SortContext {
	if r == nil {
		return SortContext{Prefix: prefix}
	}
	q := r.URL.Query()
	return SortContext{
		Params:   ParseSortPrefixed(q, prefix, defaultBy, defaultDir),
		BasePath: r.URL.EscapedPath(),
		Query:    q,
		Prefix:   prefix,
	}
}

// ByParam returns the query-parameter name holding the sort column.
func (c SortContext) ByParam() string {
	return c.Prefix + "sort_by"
}

// DirParam returns the query-parameter name holding the sort direction.
func (c SortContext) DirParam() string {
	return c.Prefix + "sort_dir"
}

// Active reports whether column is the active sort column.
func (c SortContext) Active(column string) bool {
	return column != "" && c.Params.By == column && c.Params.Dir != SortNone
}

// AriaSort returns the aria-sort value for column, or "" when the column is not
// the active sort column. An empty value must result in the attribute being
// omitted entirely, not rendered as aria-sort="".
func (c SortContext) AriaSort(column string) string {
	if !c.Active(column) {
		return ""
	}
	if c.Params.Dir == SortAsc {
		return "ascending"
	}
	return "descending"
}

// ariaLabel returns the accessible name for a sort link, including the current
// direction when the column is active.
func (p SortableTableHeaderProps) ariaLabel() string {
	label := "Sort by " + p.Label
	if p.Sort.Active(p.Column) {
		if p.Sort.Params.Dir == SortAsc {
			label += ", currently ascending"
		} else {
			label += ", currently descending"
		}
	}
	switch NextSortDirFrom(p.Sort.Params, p.Column, p.FirstDir) {
	case SortAsc:
		label += "; activate to sort ascending"
	case SortDesc:
		label += "; activate to sort descending"
	default:
		label += "; activate to clear sorting"
	}
	return label
}

// Href returns the path and query for the next sort state of column, using the
// request path captured by NewSortContext. The current query is cloned, so
// filters are preserved and the caller's query is never mutated. first is the
// direction used when column is not active. Pagination is reset by removing the
// "page" param; other params, including per_page, are preserved.
func (c SortContext) Href(column string, first SortDir) string {
	return c.HrefFor(c.BasePath, column, first)
}

// HrefFor is Href with an explicit base path. Use it when the sort link must
// target a different endpoint than the page itself, such as a fragment URL that
// only returns the table. The context query is merged with the base query (the
// base wins on conflicts) and the sort parameters are merged into the result;
// page is removed.
func (c SortContext) HrefFor(basePath string, column string, first SortDir) string {
	if basePath == "" || column == "" {
		return ""
	}
	u, err := url.Parse(basePath)
	if err != nil || u.Path == "" {
		return ""
	}
	dir := NextSortDirFrom(c.Params, column, first)
	q := cloneValues(c.Query)
	for key, values := range u.Query() {
		q[key] = values
	}
	if dir == SortNone {
		q.Del(c.ByParam())
		q.Del(c.DirParam())
	} else {
		q.Set(c.ByParam(), column)
		q.Set(c.DirParam(), string(dir))
	}
	q.Del("page")
	u.RawQuery = q.Encode()
	return u.String()
}

// ParseSortPrefixed extracts <prefix>sort_by and <prefix>sort_dir from query
// params. Defaults to defaultBy and defaultDir when the params are absent. An
// empty sort column always results in SortNone, so callers can treat "no sort"
// as "keep the natural row order".
//
// A query sort_by with an absent or invalid sort_dir uses defaultDir; if
// defaultDir is SortNone or invalid, it becomes ascending. sort_dir without a
// sort_by is ignored. A handler-supplied defaultBy with defaultDir SortNone
// keeps SortNone, so the handler's natural row order is preserved until a
// click. An invalid defaultDir is treated as SortNone.
func ParseSortPrefixed(q url.Values, prefix string, defaultBy string, defaultDir SortDir) SortParams {
	requestedBy := strings.TrimSpace(q.Get(prefix + "sort_by"))
	by := requestedBy
	if by == "" {
		by = defaultBy
	}

	dir := defaultDir
	if dir != SortAsc && dir != SortDesc {
		dir = SortNone
	}
	// A direction without a column is meaningless: it must not activate a
	// handler-supplied defaultBy.
	if requestedBy != "" {
		switch strings.ToLower(strings.TrimSpace(q.Get(prefix + "sort_dir"))) {
		case "asc":
			dir = SortAsc
		case "desc":
			dir = SortDesc
		}
	}
	if by == "" {
		dir = SortNone
	} else if requestedBy != "" && dir == SortNone {
		// A user-requested column without a usable direction means ascending,
		// not an active column with no order. A handler-supplied defaultBy with
		// SortNone keeps the natural row order.
		dir = SortAsc
	}
	return SortParams{By: by, Dir: dir}
}

// NextSortDirFrom returns the next direction for a column click:
//
//	not active      -> first
//	first direction -> the opposite direction
//	opposite        -> SortNone (clear the sort and return to the default order)
//
// first defaults to SortAsc when it is not SortAsc or SortDesc.
func NextSortDirFrom(sp SortParams, column string, first SortDir) SortDir {
	if first != SortAsc && first != SortDesc {
		first = SortAsc
	}
	if column == "" || sp.By != column || sp.Dir == SortNone {
		return first
	}
	if sp.Dir == first {
		if first == SortAsc {
			return SortDesc
		}
		return SortAsc
	}
	return SortNone
}
