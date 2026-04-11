package gallery

//go:generate go run github.com/emergent-company/go-daisy/cmd/boundarytoken -pkg github.com/emergent-company/go-daisy/components/ui    -out tokens_ui_gen.go     -out-pkg gallery ../../../../../../components/ui/boundary.go
//go:generate go run github.com/emergent-company/go-daisy/cmd/boundarytoken -pkg github.com/emergent-company/go-daisy/components/form   -out tokens_form_gen.go   -out-pkg gallery ../../../../../../components/form/boundary.go
//go:generate go run github.com/emergent-company/go-daisy/cmd/boundarytoken -pkg github.com/emergent-company/go-daisy/components/nav    -out tokens_nav_gen.go    -out-pkg gallery ../../../../../../components/nav/boundary.go
//go:generate go run github.com/emergent-company/go-daisy/cmd/boundarytoken -pkg github.com/emergent-company/go-daisy/components/modal  -out tokens_modal_gen.go  -out-pkg gallery ../../../../../../components/modal/boundary.go
//go:generate go run github.com/emergent-company/go-daisy/cmd/boundarytoken -pkg github.com/emergent-company/go-daisy/components/layout -out tokens_layout_gen.go -out-pkg gallery ../../../../../../components/layout/boundary.go
//go:generate go run github.com/emergent-company/go-daisy/cmd/boundarytoken -pkg github.com/emergent-company/go-daisy/components/logs   -out tokens_logs_gen.go   -out-pkg gallery ../../../../../../components/logs/boundary.go

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/components/form"
	"github.com/emergent-company/go-daisy/components/layout"
	"github.com/emergent-company/go-daisy/components/logs"
	"github.com/emergent-company/go-daisy/components/modal"
	"github.com/emergent-company/go-daisy/components/nav"
	"github.com/emergent-company/go-daisy/components/table"
	"github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/galleryruntime"
)

