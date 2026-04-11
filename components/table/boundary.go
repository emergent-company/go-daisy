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

// ListAreaWithBoundary wraps ListArea with a dev-mode component boundary annotation.
func ListAreaWithBoundary(props ListAreaProps) templ.Component {
	return devmode.ComponentBoundary("ListArea", props, ListArea(props))
}
