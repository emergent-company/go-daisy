package table

import "github.com/a-h/templ"

// ListAreaProps configures a self-contained filterable+paginated list area.
type ListAreaProps struct {
	ID          string // container and form ID prefix, e.g. "sessions-list"
	URL         string // hx-get URL for filter/pagination requests
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	PageSize    int
	ScrollLoad  bool // when true, renders infinite scroll instead of page buttons
	ColSpan     int  // number of columns; used by the scroll-load sentinel row
}

// ScrollSentinelProps configures the sentinel element for scroll-load mode.
type ScrollSentinelProps struct {
	ID       string
	URL      string
	NextPage int
	ColSpan  int
}

// ScrollRowsProps is the server response payload for a scroll-load fetch.
type ScrollRowsProps struct {
	ID       string
	URL      string
	NextPage int
	HasMore  bool
	ColSpan  int
	Rows     templ.Component
}
