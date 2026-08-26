package galleryruntime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

// ComponentSlugs maps component data-component names to their gallery page slugs.
// Single source of truth — used by both devoverlay.go (iframe) and pages_detail.templ (tree JS).
var ComponentSlugs = map[string]string{
	// Basics
	"Button":          "button",
	"Badge":           "badge",
	"StatusBadge":     "status-badge-real",
	"Avatar":          "avatar-real",
	"Card":            "card-real",
	"Tag":             "tag",
	"Divider":         "divider",
	"Kbd":             "kbd",
	"IconSpanColored": "button",
	// Feedback
	"Toast":           "toast-real",
	"Alert":           "alert",
	"Empty":           "empty-state-real",
	"Loader":          "loader",
	"NoPermissions":   "no-permissions",
	"SectionHeader":   "section-header",
	"Skeleton":        "skeleton",
	// Data display
	"StatCard":               "stat-card-real",
	"StatCardMinimal":        "stat-card-minimal",
	"ProgressCard":           "progress-card",
	"Timeline":               "timeline",
	"ChatBubble":             "chat-bubble",
	"ChatInput":              "chat-input",
	"AIThinkingIndicator":    "chat-bubble",
	"ChatWindow":             "chat-bubble",
	"PromptBarModelSelector": "prompt-bar-model",
	"PromptBarAbility":       "prompt-bar-ability",
	"LogsTable":              "logs-table",
	// Table
	"TableWithProps": "table",
	"Table":          "table",
	"ListArea":       "list-basic",
	// Navigation
	"ActionMenu":        "action-menu-real",
	"FilterTabs":        "filter-tabs",
	"FilterCard":        "filter-bar",
	"Pagination":        "pagination-real",
	"TabMenu":           "tabs",
	"SimpleTabs":        "tabs",
	"Tabs":              "tabs",
	"PageHeader":        "page-header-real",
	"Menu":              "menu-real",
	"TopBar":            "top-bar-real",
	"Navbar":            "navbar-real",
	"Breadcrumbs":       "breadcrumbs",
	"Dock":              "dock-nav",
	"ProfileMenu":       "profile-menu",
	"PageTitleMinimal":  "page-title-minimal",
	"PageTitleEditor":   "page-title-editor",
	"FooterMinimal":     "footer-minimal",
	// Foundation / display
	"Progress":          "progress",
	"Steps":             "steps",
	"Accordion":         "collapse",
	"Swap":              "swap",
	"Countdown":         "countdown",
	"StatusDot":         "status-dots",
	"Tooltip":           "tooltip",
	"Indicator":         "indicator",
	"Stack":             "stack",
	"Diff":              "diff",
	"Mask":              "mask",
	"Carousel":          "carousel",
	"Link":              "link-styles",
	"ThemeToggle":       "theme-toggle",
	"TypographyType":        "typography",
	"TypographyLayoutExample": "typography",
	// Layout
	"Hero":     "hero",
	"Join":     "join",
	"Fieldset": "fieldset",
	// Mockups
	"MockupBrowser": "mockup-browser",
	"MockupCode":    "mockup-code",
	"MockupPhone":   "mockup-phone",
	"MockupWindow":  "mockup-window",
	// Overlays
	"Modal":             "modal-real",
	"ConfirmPopup":      "confirm-popup",
	"FormModal":         "form-modal-real",
	"Dropdown":          "dropdown",
	"FAB":               "fab",
	"NotificationPanel": "notification-panel",
	"NotificationRow":   "notification-panel",
	// Person
	"PersonCell": "person-cell",
	// Forms
	"TextInput":      "text-input",
	"TextareaInput":  "textarea-input",
	"CheckboxInput":  "checkbox-input",
	"SelectInput":    "select-input",
	"RangeInput":     "range-input",
	"SearchInput":    "search-input-real",
	"FormField":      "form-field-real",
	"RadioGroup":     "form-radio",
	"Rating":         "form-rating",
	"FileInput":      "form-file",
	"Checkbox":       "form-checkbox",
	"Toggle":         "form-checkbox",
	"PromptBar":      "prompt-bar-minimal",
	"PromptBarAction": "prompt-bar-action",
	"InputSpinner":   "input-spinner",
	"WizardStepper":  "wizard-stepper",
}

// ComponentSlugsJSON returns ComponentSlugs as a JSON string for embedding in <script> tags.
func ComponentSlugsJSON() string {
	b, err := json.Marshal(ComponentSlugs)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// ComponentSlugsScript returns a templ.Component that renders a <script> tag
// defining COMPONENT_SLUGS as a global JS variable.
func ComponentSlugsScript() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, "<script>var COMPONENT_SLUGS=%s;</script>", ComponentSlugsJSON())
		return err
	})
}
