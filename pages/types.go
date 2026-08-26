package pages

import "github.com/emergent-company/go-daisy/shared"

// AuthPageStyle chooses the auth page layout.
type AuthPageStyle string

const (
	AuthLogin    AuthPageStyle = "login"
	AuthRegister AuthPageStyle = "register"
	AuthForgot   AuthPageStyle = "forgot"
	AuthReset    AuthPageStyle = "reset"
)

// AuthPageProps configures an authentication page.
type AuthPageProps struct {
	Style      AuthPageStyle
	LogoURL    string
	BrandName  string
	ErrorMsg   string
	SuccessMsg string
}

// DashboardStyle selects the dashboard variant.
type DashboardStyle string

const (
	DashboardCRM DashboardStyle = "crm"
)

// DashboardPageProps configures a dashboard page.
type DashboardPageProps struct {
	Style DashboardStyle
}

// DashboardProps is an alias for DashboardPageProps (used in the template).
type DashboardProps = DashboardPageProps

// ChatLayoutProps configures a chat layout page.
type ChatLayoutProps struct {
	ActiveConversation string
}

// LandingPageProps configures a landing page.
type LandingPageProps struct {
	BrandName string
	Tagline   string
}

func ternary(cond bool, a, b string) string { return shared.Ternary(cond, a, b) }
