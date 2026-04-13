package table

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// TableWithBoundary wraps Table with a dev-mode component boundary annotation.
func TableWithBoundary() templ.Component {
	return devmode.ComponentBoundary("Table", nil, Table())
}

// TableWithPropsWithBoundary wraps TableWithProps with a dev-mode component boundary annotation.
func TableWithPropsWithBoundary(props TableProps) templ.Component {
	return devmode.ComponentBoundary("TableWithProps", props, TableWithProps(props))
}

// TableHeadWithBoundary wraps TableHead with a dev-mode element boundary annotation.
// Uses ElementBoundary so the data-component attribute is placed on the <thead> element
// itself, which is required because a <div> wrapper inside a <table> is invalid HTML.
func TableHeadWithBoundary() templ.Component {
	return devmode.ElementBoundary("TableHead", nil, TableHead())
}

// TableHeadRowWithBoundary wraps TableHeadRow with a dev-mode element boundary annotation.
func TableHeadRowWithBoundary() templ.Component {
	return devmode.ElementBoundary("TableHeadRow", nil, TableHeadRow())
}

// TableHeadCellWithBoundary wraps TableHeadCell with a dev-mode element boundary annotation.
func TableHeadCellWithBoundary(label string) templ.Component {
	return devmode.ElementBoundary("TableHeadCell", map[string]any{"label": label}, TableHeadCell(label))
}

// TableBodyWithBoundary wraps TableBody with a dev-mode element boundary annotation.
func TableBodyWithBoundary() templ.Component {
	return devmode.ElementBoundary("TableBody", nil, TableBody())
}

// TableRowWithBoundary wraps TableRow with a dev-mode element boundary annotation.
func TableRowWithBoundary(id string, hover bool) templ.Component {
	return devmode.ElementBoundary("TableRow", map[string]any{"id": id, "hover": hover}, TableRow(id, hover))
}

// TableCellWithBoundary wraps TableCell with a dev-mode element boundary annotation.
func TableCellWithBoundary(class string) templ.Component {
	return devmode.ElementBoundary("TableCell", map[string]any{"class": class}, TableCell(class))
}

// ListAreaWithBoundary wraps ListArea with a dev-mode component boundary annotation.
func ListAreaWithBoundary(props ListAreaProps) templ.Component {
	return devmode.ComponentBoundary("ListArea", props, ListArea(props))
}
