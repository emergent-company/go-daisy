package compare

import (
	"fmt"
	"strings"
)

// PageMapping links a go-daisy route to its nexus-html reference file.
type PageMapping struct {
	Route    string
	Name     string
	HTMLFile string // relative path from nexus-html root
	Skip     bool
	Reason   string
}

// pages is a small constructor for a regular (non-skipped) mapping.
func page(route, name, file string) PageMapping {
	return PageMapping{Route: route, Name: name, HTMLFile: file}
}

// skip is a constructor for a skipped reference file with a reason.
func skip(name, file, reason string) PageMapping {
	return PageMapping{Name: name, HTMLFile: file, Skip: true, Reason: reason}
}

// Pages returns the full list of pages to compare: every reference file in
// nxhtml/*.html mapped to a route, with index.html and starter.html skipped.
//
// The reference set is 69 files. Two are not comparable and are marked Skip:
// index.html (a duplicate of dashboards-ecommerce.html) and starter.html
// (no route registered). The remaining 67 map 1:1 onto a route.
func Pages() []PageMapping {
	pages := make([]PageMapping, 0, 69)

	// Dashboards
	pages = append(pages,
		page("/dashboards/ecommerce", "Ecommerce Dashboard", "dashboards-ecommerce.html"),
		page("/dashboards/crm", "CRM Dashboard", "dashboards-crm.html"),
	)

	// Ecommerce CRUD: list + create (+/create) + edit (+/1) per resource.
	// Orders has no create/edit; its detail page uses a bespoke file name.
	for _, c := range []struct{ route, name, file string }{
		{"products", "Products", "apps-ecommerce-products"},
		{"sellers", "Sellers", "apps-ecommerce-sellers"},
		{"customers", "Customers", "apps-ecommerce-customers"},
		{"shops", "Shops", "apps-ecommerce-shops"},
	} {
		pages = append(pages,
			page(fmt.Sprintf("/apps/ecommerce/%s", c.route), c.name+" List", c.file+".html"),
			page(fmt.Sprintf("/apps/ecommerce/%s/create", c.route), c.name+" Create", c.file+"-create.html"),
			page(fmt.Sprintf("/apps/ecommerce/%s/1", c.route), c.name+" Edit", c.file+"-edit.html"),
		)
	}
	pages = append(pages,
		page("/apps/ecommerce/orders", "Orders List", "apps-ecommerce-orders.html"),
		page("/apps/ecommerce/orders/1", "Order Details", "apps-ecommerce-order-details.html"),
	)

	// Apps
	pages = append(pages,
		page("/apps/chat", "Chat", "apps-chat.html"),
		page("/apps/file-manager", "File Manager", "apps-file-manager.html"),
		page("/apps/gen-ai/home", "Gen AI Home", "apps-gen-ai-home.html"),
		page("/apps/gen-ai/content", "Gen AI Content", "apps-gen-ai-content.html"),
		page("/apps/gen-ai/image", "Gen AI Image", "apps-gen-ai-image.html"),
		page("/apps/gen-ai/library", "Gen AI Library", "apps-gen-ai-library.html"),
	)

	// Site: auth, settings, get-help, landing
	pages = append(pages,
		page("/auth/login", "Login", "auth-login.html"),
		page("/auth/register", "Register", "auth-register.html"),
		page("/auth/forgot-password", "Forgot Password", "auth-forgot-password.html"),
		page("/auth/reset-password", "Reset Password", "auth-reset-password.html"),
		page("/pages/settings", "Settings", "pages-settings.html"),
		page("/pages/get-help", "Get Help", "pages-get-help.html"),
		page("/landing", "Landing", "landing.html"),
	)

	// UI showcase — generated from slug lists, one source of truth.
	for _, s := range []string{
		"accordion", "alert", "avatar", "badge", "breadcrumb", "button",
		"countdown", "drawer", "dropdown", "indicator", "loading", "menu",
		"modal", "pagination", "progress", "step", "tab", "table", "timeline",
		"toast", "tooltip",
	} {
		pages = append(pages, page("/ui/components/"+s, "UI Component "+title(s), "ui-components-"+s+".html"))
	}
	for _, s := range []string{
		"checkbox", "fieldset", "file-input", "input", "label", "radio",
		"range", "rating", "select", "textarea", "toggle", "validator",
	} {
		pages = append(pages, page("/ui/forms/"+s, "UI Form "+title(s), "ui-forms-"+s+".html"))
	}
	for _, s := range []string{"area", "bar", "column", "line", "pie"} {
		pages = append(pages, page("/ui/charts/"+s, "UI Chart "+title(s), "ui-charts-apex-"+s+".html"))
	}

	// Skipped references: present in the reference set but not comparable.
	pages = append(pages,
		skip("Index (duplicate)", "index.html", "duplicate of dashboards-ecommerce.html"),
		skip("Starter", "starter.html", "no route registered"),
	)

	return pages
}

// Compared returns only the mappings that are actually fetched and diffed.
func Compared(mappings []PageMapping) []PageMapping {
	var out []PageMapping
	for _, m := range mappings {
		if !m.Skip {
			out = append(out, m)
		}
	}
	return out
}

// title converts a kebab-case slug into Title Case.
func title(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

// NexusHTMLRoot is the absolute path to the nexus-html built output directory.
// It is a var so the compare entrypoint can override it via SetConfig.
var NexusHTMLRoot = "/root/nexus-html/html"

// GoDaisyBaseURL is the base URL of the running go-daisy nexus server.
var GoDaisyBaseURL = "http://localhost:11001"

// Config holds runtime configuration for the comparator.
type Config struct {
	BaseURL string
	RefDir  string
}

// SetConfig overrides the default base URL and reference directory.
func SetConfig(c Config) {
	if c.BaseURL != "" {
		GoDaisyBaseURL = c.BaseURL
	}
	if c.RefDir != "" {
		NexusHTMLRoot = c.RefDir
	}
}

func refPath(file string) string {
	return fmt.Sprintf("%s/%s", NexusHTMLRoot, file)
}

func goDaisyURL(route string) string {
	return fmt.Sprintf("%s%s", GoDaisyBaseURL, route)
}
