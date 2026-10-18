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
// gallery:token activeConversation,searchPlaceholder,communityLabel,newContactHref
// gallery:hint activeConversation:default(Alice Johnson),searchPlaceholder:default(Search conversations...),communityLabel:default(Join a Community),newContactHref:default(#)
func ChatLayoutWithBoundary(props ChatLayoutProps) templ.Component {
	return devmode.ComponentBoundary("ChatLayout", ChatLayout(props), props)
}

// ChatContactItemWithBoundary wraps ChatContactItem with a dev-mode component boundary annotation.
// gallery:token name,time,preview,unread,active
// gallery:hint name:default(Alice Johnson),time:default(2 min ago),preview:default(Sounds good! Let me know when...),unread:range(0,9,1),active:default(true)
func ChatContactItemWithBoundary(props ChatContactProps) templ.Component {
	return devmode.ComponentBoundary("ChatContactItem", ChatContactItem(props), props)
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

// LandingTopbarWithBoundary wraps LandingTopbar with a dev-mode component boundary annotation.
// gallery:token brand,ctaLabel,showTheme
// gallery:hint brand:default(go-daisy),ctaLabel:default(Get Started),showTheme:default(true)
func LandingTopbarWithBoundary(brand, ctaLabel string, showTheme bool, links []LandingNavLink) templ.Component {
	props := LandingTopbarProps{Brand: brand, CTALabel: ctaLabel, ShowTheme: showTheme, Links: links}
	return devmode.ComponentBoundary("LandingTopbar", LandingTopbar(props), map[string]any{
		"brand": brand, "ctaLabel": ctaLabel, "showTheme": showTheme, "linkCount": len(links),
	})
}

// LandingHeroWithBoundary wraps LandingHero with a dev-mode component boundary annotation.
// gallery:token badge,title,highlight,subtitle,primaryLabel,secondaryLabel
// gallery:hint badge:default(New),title:default(Build faster with),highlight:default(go-daisy),subtitle:default(Start your next project),primaryLabel:default(Get started free),secondaryLabel:default(View on GitHub)
func LandingHeroWithBoundary(badge, title, highlight, subtitle, primaryLabel, secondaryLabel string) templ.Component {
	props := LandingHeroProps{
		Badge: badge, Title: title, Highlight: highlight, Subtitle: subtitle,
		PrimaryLabel: primaryLabel, SecondaryLabel: secondaryLabel,
	}
	return devmode.ComponentBoundary("LandingHero", LandingHero(props), map[string]any{
		"badge": badge, "title": title, "highlight": highlight, "subtitle": subtitle,
		"primaryLabel": primaryLabel, "secondaryLabel": secondaryLabel,
	})
}

// LandingFeatureGridWithBoundary wraps LandingFeatureGrid with a dev-mode component boundary annotation.
// gallery:token title,subtitle,columns
// gallery:hint title:default(Everything you need),subtitle:default(Powerful features),columns:default(3)
func LandingFeatureGridWithBoundary(id, title, subtitle string, columns int, features []LandingFeature) templ.Component {
	props := LandingFeatureGridProps{ID: id, Title: title, Subtitle: subtitle, Columns: columns, Features: features}
	return devmode.ComponentBoundary("LandingFeatureGrid", LandingFeatureGrid(props), map[string]any{
		"title": title, "subtitle": subtitle, "columns": columns, "featureCount": len(features),
	})
}

// LandingTestimonialsWithBoundary wraps LandingTestimonials with a dev-mode component boundary annotation.
// gallery:token title
// gallery:hint title:default(Loved by developers)
func LandingTestimonialsWithBoundary(id, title string, items []LandingTestimonial) templ.Component {
	props := LandingTestimonialsProps{ID: id, Title: title, Items: items}
	return devmode.ComponentBoundary("LandingTestimonials", LandingTestimonials(props), map[string]any{
		"title": title, "itemCount": len(items),
	})
}

// LandingPricingWithBoundary wraps LandingPricing with a dev-mode component boundary annotation.
// gallery:token title,subtitle
// gallery:hint title:default(Simple pricing),subtitle:default(Choose the plan that fits)
func LandingPricingWithBoundary(id, title, subtitle string, tiers []LandingPricingTier) templ.Component {
	props := LandingPricingProps{ID: id, Title: title, Subtitle: subtitle, Tiers: tiers}
	return devmode.ComponentBoundary("LandingPricing", LandingPricing(props), map[string]any{
		"title": title, "subtitle": subtitle, "tierCount": len(tiers),
	})
}

// LandingFAQWithBoundary wraps LandingFAQ with a dev-mode component boundary annotation.
// gallery:token title,subtitle
// gallery:hint title:default(Frequently asked questions),subtitle:default(Got questions? We have answers.)
func LandingFAQWithBoundary(id, title, subtitle string, items []LandingFAQItem) templ.Component {
	props := LandingFAQProps{ID: id, Title: title, Subtitle: subtitle, Items: items}
	return devmode.ComponentBoundary("LandingFAQ", LandingFAQ(props), map[string]any{
		"title": title, "subtitle": subtitle, "itemCount": len(items),
	})
}

// LandingFooterWithBoundary wraps LandingFooter with a dev-mode component boundary annotation.
// gallery:token brand,tagline,copyright
// gallery:hint brand:default(go-daisy),tagline:default(Type-safe Templ components)
func LandingFooterWithBoundary(brand, tagline, copyright string, columns []LandingFooterColumn) templ.Component {
	props := LandingFooterProps{Brand: brand, Tagline: tagline, Copyright: copyright, Columns: columns}
	return devmode.ComponentBoundary("LandingFooter", LandingFooter(props), map[string]any{
		"brand": brand, "tagline": tagline, "copyright": copyright, "columnCount": len(columns),
	})
}
