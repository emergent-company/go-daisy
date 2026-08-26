package shell

import (
	"github.com/emergent-company/go-daisy/components/layout"
	"github.com/emergent-company/go-daisy/components/nav"
)

type NexusShellData struct {
	Title         string
	ActiveSlug    string
	Breadcrumbs   []nav.BreadcrumbItem
	SidebarGroups []layout.SidebarGroup
}
