package pages

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// DataGridPageWithBoundary wraps DataGridPage with a dev-mode component boundary annotation.
// gallery:token title,subtitle,gridProps
// gallery:hint title:default(Users),subtitle:default(Manage users)
func DataGridPageWithBoundary(props DataGridPageProps) templ.Component {
	return devmode.ComponentBoundary("DataGridPage", DataGridPage(props), props)
}

// AuthPageWithBoundary wraps AuthPage with a dev-mode component boundary annotation.
// gallery:token style,brandName
// gallery:hint style:default(login),brandName:default(go-daisy)
func AuthPageWithBoundary(props AuthPageProps) templ.Component {
	return devmode.ComponentBoundary("AuthPage", AuthPage(props), props)
}

// DashboardPageWithBoundary wraps DashboardPage with a dev-mode component boundary annotation.
// gallery:token style
// gallery:hint style:default(crm)
func DashboardPageWithBoundary(props DashboardPageProps) templ.Component {
	return devmode.ComponentBoundary("DashboardPage", DashboardPage(props), props)
}

// ChatLayoutWithBoundary wraps ChatLayout with a dev-mode component boundary annotation.
// gallery:token activeConversation
// gallery:hint activeConversation:default(Alice Johnson)
func ChatLayoutWithBoundary(props ChatLayoutProps) templ.Component {
	return devmode.ComponentBoundary("ChatLayout", ChatLayout(props), props)
}

// SettingsPageWithBoundary wraps SettingsPage with a dev-mode component boundary annotation.
func SettingsPageWithBoundary() templ.Component {
	return devmode.ComponentBoundary("SettingsPage", SettingsPage())
}

// LandingPageWithBoundary wraps LandingPage with a dev-mode component boundary annotation.
// gallery:token brandName,tagline
// gallery:hint brandName:default(go-daisy)
func LandingPageWithBoundary(props LandingPageProps) templ.Component {
	return devmode.ComponentBoundary("LandingPage", LandingPage(props), props)
}
