package compare

import "fmt"

// PageMapping links a go-daisy route to its nexus-html reference file.
type PageMapping struct {
	Route      string
	Name       string
	HTMLFile   string // relative path from nexus-html root
}

// Pages returns the full list of pages to compare.
func Pages() []PageMapping {
	return []PageMapping{
		{"/dashboards/ecommerce", "Ecommerce Dashboard", "dashboards-ecommerce.html"},
		{"/dashboards/crm", "CRM Dashboard", "dashboards-crm.html"},
		{"/apps/ecommerce/products", "Products List", "apps-ecommerce-products.html"},
		{"/apps/ecommerce/products/create", "Product Create", "apps-ecommerce-products-create.html"},
		{"/apps/ecommerce/orders", "Orders List", "apps-ecommerce-orders.html"},
		{"/apps/chat", "Chat", "apps-chat.html"},
		{"/apps/file-manager", "File Manager", "apps-file-manager.html"},
		{"/apps/gen-ai/home", "Gen AI Home", "apps-gen-ai-home.html"},
		{"/apps/gen-ai/content", "Gen AI Content", "apps-gen-ai-content.html"},
		{"/apps/gen-ai/image", "Gen AI Image", "apps-gen-ai-image.html"},
		{"/apps/gen-ai/library", "Gen AI Library", "apps-gen-ai-library.html"},
		{"/pages/settings", "Settings", "pages-settings.html"},
		{"/pages/get-help", "Get Help", "pages-get-help.html"},
		{"/auth/login", "Login", "auth-login.html"},
		{"/auth/register", "Register", "auth-register.html"},
		{"/landing", "Landing", "landing.html"},
	}
}

// NexusHTMLRoot is the absolute path to the nexus-html built output directory.
const NexusHTMLRoot = "/root/nexus-html/html"

// GoDaisyBaseURL is the base URL of the running go-daisy nexus server.
const GoDaisyBaseURL = "http://localhost:11001"

func refPath(file string) string {
	return fmt.Sprintf("%s/%s", NexusHTMLRoot, file)
}

func goDaisyURL(route string) string {
	return fmt.Sprintf("%s%s", GoDaisyBaseURL, route)
}
