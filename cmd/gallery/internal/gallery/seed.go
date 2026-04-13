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

// seq renders multiple components in sequence with no wrapper element.
func seq(components ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, c := range components {
			if err := c.Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
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

func alertRenderFunc(defaultMessage string) func(params url.Values) templ.Component {
	return func(params url.Values) templ.Component {
		typ := ui.AlertType(params.Get("typ"))
		if typ == "" {
			typ = ui.AlertSuccess
		}
		icon := params.Get("icon")
		if icon == "" {
			icon = alertIconForType(typ)
		}
		message := params.Get("message")
		if message == "" {
			message = defaultMessage
		}
		return ui.AlertWithBoundary(typ, icon, message)
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
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Vertical timeline with done and pending items.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.TimelineWithBoundary([]ui.TimelineItemProps{
							{Date: "Jan 2024", Label: "Project started", Done: true},
							{Date: "Mar 2024", Label: "Beta release", Done: true},
							{Date: "Jun 2024", Label: "v1.0 launch", Done: false},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Data Display / Chat ───────────────────────────────────────────────────
		{
			Slug:        "chat-bubble",
			Name:        "Chat Bubble",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Chat conversation bubbles (sent and received).",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single chat bubble with live sent/author/message controls.",
					RenderFunc: func(params url.Values) templ.Component {
						sent := params.Get("sent") == "true"
						author := params.Get("author")
						if author == "" {
							author = "Alice"
						}
						timestamp := params.Get("timestamp")
						if timestamp == "" {
							timestamp = "10:32 AM"
						}
						message := params.Get("message")
						if message == "" {
							message = "Hey! How are you doing?"
						}
						return ui.ChatBubbleWithBoundary(sent, author, timestamp, "", message)
					},
					Tokens: ChatBubbleTokens(),
				},
				{
					Name:        "Examples",
					Description: "Sent and received bubbles together.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-col gap-2 p-4 max-w-sm mx-auto">`); err != nil {
								return err
							}
							if err := withText("Hey! How are you doing?", ui.ChatBubble(false, "Alice", "10:32 AM", "chat-bubble-primary")).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("Good thanks! How about you?", ui.ChatBubble(true, "You", "10:33 AM", "")).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Data Display / Mockups ────────────────────────────────────────────────
		{
			Slug:        "mockup-code",
			Name:        "Mockup Code",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Code block mockup with terminal-style prefix lines.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Terminal-style code block with prefix lines.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupCodeWithBoundary([]ui.MockupCodeLineProps{
							{Prefix: "$", Code: "go get github.com/emergent-company/go-daisy"},
							{Prefix: "$", Code: "task build:ui"},
							{Prefix: ">", Code: "Done in 1.2s", ColorClass: "text-success"},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "mockup-browser",
			Name:        "Mockup Browser",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Browser window mockup for UI showcasing.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Browser mockup with configurable URL.",
					RenderFunc: func(params url.Values) templ.Component {
						u := params.Get("url")
						if u == "" {
							u = "https://go-daisy.dev"
						}
						return ui.MockupBrowserWithBoundary(u)
					},
					Tokens: MockupBrowserTokens(),
				},
			},
		},

		// ── Feedback / Alerts ─────────────────────────────────────────────────────
		{
			Slug:        "alert",
			Name:        "Alert",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Alerts",
			Description: "Contextual feedback alert with configurable type, optional icon, and message.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single alert with live type and message controls.",
					RenderFunc:  alertRenderFunc("Your changes have been saved successfully."),
					Tokens:      AlertTokens(),
				},
				{
					Name:        "Examples",
					Description: "All alert types shown together.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-col gap-3 p-6">`); err != nil {
								return err
							}
							alerts := []templ.Component{
								ui.Alert(ui.AlertSuccess, "lucide--circle-check", "Your changes have been saved successfully."),
								ui.Alert(ui.AlertError, "lucide--circle-x", "Something went wrong. Please try again."),
								ui.Alert(ui.AlertWarning, "lucide--triangle-alert", "Your session will expire in 5 minutes."),
								ui.Alert(ui.AlertInfo, "lucide--info", "A new software update is available."),
							}
							for _, a := range alerts {
								if err := a.Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
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
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Filter tabs with configurable selected tab.",
					RenderFunc: func(params url.Values) templ.Component {
						selected := params.Get("selected")
						if selected == "" {
							selected = "All"
						}
						return ui.FilterTabsWithBoundary("filter", selected, []string{"All", "Active", "Pending", "Closed"})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Selected", Group: "State", Type: galleryruntime.TokenTypeSelect, Default: "All", QueryParam: "selected", Options: []galleryruntime.TokenOption{
							{Value: "All", Label: "All"},
							{Value: "Active", Label: "Active"},
							{Value: "Pending", Label: "Pending"},
							{Value: "Closed", Label: "Closed"},
						}},
					},
				},
			},
		},

		// ── Forms ─────────────────────────────────────────────────────────────────
		{
			Slug:        "form-checkbox",
			Name:        "Checkboxes & Toggles",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Toggles",
			Description: "Checkbox and toggle switch inputs.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Checkbox",
					Description: "Labeled checkbox input.",
					RenderFunc: func(params url.Values) templ.Component {
						checked := params.Get("checked") == "true"
						label := params.Get("label")
						if label == "" {
							label = "Receive email notifications"
						}
						return form.CheckboxWithBoundary("notifications", checked, label)
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Toggle",
					Description: "Toggle switch input.",
					RenderFunc: func(params url.Values) templ.Component {
						checked := params.Get("checked") == "true"
						label := params.Get("label")
						if label == "" {
							label = "Dark mode"
						}
						return form.ToggleWithBoundary("dark-mode", checked, label)
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Checkboxes and toggles together.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex flex-col gap-4 max-w-sm mx-auto">`); err != nil {
								return err
							}
							comps := []templ.Component{
								form.CheckboxWithBoundary("n1", true, "Receive email notifications"),
								form.CheckboxWithBoundary("n2", false, "Subscribe to newsletter"),
								form.ToggleWithBoundary("dark", true, "Dark mode"),
								form.ToggleWithBoundary("autosave", false, "Auto-save"),
							}
							for _, c := range comps {
								if err := c.Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "form-radio",
			Name:        "Radio Buttons",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Toggles",
			Description: "Radio button group for single-selection.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Radio group with configurable color.",
					RenderFunc: func(params url.Values) templ.Component {
						color := params.Get("color")
						if color == "" {
							color = "radio-primary"
						}
						return form.RadioGroupWithBoundary("plan", "free", [][2]string{
							{"free", "Free – $0/mo"},
							{"pro", "Pro – $12/mo"},
							{"enterprise", "Enterprise – Custom"},
						}, color)
					},
					Tokens: RadioGroupTokens(),
				},
			},
		},
		{
			Slug:        "form-rating",
			Name:        "Rating",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Star and heart rating inputs using DaisyUI rating + mask classes.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Rating input with configurable shape, value, max, and color.",
					RenderFunc: func(params url.Values) templ.Component {
						shape := form.RatingShape(params.Get("shape"))
						if shape == "" {
							shape = form.RatingStar
						}
						value := 3
						if v, err := parseInt(params.Get("value")); err == nil && v > 0 {
							value = v
						}
						max := 5
						if v, err := parseInt(params.Get("max")); err == nil && v > 0 {
							max = v
						}
						color := params.Get("color")
						if color == "" {
							color = "bg-orange-400"
						}
						return form.RatingWithBoundary("rating-demo", value, max, shape, color, "")
					},
					Tokens: RatingTokens(),
				},
			},
		},

		// ── Foundation / Display ──────────────────────────────────────────────────
		{
			Slug:        "divider",
			Name:        "Divider",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Horizontal and vertical dividers with optional label.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Divider with configurable color, orientation, and label.",
					RenderFunc: func(params url.Values) templ.Component {
						color := ui.DividerColor(params.Get("color"))
						vertical := params.Get("vertical") == "true"
						label := params.Get("label")
						return ui.DividerWithBoundary(color, vertical, label)
					},
					Tokens: DividerTokens(),
				},
				{
					Name:        "Examples",
					Description: "Horizontal and vertical dividers.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex flex-col gap-4 max-w-sm mx-auto">`); err != nil {
								return err
							}
							if err := withText("OR", ui.Divider(ui.DividerDefault, false)).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("Primary", ui.Divider(ui.DividerPrimary, false)).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div class="flex h-20 items-center gap-4"><span class="text-sm">Left</span>`); err != nil {
								return err
							}
							if err := ui.Divider(ui.DividerDefault, true).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `<span class="text-sm">Right</span></div></div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "kbd",
			Name:        "Keyboard Keys",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Keyboard shortcut display using DaisyUI kbd.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single key with configurable size and label.",
					RenderFunc: func(params url.Values) templ.Component {
						size := ui.KbdSize(params.Get("size"))
						key := params.Get("key")
						if key == "" {
							key = "⌘K"
						}
						return ui.KbdWithBoundary(size, key)
					},
					Tokens: KbdTokens(),
				},
				{
					Name:        "Examples",
					Description: "Various keyboard shortcut combinations.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-wrap gap-4 p-6 items-center justify-center">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div class="flex items-center gap-1 text-sm">Press `); err != nil {
								return err
							}
							if err := withText("⌘", ui.Kbd(ui.KbdSM)).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("K", ui.Kbd(ui.KbdSM)).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, ` to search</div>`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div class="flex items-center gap-1">`); err != nil {
								return err
							}
							if err := withText("Ctrl", ui.Kbd(ui.KbdSM)).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<span class="text-sm">+</span>`); err != nil {
								return err
							}
							if err := withText("S", ui.Kbd(ui.KbdSM)).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div>`); err != nil {
								return err
							}
							if err := withText("Enter", ui.Kbd(ui.KbdLG)).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("Esc", ui.Kbd(ui.KbdXS)).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
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
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "4-step progress indicator.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.StepsWithBoundary([]ui.StepProps{
							{Label: "Register", Done: true},
							{Label: "Choose plan", Done: true},
							{Label: "Payment", Done: false},
							{Label: "Confirm", Done: false},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "collapse",
			Name:        "Collapse / Accordion",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Collapsible sections using DaisyUI collapse.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Accordion with multiple collapsible items.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.AccordionWithBoundary([]ui.AccordionItemProps{
							{Title: "What is go-daisy?", Content: templ.Raw("go-daisy is a Go UI component library for HTMX-driven web interfaces built with DaisyUI."), Open: true},
							{Title: "How do I install it?", Content: templ.Raw("<code>go get github.com/emergent-company/go-daisy</code>"), Open: false},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "dropdown",
			Name:        "Dropdown",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Dropdown menu triggered by a button.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Dropdown menu with configurable alignment.",
					RenderFunc: func(params url.Values) templ.Component {
						align := ui.DropdownAlign(params.Get("align"))
						return ui.DropdownWithBoundary(align, ui.DropdownTrigger("Options", "btn-primary"), []ui.DropdownItemProps{
							{Label: "Profile"},
							{Label: "Settings"},
							{Label: "Help"},
							{Divider: true},
							{Label: "Sign out", Danger: true},
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Alignment", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "", QueryParam: "align", Options: []galleryruntime.TokenOption{
							{Value: "", Label: "Default (bottom)"},
							{Value: "dropdown-end", Label: "End"},
							{Value: "dropdown-top", Label: "Top"},
							{Value: "dropdown-bottom", Label: "Bottom"},
							{Value: "dropdown-left", Label: "Left"},
							{Value: "dropdown-right", Label: "Right"},
						}},
					},
				},
			},
		},
		{
			Slug:        "tooltip",
			Name:        "Tooltip",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Tooltip on hover in top, bottom, left, right positions.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Tooltip with configurable tip text and position.",
					RenderFunc: func(params url.Values) templ.Component {
						tip := params.Get("tip")
						if tip == "" {
							tip = "Helpful hint"
						}
						position := params.Get("position")
						return ui.TooltipWithBoundary(tip, position, ui.SimpleButton("Hover me", ui.ButtonPrimary, ui.ButtonSM))
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Tip", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Helpful hint", QueryParam: "tip"},
						{Label: "Position", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "", QueryParam: "position", Options: []galleryruntime.TokenOption{
							{Value: "", Label: "Default (top)"},
							{Value: "top", Label: "Top"},
							{Value: "bottom", Label: "Bottom"},
							{Value: "left", Label: "Left"},
							{Value: "right", Label: "Right"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Tooltips in all four positions.",
					RenderFunc: func(_ url.Values) templ.Component {
						return row(
							ui.TooltipWithBoundary("Default tooltip", "", ui.SimpleButton("Hover me", "", ui.ButtonSM)),
							ui.TooltipWithBoundary("Top", "top", ui.SimpleButton("Top", ui.ButtonPrimary, ui.ButtonSM)),
							ui.TooltipWithBoundary("Bottom", "bottom", ui.SimpleButton("Bottom", ui.ButtonSecondary, ui.ButtonSM)),
							ui.TooltipWithBoundary("Left", "left", ui.SimpleButton("Left", "btn-accent", ui.ButtonSM)),
							ui.TooltipWithBoundary("Right", "right", ui.SimpleButton("Right", ui.ButtonNeutral, ui.ButtonSM)),
						)
					},
				},
			},
		},
		{
			Slug:        "swap",
			Name:        "Swap",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Toggle between two visual states on click.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Swap toggle with configurable rotate animation.",
					RenderFunc: func(params url.Values) templ.Component {
						rotate := params.Get("rotate") == "true"
						return ui.SwapWithBoundary(rotate,
							ui.IconSpan("lucide--sun", "size-8"),
							ui.IconSpan("lucide--moon", "size-8"),
						)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Rotate", Group: "Animation", Type: galleryruntime.TokenTypeSelect, Default: "false", QueryParam: "rotate", Options: []galleryruntime.TokenOption{
							{Value: "false", Label: "Fade"},
							{Value: "true", Label: "Rotate"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Rotate icon swap and text button swap.",
					RenderFunc: func(_ url.Values) templ.Component {
						return row(
							ui.SwapWithBoundary(true,
								ui.IconSpan("lucide--sun", "size-8"),
								ui.IconSpan("lucide--moon", "size-8"),
							),
							ui.SwapWithBoundary(false,
								ui.SimpleButton("ON", "btn-success", ui.ButtonSM),
								ui.SimpleButton("OFF", ui.ButtonGhost, ui.ButtonSM),
							),
						)
					},
				},
			},
		},
		{
			Slug:        "hero",
			Name:        "Hero",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Layout",
			Description: "Full-width hero section with headline and CTA button.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Hero with configurable title, subtitle, and CTA label.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "go-daisy"
						}
						subtitle := params.Get("subtitle")
						if subtitle == "" {
							subtitle = "Type-safe Templ components styled with DaisyUI for HTMX apps."
						}
						ctaLabel := params.Get("ctaLabel")
						if ctaLabel == "" {
							ctaLabel = "Get Started"
						}
						return ui.HeroWithBoundary("min-h-56", title, subtitle, ctaLabel)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Title", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "go-daisy", QueryParam: "title"},
						{Label: "Subtitle", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Type-safe Templ components styled with DaisyUI for HTMX apps.", QueryParam: "subtitle"},
						{Label: "CTA Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Get Started", QueryParam: "ctaLabel"},
					},
				},
			},
		},

		// ── Data Display / List ───────────────────────────────────────────────────
		{
			Slug:        "list-basic",
			Name:        "List",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Lists",
			Description: "DaisyUI list component for vertical item groups.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Vertical list with labelled rows and status badges.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.ListWithBoundary([]ui.ListRowProps{
							{Title: "Alice Johnson", Subtitle: "alice@example.com", Badge: ui.Badge(ui.BadgeProps{Label: "Active", Variant: ui.BadgeSuccess, Style: ui.BadgeStyleSoft, Size: ui.BadgeSizeSM})},
							{Title: "Bob Smith", Subtitle: "bob@example.com", Badge: ui.Badge(ui.BadgeProps{Label: "Inactive", Variant: ui.BadgeGhost, Style: ui.BadgeStyleSoft, Size: ui.BadgeSizeSM})},
							{Title: "Carol White", Subtitle: "carol@example.com", Badge: ui.Badge(ui.BadgeProps{Label: "Pending", Variant: ui.BadgeWarning, Style: ui.BadgeStyleSoft, Size: ui.BadgeSizeSM})},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Data Display / Indicator ──────────────────────────────────────────────
		{
			Slug:        "indicator",
			Name:        "Indicator",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Badge/dot overlay indicators on components.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single indicator with configurable badge color and content.",
					RenderFunc: func(params url.Values) templ.Component {
						badgeClass := params.Get("badge")
						if badgeClass == "" {
							badgeClass = "badge-error"
						}
						count := params.Get("count")
						if count == "" {
							count = "3"
						}
						return ui.IndicatorWithBoundary(
							"badge badge-sm "+badgeClass,
							templ.Raw(count),
							ui.Button("", ui.ButtonOutline, ui.ButtonMD, ui.ButtonTypeButton, ui.ButtonShapeDefault, "lucide--bell size-4", false),
						)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Badge color", Group: "Badge", Type: galleryruntime.TokenTypeSelect, Default: "badge-error", QueryParam: "badge", Options: []galleryruntime.TokenOption{
							{Value: "badge-error", Label: "Error"},
							{Value: "badge-primary", Label: "Primary"},
							{Value: "badge-success", Label: "Success"},
							{Value: "badge-warning", Label: "Warning"},
							{Value: "badge-neutral", Label: "Neutral"},
						}},
						{Label: "Count", Group: "Badge", Type: galleryruntime.TokenTypeText, Default: "3", QueryParam: "count"},
					},
				},
				{
					Name:        "Examples",
					Description: "Badge/dot indicators on button, avatar, and card.",
					RenderFunc: func(_ url.Values) templ.Component {
						return row(
							ui.IndicatorWithBoundary("badge badge-error badge-sm",
								templ.Raw("3"),
								ui.Button("", ui.ButtonGhost, ui.ButtonSM, ui.ButtonTypeButton, ui.ButtonShapeSquare, "lucide--bell size-5", false),
							),
							ui.IndicatorWithBoundary("badge badge-primary badge-xs",
								templ.NopComponent,
								ui.Avatar("AJ", "", "", ui.AvatarMD),
							),
							ui.IndicatorWithBoundary("badge badge-success badge-sm",
								templ.Raw("New"),
								ui.CardPlaceholder("Card"),
							),
						)
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Data Display / Stack ──────────────────────────────────────────────────
		{
			Slug:        "stack",
			Name:        "Stack",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Overlapping stacked card effect.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Three cards stacked with depth effect.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.StackWithBoundary(
							ui.StackCard("Card 1", "bg-primary text-primary-content shadow-lg"),
							ui.StackCard("Card 2", "bg-secondary text-secondary-content shadow"),
							ui.StackCard("Card 3", "bg-accent text-accent-content"),
						)
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Data Display / Diff ───────────────────────────────────────────────────
		{
			Slug:        "diff",
			Name:        "Diff",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Side-by-side comparison diff panel.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Diff panel with configurable before and after content.",
					RenderFunc: func(params url.Values) templ.Component {
						before := params.Get("before")
						if before == "" {
							before = "Before: Old content here"
						}
						after := params.Get("after")
						if after == "" {
							after = "After: New content here"
						}
						return ui.DiffWithBoundary(before, after)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Before", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Before: Old content here", QueryParam: "before"},
						{Label: "After", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "After: New content here", QueryParam: "after"},
					},
				},
			},
		},

		// ── Data Display / Mask ───────────────────────────────────────────────────
		{
			Slug:        "mask",
			Name:        "Mask",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "CSS mask shapes applied to images and elements.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Mask with configurable shape.",
					RenderFunc: func(params url.Values) templ.Component {
						shape := ui.MaskShape(params.Get("shape"))
						if shape == "" {
							shape = ui.MaskSquircle
						}
						return ui.MaskWithBoundary(shape, ui.MaskSwatch("S", "bg-primary text-primary-content"))
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Shape", Group: "Appearance", Type: galleryruntime.TokenTypeSelect, Default: "mask-squircle", QueryParam: "shape", Options: []galleryruntime.TokenOption{
							{Value: "mask-squircle", Label: "Squircle"},
							{Value: "mask-heart", Label: "Heart"},
							{Value: "mask-hexagon", Label: "Hexagon"},
							{Value: "mask-triangle", Label: "Triangle"},
							{Value: "mask-circle", Label: "Circle"},
							{Value: "mask-star", Label: "Star"},
							{Value: "mask-star-2", Label: "Star 2"},
							{Value: "mask-pentagon", Label: "Pentagon"},
							{Value: "mask-diamond", Label: "Diamond"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "All mask shapes.",
					RenderFunc: func(_ url.Values) templ.Component {
						type maskEx struct {
							shape ui.MaskShape
							bg    string
							label string
						}
						examples := []maskEx{
							{ui.MaskSquircle, "bg-primary text-primary-content", "S"},
							{ui.MaskHeart, "bg-error text-error-content", "♥"},
							{ui.MaskHexagon, "bg-secondary text-secondary-content", "H"},
							{ui.MaskTriangle, "bg-accent text-accent-content", "▲"},
							{ui.MaskCircle, "bg-success text-success-content", "●"},
							{ui.MaskStar2, "bg-warning text-warning-content", "★"},
						}
						comps := make([]templ.Component, len(examples))
						for i, ex := range examples {
							e := ex
							comps[i] = ui.MaskWithBoundary(e.shape, ui.MaskSwatch(e.label, e.bg))
						}
						return row(comps...)
					},
				},
			},
		},

		// ── Data Display / Carousel ───────────────────────────────────────────────
		{
			Slug:        "carousel",
			Name:        "Carousel",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Horizontal scrolling carousel with snap items.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Three-slide horizontal carousel.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.CarouselWithBoundary(false, []ui.CarouselItemProps{
							{ID: "slide1", Content: ui.CarouselSlide("Slide 1", "bg-primary text-primary-content")},
							{ID: "slide2", Content: ui.CarouselSlide("Slide 2", "bg-secondary text-secondary-content")},
							{ID: "slide3", Content: ui.CarouselSlide("Slide 3", "bg-accent text-accent-content")},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Data Display / Countdown ──────────────────────────────────────────────
		{
			Slug:        "countdown",
			Name:        "Countdown",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Animated countdown timer display.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Countdown with configurable days, hours, minutes, and seconds.",
					RenderFunc: func(params url.Values) templ.Component {
						days := 2
						if v, err := parseInt(params.Get("days")); err == nil {
							days = v
						}
						hours := 10
						if v, err := parseInt(params.Get("hours")); err == nil {
							hours = v
						}
						minutes := 24
						if v, err := parseInt(params.Get("minutes")); err == nil {
							minutes = v
						}
						seconds := 45
						if v, err := parseInt(params.Get("seconds")); err == nil {
							seconds = v
						}
						return ui.CountdownWithBoundary(days, hours, minutes, seconds)
					},
					Tokens: CountdownTokens(),
				},
			},
		},

		// ── Data Display / Mockup Phone & Window ──────────────────────────────────
		{
			Slug:        "mockup-phone",
			Name:        "Mockup Phone",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Phone frame mockup for mobile UI display.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Phone frame with an app screen placeholder.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupPhoneWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "mockup-window",
			Name:        "Mockup Window",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Mockups",
			Description: "Desktop window frame mockup.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Desktop window frame with content placeholder.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupWindowWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Feedback / Status ─────────────────────────────────────────────────────
		{
			Slug:        "status-dots",
			Name:        "Status Dots",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "States",
			Description: "Small colored status indicator dots.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single status dot with configurable color and animation.",
					RenderFunc: func(params url.Values) templ.Component {
						color := ui.StatusColor(params.Get("color"))
						if color == "" {
							color = ui.StatusSuccess
						}
						animate := params.Get("animate") == "true"
						return ui.StatusDotWithBoundary(color, animate)
					},
					Tokens: StatusDotTokens(),
				},
				{
					Name:        "Examples",
					Description: "All status colors.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-wrap gap-6 p-6 items-center justify-center">`); err != nil {
								return err
							}
							items := []struct {
								color ui.StatusColor
								label string
							}{
								{ui.StatusSuccess, "Online"},
								{ui.StatusError, "Offline"},
								{ui.StatusWarning, "Away"},
								{ui.StatusInfo, "Busy"},
								{ui.StatusNeutral, "Unknown"},
							}
							for _, item := range items {
								if _, err := io.WriteString(w, `<div class="flex items-center gap-2 text-sm">`); err != nil {
									return err
								}
								if err := ui.StatusDot(item.color, false).Render(ctx, w); err != nil {
									return err
								}
								if _, err := io.WriteString(w, " "+item.label+`</div>`); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Overlays / Dropdown positions ─────────────────────────────────────────
		{
			Slug:        "dropdown-positions",
			Name:        "Dropdown Positions",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Dropdowns",
			Description: "Dropdown menus opening in different directions.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single dropdown with configurable open direction.",
					RenderFunc: func(params url.Values) templ.Component {
						align := ui.DropdownAlign(params.Get("align"))
						if align == "" {
							align = ui.DropdownBottom
						}
						label := "Open ▼"
						if align == ui.DropdownTop {
							label = "Open ▲"
						} else if align == ui.DropdownEnd {
							label = "Options ⋮"
						}
						return ui.DropdownWithBoundary(align, ui.DropdownTrigger(label, "btn-primary"), []ui.DropdownItemProps{
							{Label: "Edit"},
							{Label: "Duplicate"},
							{Divider: true},
							{Label: "Delete", Danger: true},
						})
					},
					FrameHeight: "220px",
					Tokens: []galleryruntime.DesignToken{
						{Label: "Position", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "", QueryParam: "align", Options: []galleryruntime.TokenOption{
							{Value: "", Label: "Bottom"},
							{Value: string(ui.DropdownTop), Label: "Top"},
							{Value: string(ui.DropdownEnd), Label: "End"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Bottom, top, and end-aligned dropdowns.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []ui.DropdownItemProps{
							{Label: "Item 1"},
							{Label: "Item 2"},
							{Label: "Item 3"},
						}
						return row(
							ui.DropdownWithBoundary("", ui.DropdownTrigger("Bottom ▼", "btn-outline"), items),
							ui.DropdownWithBoundary(ui.DropdownTop, ui.DropdownTrigger("Top ▲", "btn-outline"), items),
							ui.DropdownWithBoundary(ui.DropdownEnd, ui.DropdownTrigger("Options ⋮", "btn-primary"), []ui.DropdownItemProps{
								{Label: "Edit"},
								{Label: "Duplicate"},
								{Divider: true},
								{Label: "Delete", Danger: true},
							}),
						)
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Navigation / Breadcrumbs, Navbar, Menu, Dock ──────────────────────────
		{
			Slug:        "breadcrumbs",
			Name:        "Breadcrumbs",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Headers",
			Description: "Navigation breadcrumb trail with configurable items.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Breadcrumb trail with configurable item labels.",
					RenderFunc: func(params url.Values) templ.Component {
						item1 := params.Get("items1")
						if item1 == "" {
							item1 = "Home"
						}
						item2 := params.Get("items2")
						if item2 == "" {
							item2 = "Documents"
						}
						item3 := params.Get("items3")
						if item3 == "" {
							item3 = "Add Document"
						}
						return nav.BreadcrumbsWithBoundary([]nav.BreadcrumbItem{
							{Label: item1, Href: "#"},
							{Label: item2, Href: "#"},
							{Label: item3},
						})
					},
					Tokens: BreadcrumbsTokens(),
				},
			},
		},
		{
			Slug:        "dock-nav",
			Name:        "Dock",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Bottom dock navigation bar for mobile-style UIs.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Bottom navigation dock with configurable item labels.",
					RenderFunc: func(params url.Values) templ.Component {
						item1 := params.Get("items1")
						if item1 == "" {
							item1 = "Home"
						}
						item2 := params.Get("items2")
						if item2 == "" {
							item2 = "Search"
						}
						item3 := params.Get("items3")
						if item3 == "" {
							item3 = "Alerts"
						}
						item4 := params.Get("items4")
						if item4 == "" {
							item4 = "Profile"
						}
						return nav.DockWithBoundary([]nav.DockItem{
							{Label: item1, Icon: "lucide--home", Active: true},
							{Label: item2, Icon: "lucide--search"},
							{Label: item3, Icon: "lucide--bell"},
							{Label: item4, Icon: "lucide--user"},
						})
					},
					Tokens: DockTokens(),
				},
			},
		},

		// ── Forms / File Input ────────────────────────────────────────────────────
		{
			Slug:        "form-file",
			Name:        "File Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "File upload input field with label and accept filter.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "File input with configurable label and accept filter.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Upload file"
						}
						accept := params.Get("accept")
						return form.FileInputWithBoundary("upload", label, accept)
					},
					Tokens: FileInputTokens(),
				},
			},
		},

		// ── Foundation / Join, Link ───────────────────────────────────────────────
		{
			Slug:        "join",
			Name:        "Join",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Layout",
			Description: "Join fuses children into a single rounded group.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Input + button join with configurable orientation.",
					RenderFunc: func(params url.Values) templ.Component {
						vertical := params.Get("orientation") == "vertical"
						if vertical {
							return ui.JoinWithBoundary(true,
								ui.JoinButton("Top", ui.ButtonOutline, false),
								ui.JoinButton("Middle", ui.ButtonOutline, false),
								ui.JoinButton("Bottom", ui.ButtonOutline, false),
							)
						}
						return ui.JoinWithBoundary(false,
							ui.JoinInputPlaceholder("Search…"),
							ui.JoinButton("Go", ui.ButtonPrimary, false),
						)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Orientation", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "horizontal", QueryParam: "orientation", Options: []galleryruntime.TokenOption{
							{Value: "horizontal", Label: "Horizontal"},
							{Value: "vertical", Label: "Vertical"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Input+button, button group, and vertical join.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-col gap-4 p-6 items-center">`); err != nil {
								return err
							}
							joins := []templ.Component{
								ui.JoinWithBoundary(false,
									ui.JoinInputPlaceholder("Search…"),
									ui.JoinButton("Go", ui.ButtonPrimary, false),
								),
								ui.JoinWithBoundary(false,
									ui.JoinButton("A", ui.ButtonOutline, false),
									ui.JoinButton("B", ui.ButtonOutline, true),
									ui.JoinButton("C", ui.ButtonOutline, false),
								),
								ui.JoinWithBoundary(true,
									ui.JoinButton("Top", ui.ButtonOutline, false),
									ui.JoinButton("Middle", ui.ButtonOutline, false),
									ui.JoinButton("Bottom", ui.ButtonOutline, false),
								),
							}
							for _, c := range joins {
								if err := c.Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "link-styles",
			Name:        "Links",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "DaisyUI link styles with color variants and hover.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single link with configurable variant.",
					RenderFunc: func(params url.Values) templ.Component {
						variant := nav.LinkVariant(params.Get("variant"))
						if variant == "" {
							variant = nav.LinkDefault
						}
						return nav.LinkWithBoundary("#", variant, "Click here")
					},
					Tokens: LinkTokens(),
				},
				{
					Name:        "Examples",
					Description: "All link style variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-wrap gap-4 p-6 items-center justify-center text-sm">`); err != nil {
								return err
							}
							links := []struct {
								variant nav.LinkVariant
								label   string
							}{
								{nav.LinkDefault, "Default link"},
								{nav.LinkPrimary, "Primary"},
								{nav.LinkSecondary, "Secondary"},
								{nav.LinkAccent, "Accent"},
								{nav.LinkNeutral, "Neutral"},
								{nav.LinkHover, "Hover only"},
							}
							for _, l := range links {
								if err := withText(l.label, nav.Link("#", l.variant)).Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		{
			Slug:        "tag",
			Name:        "Tag",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Tag",
			Description: "Tag renders a removable chip badge used in multi-select fields. Clicking the × removes the tag.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single tag with configurable label and remove link.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Contract Law"
						}
						return ui.TagWithBoundary(label, "#")
					},
					Tokens: TagTokens(),
				},
				{
					Name:        "Examples",
					Description: "Multiple removable and read-only tags.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-4"><div class="flex flex-wrap gap-2">`); err != nil {
								return err
							}
							for _, label := range []string{"Contract Law", "Family Law", "Civil Litigation"} {
								if err := ui.Tag(label, "#").Render(ctx, w); err != nil {
									return err
								}
							}
							if _, err := io.WriteString(w, `</div><p class="text-xs text-base-content/50">Read-only (no remove link):</p><div class="flex flex-wrap gap-2">`); err != nil {
								return err
							}
							for _, label := range []string{"Contract Law", "Family Law", "Civil Litigation"} {
								if err := ui.Tag(label, "").Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "company-avatar",
			Name:        "Company Avatar",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Avatars",
			Description: "Circular avatar with a building icon placeholder for companies. Same sizes as Avatar. Use alongside a company name in tables and detail views.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single company avatar with configurable size.",
					RenderFunc: func(params url.Values) templ.Component {
						sizeStr := params.Get("size")
						size := ui.AvatarMD
						switch sizeStr {
						case "xs":
							size = ui.AvatarXS
						case "sm":
							size = ui.AvatarSM
						case "lg":
							size = ui.AvatarLG
						}
						return ui.AvatarWithBoundary("", "", "lucide--building-2", size)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Size", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "md", QueryParam: "size", Options: []galleryruntime.TokenOption{
							{Value: "xs", Label: "XS"},
							{Value: "sm", Label: "SM"},
							{Value: "md", Label: "MD"},
							{Value: "lg", Label: "LG"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Company avatar in all sizes plus an inline with-name usage.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex flex-wrap gap-6 items-end">`); err != nil {
								return err
							}
							sizes := []struct {
								size  ui.AvatarSize
								label string
							}{
								{ui.AvatarXS, "xs"},
								{ui.AvatarSM, "sm"},
								{ui.AvatarMD, "md"},
								{ui.AvatarLG, "lg"},
							}
							for _, s := range sizes {
								if _, err := io.WriteString(w, `<div class="flex flex-col items-center gap-2">`); err != nil {
									return err
								}
								if err := ui.Avatar("", "", "lucide--building-2", s.size).Render(ctx, w); err != nil {
									return err
								}
								if _, err := io.WriteString(w, `<span class="text-xs text-base-content/60">`+s.label+`</span></div>`); err != nil {
									return err
								}
							}
							// inline with-name example
							if _, err := io.WriteString(w, `<div class="flex flex-col items-center gap-2"><div class="flex items-center gap-2">`); err != nil {
								return err
							}
							if err := ui.Avatar("", "", "lucide--building-2", ui.AvatarXS).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<span class="text-sm font-medium">Acme Corp</span></div><span class="text-xs text-base-content/60">with name</span></div>`); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "person-avatar",
			Name:        "Person Avatar",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Avatars",
			Description: "Inline avatar + name chip with a hover card that reveals contact details. Pure CSS — no JS required.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single person chip with configurable name, role, and status.",
					RenderFunc: func(params url.Values) templ.Component {
						name := params.Get("name")
						if name == "" {
							name = "Jane Doe"
						}
						role := params.Get("role")
						if role == "" {
							role = "Senior Attorney"
						}
						status := params.Get("status")
						badgeLabel, badgeClass := "Active", "badge-success badge-soft"
						switch status {
						case "leave":
							badgeLabel, badgeClass = "On leave", "badge-warning badge-soft"
						case "closed":
							badgeLabel, badgeClass = "Closed", "badge-neutral badge-soft"
						}
						return ui.PersonChip(name, "bg-primary", "text-primary-content", "from-primary/20", "to-primary/5", ui.PersonChipContact{
							Role:        role,
							BadgeLabel:  badgeLabel,
							BadgeClass:  badgeClass,
							ProfileHref: "#",
						})
					},
					FrameHeight: "180px",
					Tokens: []galleryruntime.DesignToken{
						{Label: "Name", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Jane Doe", QueryParam: "name"},
						{Label: "Role", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Senior Attorney", QueryParam: "role"},
						{Label: "Status", Group: "Content", Type: galleryruntime.TokenTypeSelect, Default: "active", QueryParam: "status", Options: []galleryruntime.TokenOption{
							{Value: "active", Label: "Active"},
							{Value: "leave", Label: "On leave"},
							{Value: "closed", Label: "Closed"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Two person chips side-by-side: Jane Doe (primary, active) and Bob Smith (secondary, on leave).",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-8 space-y-10"><div><p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-4">Inline — initials avatar</p><div class="flex flex-wrap gap-6 items-start">`); err != nil {
								return err
							}
							if err := ui.PersonChip("Jane Doe", "bg-primary", "text-primary-content", "from-primary/20", "to-primary/5", ui.PersonChipContact{
								Email:       "jane.doe@example.com",
								Role:        "Senior Attorney",
								BadgeLabel:  "Active",
								BadgeClass:  "badge-success badge-soft",
								ProfileHref: "#",
							}).Render(ctx, w); err != nil {
								return err
							}
							if err := ui.PersonChip("Bob Smith", "bg-secondary", "text-secondary-content", "from-secondary/20", "to-secondary/5", ui.PersonChipContact{
								Email:       "bob.smith@example.com",
								Role:        "Paralegal",
								BadgeLabel:  "On leave",
								BadgeClass:  "badge-warning badge-soft",
								ProfileHref: "#",
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div></div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},

		// ── Data Display extras ────────────────────────────────────────────────
		{
			Slug:        "table-with-actions",
			Name:        "With Actions",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Full-featured table with sortable headers, status badges, avatars, and an action menu (ellipsis dropdown) per row.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Configurable row count and pagination state.",
					RenderFunc: func(params url.Values) templ.Component {
						page := 1
						if v, err := parseInt(params.Get("page")); err == nil && v > 0 {
							page = v
						}
						allRows := []table.TableWithActionsRow{
							{Name: "Alice Johnson", Status: "active", Role: "Admin", Joined: "2024-01-15"},
							{Name: "Bob Smith", Status: "pending", Role: "Employee", Joined: "2024-03-02"},
							{Name: "Carol White", Status: "closed", Role: "Employee", Joined: "2023-11-20"},
							{Name: "David Kim", Status: "active", Role: "Viewer", Joined: "2024-06-10"},
							{Name: "Eve Martinez", Status: "pending", Role: "Employee", Joined: "2024-08-22"},
						}
						rowCount := 3
						if v, err := parseInt(params.Get("rows")); err == nil && v >= 1 && v <= 5 {
							rowCount = v
						}
						rows := allRows[:rowCount]
						return table.TableWithActions(table.TableWithActionsProps{
							Rows:        rows,
							TotalCount:  47,
							CurrentPage: page,
							TotalPages:  3,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Visible rows", Group: "Data", Type: galleryruntime.TokenTypeSelect, Default: "3", QueryParam: "rows", Options: []galleryruntime.TokenOption{
							{Value: "1", Label: "1"},
							{Value: "2", Label: "2"},
							{Value: "3", Label: "3"},
							{Value: "4", Label: "4"},
							{Value: "5", Label: "5"},
						}},
						{Label: "Current page", Group: "Pagination", Type: galleryruntime.TokenTypeSelect, Default: "1", QueryParam: "page", Options: []galleryruntime.TokenOption{
							{Value: "1", Label: "Page 1"},
							{Value: "2", Label: "Page 2"},
							{Value: "3", Label: "Page 3"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Three rows with avatar, status badge, role, and ellipsis action dropdown.",
					RenderFunc: func(_ url.Values) templ.Component {
						return table.TableWithActions(table.TableWithActionsProps{
							Rows: []table.TableWithActionsRow{
								{Name: "Alice Johnson", Status: "active", Role: "Admin", Joined: "2024-01-15"},
								{Name: "Bob Smith", Status: "pending", Role: "Employee", Joined: "2024-03-02"},
								{Name: "Carol White", Status: "closed", Role: "Employee", Joined: "2023-11-20"},
							},
							TotalCount:  47,
							CurrentPage: 1,
							TotalPages:  3,
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "table-empty",
			Name:        "Table — Empty State",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Full-width empty-state row inside a tbody when the list has no items.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Empty state row with configurable message and column span.",
					RenderFunc: func(params url.Values) templ.Component {
						message := params.Get("message")
						if message == "" {
							message = "No records found."
						}
						cols, _ := parseInt(params.Get("cols"))
						if cols == 0 {
							cols = 3
						}
						tableContent := seq(
							withChildren(
								table.TableHeadWithBoundary(),
								withChildren(
									table.TableHeadRowWithBoundary(),
									func() templ.Component {
										headers := []string{"Name", "Status", "Role"}[:cols]
										comps := make([]templ.Component, 0, len(headers))
										for _, h := range headers {
											comps = append(comps, table.TableHeadCellWithBoundary(h))
										}
										return seq(comps...)
									}(),
								),
							),
							withChildren(
								table.TableBodyWithBoundary(),
								table.TableEmpty(cols, message),
							),
						)
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6">`); err != nil {
								return err
							}
							if err := withChildren(table.TableWithBoundary(), tableContent).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Message", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "No records found.", QueryParam: "message"},
						{Label: "Columns", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "3", QueryParam: "cols", Options: []galleryruntime.TokenOption{
							{Value: "2", Label: "2"},
							{Value: "3", Label: "3"},
							{Value: "4", Label: "4"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Empty state spanning 3 columns inside a zebra-striped table.",
					RenderFunc: func(_ url.Values) templ.Component {
						tableContent := seq(
							withChildren(
								table.TableHeadWithBoundary(),
								withChildren(
									table.TableHeadRowWithBoundary(),
									seq(
										table.TableHeadCellWithBoundary("Name"),
										table.TableHeadCellWithBoundary("Status"),
										table.TableHeadCellWithBoundary("Role"),
									),
								),
							),
							withChildren(
								table.TableBodyWithBoundary(),
								table.TableEmpty(3, "No records found."),
							),
						)
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6">`); err != nil {
								return err
							}
							if err := withChildren(table.TableWithBoundary(), tableContent).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "progress-card",
			Name:        "Progress Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "Card with a gradient header, a progress bar, and an optional stats row.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single progress card with configurable progress value and layout.",
					RenderFunc: func(params url.Values) templ.Component {
						progress, _ := parseInt(params.Get("progress"))
						if progress == 0 {
							progress = 72
						}
						horizontal := params.Get("layout") == "horizontal"
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 max-w-sm">`); err != nil {
								return err
							}
							if err := ui.ProgressCard(ui.ProgressCardProps{
								Title:         "Case Compliance",
								Subtitle:      "Johnson v. Smith",
								ProgressValue: progress,
								ProgressLabel: fmt.Sprintf("%d%%", progress),
								GradientClass: "bg-gradient-to-r from-primary/10 to-primary/5",
								Stats: []ui.ProgressStat{
									{Label: "Tasks", Value: "18 / 25"},
									{Label: "Documents", Value: "12 / 15"},
									{Label: "Due", Value: "Apr 30"},
								},
								Horizontal: horizontal,
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Progress", Group: "Content", Type: galleryruntime.TokenTypeRange, Default: "72", QueryParam: "progress", Min: 0, Max: 100, Step: 1},
						{Label: "Layout", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "vertical", QueryParam: "layout", Options: []galleryruntime.TokenOption{
							{Value: "vertical", Label: "Vertical"},
							{Value: "horizontal", Label: "Horizontal"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Full gradient header card and a compact horizontal inline variant.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-4 max-w-lg">`); err != nil {
								return err
							}
							if err := ui.ProgressCard(ui.ProgressCardProps{
								Title:         "Case Compliance",
								Subtitle:      "Johnson v. Smith",
								ProgressValue: 72,
								ProgressLabel: "72%",
								GradientClass: "bg-gradient-to-r from-primary/10 to-primary/5",
								Stats: []ui.ProgressStat{
									{Label: "Tasks", Value: "18 / 25"},
									{Label: "Documents", Value: "12 / 15"},
									{Label: "Due", Value: "Apr 30"},
								},
							}).Render(ctx, w); err != nil {
								return err
							}
							if err := ui.ProgressCard(ui.ProgressCardProps{
								Title:         "Document Review",
								Subtitle:      "3 of 8 complete",
								ProgressValue: 38,
								Horizontal:    true,
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "stat-card-minimal",
			Name:        "Stat Card — Minimal",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "KPI stat cards with trend indicators (↑↓). Use on dashboards to surface key metrics.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single stat card with configurable value, label, and trend.",
					RenderFunc: func(params url.Values) templ.Component {
						value := params.Get("value")
						if value == "" {
							value = "142"
						}
						label := params.Get("label")
						if label == "" {
							label = "Open Cases"
						}
						trendLabel := params.Get("trend_label")
						if trendLabel == "" {
							trendLabel = "12.3%"
						}
						trend := ui.StatTrend(params.Get("trend"))
						if trend == "" {
							trend = ui.StatTrendUp
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 max-w-xs">`); err != nil {
								return err
							}
							if err := ui.StatCardMinimal(ui.StatCardMinimalItem{
								Label:      label,
								Value:      value,
								Trend:      trend,
								TrendLabel: trendLabel,
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Value", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "142", QueryParam: "value"},
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Open Cases", QueryParam: "label"},
						{Label: "Trend", Group: "Content", Type: galleryruntime.TokenTypeSelect, Default: "up", QueryParam: "trend", Options: []galleryruntime.TokenOption{
							{Value: "up", Label: "Up ↑"},
							{Value: "down", Label: "Down ↓"},
							{Value: "", Label: "Neutral"},
						}},
						{Label: "Trend label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "12.3%", QueryParam: "trend_label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Four KPI stat cards in a responsive grid.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []ui.StatCardMinimalItem{
							{Label: "Open Cases", Value: "142", Trend: ui.StatTrendUp, TrendLabel: "12.3%"},
							{Label: "Pending Tasks", Value: "38", Trend: ui.StatTrendDown, TrendLabel: "4.1%"},
							{Label: "Clients", Value: "89", Trend: ui.StatTrendUp, TrendLabel: "7.8%"},
							{Label: "Avg. Case Days", Value: "24", Trend: ui.StatTrendUp, TrendLabel: "2.5%"},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6"><div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">`); err != nil {
								return err
							}
							for _, item := range items {
								if err := ui.StatCardMinimal(item).Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "stat-card-icon-corner",
			Name:        "Stat Card — Icon Corner",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "Stat cards with a floating icon in the top corner and a soft badge for the trend.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single stat card with configurable value, label, icon, and trend.",
					RenderFunc: func(params url.Values) templ.Component {
						value := params.Get("value")
						if value == "" {
							value = "142"
						}
						label := params.Get("label")
						if label == "" {
							label = "Open Cases"
						}
						icon := params.Get("icon")
						if icon == "" {
							icon = "lucide--briefcase"
						}
						trendLabel := params.Get("trend_label")
						if trendLabel == "" {
							trendLabel = "14.6%"
						}
						trend := ui.StatTrend(params.Get("trend"))
						if trend == "" {
							trend = ui.StatTrendUp
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 max-w-xs">`); err != nil {
								return err
							}
							if err := ui.StatCardIconCorner(ui.StatCardIconCornerItem{
								Value:      value,
								Label:      label,
								Icon:       icon,
								IconColor:  "text-primary",
								Trend:      trend,
								TrendLabel: trendLabel,
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Value", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "142", QueryParam: "value"},
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Open Cases", QueryParam: "label"},
						{Label: "Icon", Group: "Content", Type: galleryruntime.TokenTypeSelect, Default: "lucide--briefcase", QueryParam: "icon", Options: []galleryruntime.TokenOption{
							{Value: "lucide--briefcase", Label: "Briefcase"},
							{Value: "lucide--users", Label: "Users"},
							{Value: "lucide--check-square", Label: "Check square"},
							{Value: "lucide--dollar-sign", Label: "Dollar"},
							{Value: "lucide--bar-chart-2", Label: "Chart"},
						}},
						{Label: "Trend", Group: "Content", Type: galleryruntime.TokenTypeSelect, Default: "up", QueryParam: "trend", Options: []galleryruntime.TokenOption{
							{Value: "up", Label: "Up ↑"},
							{Value: "down", Label: "Down ↓"},
							{Value: "", Label: "Neutral"},
						}},
						{Label: "Trend label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "14.6%", QueryParam: "trend_label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Four stat cards with corner icons in a responsive grid.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []ui.StatCardIconCornerItem{
							{Value: "142", Label: "Open Cases", Trend: ui.StatTrendUp, TrendLabel: "14.6%", Icon: "lucide--briefcase", IconColor: "text-primary"},
							{Value: "38", Label: "Pending Tasks", Trend: ui.StatTrendDown, TrendLabel: "4.1%", Icon: "lucide--check-square", IconColor: "text-warning"},
							{Value: "89", Label: "Active Clients", Trend: ui.StatTrendUp, TrendLabel: "7.8%", Icon: "lucide--users", IconColor: "text-success"},
							{Value: "$48K", Label: "Revenue (MTD)", Trend: ui.StatTrendUp, TrendLabel: "9.2%", Icon: "lucide--dollar-sign", IconColor: "text-secondary"},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6"><div class="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-4">`); err != nil {
								return err
							}
							for _, item := range items {
								if err := ui.StatCardIconCorner(item).Render(ctx, w); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
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
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Dashboard skeleton layout placeholder.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.SkeletonDashboard(nil)
					},
					FrameHeight: "480px",
					Tokens:      []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Default 4-stat + chart/side-panel + full-width table layout.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.SkeletonDashboard(nil)
					},
					FrameHeight: "480px",
					Tokens:      []galleryruntime.DesignToken{},
				},
			},
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
			Slug:        "notification-panel",
			Name:        "Notification Panel",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Notifications",
			Description: "Tab-based notification center with All / Unread tabs and a list of notification items.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Notification panel with three sample items.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []ui.NotificationItem{
							{
								IconClass:     "bg-primary/10",
								IconTextClass: "text-primary",
								IconName:      "lucide--briefcase",
								Title:         "New case assigned",
								Body:          "Johnson v. Smith was assigned to you.",
								Time:          "2 min ago",
								Unread:        true,
							},
							{
								IconClass:     "bg-warning/10",
								IconTextClass: "text-warning",
								IconName:      "lucide--check-square",
								Title:         "Task deadline tomorrow",
								Body:          "File motion for Johnson v. Smith due soon.",
								Time:          "1 hour ago",
								Unread:        true,
							},
							{
								IconClass:     "bg-success/10",
								IconTextClass: "text-success",
								IconName:      "lucide--user",
								Title:         "Client signed in",
								Body:          "Alice Johnson accessed the client portal.",
								Time:          "Yesterday",
								Unread:        false,
							},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex justify-center">`); err != nil {
								return err
							}
							if err := ui.NotificationPanel(items, 2, "#").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Overlays extras ───────────────────────────────────────────────────
		{
			Slug:        "fab",
			Name:        "FAB",
			Category:    galleryruntime.CategoryOverlays,
			Description: "CSS-only floating action button with an expandable sub-menu of quick actions. No JS required.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "FAB appears bottom-right. Click it to expand sub-actions.",
					RenderFunc: func(_ url.Values) templ.Component {
						actions := []ui.FABAction{
							{Label: "New Case", Icon: "lucide--briefcase"},
							{Label: "Upload Doc", Icon: "lucide--file-up"},
							{Label: "Add Task", Icon: "lucide--check-square"},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="relative h-56 bg-base-100 rounded-lg border border-base-200 overflow-hidden"><p class="text-xs text-base-content/50 p-4">FAB appears bottom-right. Click it to expand sub-actions.</p>`); err != nil {
								return err
							}
							if err := ui.FAB("lucide--plus", actions).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Navigation extras ─────────────────────────────────────────────────
		{
			Slug:        "page-title-minimal",
			Name:        "Page Title — Minimal",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Page Title",
			Description: "Breadcrumb-only page header with home icon. Lightweight variant for inner pages.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Compact title bar with inline breadcrumb trail.",
					RenderFunc: func(_ url.Values) templ.Component {
						steps := []nav.PageTitleStep{
							{Label: "", Href: "#", Icon: "lucide--home"},
							{Label: "Cases", Href: "#", Icon: "lucide--briefcase"},
							{Label: "New"},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 bg-base-100">`); err != nil {
								return err
							}
							if err := nav.PageTitleMinimal("Create New Case", steps).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "page-title-editor",
			Name:        "Page Title — Editor",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Page Title",
			Description: "Full page title with DaisyUI breadcrumbs, subtitle meta, and action buttons.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Editor title with breadcrumbs, subtitle, and action buttons.",
					RenderFunc: func(_ url.Values) templ.Component {
						steps := []nav.BreadcrumbStep{
							{Label: "Dashboard", URL: "#"},
							{Label: "Cases", URL: "#"},
							{Label: "Johnson v. Smith"},
						}
						actions := []nav.PageTitleEditorAction{
							{Label: "Save Changes", Class: "btn-primary btn-sm"},
							{Label: "Preview", Class: "btn-outline btn-sm border-base-300"},
							{Icon: "lucide--ellipsis-vertical", Class: "btn-outline btn-sm border-base-300 btn-square", AriaLabel: "More options"},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 bg-base-100">`); err != nil {
								return err
							}
							if err := nav.PageTitleEditor(steps, "Johnson v. Smith", "Type: Civil Litigation", actions).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "search-dropdown",
			Name:        "Search — Dropdown",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Inline search input with a results dropdown showing recent searches and suggested items. CSS-only — no JS required.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Search input with recent and suggested result sections.",
					RenderFunc: func(_ url.Values) templ.Component {
						sections := []ui.SearchDropdownSection{
							{
								Title: "Recent",
								Items: []ui.SearchDropdownItem{
									{Icon: "lucide--briefcase", Label: "Johnson v. Smith"},
									{Icon: "lucide--user", Label: "Alice Johnson"},
								},
							},
							{
								Title: "Suggestions",
								Items: []ui.SearchDropdownItem{
									{Icon: "lucide--file", Label: "Contract_Draft_v3.pdf"},
								},
							},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex justify-center">`); err != nil {
								return err
							}
							if err := ui.SearchDropdown("Search cases, clients…", sections).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Forms extras ──────────────────────────────────────────────────────
		{
			Slug:        "filter-bar",
			Name:        "Filter Bar",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Filters",
			Description: "FilterCard wraps filter inputs in a card with Filter/Clear buttons. CompactFilterBar is the inline variant used above tables.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "FilterCard and CompactFilterBar with sample search and status inputs.",
					RenderFunc: func(_ url.Values) templ.Component {
						filterInputs := seq(
							form.SearchInputWithBoundary("q", "", "Search cases…", "", "#"),
							form.SelectInputWithBoundary("status", "Status", "", [][2]string{
								{"", "All statuses"},
								{"active", "Active"},
								{"pending", "Pending"},
								{"closed", "Closed"},
							}, "", false),
						)
						compactInputs := seq(
							form.SearchInputWithBoundary("q", "", "Search…", "", "#"),
							form.SelectInputWithBoundary("status", "", "", [][2]string{
								{"", "All statuses"},
								{"active", "Active"},
								{"closed", "Closed"},
							}, "", false),
						)
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6"><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">FilterCard</p>`); err != nil {
								return err
							}
							if err := withChildren(ui.FilterCard(ui.FilterCardProps{Action: "#"}), filterInputs).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">CompactFilterBar</p>`); err != nil {
								return err
							}
							if err := withChildren(ui.CompactFilterBar(ui.FilterCardProps{Action: "#"}), compactInputs).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "fieldset",
			Name:        "Fieldset",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Layout",
			Description: "Fieldset wrapper with an optional legend label grouping related form inputs.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Two fieldsets grouping personal information and case detail inputs.",
					RenderFunc: func(_ url.Values) templ.Component {
						personal := withChildren(
							ui.Fieldset("Personal Information"),
							seq(
								form.TextInputWithBoundary("full_name", "Full name", "Alice Johnson", "", false),
								form.TextInputWithBoundary("email", "Email", "alice@example.com", "", false),
							),
						)
						caseDetails := withChildren(
							ui.Fieldset("Case Details"),
							form.SelectInputWithBoundary("case_type", "Case type", "Civil", [][2]string{
								{"Civil", "Civil"},
								{"Criminal", "Criminal"},
								{"Family", "Family"},
							}, "", false),
						)
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 max-w-md space-y-4">`); err != nil {
								return err
							}
							if err := personal.Render(ctx, w); err != nil {
								return err
							}
							if err := caseDetails.Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "prompt-bar-minimal",
			Name:        "Prompt Bar — Minimal",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Minimal AI prompt / chat input with token counter and send button.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Prompt bar with attach, image, voice, and token counter.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex justify-center">`); err != nil {
								return err
							}
							if err := form.PromptBar(form.PromptBarProps{
								Placeholder:      "Describe what you want to generate or ask a question…",
								ShowTokenCounter: true,
								TokenCurrent:     88,
								TokenMax:         100,
								ShowAttach:       true,
								ShowImage:        true,
								ShowVoice:        true,
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "prompt-bar-action",
			Name:        "Prompt Bar — Action",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "AI prompt input with quick-action buttons (Add File, Deep Thinking, Browsing).",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Prompt bar with labelled quick-action buttons.",
					RenderFunc: func(_ url.Values) templ.Component {
						actions := []form.PromptBarActionItem{
							{Label: "Add File", Icon: "lucide--circle-plus"},
							{Label: "Deep Think", Icon: "lucide--lightbulb"},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex justify-center">`); err != nil {
								return err
							}
							if err := form.PromptBarAction("Type your request or attach files to get started…", actions).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Foundation extras ─────────────────────────────────────────────────
		{
			Slug:        "gradient-text",
			Name:        "Gradient Text",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Gradient text using Tailwind v4 bg-linear-to-r + bg-clip-text. Useful for hero headings.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Primary→secondary, success→info, and warning→error gradient examples.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 space-y-6">
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
</div>`)
					},
				},
			},
		},
		{
			Slug:        "colored-shadows",
			Name:        "Colored Shadows",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Tailwind v4 colored shadow utilities using shadow-{color}/{opacity}.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Cards and buttons with colored drop shadows.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-8 space-y-6">
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
</div>`)
					},
				},
			},
		},

		// ── Foundation extras ─────────────────────────────────────────────────────
		{
			Slug:        "typography",
			Name:        "Typography",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Typography",
			Description: "Heading and body text scale used across the application.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Heading levels, body, muted, overline, and link styles.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 space-y-3">
  <h1 class="text-3xl font-bold text-base-content">Heading 1</h1>
  <h2 class="text-2xl font-semibold text-base-content">Heading 2</h2>
  <h3 class="text-xl font-semibold text-base-content">Heading 3</h3>
  <h4 class="text-base font-semibold text-base-content">Heading 4</h4>
  <p class="text-base text-base-content/80">Body text — regular paragraph content used in cards and detail views.</p>
  <p class="text-sm text-base-content/60">Small / muted text — used for labels, hints, and secondary information.</p>
  <p class="text-xs text-base-content/50 uppercase tracking-wide font-semibold">Overline / label text</p>
  <a href="#" class="link link-primary text-sm">Link text</a>
</div>`)
					},
				},
			},
		},
		{
			Slug:        "typography-scale",
			Name:        "Typography Scale",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Typography",
			Description: "Text size scale (xs→6xl) and font weight scale (thin→black).",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Size scale from xs to 4xl and all font weights.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="space-y-6 p-6">
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
</div>`)
					},
				},
			},
		},
		{
			Slug:        "shadow-scale",
			Name:        "Shadow Scale",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Box shadows from none→2xl, colored shadows, inset shadows, and text shadows.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Box shadow, inset shadow, and text shadow scales.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="space-y-6 p-6">
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
</div>`)
					},
				},
			},
		},
		{
			Slug:        "css-filters",
			Name:        "CSS Filters",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Tailwind filter utilities: grayscale, invert, sepia, blur, brightness, contrast, saturate.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Image filter utility classes applied to sample images.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6">
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
</div>`)
					},
				},
			},
		},

		// ── Navigation extras ──────────────────────────────────────────────────────
		{
			Slug:        "footer-minimal",
			Name:        "Footer — Minimal",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Simple one-line footer with copyright text and optional links.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Footer with copyright text only.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="space-y-4 p-6"><div class="card card-border"><div class="bg-base-200/30 rounded-t-box px-5 py-3 font-medium">Copyright only</div>`); err != nil {
								return err
							}
							if err := nav.FooterMinimal("© 2025 LegalPlant. All rights reserved.", nil).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "profile-menu",
			Name:        "Profile Menu",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Misc",
			Description: "Avatar dropdown menu with grouped menu items and sign-out action.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Profile dropdown with avatar, user info, menu items, and sign-out.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []nav.ProfileMenuItem{
							{Label: "Profile", Href: "#", Icon: "lucide--user"},
							{Label: "Settings", Href: "#", Icon: "lucide--pencil"},
							{Label: "Notifications", Href: "#", Icon: "lucide--bell", Badge: 3},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex items-center justify-center p-12">`); err != nil {
								return err
							}
							if err := nav.ProfileMenu("Jane Doe", "jane@example.com", "JD", items, "#").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},

		// ── Forms extras ───────────────────────────────────────────────────────────
		{
			Slug:        "input-spinner",
			Name:        "Input Spinner",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Numeric increment/decrement input with +/- buttons. Uses vanilla JS — no library needed. Includes simple and joined variants.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Simple spinner with default styling.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6"><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Simple</p>`); err != nil {
								return err
							}
							if err := form.InputSpinner("spin1", 0, 0, 99, true, "btn-outline", "w-24").Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">With min/max (0–10)</p>`); err != nil {
								return err
							}
							if err := form.InputSpinner("spin2", 5, 0, 10, true, "btn-primary btn-sm", "w-20 input-sm").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "wizard-stepper",
			Name:        "Wizard — Stepper",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Wizard",
			Description: "Multi-step form wizard with step indicators, next/prev navigation, and a finish action. Implemented in vanilla JS — no Alpine.js needed.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Four-step case creation wizard.",
					RenderFunc: func(_ url.Values) templ.Component {
						steps := []form.WizardStep{
							{Label: "Intake"},
							{Label: "Details"},
							{Label: "Team"},
							{Label: "Review"},
						}
						panels := []form.WizardStepPanel{
							{
								Title: "Step 1 — Intake",
								Content: `<div class="form-control mb-3">
  <label class="label pb-1"><span class="label-text text-sm font-medium">Case title</span></label>
  <input type="text" placeholder="e.g. Johnson v. Smith" class="input input-bordered w-full"/>
</div>
<div class="form-control">
  <label class="label pb-1"><span class="label-text text-sm font-medium">Case type</span></label>
  <select class="select select-bordered w-full"><option>Civil</option><option>Criminal</option><option>Family</option></select>
</div>`,
							},
							{
								Title: "Step 2 — Details",
								Content: `<div class="form-control mb-3">
  <label class="label pb-1"><span class="label-text text-sm font-medium">Description</span></label>
  <textarea class="textarea textarea-bordered w-full" rows="3" placeholder="Brief description of the case…"></textarea>
</div>
<div class="form-control">
  <label class="label pb-1"><span class="label-text text-sm font-medium">Priority</span></label>
  <select class="select select-bordered w-full"><option>Normal</option><option>High</option><option>Urgent</option></select>
</div>`,
							},
							{
								Title: "Step 3 — Team",
								Content: `<div class="form-control">
  <label class="label pb-1"><span class="label-text text-sm font-medium">Lead attorney</span></label>
  <select class="select select-bordered w-full"><option>Alice Johnson</option><option>Bob Smith</option><option>Carol White</option></select>
</div>`,
							},
							{
								Title: "Step 4 — Review",
								Content: `<p class="text-sm text-base-content/60 mb-4">Review the case details before submitting.</p>
<div class="space-y-2 text-sm">
  <div class="flex gap-2"><span class="text-base-content/60 w-24">Title:</span><span class="font-medium">Johnson v. Smith</span></div>
  <div class="flex gap-2"><span class="text-base-content/60 w-24">Type:</span><span class="font-medium">Civil</span></div>
  <div class="flex gap-2"><span class="text-base-content/60 w-24">Attorney:</span><span class="font-medium">Alice Johnson</span></div>
</div>`,
							},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6">`); err != nil {
								return err
							}
							if err := form.WizardStepper("wizard-demo", steps, panels).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
		},
		{
			Slug:        "clipboard-copy",
			Name:        "Clipboard Copy",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Click-to-copy buttons with visual feedback. Uses vanilla JS navigator.clipboard — no library needed.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Copy text field, share link, and inline copy badge.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []form.ClipboardCopyItem{
							{
								ID:          "copy-input-1",
								Label:       "Copy text field",
								Value:       "CASE-2026-00142",
								ButtonLabel: "Copy",
							},
							{
								ID:          "copy-input-2",
								Label:       "Copy share link",
								Value:       "https://app.example.com/cases/CASE-2026-00142",
								ButtonLabel: "Copy Link",
								ButtonClass: "btn-primary",
							},
							{
								ID:    "copy-input-3",
								Label: "Inline copy badge",
								Value: "CASE-2026-00142",
								Mono:  true,
							},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 max-w-lg">`); err != nil {
								return err
							}
							if err := form.ClipboardCopy(items).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
			},
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
						dot := params.Get("dot") == "true"
						icon := params.Get("icon")
						label := params.Get("label")
						if label == "" {
							label = "Active"
						}
						return ui.BadgeWithBoundary(variant, style, size, dot, icon, label)
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
			Description: "DaisyUI loading spinner with centered, inline, and overlay variants.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basics",
					Description: "Configurable spinner variant: centered, inline, or overlay.",
					RenderFunc: func(params url.Values) templ.Component {
						variant := ui.LoaderVariant(params.Get("variant"))
						if variant == "" {
							variant = ui.LoaderCentered
						}
						return ui.LoaderWithBoundary(variant)
					},
					Tokens: LoaderTokens(),
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
						icon := params.Get("icon")
						size := ui.AvatarSize(params.Get("size"))
						if size == "" {
							size = ui.AvatarMD
						}
						return ui.AvatarWithBoundary(name, "", icon, size)
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
						body := seq(
							rawHTML(`<p class="text-sm text-base-content/70 mb-6">Are you sure you want to proceed? This action will be applied immediately.</p><div class="flex justify-end gap-2">`),
							withText("Cancel", ui.ButtonWithBoundary("", ui.ButtonGhost, ui.ButtonSM, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false)),
							withText("Confirm", ui.ButtonWithBoundary("", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false)),
							rawHTML(`</div>`),
						)
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
						type memberRow struct {
							Name   string
							Role   string
							Status string
							Joined string
						}
						members := []memberRow{
							{"Alice Johnson", "Admin", "Active", "Jan 2024"},
							{"Bob Martinez", "Member", "Pending", "Mar 2024"},
							{"Carol White", "Viewer", "Inactive", "Jun 2024"},
							{"David Kim", "Member", "Active", "Aug 2024"},
						}
						rowComponents := make([]templ.Component, len(members))
						for i, m := range members {
							m := m
							rowComponents[i] = withChildren(
								table.TableRowWithBoundary("", false),
								seq(
									withChildren(table.TableCellWithBoundary(""), rawHTML(m.Name)),
									withChildren(table.TableCellWithBoundary(""), rawHTML(m.Role)),
									withChildren(table.TableCellWithBoundary(""), ui.StatusBadgeWithBoundary(m.Status)),
									withChildren(table.TableCellWithBoundary(""), rawHTML(m.Joined)),
								),
							)
						}
						return withChildren(
							table.TableWithPropsWithBoundary(table.TableProps{
								Striped: zebra,
								Size:    size,
							}),
							seq(
								withChildren(
									table.TableHeadWithBoundary(),
									withChildren(
										table.TableHeadRowWithBoundary(),
										seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Status"),
											table.TableHeadCellWithBoundary("Joined"),
										),
									),
								),
								withChildren(
									table.TableBodyWithBoundary(),
									seq(rowComponents...),
								),
							),
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
				{
					Name:        "Examples",
					Description: "Default, zebra-striped, and compact (small) variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						type memberRow struct {
							Name   string
							Role   string
							Status string
							Joined string
						}
						members := []memberRow{
							{"Alice Johnson", "Admin", "Active", "Jan 2024"},
							{"Bob Martinez", "Member", "Pending", "Mar 2024"},
							{"Carol White", "Viewer", "Inactive", "Jun 2024"},
						}
						buildTable := func(props table.TableProps) templ.Component {
							rowComponents := make([]templ.Component, len(members))
							for i, m := range members {
								m := m
								rowComponents[i] = withChildren(
									table.TableRow("", false),
									seq(
										withChildren(table.TableCell(""), rawHTML(m.Name)),
										withChildren(table.TableCell(""), rawHTML(m.Role)),
										withChildren(table.TableCell(""), ui.StatusBadge(m.Status)),
										withChildren(table.TableCell(""), rawHTML(m.Joined)),
									),
								)
							}
							return withChildren(
								table.TableWithProps(props),
								seq(
									withChildren(
										table.TableHead(),
										withChildren(
											table.TableHeadRow(),
											seq(
												table.TableHeadCell("Name"),
												table.TableHeadCell("Role"),
												table.TableHeadCell("Status"),
												table.TableHeadCell("Joined"),
											),
										),
									),
									withChildren(
										table.TableBody(),
										seq(rowComponents...),
									),
								),
							)
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							variants := []struct {
								label string
								props table.TableProps
							}{
								{"Default", table.TableProps{}},
								{"Zebra stripes", table.TableProps{Striped: true}},
								{"Compact (small)", table.TableProps{Size: "sm", Striped: true}},
							}
							if _, err := io.WriteString(w, `<div class="p-6 flex flex-col gap-8">`); err != nil {
								return err
							}
							for _, v := range variants {
								if _, err := io.WriteString(w, `<div class="flex flex-col gap-2"><p class="text-xs font-semibold text-base-content/40 uppercase tracking-wide">`+v.label+`</p>`); err != nil {
									return err
								}
								if err := buildTable(v.props).Render(ctx, w); err != nil {
									return err
								}
								if _, err := io.WriteString(w, `</div>`); err != nil {
									return err
								}
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{},
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