// row renders multiple components side by side in a centred flex row.
func row(components ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="flex flex-wrap gap-4 p-6 justify-center items-center">`); err != nil {
			return err
		}
		for _, c := range components {
			if err := c.Render(ctx, w); err != nil {
				return err
			}
		}
		_, err := io.WriteString(w, `</div>`)
		return err
	})
}

// withText returns a component that renders inner with a text child injected.
func withText(text string, inner templ.Component) templ.Component {
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	})
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return inner.Render(templ.WithChildren(ctx, child), w)
	})
}

// rawHTML returns a templ.Component that writes a raw HTML string.
func rawHTML(html string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, html)
		return err
	})
}

// withChildren renders inner with the given children injected.
func withChildren(inner templ.Component, children templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return inner.Render(templ.WithChildren(ctx, children), w)
	})
}

// alertIconForType returns the canonical icon for each AlertType.
func alertIconForType(typ ui.AlertType) string {
	switch typ {
	case ui.AlertError:
		return "lucide--circle-x"
	case ui.AlertWarning:
		return "lucide--triangle-alert"
	case ui.AlertInfo:
		return "lucide--info"
	default: // AlertSuccess
		return "lucide--circle-check"
	}
}

// alertInteractiveTokens are the design tokens for the Alert Interactive story.
// Only reflects actual component props: type (enum) and message (free-text string).
var alertInteractiveTokens = []galleryruntime.DesignToken{
	{
		Label:      "Type",
		Group:      "Component",
		Type:       galleryruntime.TokenTypeSelect,
		Default:    "alert-success",
		QueryParam: "type",
		Options: []galleryruntime.TokenOption{
			{Value: "alert-success", Label: "Success"},
			{Value: "alert-error", Label: "Error"},
			{Value: "alert-warning", Label: "Warning"},
			{Value: "alert-info", Label: "Info"},
		},
	},
	{
		Label:      "Message",
		Group:      "Component",
		Type:       galleryruntime.TokenTypeText,
		Default:    "Your changes have been saved successfully.",
		QueryParam: "message",
	},
}

func alertRenderFunc(defaultMessage string) func(params url.Values) templ.Component {
	return func(params url.Values) templ.Component {
		typ := ui.AlertType(params.Get("type"))
		if typ == "" {
			typ = ui.AlertSuccess
		}
		message := params.Get("message")
		if message == "" {
			message = defaultMessage
		}
		return ui.AlertWithIconBoundary(typ, alertIconForType(typ), message)
	}
}

// Add new components here — they are immediately available in the gallery.
func AllComponents() []galleryruntime.GalleryComponent {
	return []galleryruntime.GalleryComponent{

		// ── Basics / Buttons ─────────────────────────────────────────────────────

		// ── Data Display / Timeline ───────────────────────────────────────────────
		{
			Slug:        "timeline",
			Name:        "Timeline",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Vertical timeline for activity or event history.",
			HTML: `<div class="p-4 max-w-sm mx-auto">
  <ul class="timeline timeline-vertical">
    <li>
      <div class="timeline-start text-xs text-base-content/40">Jan 2024</div>
      <div class="timeline-middle"><span class="iconify lucide--circle-check size-4 text-primary"></span></div>
      <div class="timeline-end timeline-box text-sm">Project started</div>
      <hr class="bg-primary"/>
    </li>
    <li>
      <hr class="bg-primary"/>
      <div class="timeline-start text-xs text-base-content/40">Mar 2024</div>
      <div class="timeline-middle"><span class="iconify lucide--circle-check size-4 text-primary"></span></div>
      <div class="timeline-end timeline-box text-sm">Beta release</div>
      <hr/>
    </li>
    <li>
      <hr/>
      <div class="timeline-start text-xs text-base-content/40">Jun 2024</div>
      <div class="timeline-middle"><span class="iconify lucide--circle size-4 text-base-content/30"></span></div>
      <div class="timeline-end timeline-box text-sm text-base-content/50">v1.0 launch</div>
    </li>
  </ul>
</div>`,
		},

		// ── Data Display / Chat ───────────────────────────────────────────────────
		{
			Slug:        "chat-bubble",
			Name:        "Chat Bubble",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Chat conversation bubbles (sent and received).",
			HTML: `<div class="flex flex-col gap-2 p-4 max-w-sm mx-auto">
  <div class="chat chat-start">
    <div class="chat-header text-xs text-base-content/50 mb-0.5">Alice</div>
    <div class="chat-bubble chat-bubble-primary">Hey! How are you doing?</div>
    <div class="chat-footer text-xs text-base-content/30 mt-0.5">10:32 AM</div>
  </div>
  <div class="chat chat-end">
    <div class="chat-header text-xs text-base-content/50 mb-0.5">You</div>
    <div class="chat-bubble">Good thanks! How about you?</div>
    <div class="chat-footer text-xs text-base-content/30 mt-0.5">10:33 AM</div>
  </div>
</div>`,
		},

		// ── Data Display / Mockups ────────────────────────────────────────────────
		{
			Slug:        "mockup-code",
			Name:        "Mockup Code",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Code block mockup with terminal-style prefix lines.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="mockup-code w-full max-w-lg">
    <pre data-prefix="$"><code>go get github.com/emergent-company/go-daisy</code></pre>
    <pre data-prefix="$"><code>task build:ui</code></pre>
    <pre data-prefix=">" class="text-success"><code>Done in 1.2s</code></pre>
  </div>
</div>`,
		},
		{
			Slug:        "mockup-browser",
			Name:        "Mockup Browser",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Browser window mockup for UI showcasing.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="mockup-browser border border-base-300 w-full max-w-md">
    <div class="mockup-browser-toolbar">
      <div class="input">https://go-daisy.dev</div>
    </div>
    <div class="flex justify-center px-4 py-8 border-t border-base-300 bg-base-200 text-sm text-base-content/50">
      Page content here
    </div>
  </div>
</div>`,
		},

		// ── Feedback / Alerts ─────────────────────────────────────────────────────
		{
			Slug:        "alert-variants",
			Name:        "Inline Alert With Icon",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Alerts",
			Description: "Contextual feedback alerts for success, error, warning, and info states.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single alert with live type and message controls.",
					RenderFunc:  alertRenderFunc("Your changes have been saved successfully."),
					Tokens:      alertInteractiveTokens,
				},
				{
					Name:        "Examples",
					Description: "All alert types shown together.",
					HTML: `<div class="flex flex-col gap-3 p-6">
  <div role="alert" class="alert alert-success">
    <span class="iconify lucide--circle-check size-5"></span>
    <span>Your changes have been saved successfully.</span>
  </div>
  <div role="alert" class="alert alert-error">
    <span class="iconify lucide--circle-x size-5"></span>
    <span>Something went wrong. Please try again.</span>
  </div>
  <div role="alert" class="alert alert-warning">
    <span class="iconify lucide--triangle-alert size-5"></span>
    <span>Your session will expire in 5 minutes.</span>
  </div>
  <div role="alert" class="alert alert-info">
    <span class="iconify lucide--info size-5"></span>
    <span>A new software update is available.</span>
  </div>
</div>`,
				},
			},
		},
		// ── Navigation ────────────────────────────────────────────────────────────
		{
			Slug:        "filter-tabs",
			Name:        "Filter Tabs",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Filters",
			Description: "Radio-based pill filter tabs for list filtering.",
			HTML: `<div class="p-6 flex flex-wrap gap-2 justify-center">
  <div class="join">
    <input class="join-item btn btn-sm btn-outline checked:btn-primary" type="radio" name="filter" aria-label="All" checked/>
    <input class="join-item btn btn-sm btn-outline checked:btn-primary" type="radio" name="filter" aria-label="Active"/>
    <input class="join-item btn btn-sm btn-outline checked:btn-primary" type="radio" name="filter" aria-label="Pending"/>
    <input class="join-item btn btn-sm btn-outline checked:btn-primary" type="radio" name="filter" aria-label="Closed"/>
  </div>
</div>`,
		},

		// ── Forms ─────────────────────────────────────────────────────────────────
		{
			Slug:        "form-checkbox",
			Name:        "Checkboxes & Toggles",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Toggles",
			Description: "Checkbox and toggle switch inputs.",
			HTML: `<div class="p-6 flex flex-col gap-4 max-w-sm mx-auto">
  <label class="flex items-center gap-3 cursor-pointer">
    <input type="checkbox" class="checkbox checkbox-primary" checked/>
    <span class="text-sm">Receive email notifications</span>
  </label>
  <label class="flex items-center gap-3 cursor-pointer">
    <input type="checkbox" class="checkbox checkbox-secondary"/>
    <span class="text-sm">Subscribe to newsletter</span>
  </label>
  <div class="divider my-0"></div>
  <label class="flex items-center justify-between gap-3 cursor-pointer">
    <span class="text-sm">Dark mode</span>
    <input type="checkbox" class="toggle toggle-primary" checked/>
  </label>
  <label class="flex items-center justify-between gap-3 cursor-pointer">
    <span class="text-sm">Auto-save</span>
    <input type="checkbox" class="toggle toggle-secondary"/>
  </label>
</div>`,
		},
		{
			Slug:        "form-radio",
			Name:        "Radio Buttons",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Toggles",
			Description: "Radio button group for single-selection.",
			HTML: `<div class="p-6 max-w-sm mx-auto">
  <p class="text-sm font-medium mb-3">Subscription plan</p>
  <div class="flex flex-col gap-2">
    <label class="flex items-center gap-3 cursor-pointer">
      <input type="radio" name="plan" class="radio radio-primary" checked/>
      <span class="text-sm">Free – $0/mo</span>
    </label>
    <label class="flex items-center gap-3 cursor-pointer">
      <input type="radio" name="plan" class="radio radio-primary"/>
      <span class="text-sm">Pro – $12/mo</span>
    </label>
    <label class="flex items-center gap-3 cursor-pointer">
      <input type="radio" name="plan" class="radio radio-primary"/>
      <span class="text-sm">Enterprise – Custom</span>
    </label>
  </div>
</div>`,
		},
		{
			Slug:        "form-rating",
			Name:        "Rating",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Star and heart rating inputs using DaisyUI rating + mask classes.",
			HTML: `<div class="p-6 flex flex-col gap-4 items-center">
  <div class="rating">
    <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400"/>
    <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400"/>
    <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400" checked/>
    <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400"/>
    <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400"/>
  </div>
  <div class="rating rating-sm gap-1">
    <input type="radio" name="rating-2" class="mask mask-heart bg-red-400"/>
    <input type="radio" name="rating-2" class="mask mask-heart bg-red-400"/>
    <input type="radio" name="rating-2" class="mask mask-heart bg-red-400" checked/>
    <input type="radio" name="rating-2" class="mask mask-heart bg-red-400"/>
    <input type="radio" name="rating-2" class="mask mask-heart bg-red-400"/>
  </div>
</div>`,
		},

		// ── Foundation / Display ──────────────────────────────────────────────────
		{
			Slug:        "divider",
			Name:        "Divider",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Horizontal and vertical dividers with optional label.",
			HTML: `<div class="p-6 flex flex-col gap-4 max-w-sm mx-auto">
  <div class="divider">OR</div>
  <div class="divider divider-primary">Primary</div>
  <div class="flex h-20 items-center gap-4">
    <span class="text-sm">Left</span>
    <div class="divider divider-horizontal"></div>
    <span class="text-sm">Right</span>
  </div>
</div>`,
		},
		{
			Slug:        "kbd",
			Name:        "Keyboard Keys",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Keyboard shortcut display using DaisyUI kbd.",
			HTML: `<div class="flex flex-wrap gap-4 p-6 items-center justify-center">
  <div class="flex items-center gap-1 text-sm">
    Press <kbd class="kbd kbd-sm">⌘</kbd><kbd class="kbd kbd-sm">K</kbd> to search
  </div>
  <div class="flex items-center gap-1">
    <kbd class="kbd kbd-sm">Ctrl</kbd><span class="text-sm">+</span><kbd class="kbd kbd-sm">S</kbd>
  </div>
  <kbd class="kbd kbd-lg">Enter</kbd>
  <kbd class="kbd kbd-xs">Esc</kbd>
</div>`,
		},
		{
			Slug:        "progress",
			Name:        "Progress",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "A DaisyUI linear progress bar with configurable color, value, and max.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "A DaisyUI progress bar with configurable color, value, and max.",
					RenderFunc: func(params url.Values) templ.Component {
						color := ui.ProgressColor(params.Get("color"))
						if color == "" {
							color = ui.ProgressPrimary
						}
						value := 70
						if v, err := parseInt(params.Get("value")); err == nil {
							value = v
						}
						max := 100
						if m, err := parseInt(params.Get("max")); err == nil && m > 0 {
							max = m
						}
						return ui.ProgressWithBoundary(color, value, max)
					},
					Tokens: ProgressTokens(),
				},
			},
		},
		{
			Slug:        "steps",
			Name:        "Steps",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Step progress indicator for multi-step flows.",
			HTML: `<div class="p-6 flex justify-center">
  <ul class="steps">
    <li class="step step-primary">Register</li>
    <li class="step step-primary">Choose plan</li>
    <li class="step">Payment</li>
    <li class="step">Confirm</li>
  </ul>
</div>`,
		},
		{
			Slug:        "collapse",
			Name:        "Collapse / Accordion",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Collapsible sections using DaisyUI collapse.",
			HTML: `<div class="p-4 flex flex-col gap-2 max-w-sm mx-auto">
  <div class="collapse collapse-arrow border border-base-200 bg-base-100">
    <input type="checkbox" checked/>
    <div class="collapse-title text-sm font-medium">What is go-daisy?</div>
    <div class="collapse-content text-sm text-base-content/60">
      go-daisy is a Go UI component library for HTMX-driven web interfaces built with DaisyUI.
    </div>
  </div>
  <div class="collapse collapse-arrow border border-base-200 bg-base-100">
    <input type="checkbox"/>
    <div class="collapse-title text-sm font-medium">How do I install it?</div>
    <div class="collapse-content text-sm text-base-content/60">
      <code>go get github.com/emergent-company/go-daisy</code>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "dropdown",
			Name:        "Dropdown",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Dropdown menu triggered by a button.",
			HTML: `<div class="p-6 flex justify-center min-h-40">
  <div class="dropdown">
    <div tabindex="0" role="button" class="btn btn-primary btn-sm m-1">
      Options <span class="iconify lucide--chevron-down size-3.5"></span>
    </div>
    <ul tabindex="0" class="dropdown-content z-[1] menu menu-sm p-1 bg-base-100 border border-base-200 rounded-box shadow-lg w-40 mt-1">
      <li><a>Profile</a></li>
      <li><a>Settings</a></li>
      <li><a>Help</a></li>
      <li class="divider my-0.5"></li>
      <li><a class="text-error">Sign out</a></li>
    </ul>
  </div>
</div>`,
		},
		{
			Slug:        "tooltip",
			Name:        "Tooltip",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Tooltip on hover in top, bottom, left, right positions.",
			HTML: `<div class="flex flex-wrap gap-6 p-8 justify-center items-center min-h-32">
  <div class="tooltip" data-tip="Default tooltip">
    <button class="btn btn-sm">Hover me</button>
  </div>
  <div class="tooltip tooltip-top" data-tip="Top">
    <button class="btn btn-sm btn-primary">Top</button>
  </div>
  <div class="tooltip tooltip-bottom" data-tip="Bottom">
    <button class="btn btn-sm btn-secondary">Bottom</button>
  </div>
  <div class="tooltip tooltip-left" data-tip="Left">
    <button class="btn btn-sm btn-accent">Left</button>
  </div>
  <div class="tooltip tooltip-right" data-tip="Right">
    <button class="btn btn-sm btn-neutral">Right</button>
  </div>
</div>`,
		},
		{
			Slug:        "swap",
			Name:        "Swap",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Toggle between two visual states on click.",
			HTML: `<div class="flex flex-wrap gap-6 p-6 justify-center items-center">
  <label class="swap swap-rotate">
    <input type="checkbox"/>
    <span class="swap-on iconify lucide--sun size-8"></span>
    <span class="swap-off iconify lucide--moon size-8"></span>
  </label>
  <label class="swap">
    <input type="checkbox"/>
    <span class="swap-on btn btn-sm btn-success">ON</span>
    <span class="swap-off btn btn-sm btn-ghost">OFF</span>
  </label>
</div>`,
		},
		{
			Slug:        "hero",
			Name:        "Hero",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Layout",
			Description: "Full-width hero section with headline and CTA button.",
			HTML: `<div class="hero min-h-56 bg-base-200">
  <div class="hero-content text-center">
    <div class="max-w-md">
      <h1 class="text-4xl font-bold">go-daisy</h1>
      <p class="py-4 text-base-content/60">Type-safe Templ components styled with DaisyUI for HTMX apps.</p>
      <button class="btn btn-primary">Get Started</button>
    </div>
  </div>
</div>`,
		},

		// ── Data Display / List ───────────────────────────────────────────────────
		{
			Slug:        "list-basic",
			Name:        "List",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Lists",
			Description: "DaisyUI list component for vertical item groups.",
			HTML: `<div class="p-4 max-w-sm mx-auto">
  <ul class="list bg-base-100 rounded-box border border-base-200">
    <li class="list-row">
      <div class="list-col-grow">
        <div class="font-medium text-sm">Alice Johnson</div>
        <div class="text-xs text-base-content/50">alice@example.com</div>
      </div>
      <span class="badge badge-success badge-sm">Active</span>
    </li>
    <li class="list-row">
      <div class="list-col-grow">
        <div class="font-medium text-sm">Bob Smith</div>
        <div class="text-xs text-base-content/50">bob@example.com</div>
      </div>
      <span class="badge badge-ghost badge-sm">Inactive</span>
    </li>
    <li class="list-row">
      <div class="list-col-grow">
        <div class="font-medium text-sm">Carol White</div>
        <div class="text-xs text-base-content/50">carol@example.com</div>
      </div>
      <span class="badge badge-warning badge-sm">Pending</span>
    </li>
  </ul>
</div>`,
		},

		// ── Data Display / Indicator ──────────────────────────────────────────────
		{
			Slug:        "indicator",
			Name:        "Indicator",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Badge/dot overlay indicators on components.",
			HTML: `<div class="flex flex-wrap gap-8 p-6 justify-center items-center">
  <div class="indicator">
    <span class="indicator-item badge badge-error badge-sm">3</span>
    <button class="btn btn-ghost btn-sm btn-square">
      <span class="iconify lucide--bell size-5"></span>
    </button>
  </div>
  <div class="indicator">
    <span class="indicator-item badge badge-primary badge-xs"></span>
    <div class="avatar">
      <div class="w-10 rounded-full bg-secondary flex items-center justify-center text-secondary-content font-bold text-sm">AJ</div>
    </div>
  </div>
  <div class="indicator">
    <span class="indicator-item badge badge-success badge-sm">New</span>
    <div class="card bg-base-100 border border-base-200 w-32 h-16 flex items-center justify-center text-sm text-base-content/50">Card</div>
  </div>
</div>`,
		},

		// ── Data Display / Stack ──────────────────────────────────────────────────
		{
			Slug:        "stack",
			Name:        "Stack",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Overlapping stacked card effect.",
			HTML: `<div class="p-8 flex justify-center">
  <div class="stack w-48">
    <div class="card bg-primary text-primary-content shadow-lg h-24 flex items-center justify-center text-sm font-medium">Card 1</div>
    <div class="card bg-secondary text-secondary-content shadow h-24 flex items-center justify-center text-sm font-medium">Card 2</div>
    <div class="card bg-accent text-accent-content h-24 flex items-center justify-center text-sm font-medium">Card 3</div>
  </div>
</div>`,
		},

		// ── Data Display / Diff ───────────────────────────────────────────────────
		{
			Slug:        "diff",
			Name:        "Diff",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Side-by-side comparison diff panel.",
			HTML: `<div class="p-4 flex justify-center">
  <div class="diff aspect-[16/6] w-full max-w-lg rounded-box overflow-hidden border border-base-200">
    <div class="diff-item-1 bg-base-100 flex items-center justify-center p-4 text-sm text-base-content/60">
      Before: Old content here
    </div>
    <div class="diff-item-2 bg-base-200 flex items-center justify-center p-4 text-sm font-medium">
      After: New content here
    </div>
    <div class="diff-resizer"></div>
  </div>
</div>`,
		},

		// ── Data Display / Mask ───────────────────────────────────────────────────
		{
			Slug:        "mask",
			Name:        "Mask",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "CSS mask shapes applied to images and elements.",
			HTML: `<div class="flex flex-wrap gap-4 p-6 justify-center items-center">
  <div class="mask mask-squircle bg-primary w-16 h-16 flex items-center justify-center text-primary-content font-bold text-lg">S</div>
  <div class="mask mask-heart bg-error w-16 h-16 flex items-center justify-center text-error-content font-bold text-lg">♥</div>
  <div class="mask mask-hexagon bg-secondary w-16 h-16 flex items-center justify-center text-secondary-content font-bold text-lg">H</div>
  <div class="mask mask-triangle bg-accent w-16 h-16 flex items-center justify-center text-accent-content font-bold text-lg">▲</div>
  <div class="mask mask-circle bg-success w-16 h-16 flex items-center justify-center text-success-content font-bold text-lg">●</div>
  <div class="mask mask-star-2 bg-warning w-16 h-16 flex items-center justify-center text-warning-content font-bold text-lg">★</div>
</div>`,
		},

		// ── Data Display / Carousel ───────────────────────────────────────────────
		{
			Slug:        "carousel",
			Name:        "Carousel",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Horizontal scrolling carousel with snap items.",
			HTML: `<div class="p-4 flex justify-center">
  <div class="carousel w-full max-w-sm rounded-box gap-2">
    <div id="slide1" class="carousel-item w-full">
      <div class="bg-primary h-32 w-full rounded-box flex items-center justify-center text-primary-content font-bold">Slide 1</div>
    </div>
    <div id="slide2" class="carousel-item w-full">
      <div class="bg-secondary h-32 w-full rounded-box flex items-center justify-center text-secondary-content font-bold">Slide 2</div>
    </div>
    <div id="slide3" class="carousel-item w-full">
      <div class="bg-accent h-32 w-full rounded-box flex items-center justify-center text-accent-content font-bold">Slide 3</div>
    </div>
  </div>
</div>`,
		},

		// ── Data Display / Countdown ──────────────────────────────────────────────
		{
			Slug:        "countdown",
			Name:        "Countdown",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Animated countdown timer display.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="flex gap-5 text-center">
    <div class="flex flex-col items-center">
      <span class="countdown font-mono text-4xl text-primary"><span style="--value:02"></span></span>
      <span class="text-xs text-base-content/50 mt-1">days</span>
    </div>
    <div class="flex flex-col items-center">
      <span class="countdown font-mono text-4xl"><span style="--value:10"></span></span>
      <span class="text-xs text-base-content/50 mt-1">hours</span>
    </div>
    <div class="flex flex-col items-center">
      <span class="countdown font-mono text-4xl"><span style="--value:24"></span></span>
      <span class="text-xs text-base-content/50 mt-1">min</span>
    </div>
    <div class="flex flex-col items-center">
      <span class="countdown font-mono text-4xl"><span style="--value:45"></span></span>
      <span class="text-xs text-base-content/50 mt-1">sec</span>
    </div>
  </div>
</div>`,
		},

		// ── Data Display / Mockup Phone & Window ──────────────────────────────────
		{
			Slug:        "mockup-phone",
			Name:        "Mockup Phone",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Phone frame mockup for mobile UI display.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="mockup-phone">
    <div class="mockup-phone-camera"></div>
    <div class="mockup-phone-display">
      <div class="artboard artboard-demo phone-1 bg-base-200 flex items-center justify-center text-sm text-base-content/50">
        App screen
      </div>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "mockup-window",
			Name:        "Mockup Window",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Desktop window frame mockup.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="mockup-window border border-base-300 bg-base-100 w-full max-w-md">
    <div class="flex justify-center px-4 py-8 border-t border-base-300 bg-base-200 text-sm text-base-content/50">
      Window content here
    </div>
  </div>
</div>`,
		},

		// ── Feedback / Status ─────────────────────────────────────────────────────
		{
			Slug:        "status-dots",
			Name:        "Status Dots",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "States",
			Description: "Small colored status indicator dots.",
			HTML: `<div class="flex flex-wrap gap-6 p-6 items-center justify-center">
  <div class="flex items-center gap-2 text-sm">
    <span class="status status-success"></span> Online
  </div>
  <div class="flex items-center gap-2 text-sm">
    <span class="status status-error"></span> Offline
  </div>
  <div class="flex items-center gap-2 text-sm">
    <span class="status status-warning"></span> Away
  </div>
  <div class="flex items-center gap-2 text-sm">
    <span class="status status-info"></span> Busy
  </div>
  <div class="flex items-center gap-2 text-sm">
    <span class="status status-neutral"></span> Unknown
  </div>
</div>`,
		},

		// ── Overlays / Dropdown positions ─────────────────────────────────────────
		{
			Slug:        "dropdown-positions",
			Name:        "Dropdown Positions",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Dropdowns",
			Description: "Dropdown menus opening in different directions.",
			HTML: `<div class="flex flex-wrap gap-4 p-8 justify-center min-h-52">
  <div class="dropdown">
    <div tabindex="0" role="button" class="btn btn-sm btn-outline">Bottom ▼</div>
    <ul tabindex="0" class="dropdown-content menu menu-sm p-1 bg-base-100 border border-base-200 rounded-box shadow w-36 mt-1 z-[1]">
      <li><a>Item 1</a></li><li><a>Item 2</a></li><li><a>Item 3</a></li>
    </ul>
  </div>
  <div class="dropdown dropdown-top">
    <div tabindex="0" role="button" class="btn btn-sm btn-outline">Top ▲</div>
    <ul tabindex="0" class="dropdown-content menu menu-sm p-1 bg-base-100 border border-base-200 rounded-box shadow w-36 mb-1 z-[1]">
      <li><a>Item 1</a></li><li><a>Item 2</a></li><li><a>Item 3</a></li>
    </ul>
  </div>
  <div class="dropdown dropdown-end">
    <div tabindex="0" role="button" class="btn btn-sm btn-primary">Options ⋮</div>
    <ul tabindex="0" class="dropdown-content menu menu-sm p-1 bg-base-100 border border-base-200 rounded-box shadow w-36 mt-1 z-[1]">
      <li><a>Edit</a></li><li><a>Duplicate</a></li><li class="divider my-0.5"></li><li><a class="text-error">Delete</a></li>
    </ul>
  </div>
</div>`,
		},

		// ── Navigation / Breadcrumbs, Navbar, Menu, Dock ──────────────────────────
		{
			Slug:        "breadcrumbs",
			Name:        "Breadcrumbs",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Headers",
			Description: "Navigation breadcrumb trail with icon support.",
			HTML: `<div class="p-6 flex flex-col gap-4">
  <div class="breadcrumbs text-sm">
    <ul>
      <li><a href="#"><span class="iconify lucide--home size-4"></span> Home</a></li>
      <li><a href="#">Documents</a></li>
      <li>Add Document</li>
    </ul>
  </div>
  <div class="breadcrumbs text-xs text-base-content/50">
    <ul>
      <li><a href="#">Dashboard</a></li>
      <li><a href="#">Settings</a></li>
      <li>Profile</li>
    </ul>
  </div>
</div>`,
		},
		{
			Slug:        "dock-nav",
			Name:        "Dock",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Bottom dock navigation bar for mobile-style UIs.",
			HTML: `<div class="relative min-h-32 bg-base-100 border border-base-200 rounded-box overflow-hidden">
  <div class="dock">
    <button class="dock-active">
      <span class="iconify lucide--home size-5"></span>
      <span class="dock-label">Home</span>
    </button>
    <button>
      <span class="iconify lucide--search size-5"></span>
      <span class="dock-label">Search</span>
    </button>
    <button>
      <span class="iconify lucide--bell size-5"></span>
      <span class="dock-label">Alerts</span>
    </button>
    <button>
      <span class="iconify lucide--user size-5"></span>
      <span class="dock-label">Profile</span>
    </button>
  </div>
</div>`,
		},

		// ── Forms / File Input ────────────────────────────────────────────────────
		{
			Slug:        "form-file",
			Name:        "File Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "File upload input field with label and accept filter.",
			HTML: `<div class="p-6 max-w-sm mx-auto flex flex-col gap-4">
  <fieldset class="fieldset">
    <legend class="fieldset-legend">Upload file</legend>
    <input type="file" class="file-input file-input-bordered w-full"/>
  </fieldset>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">Profile image</legend>
    <input type="file" accept="image/*" class="file-input file-input-bordered file-input-primary w-full"/>
    <span class="fieldset-label text-base-content/50">PNG, JPG up to 2MB</span>
  </fieldset>
</div>`,
		},

		// ── Foundation / Join, Link ───────────────────────────────────────────────
		{
			Slug:        "join",
			Name:        "Join",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Layout",
			Description: "Join fuses children into a single rounded group.",
			HTML: `<div class="flex flex-col gap-4 p-6 items-center">
  <div class="join">
    <input class="join-item input input-bordered input-sm" placeholder="Search…"/>
    <button class="join-item btn btn-sm btn-primary">Go</button>
  </div>
  <div class="join">
    <button class="join-item btn btn-sm btn-outline">A</button>
    <button class="join-item btn btn-sm btn-outline btn-active">B</button>
    <button class="join-item btn btn-sm btn-outline">C</button>
  </div>
  <div class="join join-vertical w-48">
    <button class="join-item btn btn-sm btn-outline">Top</button>
    <button class="join-item btn btn-sm btn-outline">Middle</button>
    <button class="join-item btn btn-sm btn-outline">Bottom</button>
  </div>
</div>`,
		},
		{
			Slug:        "link-styles",
			Name:        "Links",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "DaisyUI link styles with color variants and hover.",
			HTML: `<div class="flex flex-wrap gap-4 p-6 items-center justify-center text-sm">
  <a href="#" class="link">Default link</a>
  <a href="#" class="link link-primary">Primary</a>
  <a href="#" class="link link-secondary">Secondary</a>
  <a href="#" class="link link-accent">Accent</a>
  <a href="#" class="link link-neutral">Neutral</a>
  <a href="#" class="link link-hover">Hover only</a>
</div>`,
		},

		{
			Slug:        "tag",
			Name:        "Tag",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Tag",
			Description: "Tag renders a removable chip badge used in multi-select fields. Clicking the × removes the tag.",
			HTML: `<div class="p-6 space-y-4">
  <div class="flex flex-wrap gap-2">
    <span class="badge badge-sm gap-1">
      Contract Law
      <a href="#" class="opacity-60 hover:opacity-100" aria-label="Remove Contract Law">✕</a>
    </span>
    <span class="badge badge-sm gap-1">
      Family Law
      <a href="#" class="opacity-60 hover:opacity-100" aria-label="Remove Family Law">✕</a>
    </span>
    <span class="badge badge-sm gap-1">
      Civil Litigation
      <a href="#" class="opacity-60 hover:opacity-100" aria-label="Remove Civil Litigation">✕</a>
    </span>
    <span class="badge badge-sm gap-1">
      Immigration
      <a href="#" class="opacity-60 hover:opacity-100" aria-label="Remove Immigration">✕</a>
    </span>
  </div>
  <p class="text-xs text-base-content/50">Tags without a remove link (read-only):</p>
  <div class="flex flex-wrap gap-2">
    <span class="badge badge-sm">Contract Law</span>
    <span class="badge badge-sm">Family Law</span>
    <span class="badge badge-sm">Civil Litigation</span>
  </div>
</div>`,
		},
		{
			Slug:        "company-avatar",
			Name:        "Company Avatar",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Avatars",
			Description: "Circular avatar with a building icon placeholder for companies. Same sizes as Avatar. Use alongside a company name in tables and detail views.",
			HTML: `<div class="p-6 flex flex-wrap gap-6 items-end">
  <div class="flex flex-col items-center gap-2">
    <div class="avatar avatar-xs">
      <div class="avatar-placeholder rounded-full bg-base-300 flex items-center justify-center text-base-content/60 w-6 h-6">
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 22V4a2 2 0 012-2h8a2 2 0 012 2v18"/><path d="M6 12H4a2 2 0 00-2 2v8h20v-8a2 2 0 00-2-2h-2"/></svg>
      </div>
    </div>
    <span class="text-xs text-base-content/60">xs</span>
  </div>
  <div class="flex flex-col items-center gap-2">
    <div class="avatar avatar-sm">
      <div class="avatar-placeholder rounded-full bg-base-300 flex items-center justify-center text-base-content/60 w-8 h-8">
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 22V4a2 2 0 012-2h8a2 2 0 012 2v18"/><path d="M6 12H4a2 2 0 00-2 2v8h20v-8a2 2 0 00-2-2h-2"/></svg>
      </div>
    </div>
    <span class="text-xs text-base-content/60">sm</span>
  </div>
  <div class="flex flex-col items-center gap-2">
    <div class="avatar avatar-md">
      <div class="avatar-placeholder rounded-full bg-base-300 flex items-center justify-center text-base-content/60 w-10 h-10">
        <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 22V4a2 2 0 012-2h8a2 2 0 012 2v18"/><path d="M6 12H4a2 2 0 00-2 2v8h20v-8a2 2 0 00-2-2h-2"/></svg>
      </div>
    </div>
    <span class="text-xs text-base-content/60">md</span>
  </div>
  <div class="flex flex-col items-center gap-2">
    <div class="avatar avatar-lg">
      <div class="avatar-placeholder rounded-full bg-base-300 flex items-center justify-center text-base-content/60 w-14 h-14">
        <svg class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 22V4a2 2 0 012-2h8a2 2 0 012 2v18"/><path d="M6 12H4a2 2 0 00-2 2v8h20v-8a2 2 0 00-2-2h-2"/></svg>
      </div>
    </div>
    <span class="text-xs text-base-content/60">lg</span>
  </div>
  <div class="flex flex-col items-center gap-2">
    <div class="flex items-center gap-2">
      <div class="avatar avatar-xs">
        <div class="avatar-placeholder rounded-full bg-base-300 flex items-center justify-center text-base-content/60 w-6 h-6">
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 22V4a2 2 0 012-2h8a2 2 0 012 2v18"/><path d="M6 12H4a2 2 0 00-2 2v8h20v-8a2 2 0 00-2-2h-2"/></svg>
        </div>
      </div>
      <span class="text-sm font-medium">Acme Corp</span>
    </div>
    <span class="text-xs text-base-content/60">with name</span>
  </div>
</div>`,
		},
		{
			Slug:        "person-avatar",
			Name:        "Person Avatar",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Avatars",
			Description: "Inline avatar + name chip with a hover card that reveals contact details. Pure CSS — no JS required.",
			HTML: `<style>
  .person-chip .person-card {
    opacity: 0;
    pointer-events: none;
    transform: translateY(4px) scale(0.98);
    transition: opacity 120ms ease, transform 120ms ease;
  }
  /* Invisible bridge above the card keeps hover active across the gap */
  .person-chip .person-card::before {
    content: '';
    display: block;
    position: absolute;
    top: -8px;
    left: 0;
    right: 0;
    height: 8px;
  }
  .person-chip:hover .person-card,
  .person-chip:focus-within .person-card {
    opacity: 1;
    pointer-events: auto;
    transform: translateY(0) scale(1);
  }
</style>
<div class="p-8 space-y-10">
  <div>
    <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-4">Inline — initials avatar</p>
    <div class="flex flex-wrap gap-6 items-start">
      <div class="relative person-chip inline-flex items-center gap-1.5 rounded-full px-2 py-1 cursor-default select-none hover:bg-base-200 transition-colors">
        <div class="size-6 rounded-full bg-primary flex items-center justify-center text-primary-content text-[10px] font-semibold shrink-0">JD</div>
        <span class="text-sm font-medium text-base-content leading-none">Jane Doe</span>
        <div class="person-card absolute left-[-8px] top-full z-50 w-64 rounded-2xl border border-base-200 bg-base-100 shadow-xl overflow-hidden">
          <div class="h-10 bg-gradient-to-r from-primary/20 to-primary/5"></div>
          <div class="px-4 pb-3 -mt-5 flex items-end gap-3">
            <div class="size-10 rounded-full bg-primary flex items-center justify-center text-primary-content text-sm font-bold ring-2 ring-base-100 shrink-0">JD</div>
            <div class="pb-0.5 min-w-0">
              <p class="text-sm font-semibold text-base-content leading-tight truncate">Jane Doe</p>
              <p class="text-xs text-base-content/50 leading-tight">Senior Attorney</p>
            </div>
          </div>
          <div class="px-4 pb-4 space-y-2">
            <div class="flex items-center gap-2 text-xs text-base-content/60">
              <svg class="size-3.5 shrink-0 text-base-content/30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
              <span class="truncate">jane.doe@example.com</span>
            </div>
            <div class="flex items-center justify-between pt-1">
              <span class="badge badge-success badge-sm badge-soft">Active</span>
              <a class="btn btn-xs btn-ghost text-primary gap-1 px-2">View profile <svg class="size-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg></a>
            </div>
          </div>
        </div>
      </div>
      <div class="relative person-chip inline-flex items-center gap-2 rounded-full px-2.5 py-1.5 cursor-default select-none hover:bg-base-200 transition-colors">
        <div class="size-8 rounded-full bg-secondary flex items-center justify-center text-secondary-content text-xs font-semibold shrink-0">BS</div>
        <span class="text-sm font-medium text-base-content leading-none">Bob Smith</span>
        <div class="person-card absolute left-[-10px] top-full z-50 w-64 rounded-2xl border border-base-200 bg-base-100 shadow-xl overflow-hidden">
          <div class="h-10 bg-gradient-to-r from-secondary/20 to-secondary/5"></div>
          <div class="px-4 pb-3 -mt-5 flex items-end gap-3">
            <div class="size-10 rounded-full bg-secondary flex items-center justify-center text-secondary-content text-sm font-bold ring-2 ring-base-100 shrink-0">BS</div>
            <div class="pb-0.5 min-w-0">
              <p class="text-sm font-semibold text-base-content leading-tight truncate">Bob Smith</p>
              <p class="text-xs text-base-content/50 leading-tight">Paralegal</p>
            </div>
          </div>
          <div class="px-4 pb-4 space-y-2">
            <div class="flex items-center gap-2 text-xs text-base-content/60">
              <svg class="size-3.5 shrink-0 text-base-content/30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
              <span class="truncate">bob.smith@example.com</span>
            </div>
            <div class="flex items-center justify-between pt-1">
              <span class="badge badge-warning badge-sm badge-soft">On leave</span>
              <a class="btn btn-xs btn-ghost text-primary gap-1 px-2">View profile <svg class="size-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg></a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`,
		},

		// ── Data Display extras ────────────────────────────────────────────────
		{
			Slug:        "table-with-actions",
			Name:        "With Actions",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Full-featured table with sortable headers, status badges, avatars, and an action menu (ellipsis dropdown) per row.",
			HTML: `<div class="p-6">
  <div class="overflow-x-auto rounded-md bg-base-100">
    <table class="table table-zebra w-full">
      <thead>
        <tr>
          <th class="text-xs font-semibold text-base-content/60 uppercase">
            <a class="flex items-center gap-1 hover:text-base-content cursor-pointer">
              Name
              <svg class="w-3 h-3" viewBox="0 0 16 16" fill="currentColor"><path d="M8 4l4 4H4l4-4zm0 8L4 8h8l-4 4z"/></svg>
            </a>
          </th>
          <th class="text-xs font-semibold text-base-content/60 uppercase">Status</th>
          <th class="text-xs font-semibold text-base-content/60 uppercase">Role</th>
          <th class="text-xs font-semibold text-base-content/60 uppercase">Joined</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr class="hover">
          <td>
            <div class="flex items-center gap-3">
              <div class="avatar avatar-sm">
                <div class="avatar-placeholder rounded-full bg-primary flex items-center justify-center font-semibold text-primary-content text-xs w-8 h-8">AJ</div>
              </div>
              <span class="font-medium">Alice Johnson</span>
            </div>
          </td>
          <td><span class="badge badge-sm badge-success">active</span></td>
          <td class="text-sm text-base-content/70">Admin</td>
          <td class="text-sm text-base-content/60">2024-01-15</td>
          <td class="text-right">
            <div class="dropdown dropdown-end">
              <button tabindex="0" type="button" class="btn btn-ghost btn-sm btn-square">
                <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="1.5"/><circle cx="12" cy="12" r="1.5"/><circle cx="12" cy="19" r="1.5"/></svg>
              </button>
              <ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box shadow-lg border border-base-200 z-50 w-44 p-1 mt-1">
                <li><a class="flex items-center gap-2 text-sm">View</a></li>
                <li><a class="flex items-center gap-2 text-sm">Edit</a></li>
                <li><a class="flex items-center gap-2 text-sm text-error">Delete</a></li>
              </ul>
            </div>
          </td>
        </tr>
        <tr class="hover">
          <td>
            <div class="flex items-center gap-3">
              <div class="avatar avatar-sm">
                <div class="avatar-placeholder rounded-full bg-secondary flex items-center justify-center font-semibold text-secondary-content text-xs w-8 h-8">BS</div>
              </div>
              <span class="font-medium">Bob Smith</span>
            </div>
          </td>
          <td><span class="badge badge-sm badge-warning">pending</span></td>
          <td class="text-sm text-base-content/70">Employee</td>
          <td class="text-sm text-base-content/60">2024-03-02</td>
          <td class="text-right">
            <div class="dropdown dropdown-end">
              <button tabindex="0" type="button" class="btn btn-ghost btn-sm btn-square">
                <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="1.5"/><circle cx="12" cy="12" r="1.5"/><circle cx="12" cy="19" r="1.5"/></svg>
              </button>
              <ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box shadow-lg border border-base-200 z-50 w-44 p-1 mt-1">
                <li><a class="flex items-center gap-2 text-sm">View</a></li>
                <li><a class="flex items-center gap-2 text-sm">Edit</a></li>
                <li><a class="flex items-center gap-2 text-sm text-error">Delete</a></li>
              </ul>
            </div>
          </td>
        </tr>
        <tr class="hover">
          <td>
            <div class="flex items-center gap-3">
              <div class="avatar avatar-sm">
                <div class="avatar-placeholder rounded-full bg-accent flex items-center justify-center font-semibold text-accent-content text-xs w-8 h-8">CW</div>
              </div>
              <span class="font-medium">Carol White</span>
            </div>
          </td>
          <td><span class="badge badge-sm badge-error">closed</span></td>
          <td class="text-sm text-base-content/70">Employee</td>
          <td class="text-sm text-base-content/60">2023-11-20</td>
          <td class="text-right">
            <div class="dropdown dropdown-end">
              <button tabindex="0" type="button" class="btn btn-ghost btn-sm btn-square">
                <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="1.5"/><circle cx="12" cy="12" r="1.5"/><circle cx="12" cy="19" r="1.5"/></svg>
              </button>
              <ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box shadow-lg border border-base-200 z-50 w-44 p-1 mt-1">
                <li><a class="flex items-center gap-2 text-sm">View</a></li>
                <li><a class="flex items-center gap-2 text-sm">Edit</a></li>
                <li><a class="flex items-center gap-2 text-sm text-error">Delete</a></li>
              </ul>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <div class="flex items-center justify-between mt-4">
    <p class="text-sm text-base-content/60">Showing 1 to 3 of 47 entries</p>
    <div class="join">
      <button class="join-item btn btn-sm btn-active">1</button>
      <button class="join-item btn btn-sm">2</button>
      <button class="join-item btn btn-sm">3</button>
      <button class="join-item btn btn-sm">»</button>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "table-empty",
			Name:        "Table — Empty State",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Full-width empty-state row inside a tbody when the list has no items.",
			HTML: `<div class="p-6">
  <div class="overflow-x-auto rounded-md bg-base-100">
    <table class="table table-zebra w-full">
      <thead>
        <tr>
          <th class="text-xs font-semibold text-base-content/60 uppercase">Name</th>
          <th class="text-xs font-semibold text-base-content/60 uppercase">Status</th>
          <th class="text-xs font-semibold text-base-content/60 uppercase">Role</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td colspan="3" class="text-center text-base-content/50 py-8">No records found.</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>`,
		},
		{
			Slug:        "progress-card",
			Name:        "Progress Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "Card with a gradient header, a progress bar, and an optional stats row.",
			HTML: `<div class="p-6 space-y-4 max-w-lg">
  <div class="card bg-base-100 shadow-sm border border-base-200 overflow-hidden">
    <div class="p-4 bg-gradient-to-r from-primary/10 to-primary/5">
      <h3 class="text-base font-semibold text-base-content">Case Compliance</h3>
      <p class="text-sm text-base-content/60 mt-0.5">Johnson v. Smith</p>
    </div>
    <div class="card-body p-4 pt-3 gap-3">
      <div class="flex items-center gap-3">
        <progress class="progress progress-primary flex-1" value="72" max="100"></progress>
        <span class="text-sm font-medium text-base-content/70 whitespace-nowrap">72%</span>
      </div>
      <div class="flex flex-wrap gap-4">
        <div class="flex flex-col">
          <span class="text-xs text-base-content/60">Tasks</span>
          <span class="text-sm font-semibold text-base-content">18 / 25</span>
        </div>
        <div class="flex flex-col">
          <span class="text-xs text-base-content/60">Documents</span>
          <span class="text-sm font-semibold text-base-content">12 / 15</span>
        </div>
        <div class="flex flex-col">
          <span class="text-xs text-base-content/60">Due</span>
          <span class="text-sm font-semibold text-base-content">Apr 30</span>
        </div>
      </div>
    </div>
  </div>
  <div class="card bg-base-100 shadow-sm border border-base-200 overflow-hidden">
    <div class="px-4 py-3 flex items-center justify-between gap-3">
      <div>
        <p class="text-sm font-semibold text-base-content">Document Review</p>
        <p class="text-xs text-base-content/60">3 of 8 complete</p>
      </div>
      <div class="flex items-center gap-2 min-w-0">
        <progress class="progress progress-primary w-24" value="38" max="100"></progress>
        <span class="text-xs font-medium text-base-content/70 whitespace-nowrap">38%</span>
      </div>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "stat-card-minimal",
			Name:        "Stat Card — Minimal",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "KPI stat cards with trend indicators (↑↓). Use on dashboards to surface key metrics.",
			HTML: `<div class="p-6">
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
    <div class="card bg-base-100 card-border">
      <div class="card-body">
        <p class="text-base-content/60 text-xs font-medium tracking-wide uppercase">Open Cases</p>
        <div class="mt-4 flex items-end justify-between gap-2 text-sm">
          <p class="text-2xl/none font-semibold">142</p>
          <div class="text-success flex items-center gap-0.5 px-1 font-medium text-xs">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>
            12.3%
          </div>
        </div>
      </div>
    </div>
    <div class="card bg-base-100 card-border">
      <div class="card-body">
        <p class="text-base-content/60 text-xs font-medium tracking-wide uppercase">Pending Tasks</p>
        <div class="mt-4 flex items-end justify-between gap-2 text-sm">
          <p class="text-2xl/none font-semibold">38</p>
          <div class="text-error flex items-center gap-0.5 px-1 font-medium text-xs">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>
            4.1%
          </div>
        </div>
      </div>
    </div>
    <div class="card bg-base-100 card-border">
      <div class="card-body">
        <p class="text-base-content/60 text-xs font-medium tracking-wide uppercase">Clients</p>
        <div class="mt-4 flex items-end justify-between gap-2 text-sm">
          <p class="text-2xl/none font-semibold">89</p>
          <div class="text-success flex items-center gap-0.5 px-1 font-medium text-xs">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>
            7.8%
          </div>
        </div>
      </div>
    </div>
    <div class="card bg-base-100 card-border">
      <div class="card-body">
        <p class="text-base-content/60 text-xs font-medium tracking-wide uppercase">Avg. Case Days</p>
        <div class="mt-4 flex items-end justify-between gap-2 text-sm">
          <p class="text-2xl/none font-semibold">24</p>
          <div class="text-success flex items-center gap-0.5 px-1 font-medium text-xs">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>
            2.5%
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "stat-card-icon-corner",
			Name:        "Stat Card — Icon Corner",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "Stat cards with a floating icon in the top corner and a soft badge for the trend.",
			HTML: `<div class="p-6">
  <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-4">
    <div class="card bg-base-100 relative shadow-sm border border-base-200">
      <div class="card-body gap-2">
        <p class="text-2xl/none font-semibold">142</p>
        <p class="text-base-content/60 mt-1 text-sm font-medium">Open Cases</p>
        <div class="mt-5 flex items-center gap-2">
          <span class="badge badge-soft badge-success badge-sm gap-0.5 px-1 font-medium">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>
            14.6%
          </span>
          <p class="text-base-content/60 text-sm">vs last month</p>
        </div>
      </div>
      <div class="absolute -end-3 -top-3 rounded-full bg-gradient-to-bl from-transparent to-base-200/60 p-1.5">
        <div class="bg-base-100 flex items-center justify-center rounded-full p-2 shadow">
          <svg class="size-5 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
        </div>
      </div>
    </div>
    <div class="card bg-base-100 relative shadow-sm border border-base-200">
      <div class="card-body gap-2">
        <p class="text-2xl/none font-semibold">38</p>
        <p class="text-base-content/60 mt-1 text-sm font-medium">Pending Tasks</p>
        <div class="mt-5 flex items-center gap-2">
          <span class="badge badge-soft badge-error badge-sm gap-0.5 px-1 font-medium">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>
            4.1%
          </span>
          <p class="text-base-content/60 text-sm">vs last month</p>
        </div>
      </div>
      <div class="absolute -end-3 -top-3 rounded-full bg-gradient-to-bl from-transparent to-base-200/60 p-1.5">
        <div class="bg-base-100 flex items-center justify-center rounded-full p-2 shadow">
          <svg class="size-5 text-warning" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>
        </div>
      </div>
    </div>
    <div class="card bg-base-100 relative shadow-sm border border-base-200">
      <div class="card-body gap-2">
        <p class="text-2xl/none font-semibold">89</p>
        <p class="text-base-content/60 mt-1 text-sm font-medium">Active Clients</p>
        <div class="mt-5 flex items-center gap-2">
          <span class="badge badge-soft badge-success badge-sm gap-0.5 px-1 font-medium">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>
            7.8%
          </span>
          <p class="text-base-content/60 text-sm">vs last month</p>
        </div>
      </div>
      <div class="absolute -end-3 -top-3 rounded-full bg-gradient-to-bl from-transparent to-base-200/60 p-1.5">
        <div class="bg-base-100 flex items-center justify-center rounded-full p-2 shadow">
          <svg class="size-5 text-success" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/></svg>
        </div>
      </div>
    </div>
    <div class="card bg-base-100 relative shadow-sm border border-base-200">
      <div class="card-body gap-2">
        <p class="text-2xl/none font-semibold">$48K</p>
        <p class="text-base-content/60 mt-1 text-sm font-medium">Revenue (MTD)</p>
        <div class="mt-5 flex items-center gap-2">
          <span class="badge badge-soft badge-success badge-sm gap-0.5 px-1 font-medium">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>
            9.2%
          </span>
          <p class="text-base-content/60 text-sm">vs last month</p>
        </div>
      </div>
      <div class="absolute -end-3 -top-3 rounded-full bg-gradient-to-bl from-transparent to-base-200/60 p-1.5">
        <div class="bg-base-100 flex items-center justify-center rounded-full p-2 shadow">
          <svg class="size-5 text-secondary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/></svg>
        </div>
      </div>
    </div>
  </div>
</div>`,
		},

		// ── Feedback extras ───────────────────────────────────────────────────
		{
			Slug:        "skeleton",
			Name:        "Skeleton",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Loading",
			Description: "A DaisyUI skeleton placeholder block. Use the classes token to control size and shape.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "A skeleton placeholder with configurable Tailwind size classes.",
					RenderFunc: func(params url.Values) templ.Component {
						classes := params.Get("classes")
						if classes == "" {
							classes = "h-4 w-full"
						}
						return ui.SkeletonWithBoundary(classes)
					},
					Tokens: SkeletonTokens(),
				},
			},
		},
		{
			Slug:        "skeleton-dashboard",
			Name:        "Skeleton — Dashboard Layout",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Loading",
			Description: "Dashed placeholder grid for planning dashboard layouts before content loads.",
			HTML: `<div class="p-6 space-y-4">
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
    <div class="border-base-300/80 bg-base-200/20 rounded-box flex h-28 flex-col items-center justify-center border border-dashed text-center">
      <p class="text-base-content/80 font-medium text-sm">Stats</p>
      <p class="text-base-content/50 text-xs">1 / 4</p>
    </div>
    <div class="border-base-300/80 bg-base-200/20 rounded-box flex h-28 flex-col items-center justify-center border border-dashed text-center">
      <p class="text-base-content/80 font-medium text-sm">Stats</p>
      <p class="text-base-content/50 text-xs">1 / 4</p>
    </div>
    <div class="border-base-300/80 bg-base-200/20 rounded-box flex h-28 flex-col items-center justify-center border border-dashed text-center">
      <p class="text-base-content/80 font-medium text-sm">Stats</p>
      <p class="text-base-content/50 text-xs">1 / 4</p>
    </div>
    <div class="border-base-300/80 bg-base-200/20 rounded-box flex h-28 flex-col items-center justify-center border border-dashed text-center">
      <p class="text-base-content/80 font-medium text-sm">Stats</p>
      <p class="text-base-content/50 text-xs">1 / 4</p>
    </div>
  </div>
  <div class="grid grid-cols-12 gap-4">
    <div class="border-base-300/80 bg-base-200/20 rounded-box col-span-12 flex h-64 flex-col items-center justify-center border border-dashed text-center lg:col-span-8">
      <p class="text-base-content/80 font-medium text-sm">Primary Chart</p>
      <p class="text-base-content/50 text-xs">8 / 12</p>
    </div>
    <div class="border-base-300/80 bg-base-200/20 rounded-box col-span-12 flex h-64 flex-col items-center justify-center border border-dashed text-center lg:col-span-4">
      <p class="text-base-content/80 font-medium text-sm">Side Panel</p>
      <p class="text-base-content/50 text-xs">4 / 12</p>
    </div>
  </div>
  <div class="border-base-300/80 bg-base-200/20 rounded-box flex h-48 flex-col items-center justify-center border border-dashed text-center">
    <p class="text-base-content/80 font-medium text-sm">Full-Width Table</p>
    <p class="text-base-content/50 text-xs">12 / 12</p>
  </div>
</div>`,
		},
		{
			Slug:        "section-header",
			Name:        "Section Header",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Indicators",
			Description: "Divider with a label — used to separate logical groups within a form or detail panel.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "A section divider label with configurable title text.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Personal Information"
						}
						return ui.SectionHeaderWithBoundary(title)
					},
					Tokens: SectionHeaderTokens(),
				},
			},
		},
		{
			Slug:        "no-permissions",
			Name:        "No Permissions",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "States",
			Description: "Permission-denied placeholder shown when the current user lacks access to a section.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "A fixed permission-denied placeholder with no configurable props.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.NoPermissionsWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "log-status-dot",
			Name:        "Log Status Dot",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "States",
			Description: "DaisyUI status dot for workflow log entries. Color indicates the entry state.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "A small colored status indicator dot with configurable color.",
					RenderFunc: func(params url.Values) templ.Component {
						color := ui.StatusColor(params.Get("color"))
						if color == "" {
							color = ui.StatusSuccess
						}
						return ui.StatusWithBoundary(color)
					},
					Tokens: StatusTokens(),
				},
			},
		},
		{
			Slug:        "notification-panel",
			Name:        "Notification Panel",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Notifications",
			Description: "Tab-based notification center with All / Unread tabs and a list of notification items.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="bg-base-100 border border-base-200 rounded-box shadow-lg w-80">
    <div class="flex items-center justify-between px-4 pt-4 pb-2">
      <p class="font-semibold text-sm">Notifications</p>
      <button class="btn btn-ghost btn-xs text-primary">Mark all read</button>
    </div>
    <div role="tablist" class="tabs tabs-border px-4">
      <a role="tab" class="tab tab-active tab-sm">All</a>
      <a role="tab" class="tab tab-sm">Unread <span class="badge badge-xs badge-primary ml-1">3</span></a>
    </div>
    <ul class="divide-y divide-base-200">
      <li class="flex gap-3 px-4 py-3 bg-primary/5">
        <div class="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
          <svg class="size-4 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-base-content">New case assigned</p>
          <p class="text-xs text-base-content/60 mt-0.5">Johnson v. Smith was assigned to you.</p>
          <p class="text-xs text-base-content/50 mt-1">2 min ago</p>
        </div>
        <span class="mt-1 size-2 shrink-0 rounded-full bg-primary"></span>
      </li>
      <li class="flex gap-3 px-4 py-3 bg-primary/5">
        <div class="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-full bg-warning/10">
          <svg class="size-4 text-warning" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-base-content">Task deadline tomorrow</p>
          <p class="text-xs text-base-content/60 mt-0.5">File motion for Johnson v. Smith due soon.</p>
          <p class="text-xs text-base-content/50 mt-1">1 hour ago</p>
        </div>
        <span class="mt-1 size-2 shrink-0 rounded-full bg-primary"></span>
      </li>
      <li class="flex gap-3 px-4 py-3">
        <div class="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-full bg-success/10">
          <svg class="size-4 text-success" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-base-content">Client signed in</p>
          <p class="text-xs text-base-content/60 mt-0.5">Alice Johnson accessed the client portal.</p>
          <p class="text-xs text-base-content/50 mt-1">Yesterday</p>
        </div>
      </li>
    </ul>
    <div class="px-4 py-3 border-t border-base-200">
      <a href="#" class="text-xs text-primary hover:underline">View all notifications</a>
    </div>
  </div>
</div>`,
		},

		// ── Overlays extras ───────────────────────────────────────────────────
		{
			Slug:        "fab",
			Name:        "FAB",
			Category:    galleryruntime.CategoryOverlays,
			Description: "CSS-only floating action button with an expandable sub-menu of quick actions. No JS required.",
			HTML: `<div class="p-6 relative h-56 bg-base-100 rounded-lg border border-base-200 overflow-hidden">
  <p class="text-xs text-base-content/50 p-2">FAB appears bottom-right. Click it to expand sub-actions.</p>
  <div class="absolute bottom-4 end-4 z-10">
    <div class="dropdown dropdown-top dropdown-end">
      <div tabindex="0" role="button" class="btn btn-primary btn-lg btn-circle shadow-lg">
        <svg class="size-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </div>
      <div tabindex="0" class="dropdown-content mb-3 space-y-2 flex flex-col items-end">
        <div class="flex items-center gap-2">
          <span class="bg-base-100 text-xs font-medium px-2 py-1 rounded shadow text-base-content/70">New Case</span>
          <div class="btn btn-sm btn-circle btn-outline bg-base-100 shadow">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="bg-base-100 text-xs font-medium px-2 py-1 rounded shadow text-base-content/70">Upload Doc</span>
          <div class="btn btn-sm btn-circle btn-outline bg-base-100 shadow">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="12" y2="12"/><line x1="15" y1="15" x2="12" y2="12"/></svg>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="bg-base-100 text-xs font-medium px-2 py-1 rounded shadow text-base-content/70">Add Task</span>
          <div class="btn btn-sm btn-circle btn-outline bg-base-100 shadow">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`,
		},

		// ── Navigation extras ─────────────────────────────────────────────────
		{
			Slug:        "page-title-minimal",
			Name:        "Page Title — Minimal",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Page Title",
			Description: "Breadcrumb-only page header with home icon. Lightweight variant for inner pages.",
			HTML: `<div class="p-6 bg-base-100">
  <div class="bg-base-100 rounded-lg border border-base-200 px-4 py-3">
    <div class="flex w-full items-center justify-between gap-3">
      <p class="font-semibold text-base-content">Create New Case</p>
      <div class="text-base-content/80 flex items-center gap-3 text-sm">
        <a href="#" aria-label="Home">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
        </a>
        <span>/</span>
        <a class="inline-flex items-center gap-1 hover:text-primary" href="#">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
          Cases
        </a>
        <span>/</span>
        <span class="text-base-content font-medium">New</span>
      </div>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "page-title-editor",
			Name:        "Page Title — Editor",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Page Title",
			Description: "Full page title with DaisyUI breadcrumbs, subtitle meta, and action buttons.",
			HTML: `<div class="p-6 bg-base-100">
  <div class="bg-base-100 rounded-lg border border-base-200 px-4 py-3">
    <div class="breadcrumbs text-sm mb-1">
      <ul>
        <li>
          <a href="#" class="flex items-center gap-1">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
            Dashboard
          </a>
        </li>
        <li>
          <a href="#" class="flex items-center gap-1">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
            Cases
          </a>
        </li>
        <li>Johnson v. Smith</li>
      </ul>
    </div>
    <div class="flex items-end justify-between gap-3">
      <div>
        <p class="font-semibold text-base-content sm:text-lg">Johnson v. Smith</p>
        <div class="text-base-content/60 mt-0.5 flex items-center gap-3 text-sm">
          <div class="flex items-center gap-1">
            <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
            Type: Civil Litigation
          </div>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <button class="btn btn-primary btn-sm">Save Changes</button>
        <button class="btn btn-outline btn-sm border-base-300">Preview</button>
        <button class="btn btn-outline btn-sm border-base-300 btn-square" aria-label="More options">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="5" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="12" cy="19" r="1"/></svg>
        </button>
      </div>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "search-dropdown",
			Name:        "Search — Dropdown",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Inline search input with a results dropdown showing recent searches and suggested items. CSS-only — no JS required.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="dropdown dropdown-bottom w-full max-w-sm">
    <label tabindex="0" class="input input-sm input-bordered flex items-center gap-2 px-2 w-full">
      <svg class="size-4 text-base-content/50" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
      <input placeholder="Search cases, clients…" class="grow bg-transparent text-sm focus:outline-none" type="text"/>
    </label>
    <div tabindex="0" class="dropdown-content bg-base-100 rounded-box mt-1 w-full shadow-lg border border-base-200 z-50">
      <ul class="menu menu-sm w-full p-1">
        <li><p class="menu-title text-xs text-base-content/50 px-3 py-1">Recent</p></li>
        <li>
          <a class="flex items-center gap-2 text-sm">
            <svg class="size-4 text-base-content/50" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
            Johnson v. Smith
          </a>
        </li>
        <li>
          <a class="flex items-center gap-2 text-sm">
            <svg class="size-4 text-base-content/50" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
            Alice Johnson
          </a>
        </li>
        <li><p class="menu-title text-xs text-base-content/50 px-3 py-1 mt-1">Suggestions</p></li>
        <li>
          <a class="flex items-center gap-2 text-sm">
            <svg class="size-4 text-base-content/50" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            Contract_Draft_v3.pdf
          </a>
        </li>
      </ul>
    </div>
  </div>
</div>`,
		},

		// ── Forms extras ──────────────────────────────────────────────────────
		{
			Slug:        "filter-bar",
			Name:        "Filter Bar",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Filters",
			Description: "FilterCard wraps filter inputs in a card with Filter/Clear buttons. CompactFilterBar is the inline variant used above tables.",
			HTML: `<div class="p-6 space-y-6">
  <div>
    <p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">FilterCard</p>
    <div class="card bg-base-100 shadow-sm border border-base-200 mb-4">
      <div class="card-body p-4">
        <form class="flex flex-wrap gap-3 items-end">
          <div class="form-control">
            <label class="label pb-1"><span class="label-text text-sm font-medium text-base-content/80">Search</span></label>
            <input type="search" placeholder="Search cases…" class="input input-bordered input-sm w-full max-w-xs"/>
          </div>
          <div class="form-control">
            <label class="label pb-1"><span class="label-text text-sm font-medium text-base-content/80">Status</span></label>
            <select class="select select-bordered select-sm">
              <option value="">All statuses</option>
              <option>Active</option>
              <option>Pending</option>
              <option>Closed</option>
            </select>
          </div>
          <div class="flex gap-2">
            <button type="submit" class="btn btn-primary btn-sm">Filter</button>
            <a href="#" class="btn btn-ghost btn-sm">Clear</a>
          </div>
        </form>
      </div>
    </div>
  </div>
  <div>
    <p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">CompactFilterBar</p>
    <form class="flex flex-wrap gap-3 items-end mb-4">
      <input type="search" placeholder="Search…" class="input input-bordered input-sm w-full max-w-xs"/>
      <select class="select select-bordered select-sm">
        <option value="">All statuses</option>
        <option>Active</option>
        <option>Closed</option>
      </select>
      <div class="flex gap-2">
        <button type="submit" class="btn btn-primary btn-sm">Filter</button>
        <a href="#" class="btn btn-ghost btn-sm">Clear</a>
      </div>
    </form>
  </div>
</div>`,
		},
		{
			Slug:        "fieldset",
			Name:        "Fieldset",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Layout",
			Description: "Fieldset wrapper with an optional legend label grouping related form inputs.",
			HTML: `<div class="p-6 max-w-md space-y-4">
  <fieldset class="fieldset">
    <legend class="fieldset-legend">Personal Information</legend>
    <div class="form-control mb-3">
      <label class="label pb-1"><span class="label-text font-medium text-sm text-base-content/80">Full name</span></label>
      <input type="text" placeholder="Alice Johnson" class="input input-bordered w-full"/>
    </div>
    <div class="form-control mb-3">
      <label class="label pb-1"><span class="label-text font-medium text-sm text-base-content/80">Email</span></label>
      <input type="email" placeholder="alice@example.com" class="input input-bordered w-full"/>
    </div>
  </fieldset>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">Case Details</legend>
    <div class="form-control mb-3">
      <label class="label pb-1"><span class="label-text font-medium text-sm text-base-content/80">Case type</span></label>
      <select class="select select-bordered w-full">
        <option>Civil</option>
        <option>Criminal</option>
        <option>Family</option>
      </select>
    </div>
  </fieldset>
</div>`,
		},
		{
			Slug:        "prompt-bar-minimal",
			Name:        "Prompt Bar — Minimal",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Minimal AI prompt / chat input with token counter and send button.",
			HTML: `<div class="p-6 flex justify-center">
  <div class="card bg-base-100 card-border w-full max-w-2xl">
    <div class="card-body p-3">
      <textarea
        class="textarea w-full h-28 resize-none border-0 p-1 text-base focus:outline-none m-0"
        placeholder="Describe what you want to generate or ask a question…"></textarea>
      <div class="mt-2 flex items-end justify-between">
        <div class="flex items-center gap-0.5">
          <button type="button" class="btn btn-sm btn-square btn-ghost" aria-label="Attach file">
            <svg class="size-4 text-base-content/60" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21.44 11.05l-9.19 9.19a6 6 0 01-8.49-8.49l9.19-9.19a4 4 0 015.66 5.66l-9.2 9.19a2 2 0 01-2.83-2.83l8.49-8.48"/></svg>
          </button>
          <button type="button" class="btn btn-sm btn-square btn-ghost" aria-label="Insert image">
            <svg class="size-4 text-base-content/60" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
          </button>
        </div>
        <span class="text-xs">
          <span class="text-base-content/80">Tokens:</span>
          <span class="text-error font-medium">88</span>
          <span class="text-base-content/60">/100</span>
        </span>
        <div class="flex items-center gap-1">
          <button type="button" class="btn btn-sm btn-square btn-ghost" aria-label="Voice input">
            <svg class="size-4 text-base-content/60" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 1a3 3 0 00-3 3v8a3 3 0 006 0V4a3 3 0 00-3-3z"/><path d="M19 10v2a7 7 0 01-14 0v-2"/><line x1="12" y1="19" x2="12" y2="23"/><line x1="8" y1="23" x2="16" y2="23"/></svg>
          </button>
          <button type="button" class="btn btn-primary btn-square btn-sm" aria-label="Send">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="19" x2="12" y2="5"/><polyline points="5 12 12 5 19 12"/></svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "prompt-bar-action",
			Name:        "Prompt Bar — Action",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "AI prompt input with quick-action buttons (Add File, Deep Thinking, Browsing).",
			HTML: `<div class="p-6 flex justify-center">
  <div class="card bg-base-100 card-border w-full max-w-2xl">
    <div class="card-body p-3">
      <textarea
        class="textarea w-full h-28 resize-none border-0 p-1 text-base focus:outline-none m-0"
        placeholder="Type your request or attach files to get started…"></textarea>
      <div class="mt-2 flex items-end justify-between">
        <div class="flex items-center gap-2">
          <button type="button" class="btn btn-sm btn-outline border-base-300 text-base-content/80 gap-1">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
            Add File
          </button>
          <button type="button" class="btn btn-sm btn-outline border-base-300 text-base-content/80 gap-1">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/></svg>
            Deep Think
          </button>
        </div>
        <button type="button" class="btn btn-primary btn-square btn-sm" aria-label="Send">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="19" x2="12" y2="5"/><polyline points="5 12 12 5 19 12"/></svg>
        </button>
      </div>
    </div>
  </div>
</div>`,
		},

		// ── Foundation extras ─────────────────────────────────────────────────
		{
			Slug:        "gradient-text",
			Name:        "Gradient Text",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Gradient text using Tailwind v4 bg-linear-to-r + bg-clip-text. Useful for hero headings.",
			HTML: `<div class="p-6 space-y-6">
  <p class="inline-block bg-linear-to-r from-primary to-secondary bg-clip-text text-3xl font-black text-transparent">
    go-daisy — UI Component Library
  </p>
  <p class="inline-block bg-linear-to-r from-success to-info bg-clip-text text-2xl font-bold text-transparent">
    Powered by DaisyUI + HTMX
  </p>
  <p class="inline-block bg-linear-to-r from-warning to-error bg-clip-text text-xl font-semibold text-transparent">
    Deadline approaching — 3 days left
  </p>
  <p class="text-sm text-base-content/60">Uses <code class="bg-base-200 px-1 rounded text-xs">bg-linear-to-r from-X to-Y bg-clip-text text-transparent</code></p>
</div>`,
		},
		{
			Slug:        "colored-shadows",
			Name:        "Colored Shadows",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Tailwind v4 colored shadow utilities using shadow-{color}/{opacity}.",
			HTML: `<div class="p-8 space-y-6">
  <div class="flex flex-wrap gap-6">
    <div class="card bg-base-100 rounded-box shadow-lg shadow-primary/20 p-4 w-36 text-center">
      <p class="text-sm font-semibold">Primary</p>
      <p class="text-xs text-base-content/60 mt-1">shadow-primary/20</p>
    </div>
    <div class="card bg-base-100 rounded-box shadow-lg shadow-secondary/20 p-4 w-36 text-center">
      <p class="text-sm font-semibold">Secondary</p>
      <p class="text-xs text-base-content/60 mt-1">shadow-secondary/20</p>
    </div>
    <div class="card bg-base-100 rounded-box shadow-lg shadow-success/20 p-4 w-36 text-center">
      <p class="text-sm font-semibold">Success</p>
      <p class="text-xs text-base-content/60 mt-1">shadow-success/20</p>
    </div>
    <div class="card bg-base-100 rounded-box shadow-lg shadow-error/20 p-4 w-36 text-center">
      <p class="text-sm font-semibold">Error</p>
      <p class="text-xs text-base-content/60 mt-1">shadow-error/20</p>
    </div>
    <div class="card bg-base-100 rounded-box shadow-lg shadow-warning/20 p-4 w-36 text-center">
      <p class="text-sm font-semibold">Warning</p>
      <p class="text-xs text-base-content/60 mt-1">shadow-warning/20</p>
    </div>
  </div>
  <div class="flex flex-wrap gap-4">
    <button class="btn btn-primary shadow-lg shadow-primary/30">Primary Button</button>
    <button class="btn btn-success shadow-lg shadow-success/30">Success Button</button>
    <button class="btn btn-error shadow-lg shadow-error/30">Danger Button</button>
  </div>
</div>`,
		},

		// ── Foundation extras ─────────────────────────────────────────────────────
		{
			Slug:        "typography",
			Name:        "Typography",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Typography",
			Description: "Heading and body text scale used across the application.",
			HTML: `<div class="p-6 space-y-3">
  <h1 class="text-3xl font-bold text-base-content">Heading 1</h1>
  <h2 class="text-2xl font-semibold text-base-content">Heading 2</h2>
  <h3 class="text-xl font-semibold text-base-content">Heading 3</h3>
  <h4 class="text-base font-semibold text-base-content">Heading 4</h4>
  <p class="text-base text-base-content/80">Body text — regular paragraph content used in cards and detail views.</p>
  <p class="text-sm text-base-content/60">Small / muted text — used for labels, hints, and secondary information.</p>
  <p class="text-xs text-base-content/50 uppercase tracking-wide font-semibold">Overline / label text</p>
  <a href="#" class="link link-primary text-sm">Link text</a>
</div>`,
		},
		{
			Slug:        "typography-scale",
			Name:        "Typography Scale",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Typography",
			Description: "Text size scale (xs→6xl) and font weight scale (thin→black).",
			HTML: `<div class="space-y-6 p-6">
  <div class="card card-border">
    <div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Sizes</div>
    <div class="flex flex-col gap-3 p-6">
      <p class="text-xs">The quick brown fox jumps… <span class="text-base-content/40">text-xs</span></p>
      <p class="text-sm">The quick brown fox jumps… <span class="text-base-content/40">text-sm</span></p>
      <p class="text-base">The quick brown fox jumps… <span class="text-base-content/40">text-base</span></p>
      <p class="text-lg">The quick brown fox jumps… <span class="text-base-content/40">text-lg</span></p>
      <p class="text-xl">The quick brown fox jumps… <span class="text-base-content/40">text-xl</span></p>
      <p class="text-2xl">The quick brown fox jumps… <span class="text-base-content/40">text-2xl</span></p>
      <p class="text-3xl">The quick brown fox jumps… <span class="text-base-content/40">text-3xl</span></p>
      <p class="text-4xl">The quick brown fox…  <span class="text-base-content/40">text-4xl</span></p>
    </div>
  </div>
  <div class="card card-border">
    <div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Weights</div>
    <div class="flex flex-col gap-3 p-6">
      <p class="font-thin">The quick brown fox jumps… <span class="text-base-content/40 font-normal">font-thin</span></p>
      <p class="font-light">The quick brown fox jumps… <span class="text-base-content/40 font-normal">font-light</span></p>
      <p class="font-normal">The quick brown fox jumps… <span class="text-base-content/40">font-normal</span></p>
      <p class="font-medium">The quick brown fox jumps… <span class="text-base-content/40">font-medium</span></p>
      <p class="font-semibold">The quick brown fox jumps… <span class="text-base-content/40">font-semibold</span></p>
      <p class="font-bold">The quick brown fox jumps… <span class="text-base-content/40">font-bold</span></p>
      <p class="font-extrabold">The quick brown fox jumps… <span class="text-base-content/40">font-extrabold</span></p>
      <p class="font-black">The quick brown fox jumps… <span class="text-base-content/40">font-black</span></p>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "shadow-scale",
			Name:        "Shadow Scale",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Box shadows from none→2xl, colored shadows, inset shadows, and text shadows.",
			HTML: `<div class="space-y-6 p-6">
  <div class="card card-border bg-base-200/20">
    <div class="bg-base-200/40 rounded-t-box px-5 py-3 font-medium">Box Shadow</div>
    <div class="grid grid-cols-2 gap-6 p-6 lg:grid-cols-4">
      <div class="bg-base-100 rounded-box text-base-content/60 flex h-24 items-center justify-center text-sm shadow-none">shadow-none</div>
      <div class="bg-base-100 rounded-box text-base-content/60 flex h-24 items-center justify-center text-sm shadow-sm">shadow-sm</div>
      <div class="bg-base-100 rounded-box text-base-content/60 flex h-24 items-center justify-center text-sm shadow-md">shadow-md</div>
      <div class="bg-base-100 rounded-box text-base-content/60 flex h-24 items-center justify-center text-sm shadow-lg">shadow-lg</div>
      <div class="bg-base-100 rounded-box text-base-content/60 flex h-24 items-center justify-center text-sm shadow-xl">shadow-xl</div>
      <div class="bg-base-100 rounded-box text-base-content/60 flex h-24 items-center justify-center text-sm shadow-2xl">shadow-2xl</div>
      <div class="bg-base-100 rounded-box text-base-content/60 shadow-primary/20 flex h-24 items-center justify-center text-sm shadow-lg">shadow-primary</div>
      <div class="bg-base-100 rounded-box text-base-content/60 shadow-error/20 flex h-24 items-center justify-center text-sm shadow-lg">shadow-error</div>
    </div>
  </div>
  <div class="card card-border">
    <div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Inset Shadow</div>
    <div class="grid grid-cols-2 gap-6 p-6 lg:grid-cols-4">
      <div class="bg-base-100 border-base-200 rounded-box text-base-content/60 flex h-24 items-center justify-center border text-xs inset-shadow-none">inset-none</div>
      <div class="bg-base-100 border-base-200 rounded-box text-base-content/60 flex h-24 items-center justify-center border text-xs inset-shadow-xs">inset-xs</div>
      <div class="bg-base-100 border-base-200 rounded-box text-base-content/60 flex h-24 items-center justify-center border text-xs inset-shadow-sm">inset-sm</div>
      <div class="bg-base-100 border-base-200 rounded-box text-base-content/60 inset-shadow-primary/15 flex h-24 items-center justify-center border text-xs inset-shadow-sm">inset-primary</div>
    </div>
  </div>
  <div class="card card-border">
    <div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Text Shadow</div>
    <div class="flex flex-col gap-3 p-6">
      <p class="font-semibold capitalize text-shadow-none sm:text-lg">text-shadow-none</p>
      <p class="font-semibold capitalize text-shadow-sm sm:text-lg">text-shadow-sm</p>
      <p class="font-semibold capitalize text-shadow-md sm:text-lg">text-shadow-md</p>
      <p class="font-semibold capitalize text-shadow-lg sm:text-lg">text-shadow-lg</p>
      <p class="text-primary text-shadow-primary/20 font-semibold capitalize text-shadow-lg sm:text-lg">text-shadow-primary</p>
      <p class="text-error text-shadow-error/20 font-semibold capitalize text-shadow-lg sm:text-lg">text-shadow-error</p>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "css-filters",
			Name:        "CSS Filters",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Tailwind filter utilities: grayscale, invert, sepia, blur, brightness, contrast, saturate.",
			HTML: `<div class="p-6">
  <div class="card card-border">
    <div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Image Filters</div>
    <div class="grid grid-cols-3 gap-6 p-6 lg:grid-cols-4">
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">Normal</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 grayscale bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">grayscale</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 invert bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">invert</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 sepia bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">sepia</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 blur-sm bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">blur-sm</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 brightness-125 bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">brightness-125</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 contrast-200 bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">contrast-200</p>
      </div>
      <div class="flex flex-col items-center gap-2">
        <div class="bg-primary rounded-box size-24 saturate-200 bg-cover bg-center" style="background-image:url(https://picsum.photos/seed/a/96/96)"></div>
        <p class="text-base-content/60 text-xs">saturate-200</p>
      </div>
    </div>
  </div>
</div>`,
		},

		// ── Navigation extras ──────────────────────────────────────────────────────
		{
			Slug:        "footer-minimal",
			Name:        "Footer — Minimal",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Simple one-line footer with copyright text and optional links.",
			HTML: `<div class="space-y-4 p-6">
  <div class="card card-border">
    <div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Copyright only</div>
    <div class="flex w-full items-center justify-center px-4 py-3 border-t border-base-200">
      <span class="text-base-content/80 text-sm">© 2025 LegalPlant. All rights reserved.</span>
    </div>
  </div>
</div>`,
		},
		{
			Slug:        "profile-menu",
			Name:        "Profile Menu",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Avatar dropdown menu with grouped menu items and sign-out action.",
			HTML: `<div class="flex items-center justify-center p-12">
  <div class="dropdown dropdown-bottom dropdown-end">
    <div tabindex="0" class="avatar bg-base-200 size-12 cursor-pointer overflow-hidden rounded-full">
      <div class="bg-primary text-primary-content flex size-full items-center justify-center font-semibold text-lg">JD</div>
    </div>
    <div tabindex="0" class="dropdown-content bg-base-100 rounded-box mt-2 w-56 shadow-lg">
      <div class="px-4 py-3 border-b border-base-200">
        <p class="font-semibold text-sm">Jane Doe</p>
        <p class="text-base-content/60 text-xs mt-0.5">jane@example.com</p>
      </div>
      <ul class="menu w-full p-2">
        <li>
          <a href="#">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0012 15.75a7.488 7.488 0 00-5.982 2.975m11.963 0a9 9 0 10-11.963 0m11.963 0A8.966 8.966 0 0112 21a8.966 8.966 0 01-5.982-2.275M15 9.75a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
            <span>Profile</span>
          </a>
        </li>
        <li>
          <a href="#">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"/></svg>
            <span>Settings</span>
          </a>
        </li>
        <li>
          <a href="#">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0"/></svg>
            <span>Notifications</span>
            <span class="badge badge-sm badge-primary">3</span>
          </a>
        </li>
      </ul>
      <hr class="border-base-200" />
      <ul class="menu w-full p-2">
        <li>
          <a class="text-error hover:bg-error/10" href="#">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15m3 0l3-3m0 0l-3-3m3 3H9"/></svg>
            <span>Sign Out</span>
          </a>
        </li>
      </ul>
    </div>
  </div>
</div>`,
		},

		// ── Forms extras ───────────────────────────────────────────────────────────
		{
			Slug:        "input-spinner",
			Name:        "Input Spinner",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Numeric increment/decrement input with +/- buttons. Uses vanilla JS — no library needed. Includes simple and joined variants.",
			HTML: `<div class="p-6 space-y-6">
  <div>
    <p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Simple</p>
    <div class="flex items-center gap-3">
      <button type="button" class="btn btn-square btn-outline" onclick="spinnerDecrement('spin1')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
      <input id="spin1" type="number" value="0" min="0" max="99"
        class="input input-bordered w-24 text-center [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"/>
      <button type="button" class="btn btn-square btn-outline" onclick="spinnerIncrement('spin1')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
    </div>
  </div>
  <div>
    <p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Join variant</p>
    <div class="join">
      <button type="button" class="btn btn-outline btn-square join-item" onclick="spinnerDecrement('spin2', 5)">-5</button>
      <button type="button" class="btn btn-outline btn-square join-item" onclick="spinnerDecrement('spin2')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
      <input id="spin2" type="number" value="0"
        class="input input-bordered join-item w-20 text-center [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"/>
      <button type="button" class="btn btn-outline btn-square join-item" onclick="spinnerIncrement('spin2')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
      <button type="button" class="btn btn-outline btn-square join-item" onclick="spinnerIncrement('spin2', 5)">+5</button>
    </div>
  </div>
  <div>
    <p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">With min/max (0–10)</p>
    <div class="flex items-center gap-3">
      <button type="button" class="btn btn-square btn-primary btn-sm" onclick="spinnerDecrement('spin3')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
      <input id="spin3" type="number" value="5" min="0" max="10"
        class="input input-bordered input-sm w-20 text-center [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"/>
      <button type="button" class="btn btn-square btn-primary btn-sm" onclick="spinnerIncrement('spin3')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
    </div>
  </div>
</div>
<script>
function spinnerGet(id) { return document.getElementById(id); }
function spinnerIncrement(id, step) {
  var el = spinnerGet(id); if (!el) return;
  var v = parseFloat(el.value)||0, s = step||1;
  var max = el.max !== '' ? parseFloat(el.max) : Infinity;
  el.value = Math.min(v + s, max);
}
function spinnerDecrement(id, step) {
  var el = spinnerGet(id); if (!el) return;
  var v = parseFloat(el.value)||0, s = step||1;
  var min = el.min !== '' ? parseFloat(el.min) : -Infinity;
  el.value = Math.max(v - s, min);
}
</script>`,
		},
		{
			Slug:        "wizard-stepper",
			Name:        "Wizard — Stepper",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Wizard",
			Description: "Multi-step form wizard with step indicators, next/prev navigation, and a finish action. Implemented in vanilla JS — no Alpine.js needed.",
			HTML: `<div class="p-6 max-w-lg mx-auto" id="wizard-demo">
  <!-- Step indicator -->
  <div class="flex items-center gap-2 mb-6" id="wizard-steps">
    <div class="wizard-step-indicator flex items-center gap-2 cursor-pointer" data-step="1" onclick="wizardGoTo(1)">
      <div id="wizard-dot-1" class="flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-primary text-primary-content">1</div>
      <span class="font-medium text-sm max-lg:hidden">Intake</span>
    </div>
    <div class="h-px flex-1 bg-base-300"></div>
    <div class="wizard-step-indicator flex items-center gap-2 cursor-pointer" data-step="2" onclick="wizardGoTo(2)">
      <div id="wizard-dot-2" class="flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-base-200 text-base-content/50">2</div>
      <span class="font-medium text-sm text-base-content/50 max-lg:hidden">Details</span>
    </div>
    <div class="h-px flex-1 bg-base-300"></div>
    <div class="wizard-step-indicator flex items-center gap-2 cursor-pointer" data-step="3" onclick="wizardGoTo(3)">
      <div id="wizard-dot-3" class="flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-base-200 text-base-content/50">3</div>
      <span class="font-medium text-sm text-base-content/50 max-lg:hidden">Team</span>
    </div>
    <div class="h-px flex-1 bg-base-300"></div>
    <div class="wizard-step-indicator flex items-center gap-2 cursor-pointer" data-step="4" onclick="wizardGoTo(4)">
      <div id="wizard-dot-4" class="flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-base-200 text-base-content/50">4</div>
      <span class="font-medium text-sm text-base-content/50 max-lg:hidden">Review</span>
    </div>
  </div>
  <!-- Step panels -->
  <div class="bg-base-100 border border-base-300 rounded-box p-5 min-h-40">
    <div id="wizard-panel-1">
      <h3 class="font-semibold mb-3">Step 1 — Intake</h3>
      <div class="form-control mb-3">
        <label class="label pb-1"><span class="label-text text-sm font-medium">Case title</span></label>
        <input type="text" placeholder="e.g. Johnson v. Smith" class="input input-bordered w-full"/>
      </div>
      <div class="form-control">
        <label class="label pb-1"><span class="label-text text-sm font-medium">Case type</span></label>
        <select class="select select-bordered w-full"><option>Civil</option><option>Criminal</option><option>Family</option></select>
      </div>
    </div>
    <div id="wizard-panel-2" class="hidden">
      <h3 class="font-semibold mb-3">Step 2 — Details</h3>
      <div class="form-control mb-3">
        <label class="label pb-1"><span class="label-text text-sm font-medium">Description</span></label>
        <textarea class="textarea textarea-bordered w-full" rows="3" placeholder="Brief description of the case…"></textarea>
      </div>
      <div class="form-control">
        <label class="label pb-1"><span class="label-text text-sm font-medium">Priority</span></label>
        <select class="select select-bordered w-full"><option>Normal</option><option>High</option><option>Urgent</option></select>
      </div>
    </div>
    <div id="wizard-panel-3" class="hidden">
      <h3 class="font-semibold mb-3">Step 3 — Team</h3>
      <div class="form-control">
        <label class="label pb-1"><span class="label-text text-sm font-medium">Lead attorney</span></label>
        <select class="select select-bordered w-full"><option>Alice Johnson</option><option>Bob Smith</option><option>Carol White</option></select>
      </div>
    </div>
    <div id="wizard-panel-4" class="hidden">
      <h3 class="font-semibold mb-3">Step 4 — Review</h3>
      <p class="text-sm text-base-content/60 mb-4">Review the case details before submitting.</p>
      <div class="space-y-2 text-sm">
        <div class="flex gap-2"><span class="text-base-content/60 w-24">Title:</span><span class="font-medium">Johnson v. Smith</span></div>
        <div class="flex gap-2"><span class="text-base-content/60 w-24">Type:</span><span class="font-medium">Civil</span></div>
        <div class="flex gap-2"><span class="text-base-content/60 w-24">Attorney:</span><span class="font-medium">Alice Johnson</span></div>
      </div>
    </div>
  </div>
  <!-- Navigation -->
  <div class="mt-4 flex justify-between">
    <button id="wizard-prev" class="btn btn-ghost btn-sm" onclick="wizardPrev()" disabled>
      <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
      Back
    </button>
    <button id="wizard-next" class="btn btn-primary btn-sm" onclick="wizardNext()">
      Next
      <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
    </button>
  </div>
</div>
<script>
(function() {
  var current = 1, total = 4;
  window.wizardGoTo = function(n) {
    document.getElementById('wizard-panel-' + current).classList.add('hidden');
    current = Math.max(1, Math.min(total, n));
    document.getElementById('wizard-panel-' + current).classList.remove('hidden');
    for (var i = 1; i <= total; i++) {
      var dot = document.getElementById('wizard-dot-' + i);
      if (i < current) {
        dot.className = 'flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-primary text-primary-content';
        dot.innerHTML = '<svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>';
      } else if (i === current) {
        dot.className = 'flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-primary text-primary-content';
        dot.textContent = i;
      } else {
        dot.className = 'flex size-7 items-center justify-center rounded-full font-semibold text-sm bg-base-200 text-base-content/50';
        dot.textContent = i;
      }
    }
    document.getElementById('wizard-prev').disabled = current === 1;
    var nextBtn = document.getElementById('wizard-next');
    if (current === total) {
      nextBtn.textContent = 'Finish';
      nextBtn.className = 'btn btn-success btn-sm';
    } else {
      nextBtn.innerHTML = 'Next <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>';
      nextBtn.className = 'btn btn-primary btn-sm';
    }
  };
  window.wizardNext = function() { if (current < total) wizardGoTo(current + 1); };
  window.wizardPrev = function() { if (current > 1) wizardGoTo(current - 1); };
})();
</script>`,
		},
		{
			Slug:        "clipboard-copy",
			Name:        "Clipboard Copy",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Click-to-copy buttons with visual feedback. Uses vanilla JS navigator.clipboard — no library needed.",
			HTML: `<div class="p-6 space-y-4 max-w-lg">
  <div>
    <p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Copy text field</p>
    <div class="flex items-center gap-2">
      <input id="copy-input-1" type="text" value="CASE-2026-00142" readonly
        class="input input-bordered input-sm flex-1 font-mono text-sm"/>
      <button class="btn btn-sm btn-outline" onclick="copyToClipboard('copy-input-1', this)">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
        Copy
      </button>
    </div>
  </div>
  <div>
    <p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Copy share link</p>
    <div class="flex items-center gap-2">
      <input id="copy-input-2" type="text" value="https://app.example.com/cases/CASE-2026-00142" readonly
        class="input input-bordered input-sm flex-1 text-sm text-base-content/60"/>
      <button class="btn btn-sm btn-primary" onclick="copyToClipboard('copy-input-2', this)">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
        Copy Link
      </button>
    </div>
  </div>
  <div>
    <p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Inline copy badge</p>
    <div class="flex items-center gap-2">
      <code class="bg-base-200 text-base-content px-3 py-1.5 rounded text-sm font-mono">CASE-2026-00142</code>
      <button class="btn btn-ghost btn-xs gap-1" onclick="copyText('CASE-2026-00142', this)">
        <svg class="size-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
        Copy
      </button>
    </div>
  </div>
</div>
<script>
function copyToClipboard(inputId, btn) {
  var el = document.getElementById(inputId);
  if (!el) return;
  navigator.clipboard.writeText(el.value).then(function() {
    var orig = btn.innerHTML;
    btn.innerHTML = '✓ Copied!';
    btn.disabled = true;
    setTimeout(function() { btn.innerHTML = orig; btn.disabled = false; }, 2000);
  });
}
function copyText(text, btn) {
  navigator.clipboard.writeText(text).then(function() {
    var orig = btn.innerHTML;
    btn.innerHTML = '✓';
    setTimeout(function() { btn.innerHTML = orig; }, 2000);
  });
}
</script>`,
		},

		// ── Real component entries (WithBoundary + RenderFunc) ───────────────────

		// ui.Button
		{
			Slug:        "button",
			Name:        "Button",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Buttons",
			Description: "A DaisyUI button with configurable variant, size, type, shape, icon, and loading state.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Default",
					Description: "Standard button with live controls.",
					RenderFunc: func(params url.Values) templ.Component {
						variant := ui.ButtonVariant(params.Get("variant"))
						if variant == "" {
							variant = ui.ButtonPrimary
						}
						size := ui.ButtonSize(params.Get("size"))
						typ := ui.ButtonType(params.Get("typ"))
						if typ == "" {
							typ = ui.ButtonTypeButton
						}
						href := params.Get("href")
						if href == "" {
							href = "#"
						}
						shape := ui.ButtonShape(params.Get("shape"))
						icon := params.Get("icon")
						loading := params.Get("loading") == "true"
						if shape == ui.ButtonShapeDefault && typ != ui.ButtonTypeLink {
							return withText("Save changes", ui.ButtonWithBoundary(href, variant, size, typ, shape, icon, loading))
						}
						if typ == ui.ButtonTypeLink {
							return withText("Go to dashboard", ui.ButtonWithBoundary(href, variant, size, typ, shape, icon, loading))
						}
						return ui.ButtonWithBoundary(href, variant, size, typ, shape, icon, loading)
					},
					Tokens: ButtonTokens(),
				},
			},
		},

		// ui.Badge
		{
			Slug:        "badge",
			Name:        "Badge",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Badges",
			Description: "A DaisyUI badge with configurable intent, style, size, and optional icon.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live intent, style, size, and icon controls.",
					RenderFunc: func(params url.Values) templ.Component {
						variant := ui.BadgeIntent(params.Get("variant"))
						if variant == "" {
							variant = ui.BadgePrimary
						}
						style := ui.BadgeStyle(params.Get("style"))
						size := ui.BadgeSize(params.Get("size"))
						icon := params.Get("icon")
						label := params.Get("label")
						if label == "" {
							label = "Active"
						}
						return ui.BadgeWithBoundary(variant, style, size, icon, label)
					},
					Tokens: BadgeTokens(),
				},
			},
		},

		// ui.StatusBadge
		{
			Slug:        "status-badge-real",
			Name:        "Status Badge",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Badges",
			Description: "Maps a string status to an appropriate intent badge automatically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live status control.",
					RenderFunc: func(params url.Values) templ.Component {
						status := params.Get("status")
						if status == "" {
							status = "active"
						}
						return ui.StatusBadgeWithBoundary(status)
					},
					Tokens: StatusBadgeTokens(),
				},
			},
		},

		// ui.Card
		{
			Slug:        "card-real",
			Name:        "Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "A DaisyUI card container with a title.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live title control.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Card Title"
						}
						return withChildren(ui.CardWithBoundary(title), rawHTML(`<p class="text-sm text-base-content/70">Card body content goes here.</p>`))
					},
					Tokens: CardTokens(),
				},
			},
		},

		// ui.InlineAlert (no icon)
		{
			Slug:        "inline-alert",
			Name:        "Inline Alert",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Alerts",
			Description: "A plain inline alert without an icon.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live type and message controls.",
					RenderFunc: func(params url.Values) templ.Component {
						typ := ui.AlertType(params.Get("typ"))
						if typ == "" {
							typ = ui.AlertSuccess
						}
						message := params.Get("message")
						if message == "" {
							message = "Operation completed successfully."
						}
						return ui.AlertWithBoundary(typ, message)
					},
					Tokens: AlertTokens(),
				},
			},
		},

		// ui.Toast
		{
			Slug:        "toast-real",
			Name:        "Toast",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Toasts",
			Description: "A toast notification with type and message.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live type and message controls.",
					RenderFunc: func(params url.Values) templ.Component {
						typ := ui.ToastType(params.Get("typ"))
						if typ == "" {
							typ = ui.ToastSuccess
						}
						message := params.Get("message")
						if message == "" {
							message = "Action completed successfully."
						}
						return ui.ToastWithBoundary(typ, message)
					},
					Tokens: ToastTokens(),
				},
			},
		},

		// ui.Pagination
		{
			Slug:        "pagination-real",
			Name:        "Pagination",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Pagination",
			Description: "A DaisyUI pagination control.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live current page control.",
					RenderFunc: func(params url.Values) templ.Component {
						page := 1
						if p := params.Get("currentPage"); p != "" {
							if v, err := parseInt(p); err == nil && v > 0 {
								page = v
							}
						}
						totalPages := 10
						if p := params.Get("totalPages"); p != "" {
							if v, err := parseInt(p); err == nil && v > 0 {
								totalPages = v
							}
						}
						return ui.PaginationWithBoundary(page, totalPages, "#", "main-content")
					},
					Tokens: PaginationTokens(),
				},
			},
		},

		// ui.Empty
		{
			Slug:        "empty-state-real",
			Name:        "Empty State",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "States",
			Description: "An empty state placeholder with icon, title, and description.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live title and description controls.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "No results found"
						}
						desc := params.Get("description")
						if desc == "" {
							desc = "Try adjusting your search or filters."
						}
						return ui.EmptyWithBoundary("lucide--search", title, desc)
					},
					Tokens: EmptyTokens(),
				},
			},
		},

		// ui.Loader
		{
			Slug:        "loader",
			Name:        "Loader",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Loading",
			Description: "Centered loading spinner for async content areas.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "A centered DaisyUI loading spinner with no configurable props.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.LoaderWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ui.StatCard
		{
			Slug:        "stat-card-real",
			Name:        "Stat Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "A compact summary stat widget with icon, value, and label.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live label, value, and icon controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Active Sessions"
						}
						value := params.Get("value")
						if value == "" {
							value = "42"
						}
						icon := params.Get("icon")
						if icon == "" {
							icon = "lucide--users"
						}
						iconColor := params.Get("iconColor")
						if iconColor == "" {
							iconColor = "bg-primary/10 text-primary"
						}
						return ui.StatCardWithBoundary(ui.StatCardProps{
							Label:     label,
							Value:     value,
							Icon:      icon,
							IconColor: iconColor,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{
							Label:      "Label",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "Active Sessions",
							QueryParam: "label",
						},
						{
							Label:      "Value",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "42",
							QueryParam: "value",
						},
						{
							Label:      "Icon",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "lucide--users",
							QueryParam: "icon",
						},
						{
							Label:      "Icon Color",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "bg-primary/10 text-primary",
							QueryParam: "iconColor",
							Options: []galleryruntime.TokenOption{
								{Value: "bg-primary/10 text-primary", Label: "Primary"},
								{Value: "bg-secondary/10 text-secondary", Label: "Secondary"},
								{Value: "bg-success/10 text-success", Label: "Success"},
								{Value: "bg-error/10 text-error", Label: "Error"},
								{Value: "bg-warning/10 text-warning", Label: "Warning"},
								{Value: "bg-info/10 text-info", Label: "Info"},
							},
						},
					},
				},
			},
		},

		// ui.ActionMenu
		{
			Slug:        "action-menu-real",
			Name:        "Action Menu",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Menus",
			Description: "A dropdown action menu with configurable items.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sample action menu with three items.",
					RenderFunc: func(params url.Values) templ.Component {
						item1 := params.Get("items1")
						if item1 == "" {
							item1 = "Edit"
						}
						item2 := params.Get("items2")
						if item2 == "" {
							item2 = "Duplicate"
						}
						item3 := params.Get("items3")
						if item3 == "" {
							item3 = "Delete"
						}
						return ui.ActionMenuWithBoundary([]ui.ActionMenuItem{
							{Label: item1, Icon: "lucide--pencil", HXGet: "#"},
							{Label: item2, Icon: "lucide--copy", HXGet: "#"},
							{Label: item3, Icon: "lucide--trash-2", HXGet: "#", Danger: true},
						})
					},
					Tokens: ActionMenuTokens(),
				},
			},
		},

		// ui.Avatar
		{
			Slug:        "avatar-real",
			Name:        "Avatar",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Avatars",
			Description: "An avatar with initials fallback and configurable size.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live name and size controls.",
					RenderFunc: func(params url.Values) templ.Component {
						name := params.Get("name")
						if name == "" {
							name = "Jane Smith"
						}
						size := ui.AvatarSize(params.Get("size"))
						if size == "" {
							size = ui.AvatarMD
						}
						return ui.AvatarWithBoundary(name, "", size)
					},
					Tokens: AvatarTokens(),
				},
			},
		},

		// form.TextInput
		{
			Slug:        "text-input",
			Name:        "Text Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "A labelled text input field.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live label and required controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Email address"
						}
						required := params.Get("required") == "true"
						return form.TextInputWithBoundary("email", label, "", "", required)
					},
					Tokens: TextInputTokens(),
				},
			},
		},

		// form.TextareaInput
		{
			Slug:        "textarea-input",
			Name:        "Textarea Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "A labelled textarea input field.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live label and rows controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Description"
						}
						rows := 4
						if r := params.Get("rows"); r != "" {
							if v, err := parseInt(r); err == nil && v > 0 {
								rows = v
							}
						}
						required := params.Get("required") == "true"
						return form.TextareaInputWithBoundary("description", label, "", "", rows, required)
					},
					Tokens: TextareaInputTokens(),
				},
			},
		},

		// form.CheckboxInput
		{
			Slug:        "checkbox-input",
			Name:        "Checkbox Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Toggles",
			Description: "A labelled checkbox input.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live label and checked state controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "I agree to the terms"
						}
						checked := params.Get("checked") == "true"
						return form.CheckboxInputWithBoundary("agree", label, checked, "")
					},
					Tokens: CheckboxInputTokens(),
				},
			},
		},

		// form.SelectInput
		{
			Slug:        "select-input",
			Name:        "Select Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "A labelled select dropdown.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live label and selected value controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Country"
						}
						selected := params.Get("selected")
						required := params.Get("required") == "true"
						return form.SelectInputWithBoundary("country", label, selected, [][2]string{
							{"us", "United States"},
							{"gb", "United Kingdom"},
							{"ca", "Canada"},
							{"au", "Australia"},
						}, "", required)
					},
					Tokens: SelectInputTokens(),
				},
			},
		},

		// form.RangeInput
		{
			Slug:        "range-input",
			Name:        "Range Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "A labelled range slider input.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live value and color controls.",
					RenderFunc: func(params url.Values) templ.Component {
						val := 50
						if v := params.Get("value"); v != "" {
							if n, err := parseInt(v); err == nil {
								val = n
							}
						}
						color := params.Get("color")
						if color == "" {
							color = "range-primary"
						}
						return form.RangeInputWithBoundary("volume", "Volume", val, 0, 100, 1, color)
					},
					Tokens: RangeInputTokens(),
				},
			},
		},

		// form.FormField (unified)
		{
			Slug:        "form-field-real",
			Name:        "Form Field",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "A unified form field that renders the appropriate input based on type.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live type, label, placeholder, required, disabled, and error controls.",
					RenderFunc: func(params url.Values) templ.Component {
						typ := form.FormFieldType(params.Get("typ"))
						if typ == "" {
							typ = form.FieldText
						}
						label := params.Get("label")
						if label == "" {
							label = "Full name"
						}
						placeholder := params.Get("placeholder")
						if placeholder == "" {
							placeholder = "Enter value..."
						}
						required := params.Get("required") == "true"
						disabled := params.Get("disabled") == "true"
						errMsg := params.Get("error")
						return form.FormFieldWithBoundary(form.FormFieldProps{
							Type:        typ,
							Name:        "demo",
							Label:       label,
							Placeholder: placeholder,
							Required:    required,
							Disabled:    disabled,
							Error:       errMsg,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{
							Label:      "Type",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "text",
							QueryParam: "typ",
							Options: []galleryruntime.TokenOption{
								{Value: "text", Label: "Text"},
								{Value: "textarea", Label: "Textarea"},
								{Value: "email", Label: "Email"},
								{Value: "number", Label: "Number"},
								{Value: "date", Label: "Date"},
								{Value: "checkbox", Label: "Checkbox"},
								{Value: "select", Label: "Select"},
							},
						},
						{
							Label:      "Label",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "Full name",
							QueryParam: "label",
						},
						{
							Label:      "Placeholder",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "Enter value...",
							QueryParam: "placeholder",
						},
						{
							Label:      "Required",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "false",
							QueryParam: "required",
							Options: []galleryruntime.TokenOption{
								{Value: "false", Label: "No"},
								{Value: "true", Label: "Yes"},
							},
						},
						{
							Label:      "Disabled",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "false",
							QueryParam: "disabled",
							Options: []galleryruntime.TokenOption{
								{Value: "false", Label: "No"},
								{Value: "true", Label: "Yes"},
							},
						},
						{
							Label:      "Error",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "",
							QueryParam: "error",
						},
					},
				},
			},
		},

		// form.SearchInput
		{
			Slug:        "search-input-real",
			Name:        "Search Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "A search input field with a magnifier icon.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live value and placeholder controls.",
					RenderFunc: func(params url.Values) templ.Component {
						value := params.Get("value")
						placeholder := params.Get("placeholder")
						if placeholder == "" {
							placeholder = "Search..."
						}
						return form.SearchInputWithBoundary("q", value, placeholder, "", "")
					},
					Tokens: []galleryruntime.DesignToken{
						{
							Label:      "Value",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "",
							QueryParam: "value",
						},
						{
							Label:      "Placeholder",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "Search...",
							QueryParam: "placeholder",
						},
					},
				},
			},
		},

		// nav.TopBar
		{
			Slug:        "top-bar-real",
			Name:        "Top Bar",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Page Title",
			Description: "A top navigation bar with a title.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live title control.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Dashboard"
						}
						return nav.TopBarWithBoundary(title)
					},
					Tokens: TopBarTokens(),
				},
			},
		},

		// nav.TabMenu
		{
			Slug:        "tab-menu-real",
			Name:        "Tab Menu",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Tabs",
			Description: "A full-page HTMX tab strip.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sample tab menu with three tabs.",
					RenderFunc: func(params url.Values) templ.Component {
						tab1 := params.Get("tabs1")
						if tab1 == "" {
							tab1 = "Overview"
						}
						tab2 := params.Get("tabs2")
						if tab2 == "" {
							tab2 = "Activity"
						}
						tab3 := params.Get("tabs3")
						if tab3 == "" {
							tab3 = "Settings"
						}
						tabs := []nav.Tab{
							{Label: tab1, Href: "#", Active: true},
							{Label: tab2, Href: "#"},
							{Label: tab3, Href: "#"},
						}
						return nav.TabMenuWithBoundary(tabs)
					},
					Tokens: TabMenuTokens(),
				},
			},
		},

		// nav.SimpleTabs
		{
			Slug:        "simple-tabs-real",
			Name:        "Simple Tabs",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Tabs",
			Description: "A simple tab strip without HTMX.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sample simple tabs.",
					RenderFunc: func(params url.Values) templ.Component {
						tab1 := params.Get("tabs1")
						if tab1 == "" {
							tab1 = "All"
						}
						tab2 := params.Get("tabs2")
						if tab2 == "" {
							tab2 = "Open"
						}
						tab3 := params.Get("tabs3")
						if tab3 == "" {
							tab3 = "Closed"
						}
						tabs := []nav.Tab{
							{Label: tab1, Href: "#", Active: true},
							{Label: tab2, Href: "#"},
							{Label: tab3, Href: "#"},
						}
						return nav.SimpleTabsWithBoundary(tabs)
					},
					Tokens: SimpleTabsTokens(),
				},
			},
		},

		// nav.PageHeader
		{
			Slug:        "page-header-real",
			Name:        "Page Header",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Headers",
			Description: "A breadcrumb-based page header.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sample page header with breadcrumb trail.",
					RenderFunc: func(params url.Values) templ.Component {
						step1 := params.Get("steps1")
						if step1 == "" {
							step1 = "Home"
						}
						step2 := params.Get("steps2")
						if step2 == "" {
							step2 = "Cases"
						}
						step3 := params.Get("steps3")
						if step3 == "" {
							step3 = "Edit Record"
						}
						return nav.PageHeaderWithBoundary(nav.Crumbs(step1, "/", step2, "/cases", step3))
					},
					Tokens: PageHeaderTokens(),
				},
			},
		},

		// nav.Menu
		{
			Slug:        "menu-real",
			Name:        "Menu",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Menus",
			Description: "A vertical navigation menu.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live size control.",
					RenderFunc: func(params url.Values) templ.Component {
						size := nav.MenuSize(params.Get("size"))
						item1 := params.Get("items1")
						if item1 == "" {
							item1 = "Dashboard"
						}
						item2 := params.Get("items2")
						if item2 == "" {
							item2 = "Cases"
						}
						item3 := params.Get("items3")
						if item3 == "" {
							item3 = "Contacts"
						}
						item4 := params.Get("items4")
						if item4 == "" {
							item4 = "Settings"
						}
						return nav.MenuWithBoundary(size, []nav.MenuItem{
							{Label: item1, Icon: "lucide--layout-dashboard", Href: "#", Active: true},
							{Label: item2, Icon: "lucide--folder-open", Href: "#"},
							{Label: item3, Icon: "lucide--users", Href: "#"},
							{Label: item4, Icon: "lucide--settings", Href: "#"},
						})
					},
					Tokens: MenuTokens(),
				},
			},
		},

		// modal.Modal
		{
			Slug:        "modal-real",
			Name:        "Modal",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Modals",
			Description: "A DaisyUI modal dialog.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live title and size controls.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Confirm Action"
						}
						size := modal.ModalSize(params.Get("size"))
						body := rawHTML(`<p class="text-sm text-base-content/70 mb-6">Are you sure you want to proceed? This action will be applied immediately.</p><div class="flex justify-end gap-2"><button type="button" class="btn btn-ghost btn-sm" onclick="this.closest('dialog').remove()">Cancel</button><button type="button" class="btn btn-primary btn-sm">Confirm</button></div>`)
						inner := withChildren(modal.ModalWithBoundary(title, size), body)
						// Wrap in a min-height container so the iframe auto-resize picks up the dialog height.
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div style="min-height:280px;position:relative;">`); err != nil {
								return err
							}
							if err := inner.Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: ModalTokens(),
				},
			},
		},

		// modal.ConfirmPopup
		{
			Slug:        "confirm-popup",
			Name:        "Confirm Popup",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Modals",
			Description: "A confirmation dialog popup.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live title and message controls.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Delete record?"
						}
						message := params.Get("message")
						if message == "" {
							message = "This action cannot be undone."
						}
						inner := modal.ConfirmPopupWithBoundary(title, message, "Delete", "#", "DELETE")
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div style="min-height:240px;position:relative;">`); err != nil {
								return err
							}
							if err := inner.Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: ConfirmPopupTokens(),
				},
			},
		},

		// modal.FormModal
		{
			Slug:        "form-modal-real",
			Name:        "Form Modal",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Modals",
			Description: "A modal dialog containing an HTMX form.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live title, size, and submit label controls.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Edit Record"
						}
						size := modal.ModalSize(params.Get("size"))
						submitText := params.Get("submitText")
						if submitText == "" {
							submitText = "Save"
						}
						inner := modal.FormModalWithBoundary(modal.FormModalProps{
							ID:         "gallery-form-modal",
							Title:      title,
							Size:       size,
							SubmitText: submitText,
							Action:     "#",
							Method:     "post",
						})
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div style="min-height:300px;position:relative;">`); err != nil {
								return err
							}
							if err := inner.Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div><script>document.addEventListener('DOMContentLoaded',function(){var d=document.getElementById('gallery-form-modal');if(d&&d.showModal)d.showModal();});</script>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{
							Label:      "Title",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "Edit Record",
							QueryParam: "title",
						},
						{
							Label:      "Size",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "",
							QueryParam: "size",
							Options: []galleryruntime.TokenOption{
								{Value: "modal-sm", Label: "SM"},
								{Value: "", Label: "MD"},
								{Value: "modal-lg", Label: "LG"},
								{Value: "modal-xl", Label: "XL"},
							},
						},
						{
							Label:      "Submit Label",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "Save",
							QueryParam: "submitText",
						},
					},
				},
			},
		},

		// table.TableWithProps
		{
			Slug:        "table-real",
			Name:        "Table",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "A configurable DaisyUI data table.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live zebra and pinned controls.",
					RenderFunc: func(params url.Values) templ.Component {
						zebra := params.Get("zebra") == "true"
						size := ""
						if params.Get("size") == "sm" {
							size = "sm"
						}
						rows := rawHTML(`<thead><tr><th>Name</th><th>Role</th><th>Status</th><th>Joined</th></tr></thead><tbody><tr><td>Alice Johnson</td><td>Admin</td><td><span class="badge badge-success badge-sm">Active</span></td><td>Jan 2024</td></tr><tr><td>Bob Martinez</td><td>Member</td><td><span class="badge badge-warning badge-sm">Pending</span></td><td>Mar 2024</td></tr><tr><td>Carol White</td><td>Viewer</td><td><span class="badge badge-error badge-sm">Inactive</span></td><td>Jun 2024</td></tr><tr><td>David Kim</td><td>Member</td><td><span class="badge badge-success badge-sm">Active</span></td><td>Aug 2024</td></tr></tbody>`)
						return withChildren(
							table.TableWithPropsWithBoundary(table.TableProps{
								Striped: zebra,
								Size:    size,
							}),
							rows,
						)
					},
					Tokens: []galleryruntime.DesignToken{
						{
							Label:      "Zebra Stripes",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "false",
							QueryParam: "zebra",
							Options: []galleryruntime.TokenOption{
								{Value: "false", Label: "Off"},
								{Value: "true", Label: "On"},
							},
						},
						{
							Label:      "Size",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeSelect,
							Default:    "",
							QueryParam: "size",
							Options: []galleryruntime.TokenOption{
								{Value: "", Label: "Default"},
								{Value: "sm", Label: "Small"},
							},
						},
					},
				},
			},
		},

		// logs.LogsTable
		{
			Slug:        "logs-table",
			Name:        "Logs Table",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "A workflow/event log display table.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sample log entries.",
					RenderFunc: func(params url.Values) templ.Component {
						now := time.Now()
						msg1 := params.Get("entries1")
						if msg1 == "" {
							msg1 = "Record created successfully."
						}
						msg2 := params.Get("entries2")
						if msg2 == "" {
							msg2 = "Workflow triggered."
						}
						msg3 := params.Get("entries3")
						if msg3 == "" {
							msg3 = "Rate limit approaching threshold."
						}
						msg4 := params.Get("entries4")
						if msg4 == "" {
							msg4 = "Integration sync failed."
						}
						return logs.LogsTableWithBoundary([]logs.LogEntry{
							{Type: "success", Message: msg1, CreatedAt: now.Add(-1 * time.Minute)},
							{Type: "info", Message: msg2, CreatedAt: now.Add(-3 * time.Minute)},
							{Type: "warn", Message: msg3, CreatedAt: now.Add(-10 * time.Minute)},
							{Type: "error", Message: msg4, CreatedAt: now.Add(-30 * time.Minute)},
						})
					},
					Tokens: LogsTableTokens(),
				},
			},
		},

		// layout.Navbar
		{
			Slug:        "navbar-real",
			Name:        "Navbar",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Headers",
			Description: "The application top navigation bar.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live app name control.",
					RenderFunc: func(params url.Values) templ.Component {
						appName := params.Get("appName")
						if appName == "" {
							appName = "MyApp"
						}
						return layout.NavbarWithBoundary(appName)
					},
					Tokens: NavbarTokens(),
				},
			},
		},
	}
}

// ── helpers used by new real-component entries ────────────────────────────────

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
