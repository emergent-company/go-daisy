package shell

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/components/layout"
	"github.com/emergent-company/go-daisy/components/nav"
)

func mi(activeSlug, label, href, icon string) templ.Component {
	if href == "#" || href == "" {
		return layout.SidebarMenuItemNexus(label, href, icon, false)
	}
	active := strings.HasPrefix(activeSlug, href) && href != "/"
	return layout.SidebarMenuItemNexus(label, href, icon, active)
}

func sc(activeSlug, icon, label string, children []templ.Component, badge ...string) templ.Component {
	open := collapsibleOpen(activeSlug, icon)
	props := layout.SidebarCollapsibleNexusProps{Icon: icon, Label: label, Children: children, Open: open}
	if len(badge) > 0 {
		props.Badge = badge[0]
	}
	return layout.SidebarCollapsibleNexusWithProps(props)
}

func emptyItem(activeSlug, label, href string) templ.Component {
	return mi(activeSlug, label, href, "")
}

func emptyItemBadge(activeSlug, label, href, badge string) templ.Component {
	active := strings.HasPrefix(activeSlug, href) && href != "/"
	return layout.SidebarMenuItemNexusWithProps(layout.SidebarMenuItemNexusProps{Label: label, Href: href, Active: active, Badge: badge})
}

func collapsibleOpen(activeSlug, icon string) bool {
	switch icon {
	case "lucide--monitor-dot":
		return strings.HasPrefix(activeSlug, "/dashboards/")
	case "lucide--store":
		return strings.HasPrefix(activeSlug, "/apps/ecommerce/")
	case "lucide--brain-circuit":
		return strings.HasPrefix(activeSlug, "/apps/gen-ai/")
	case "lucide--shield-check":
		return strings.HasPrefix(activeSlug, "/auth/")
	case "lucide--package":
		return strings.HasPrefix(activeSlug, "/ui/components/")
	case "lucide--file-text":
		return strings.HasPrefix(activeSlug, "/ui/forms/")
	case "lucide--bar-chart-big":
		return strings.HasPrefix(activeSlug, "/ui/charts/")
	default:
		return false
	}
}

func nexusNotifications() []nav.NotificationItem {
	return []nav.NotificationItem{
		{AvatarSrc: "https://img.daisyui.com/images/profile/demo/1@94.webp", AvatarAlt: "User", Name: "Customer", Message: "Customer has requested a return for item", TimeAgo: "1 Hour ago"},
		{AvatarSrc: "https://img.daisyui.com/images/profile/demo/2@94.webp", AvatarAlt: "User", Name: "Review", Message: "A new review has been submitted for product", TimeAgo: "1 Hour ago"},
		{AvatarSrc: "https://img.daisyui.com/images/profile/demo/3@94.webp", AvatarAlt: "User", Name: "Promotion", Message: "Prepare for the upcoming weekend promotion", TimeAgo: "2 Days ago"},
		{AvatarSrc: "https://img.daisyui.com/images/profile/demo/4@94.webp", AvatarAlt: "User", Name: "Stock", Message: "Product 'ABC123' is running low in stock.", TimeAgo: "3 Days ago"},
		{AvatarSrc: "https://img.daisyui.com/images/profile/demo/5@94.webp", AvatarAlt: "User", Name: "Payment", Message: "Payment received for Order ID: #67890", TimeAgo: "Week ago"},
	}
}

// apexGroup renders the reference's nested "Apex" group under Charts.
func apexGroup(activeSlug string) templ.Component {
	open := strings.HasPrefix(activeSlug, "/ui/charts/")
	children := []templ.Component{
		emptyItem(activeSlug, "Area", "/ui/charts/area"),
		emptyItem(activeSlug, "Bar", "/ui/charts/bar"),
		emptyItem(activeSlug, "Column", "/ui/charts/column"),
		emptyItem(activeSlug, "Line", "/ui/charts/line"),
		emptyItem(activeSlug, "Pie", "/ui/charts/pie"),
	}
	return layout.SidebarCollapsibleNexusWithProps(layout.SidebarCollapsibleNexusProps{Label: "Apex", Children: children, Open: open})
}
