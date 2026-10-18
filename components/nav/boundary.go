package nav

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
	"github.com/emergent-company/go-daisy/shared"
)

// PageHeaderWithBoundary wraps PageHeader with a dev-mode component boundary annotation.
// gallery:token steps
// gallery:hint steps:slice(3)
func PageHeaderWithBoundary(steps []BreadcrumbItem) templ.Component {
	return devmode.ComponentBoundary("PageHeader", PageHeader(steps, nil), map[string]any{"stepCount": len(steps)})
}

// TabMenuWithBoundary wraps TabMenu with a dev-mode component boundary annotation.
// gallery:token tabs
// gallery:hint tabs:slice(3)
func TabMenuWithBoundary(tabs []Tab, target ...string) templ.Component {
	if len(target) > 0 {
		return devmode.ComponentBoundary("TabMenu", TabMenu(tabs, nil, target[0]), map[string]any{"tabCount": len(tabs)})
	}
	return devmode.ComponentBoundary("TabMenu", TabMenu(tabs, nil, ""), map[string]any{"tabCount": len(tabs)})
}

// SimpleTabsWithBoundary wraps SimpleTabs with a dev-mode component boundary annotation.
// gallery:token tabs
// gallery:hint tabs:slice(3)
func SimpleTabsWithBoundary(tabs []Tab) templ.Component {
	return devmode.ComponentBoundary("SimpleTabs", SimpleTabs(tabs), map[string]any{"tabCount": len(tabs)})
}

// TopBarWithBoundary wraps TopBar with a dev-mode component boundary annotation.
// gallery:token title,scrollAware
// gallery:hint title:default(My Application)
// gallery:hint scrollAware:default(false)
func TopBarWithBoundary(title string, scrollAware bool) templ.Component {
	props := TopBarProps{Title: title, ScrollAware: scrollAware}
	return devmode.ComponentBoundary("TopBar", TopBar(props), map[string]any{
		"title":       title,
		"scrollAware": scrollAware,
	})
}

// MenuWithBoundary wraps Menu with a dev-mode component boundary annotation.
// gallery:token size,items
// gallery:hint items:slice(4)
func MenuWithBoundary(size MenuSize, items []MenuItem) templ.Component {
	return devmode.ComponentBoundary("Menu", Menu(size, items), map[string]any{
		"size":      string(size),
		"itemCount": len(items),
	})
}

// BreadcrumbsWithBoundary wraps Breadcrumbs with a dev-mode component boundary annotation.
// gallery:token items
// gallery:hint items:slice(3)
func BreadcrumbsWithBoundary(items []BreadcrumbItem, divider BreadcrumbsDivider) templ.Component {
	return devmode.ComponentBoundary("Breadcrumbs", Breadcrumbs(items, divider), map[string]any{"itemCount": len(items), "divider": string(divider)})
}

// DockWithBoundary wraps Dock with a dev-mode component boundary annotation.
// gallery:token items
// gallery:hint items:slice(4)
func DockWithBoundary(items []DockItem) templ.Component {
	return devmode.ComponentBoundary("Dock", Dock(items), map[string]any{"itemCount": len(items)})
}

// LinkWithBoundary wraps Link with a dev-mode component boundary annotation.
// gallery:token variant
// gallery:hint variant:default(link)
func LinkWithBoundary(href string, variant LinkVariant, label string) templ.Component {
	inner := shared.RenderInto(Link(href, variant, "", nil), shared.StrComp(label))
	return devmode.ComponentBoundary("Link", inner, map[string]any{
		"href":    href,
		"variant": string(variant),
		"label":   label,
	})
}

// PageHeadingWithBoundary wraps PageHeading with a dev-mode component boundary annotation.
func PageHeadingWithBoundary(props PageHeadingProps) templ.Component {
	return devmode.ComponentBoundary("PageHeading", PageHeading(props), map[string]any{
		"title":           props.Title,
		"breadcrumbCount": len(props.Breadcrumbs),
		"kicker":          props.Kicker,
		"subtitle":        props.Subtitle,
		"margin":          props.Margin,
		"dashboard":       props.Dashboard,
		"bare":            props.Bare,
		"subtitleFull":    props.SubtitleFull,
		"hideBreadcrumbs": props.HideBreadcrumbs,
		"noTopMargin":     props.NoTopMargin,
		"flat":            props.Flat,
	})
}

