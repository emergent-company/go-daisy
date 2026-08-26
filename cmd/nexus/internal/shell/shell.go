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

func sc(activeSlug, icon, label string, children []templ.Component) templ.Component {
	open := collapsibleOpen(activeSlug, icon)
	return layout.SidebarCollapsibleNexus(icon, label, children, open)
}

func emptyItem(activeSlug, label, href string) templ.Component {
	return mi(activeSlug, label, href, "")
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
