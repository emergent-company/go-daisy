// Package galleryruntime provides a reusable component gallery server that any
// Go/Templ project can embed. Projects supply a []GalleryComponent registry;
// galleryruntime handles the HTTP server, sandboxed iframe previews, design-token
// panel, and feedback loop.
package galleryruntime

import "github.com/a-h/templ"

// TokenType describes the kind of control rendered for a design token.
type TokenType string

const (
	TokenTypeRange  TokenType = "range"  // numeric slider (produces px / rem / % values)
	TokenTypeColor  TokenType = "color"  // colour picker
	TokenTypeSelect TokenType = "select" // named option list
)

// DesignToken describes a single manipulable design parameter for a component.
// Each token maps to a CSS custom property (CSSVar) that is injected as an
// override into the component's sandboxed iframe preview.
type DesignToken struct {
	// CSSVar is the CSS custom property name, e.g. "--btn-p".
	CSSVar string
	// Selector is the CSS selector to inject the override on.
	// Defaults to ":root" if empty. Use ".btn", ".badge" etc. for DaisyUI
	// component-level vars that are set on the element, not :root.
	Selector string
	// Label is the human-readable name shown in the panel, e.g. "Padding scale".
	Label string
	// Group is the section header in the panel (e.g. "Spacing", "Color", "Typography").
	Group string
	// Type determines the input control rendered.
	Type TokenType
	// Default is the initial value as a CSS string, e.g. "0.75rem", "#6d62d4", "md".
	Default string
	// Unit is appended to numeric range values, e.g. "rem", "px", "%".
	// Ignored for color and select types.
	Unit string
	// Min / Max / Step apply to range tokens.
	Min, Max, Step float64
	// Options is the ordered list of (value, label) pairs for select tokens.
	Options []TokenOption
}

// TokenOption is a single choice in a select-type token.
type TokenOption struct {
	Value string // CSS value emitted as the custom property value
	Label string // Human-readable label shown in the dropdown
}

// Category groups components into logical sections in the sidebar.
type Category string

const (
	CategoryFoundation  Category = "Foundation"
	CategoryBasics      Category = "Basics"
	CategoryDataDisplay Category = "Data Display"
	CategoryFeedback    Category = "Feedback"
	CategoryOverlays    Category = "Overlays"
	CategoryNavigation  Category = "Navigation"
	CategoryForms       Category = "Forms"
	CategoryLayout      Category = "Layout"
)

// GalleryStory is a single named variant/story of a component preview.
// Stories let you show the same component in different states side-by-side
// (e.g. "Default", "Loading", "Error", "Empty").
type GalleryStory struct {
	// Name is the display label shown in the variant tab strip, e.g. "Default".
	Name string
	// Description is optional extra context shown below the tab name.
	Description string
	// Templ is a pre-instantiated templ.Component for this story.
	// Mutually exclusive with HTML.
	Templ templ.Component
	// HTML is a raw HTML snippet for this story.
	// Mutually exclusive with Templ.
	HTML string
	// FrameHeight overrides the default iframe height for this story only.
	FrameHeight string
}

// GalleryComponent describes a single component entry in the gallery.
type GalleryComponent struct {
	Slug        string
	Name        string
	Category    Category
	Subcategory string // optional grouping within a category (e.g. "Buttons", "Modal")
	Description string

	// Root-level preview fields (backward-compatible; used when Variants is empty).
	HTML        string          // raw HTML snippet — rendered via srcdoc iframe
	Templ       templ.Component // full-page templ component — rendered via /gallery/render/{slug}
	FrameHeight string          // optional iframe height (e.g. "100vh"); defaults to "400px"

	// Tokens are optional design token controls shown in the token panel.
	Tokens []DesignToken

	// Variants lists named stories for this component. When non-empty, the
	// detail page shows a tab strip to switch between them. The root HTML/Templ
	// fields serve as a fallback "Default" story when Variants is nil.
	Variants []GalleryStory
}

// EffectiveVariants returns the list of stories to display for this component.
// If Variants is set, it is returned as-is. Otherwise a single "Default" story
// is synthesised from the root HTML/Templ fields (backward compat).
func (c GalleryComponent) EffectiveVariants() []GalleryStory {
	if len(c.Variants) > 0 {
		return c.Variants
	}
	return []GalleryStory{{
		Name:        "Default",
		Templ:       c.Templ,
		HTML:        c.HTML,
		FrameHeight: c.FrameHeight,
	}}
}

// SubcategoryGroup groups components under a named subcategory within a category.
type SubcategoryGroup struct {
	Name       string // subcategory label; empty string means ungrouped (top-level in category)
	Components []GalleryComponent
}

// CategoryGroup groups components under a named category for sidebar rendering.
// Subcategories provides 3-level hierarchy: Category > Subcategory > Component.
type CategoryGroup struct {
	Name          Category
	Components    []GalleryComponent // flat list (kept for backward compat)
	Subcategories []SubcategoryGroup // ordered subcategory groups
}
