package table

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// TableWithBoundary wraps Table with a dev-mode component boundary annotation.
func TableWithBoundary() templ.Component {
	return devmode.ComponentBoundary("Table", Table())
}

// TableWithPropsWithBoundary wraps TableWithProps with a dev-mode component boundary annotation.
func TableWithPropsWithBoundary(props TableProps) templ.Component {
	return devmode.ComponentBoundary("TableWithProps", TableWithProps(props), props)
}

// TableHeadWithBoundary wraps TableHead with a dev-mode element boundary annotation.
// Uses ElementBoundary so the data-component attribute is placed on the <thead> element
// itself, which is required because a <div> wrapper inside a <table> is invalid HTML.
func TableHeadWithBoundary() templ.Component {
	return devmode.ElementBoundary("TableHead", TableHead())
}

// TableHeadRowWithBoundary wraps TableHeadRow with a dev-mode element boundary annotation.
func TableHeadRowWithBoundary() templ.Component {
	return devmode.ElementBoundary("TableHeadRow", TableHeadRow())
}

// TableHeadCellWithBoundary wraps TableHeadCell with a dev-mode element boundary annotation.
func TableHeadCellWithBoundary(label string) templ.Component {
	return devmode.ElementBoundary("TableHeadCell", TableHeadCell(label), map[string]any{"label": label})
}

// TableHeaderWithBoundary wraps TableHeader with a dev-mode element boundary annotation.
// Uses ElementBoundary because TableHeader renders a <th> element directly.
//
// The canonical wrapper stays available even though TableHeader is deprecated;
// new code should use SortableTableHeaderWithBoundary.
func TableHeaderWithBoundary(label string, sortKey string, currentSortKey string, currentDir SortDir, baseURL string) templ.Component {
	return devmode.ElementBoundary("TableHeader", TableHeader(label, sortKey, currentSortKey, currentDir, baseURL), map[string]any{"label": label, "sortKey": sortKey})
}

// SortableTableHeaderWithBoundary wraps SortableTableHeader with a dev-mode
// element boundary annotation. Uses ElementBoundary because SortableTableHeader
// renders a <th> element directly.
func SortableTableHeaderWithBoundary(props SortableTableHeaderProps) templ.Component {
	return devmode.ElementBoundary("SortableTableHeader", SortableTableHeader(props), props)
}

// TableBodyWithBoundary wraps TableBody with a dev-mode element boundary annotation.
func TableBodyWithBoundary() templ.Component {
	return devmode.ElementBoundary("TableBody", TableBody())
}

// TableRowWithBoundary wraps TableRow with a dev-mode element boundary annotation.
func TableRowWithBoundary(id string, hover bool) templ.Component {
	return devmode.ElementBoundary("TableRow", TableRow(id, hover, nil), map[string]any{"id": id, "hover": hover})
}

// TableCellWithBoundary wraps TableCell with a dev-mode element boundary annotation.
func TableCellWithBoundary(class string) templ.Component {
	return devmode.ElementBoundary("TableCell", TableCell(class, nil), map[string]any{"class": class})
}

// TableSelectAllWithBoundary wraps TableSelectAll with a dev-mode element boundary annotation.
// gallery:token label
// gallery:hint label:default(Select all)
func TableSelectAllWithBoundary(label string) templ.Component {
	return devmode.ElementBoundary("TableSelectAll", TableSelectAll(label), map[string]any{"label": label})
}

// TableSelectRowWithBoundary wraps TableSelectRow with a dev-mode element boundary annotation.
// gallery:token name,checked
// gallery:hint name:default(row)
func TableSelectRowWithBoundary(name string, checked bool) templ.Component {
	return devmode.ElementBoundary("TableSelectRow", TableSelectRow(name, checked), map[string]any{"name": name, "checked": checked})
}

// ListAreaWithBoundary wraps ListArea with a dev-mode component boundary annotation.
func ListAreaWithBoundary(props ListAreaProps) templ.Component {
	return devmode.ComponentBoundary("ListArea", ListArea(props), props)
}

// TableEmptyWithBoundary wraps TableEmpty with a dev-mode element boundary annotation.
func TableEmptyWithBoundary(colspan int, message string) templ.Component {
	return devmode.ElementBoundary("TableEmpty", TableEmpty(colspan, message), map[string]any{"colspan": colspan, "message": message})
}

// TableCardWithBoundary wraps TableCard with a dev-mode component boundary annotation.
// Toolbar and Footer are templ.Component slots, so there is no gallery token
// surface to expose.
func TableCardWithBoundary(props TableCardProps) templ.Component {
	return devmode.ComponentBoundary("TableCard", TableCard(props), map[string]any{
		"toolbar": props.Toolbar != nil,
		"footer":  props.Footer != nil,
	})
}

// TableCardWrapperWithBoundary wraps TableCardWrapper with a dev-mode component boundary annotation.
func TableCardWrapperWithBoundary() templ.Component {
	return devmode.ComponentBoundary("TableCardWrapper", TableCardWrapper())
}

// TableCardHeaderWithBoundary wraps TableCardHeader with a dev-mode component boundary annotation.
func TableCardHeaderWithBoundary() templ.Component {
	return devmode.ComponentBoundary("TableCardHeader", TableCardHeader())
}

// TableCardFooterWithBoundary wraps TableCardFooter with a dev-mode component boundary annotation.
func TableCardFooterWithBoundary(props TableCardFooterProps) templ.Component {
	return devmode.ComponentBoundary("TableCardFooter", TableCardFooter(props), props)
}

// TableSearchWithBoundary wraps TableSearch with a dev-mode component boundary annotation.
// gallery:token placeholder
// gallery:hint placeholder:default(Search orders)
func TableSearchWithBoundary(name string, value string, placeholder string, hxTarget string, hxGet string) templ.Component {
	return devmode.ComponentBoundary("TableSearch", TableSearch(name, value, placeholder, hxTarget, hxGet), map[string]any{
		"name":        name,
		"placeholder": placeholder,
	})
}

// TableSelectWithBoundary wraps TableSelect with a dev-mode component boundary annotation.
func TableSelectWithBoundary(name string, hxTarget string, hxGet string) templ.Component {
	return devmode.ComponentBoundary("TableSelect", TableSelect(name, hxTarget, hxGet), map[string]any{"name": name})
}

// DataTableWithBoundary wraps DataTable with a dev-mode component boundary annotation.
// gallery:token searchable,sortable,striped,compact,pageSize
// gallery:hint pageSize:range(1,100,1)
func DataTableWithBoundary(props DataTableProps) templ.Component {
	return devmode.ComponentBoundary("DataTable", DataTable(props), map[string]any{
		"searchable": props.Searchable,
		"sortable":   props.Sortable,
		"striped":    props.Striped,
		"compact":    props.Compact,
		"pageSize":   props.PageSize,
	})
}

// DataGridWithBoundary wraps DataGrid with a dev-mode component boundary annotation.
// gallery:token striped,bordered,compact,selectable,pageSize
// gallery:hint pageSize:range(1,100,1)
func DataGridWithBoundary(props DataGridProps) templ.Component {
	return devmode.ComponentBoundary("DataGrid", DataGrid(props), map[string]any{
		"searchable": props.SearchName != "",
		"striped":    props.Striped,
		"compact":    props.Compact,
		"selectable": props.Selectable,
		"pageSize":   props.PageSize,
	})
}
