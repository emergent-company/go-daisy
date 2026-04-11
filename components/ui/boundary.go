package ui

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// ButtonWithBoundary wraps Button with a dev-mode component boundary annotation.
// gallery:token variant,size,typ,shape,icon,loading
// gallery:hint href:default(#)
func ButtonWithBoundary(href string, variant ButtonVariant, size ButtonSize, typ ButtonType, shape ButtonShape, icon string, loading bool) templ.Component {
	return devmode.ComponentBoundary("Button", map[string]any{
		"href":    href,
		"variant": string(variant),
		"size":    string(size),
		"type":    string(typ),
		"shape":   string(shape),
		"icon":    icon,
		"loading": loading,
	}, Button(href, variant, size, typ, shape, icon, loading))
}

// BadgeWithBoundary wraps Badge with a dev-mode component boundary annotation.
// gallery:token variant,style,size,icon,label
// gallery:hint label:default(Active)
func BadgeWithBoundary(variant BadgeIntent, style BadgeStyle, size BadgeSize, icon string, label string) templ.Component {
	props := BadgeProps{Label: label, Variant: variant, Style: style, Size: size, Icon: icon}
	return devmode.ComponentBoundary("Badge", props, Badge(props))
}

// StatusBadgeWithBoundary wraps StatusBadge with a dev-mode component boundary annotation.
// gallery:token status
// gallery:hint status:default(active)
func StatusBadgeWithBoundary(status string) templ.Component {
	return devmode.ComponentBoundary("StatusBadge", map[string]any{"status": status}, StatusBadge(status))
}

// AvatarWithBoundary wraps Avatar with a dev-mode component boundary annotation.
// gallery:token name,size
// gallery:hint name:default(Jane Smith)
func AvatarWithBoundary(name string, src string, size AvatarSize) templ.Component {
	return devmode.ComponentBoundary("Avatar", map[string]any{
		"name": name,
		"src":  src,
		"size": string(size),
	}, Avatar(name, src, size))
}

// CardWithBoundary wraps Card with a dev-mode component boundary annotation.
// gallery:token title
// gallery:hint title:default(Card Title)
func CardWithBoundary(title string) templ.Component {
	return devmode.ComponentBoundary("Card", map[string]any{"title": title}, Card(title))
}

// AlertWithBoundary wraps InlineAlert with a dev-mode component boundary annotation.
// gallery:token typ,message
// gallery:hint message:default(Operation completed successfully.)
func AlertWithBoundary(typ AlertType, message string) templ.Component {
	return devmode.ComponentBoundary("InlineAlert", map[string]any{
		"type":    string(typ),
		"message": message,
	}, InlineAlert(typ, message))
}

// AlertWithIconBoundary wraps InlineAlertWithIcon with a dev-mode component boundary annotation.
// gallery:token typ,message
// gallery:hint message:default(Operation completed successfully.)
func AlertWithIconBoundary(typ AlertType, icon string, message string) templ.Component {
	return devmode.ComponentBoundary("InlineAlertWithIcon", map[string]any{
		"type":    string(typ),
		"icon":    icon,
		"message": message,
	}, InlineAlertWithIcon(typ, icon, message))
}

// ToastWithBoundary wraps Toast with a dev-mode component boundary annotation.
// gallery:token typ,message
// gallery:hint message:default(Action completed successfully.)
func ToastWithBoundary(typ ToastType, message string) templ.Component {
	return devmode.ComponentBoundary("Toast", map[string]any{
		"type":    string(typ),
		"message": message,
	}, Toast(typ, message))
}

// PaginationWithBoundary wraps Pagination with a dev-mode component boundary annotation.
// gallery:token currentPage,totalPages
// gallery:hint currentPage:range(1,20,1)
// gallery:hint totalPages:range(1,20,1)
func PaginationWithBoundary(currentPage int, totalPages int, baseURL string, targetID string) templ.Component {
	return devmode.ComponentBoundary("Pagination", map[string]any{
		"currentPage": currentPage,
		"totalPages":  totalPages,
		"baseURL":     baseURL,
		"targetID":    targetID,
	}, Pagination(currentPage, totalPages, baseURL, targetID))
}

// StatCardWithBoundary wraps StatCard with a dev-mode component boundary annotation.
func StatCardWithBoundary(p StatCardProps) templ.Component {
	return devmode.ComponentBoundary("StatCard", p, StatCard(p))
}

// EmptyWithBoundary wraps Empty with a dev-mode component boundary annotation.
// gallery:token title,description
// gallery:hint title:default(Nothing here yet)
// gallery:hint description:default(Add some items to get started.)
func EmptyWithBoundary(icon string, title string, description string) templ.Component {
	return devmode.ComponentBoundary("Empty", map[string]any{
		"icon":        icon,
		"title":       title,
		"description": description,
	}, Empty(icon, title, description))
}

// LoaderWithBoundary wraps Loader with a dev-mode component boundary annotation.
func LoaderWithBoundary() templ.Component {
	return devmode.ComponentBoundary("Loader", nil, Loader())
}

// ActionMenuWithBoundary wraps ActionMenu with a dev-mode component boundary annotation.
// gallery:token items
// gallery:hint items:slice(3)
func ActionMenuWithBoundary(items []ActionMenuItem) templ.Component {
	return devmode.ComponentBoundary("ActionMenu", map[string]any{"itemCount": len(items)}, ActionMenu(items))
}

// FilterCardWithBoundary wraps FilterCard with a dev-mode component boundary annotation.
func FilterCardWithBoundary(props FilterCardProps) templ.Component {
	return devmode.ComponentBoundary("FilterCard", props, FilterCard(props))
}

// ProgressWithBoundary wraps Progress with a dev-mode component boundary annotation.
// gallery:token color,value,max
// gallery:hint value:range(0,100,1)
// gallery:hint value:default(70)
// gallery:hint max:range(1,200,1)
// gallery:hint max:default(100)
func ProgressWithBoundary(color ProgressColor, value int, max int) templ.Component {
	return devmode.ComponentBoundary("Progress", map[string]any{
		"color": string(color),
		"value": value,
		"max":   max,
	}, Progress(color, value, max))
}

// SkeletonWithBoundary wraps Skeleton with a dev-mode component boundary annotation.
// gallery:token classes
// gallery:hint classes:default(h-4 w-full)
func SkeletonWithBoundary(classes string) templ.Component {
	return devmode.ComponentBoundary("Skeleton", map[string]any{"classes": classes}, Skeleton(classes))
}

// SectionHeaderWithBoundary wraps SectionHeader with a dev-mode component boundary annotation.
// gallery:token title
// gallery:hint title:default(Personal Information)
func SectionHeaderWithBoundary(title string) templ.Component {
	return devmode.ComponentBoundary("SectionHeader", map[string]any{"title": title}, SectionHeader(title))
}

// NoPermissionsWithBoundary wraps NoPermissions with a dev-mode component boundary annotation.
func NoPermissionsWithBoundary() templ.Component {
	return devmode.ComponentBoundary("NoPermissions", nil, NoPermissions())
}

// StatusWithBoundary wraps Status with a dev-mode component boundary annotation.
// gallery:token color
// gallery:hint color:default(status-success)
func StatusWithBoundary(color StatusColor) templ.Component {
	return devmode.ComponentBoundary("Status", map[string]any{"color": string(color)}, Status(color))
}
