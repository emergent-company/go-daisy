package nav

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// PageHeaderWithBoundary wraps PageHeader with a dev-mode component boundary annotation.
// gallery:token steps
// gallery:hint steps:slice(3)
func PageHeaderWithBoundary(steps []BreadcrumbStep) templ.Component {
	return devmode.ComponentBoundary("PageHeader", map[string]any{"stepCount": len(steps)}, PageHeader(steps))
}

// TabMenuWithBoundary wraps TabMenu with a dev-mode component boundary annotation.
// gallery:token tabs
// gallery:hint tabs:slice(3)
func TabMenuWithBoundary(tabs []Tab, target ...string) templ.Component {
	return devmode.ComponentBoundary("TabMenu", map[string]any{"tabCount": len(tabs)}, TabMenu(tabs, target...))
}

// SimpleTabsWithBoundary wraps SimpleTabs with a dev-mode component boundary annotation.
// gallery:token tabs
// gallery:hint tabs:slice(3)
func SimpleTabsWithBoundary(tabs []Tab) templ.Component {
	return devmode.ComponentBoundary("SimpleTabs", map[string]any{"tabCount": len(tabs)}, SimpleTabs(tabs))
}

// TopBarWithBoundary wraps TopBar with a dev-mode component boundary annotation.
// gallery:token title
// gallery:hint title:default(My Application)
func TopBarWithBoundary(title string) templ.Component {
	return devmode.ComponentBoundary("TopBar", map[string]any{"title": title}, TopBar(title))
}

// MenuWithBoundary wraps Menu with a dev-mode component boundary annotation.
// gallery:token size,items
// gallery:hint items:slice(4)
func MenuWithBoundary(size MenuSize, items []MenuItem) templ.Component {
	return devmode.ComponentBoundary("Menu", map[string]any{
		"size":      string(size),
		"itemCount": len(items),
	}, Menu(size, items))
}
