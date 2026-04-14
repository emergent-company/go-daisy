package nav

import (
	"context"
	"io"

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

// BreadcrumbsWithBoundary wraps Breadcrumbs with a dev-mode component boundary annotation.
// gallery:token items
// gallery:hint items:slice(3)
func BreadcrumbsWithBoundary(items []BreadcrumbItem) templ.Component {
	return devmode.ComponentBoundary("Breadcrumbs", map[string]any{"itemCount": len(items)}, Breadcrumbs(items))
}

// DockWithBoundary wraps Dock with a dev-mode component boundary annotation.
// gallery:token items
// gallery:hint items:slice(4)
func DockWithBoundary(items []DockItem) templ.Component {
	return devmode.ComponentBoundary("Dock", map[string]any{"itemCount": len(items)}, Dock(items))
}

// LinkWithBoundary wraps Link with a dev-mode component boundary annotation.
// gallery:token variant
// gallery:hint variant:default(link)
func LinkWithBoundary(href string, variant LinkVariant, label string) templ.Component {
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, label)
		return err
	})
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Link(href, variant).Render(templ.WithChildren(ctx, child), w)
	})
	return devmode.ComponentBoundary("Link", map[string]any{
		"href":    href,
		"variant": string(variant),
		"label":   label,
	}, inner)
}

// PageTitleMinimalWithBoundary wraps PageTitleMinimal with a dev-mode component boundary annotation.
func PageTitleMinimalWithBoundary(title string, steps []PageTitleStep) templ.Component {
	return devmode.ComponentBoundary("PageTitleMinimal", map[string]any{
		"title":     title,
		"stepCount": len(steps),
	}, PageTitleMinimal(title, steps))
}

// PageTitleEditorWithBoundary wraps PageTitleEditor with a dev-mode component boundary annotation.
func PageTitleEditorWithBoundary(steps []BreadcrumbStep, title, subtitle string, actions []PageTitleEditorAction) templ.Component {
	return devmode.ComponentBoundary("PageTitleEditor", map[string]any{
		"title":       title,
		"subtitle":    subtitle,
		"stepCount":   len(steps),
		"actionCount": len(actions),
	}, PageTitleEditor(steps, title, subtitle, actions))
}

// FooterMinimalWithBoundary wraps FooterMinimal with a dev-mode component boundary annotation.
func FooterMinimalWithBoundary(copyright string, links []FooterLink) templ.Component {
	return devmode.ComponentBoundary("FooterMinimal", map[string]any{
		"copyright": copyright,
		"linkCount": len(links),
	}, FooterMinimal(copyright, links))
}

// ProfileMenuWithBoundary wraps ProfileMenu with a dev-mode component boundary annotation.
func ProfileMenuWithBoundary(name, email, initials string, items []ProfileMenuItem, signOutHref string) templ.Component {
	return devmode.ComponentBoundary("ProfileMenu", map[string]any{
		"name":        name,
		"email":       email,
		"initials":    initials,
		"itemCount":   len(items),
		"signOutHref": signOutHref,
	}, ProfileMenu(name, email, initials, items, signOutHref))
}