// FooterMinimalWithBoundary wraps FooterMinimal with a dev-mode component boundary annotation.
func FooterMinimalWithBoundary(copyright string, links []FooterLink) templ.Component {
	return devmode.ComponentBoundary("FooterMinimal", FooterMinimal(copyright, links), map[string]any{
		"copyright": copyright,
		"linkCount": len(links),
	})
}

// ProfileMenuWithBoundary wraps ProfileMenu with a dev-mode component boundary annotation.
func ProfileMenuWithBoundary(name, email, initials string, items []ProfileMenuItem, signOutHref string) templ.Component {
	return devmode.ComponentBoundary("ProfileMenu", ProfileMenu(name, email, initials, items, signOutHref), map[string]any{
		"name":        name,
		"email":       email,
		"initials":    initials,
		"itemCount":   len(items),
		"signOutHref": signOutHref,
	})
}

// FooterVariantWithBoundary wraps FooterVariant with a dev-mode boundary.
func FooterVariantWithBoundary(style string, opts FooterVariantOpts) templ.Component {
	return devmode.ComponentBoundary("FooterVariant", FooterVariant(style, opts), map[string]any{
		"style":   style,
		"version": 1,
	})
}

// NotificationVariantWithBoundary wraps NotificationVariant with a dev-mode boundary.
func NotificationVariantWithBoundary(style string, opts NotificationVariantOpts) templ.Component {
	return devmode.ComponentBoundary("NotificationVariant", NotificationVariant(style, opts), map[string]any{
		"style": style,
	})
}

// SearchModalWithBoundary wraps SearchModal with a dev-mode boundary.
func SearchModalWithBoundary(style string, opts SearchModalOpts) templ.Component {
	return devmode.ComponentBoundary("SearchModal", SearchModal(style, opts), map[string]any{
		"style": style,
	})
}

// SearchTriggerWithBoundary wraps SearchTrigger with a dev-mode component boundary annotation.
// gallery:token dialogId,placeholder
// gallery:hint dialogId:default(search-modal)
// gallery:hint placeholder:default(Search)
func SearchTriggerWithBoundary(dialogID, placeholder string) templ.Component {
	props := SearchTriggerProps{DialogID: dialogID, Placeholder: placeholder}
	return devmode.ComponentBoundary("SearchTrigger", SearchTrigger(props), map[string]any{
		"dialogId":    dialogID,
		"placeholder": placeholder,
	})
}

// ProfileMenuVariantWithBoundary wraps ProfileMenuVariant with a dev-mode boundary.
func ProfileMenuVariantWithBoundary(style string, opts ProfileMenuVariantOpts) templ.Component {
	return devmode.ComponentBoundary("ProfileMenuVariant", ProfileMenuVariant(style, opts), map[string]any{
		"style": style,
	})
}

// PageTitleVariantWithBoundary wraps PageTitleVariant with a dev-mode boundary.
func PageTitleVariantWithBoundary(style string, opts PageTitleVariantOpts) templ.Component {
	return devmode.ComponentBoundary("PageTitleVariant", PageTitleVariant(style, opts), map[string]any{
		"style": style,
	})
}

// MegamenuWithBoundary wraps Megamenu with a dev-mode component boundary annotation.
func MegamenuWithBoundary(items []MegamenuItem) templ.Component {
	return devmode.ComponentBoundary("Megamenu", Megamenu(items), map[string]any{
		"itemCount": len(items),
	})
}

// MenuSectionWithBoundary wraps MenuSection with a dev-mode component boundary annotation.
// gallery:token title,items
// gallery:hint title:default(Section)
// gallery:hint items:slice(3)
func MenuSectionWithBoundary(title string, items []MenuItem) templ.Component {
	return devmode.ComponentBoundary("MenuSection", MenuSection(title, items), map[string]any{
		"title":     title,
		"itemCount": len(items),
	})
}
