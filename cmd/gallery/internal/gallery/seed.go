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
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/components/form"
	"github.com/emergent-company/go-daisy/components/layout"
	"github.com/emergent-company/go-daisy/components/logs"
	"github.com/emergent-company/go-daisy/pages"
	"github.com/emergent-company/go-daisy/components/modal"
	"github.com/emergent-company/go-daisy/components/nav"
	"github.com/emergent-company/go-daisy/components/table"
	"github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/devmode"
	"github.com/emergent-company/go-daisy/galleryruntime"
	"github.com/emergent-company/go-daisy/shared"
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
	components := []galleryruntime.GalleryComponent{

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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Vertical timeline with mixed done/pending states.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.TimelineWithBoundary([]ui.TimelineItemProps{
							{Date: "Day 1", Label: "Order placed", Done: true},
							{Date: "Day 2", Label: "Processing", Done: true},
							{Date: "Day 3", Label: "Shipped", Done: false},
							{Date: "Day 4", Label: "Delivered", Done: false},
						})
					},
				},
			},
		},

		// ── Data Display / Chat ───────────────────────────────────────────────────
		{
			Slug:        "chat-bubble",
			Name:        "Chat Bubble",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Chat conversation bubbles (sent and received) with avatar, bot-icon, and hover action toolbar.",
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
						return ui.ChatBubbleWithBoundary(sent, author, timestamp, "", false, "", false, message)
					},
					Tokens: ChatBubbleTokens(),
				},
				{
					Name:        "Examples",
					Description: "Sent and received bubbles with avatars.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-col gap-2 p-4 max-w-lg mx-auto">`); err != nil {
								return err
							}
							if err := withText("Hey! How are you doing?", ui.ChatBubble(false, "Alice", "10:32 AM", "./images/avatars/2.png", false, "bg-base-200", false, nil)).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("Good thanks! How about you?", ui.ChatBubble(true, "You", "10:33 AM", "./images/avatars/1.png", false, "bg-base-200", false, nil)).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("Just finished a great book. Have any recommendations?", ui.ChatBubble(false, "Alice", "10:34 AM", "./images/avatars/2.png", false, "bg-base-200", false, nil)).Render(ctx, w); err != nil {
								return err
							}
							if err := withText("I'd recommend 'The Silent Observer' — it's a must-read!", ui.ChatBubble(true, "You", "10:36 AM", "./images/avatars/1.png", false, "bg-base-200", false, nil)).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
				{
					Name:        "AI Conversation",
					Description: "AI bot bubble with bot icon, hover actions, and thinking indicator.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-col gap-3 p-4 max-w-xl mx-auto">`); err != nil {
								return err
							}
							// User message
							if err := withText("Can you provide an estimated timeline for completion?", ui.ChatBubble(true, "", "Weeks ago", "", false, "bg-base-200", false, nil)).Render(ctx, w); err != nil {
								return err
							}
							// AI response with bot icon and hover actions
							if err := withText("Certainly! Based on our current progress, we estimate the project will be completed within 4–6 weeks. Let me know if you'd like a detailed breakdown.", ui.ChatBubble(false, "", "Week ago", "", true, "bg-base-200", true, nil)).Render(ctx, w); err != nil {
								return err
							}
							// Another user message
							if err := withText("Can you generate a follow-up summary?", ui.ChatBubble(true, "", "now", "", false, "bg-base-200", false, nil)).Render(ctx, w); err != nil {
								return err
							}
							// AI thinking indicator
							if err := ui.AIThinkingIndicator().Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
				},
				{
					Name:        "Bubble Colors",
					Description: "Chat bubbles using DaisyUI color modifiers.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-col gap-2 p-4 max-w-sm mx-auto">`); err != nil {
								return err
							}
							type bubble struct {
								text  string
								class string
								sent  bool
							}
							bubbles := []bubble{
								{"Primary message", "chat-bubble-primary", false},
								{"Secondary reply", "chat-bubble-secondary", true},
								{"Accent highlight", "chat-bubble-accent", false},
								{"Success confirmation", "chat-bubble-success", true},
								{"Warning notice", "chat-bubble-warning", false},
								{"Error alert", "chat-bubble-error", true},
							}
							for _, b := range bubbles {
								if err := withText(b.text, ui.ChatBubble(b.sent, "", "", "", false, b.class, false, nil)).Render(ctx, w); err != nil {
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
			Slug:        "chat-input",
			Name:        "Chat Input",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Compact single-line chat input bar with optional attach button and send button.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Chat input with attach button.",
					RenderFunc: func(params url.Values) templ.Component {
						placeholder := params.Get("placeholder")
						if placeholder == "" {
							placeholder = "Type a message..."
						}
						return ui.ChatInputWithBoundary(true, placeholder)
					},
					Tokens: ChatInputTokens(),
				},
				{
					Name:        "Examples",
					Description: "Chat input variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="space-y-4 max-w-lg mx-auto p-4">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<p class="text-xs text-base-content/60 font-semibold uppercase">With attach button</p>`); err != nil {
								return err
							}
							if err := ui.ChatInput(true, "Type a message...", nil).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<p class="text-xs text-base-content/60 font-semibold uppercase">Without attach button</p>`); err != nil {
								return err
							}
							if err := ui.ChatInput(false, "Send a reply...", nil).Render(ctx, w); err != nil {
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Terminal code block with multiple lines and colors.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupCodeWithBoundary([]ui.MockupCodeLineProps{
							{Prefix: "$", Code: "npm install go-daisy"},
							{Prefix: ">", Code: "Installing packages...", ColorClass: "text-warning"},
							{Prefix: "", Code: "added 42 packages"},
							{Prefix: "", Code: "Done!", ColorClass: "text-success"},
						})
					},
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
				{
					Name:        "Examples",
					Description: "Browser mockup with a custom URL.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupBrowserWithBoundary("https://app.example.com/dashboard")
					},
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
					Description: "All alert types shown individually.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Success",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AlertWithBoundary(ui.AlertSuccess, "lucide--circle-check", "Your changes have been saved successfully.")
							},
						},
						{
							Label: "Error",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AlertWithBoundary(ui.AlertError, "lucide--circle-x", "Something went wrong. Please try again.")
							},
						},
						{
							Label: "Warning",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AlertWithBoundary(ui.AlertWarning, "lucide--triangle-alert", "Your session will expire in 5 minutes.")
							},
						},
						{
							Label: "Info",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AlertWithBoundary(ui.AlertInfo, "lucide--info", "A new software update is available.")
							},
						},
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
				{
					Name:        "Examples",
					Description: "Filter tabs with different selections.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "All selected",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.FilterTabsWithBoundary("filter1", "All", []string{"All", "Active", "Pending", "Closed"})
							},
						},
						{
							Label: "Active selected",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.FilterTabsWithBoundary("filter2", "Active", []string{"All", "Active", "Pending", "Closed"})
							},
						},
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
					Description: "Checkboxes and toggles in checked and unchecked states.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Checkbox (checked)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CheckboxWithBoundary("n1", true, "Receive email notifications")
							},
						},
						{
							Label: "Checkbox (unchecked)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CheckboxWithBoundary("n2", false, "Subscribe to newsletter")
							},
						},
						{
							Label: "Toggle (on)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.ToggleWithBoundary("dark", true, "Dark mode")
							},
						},
						{
							Label: "Toggle (off)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.ToggleWithBoundary("autosave", false, "Auto-save")
							},
						},
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
						return form.RadioGroupWithBoundary("plan", "free", []form.SelectOption{
							{Value: "free", Label: "Free – $0/mo"},
							{Value: "pro", Label: "Pro – $12/mo"},
							{Value: "enterprise", Label: "Enterprise – Custom"},
						}, color)
					},
					Tokens: RadioGroupTokens(),
				},
				{
					Name:        "Examples",
					Description: "Radio groups with different colors.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Primary",
							RenderFunc: func(_ url.Values) templ.Component {
								opts := []form.SelectOption{{Value: "opt1", Label: "Option 1"}, {Value: "opt2", Label: "Option 2"}, {Value: "opt3", Label: "Option 3"}}
								return form.RadioGroupWithBoundary("radio-primary", "opt1", opts, "radio-primary")
							},
						},
						{
							Label: "Secondary",
							RenderFunc: func(_ url.Values) templ.Component {
								opts := []form.SelectOption{{Value: "opt1", Label: "Option 1"}, {Value: "opt2", Label: "Option 2"}, {Value: "opt3", Label: "Option 3"}}
								return form.RadioGroupWithBoundary("radio-secondary", "opt1", opts, "radio-secondary")
							},
						},
						{
							Label: "Accent",
							RenderFunc: func(_ url.Values) templ.Component {
								opts := []form.SelectOption{{Value: "opt1", Label: "Option 1"}, {Value: "opt2", Label: "Option 2"}, {Value: "opt3", Label: "Option 3"}}
								return form.RadioGroupWithBoundary("radio-accent", "opt1", opts, "radio-accent")
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Star and heart rating shapes at different values.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Stars (3/5)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.RatingWithBoundary("r1", 3, 5, form.RatingStar, "rating-warning", "")
							},
						},
						{
							Label: "Hearts (4/5)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.RatingWithBoundary("r2", 4, 5, form.RatingHeart, "rating-error", "")
							},
						},
					},
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
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Default (OR)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.DividerWithBoundary(ui.DividerDefault, false, "OR")
							},
						},
						{
							Label: "Primary",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.DividerWithBoundary(ui.DividerPrimary, false, "Primary")
							},
						},
						{
							Label: "Vertical",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.DividerWithBoundary(ui.DividerDefault, true, "")
							},
						},
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
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "⌘K search",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.KbdWithBoundary(ui.KbdSM, "⌘K")
							},
						},
						{
							Label: "Ctrl+S",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.KbdWithBoundary(ui.KbdSM, "Ctrl+S")
							},
						},
						{
							Label: "Enter (large)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.KbdWithBoundary(ui.KbdLG, "Enter")
							},
						},
						{
							Label: "Esc (xs)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.KbdWithBoundary(ui.KbdXS, "Esc")
							},
						},
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Progress bars in all DaisyUI color variants.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Primary (40%)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressWithBoundary(ui.ProgressPrimary, 40, 100)
							},
						},
						{
							Label: "Secondary (60%)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressWithBoundary(ui.ProgressSecondary, 60, 100)
							},
						},
						{
							Label: "Success (75%)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressWithBoundary(ui.ProgressSuccess, 75, 100)
							},
						},
						{
							Label: "Success (90%)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressWithBoundary(ui.ProgressSuccess, 90, 100)
							},
						},
						{
							Label: "Error (25%)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressWithBoundary(ui.ProgressError, 25, 100)
							},
						},
						{
							Label: "Warning (50%)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressWithBoundary(ui.ProgressWarning, 50, 100)
							},
						},
					},
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Horizontal and vertical step trackers.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "2 of 4 done",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StepsWithBoundary([]ui.StepProps{
									{Label: "Register", Done: true},
									{Label: "Profile", Done: true},
									{Label: "Billing", Done: false},
									{Label: "Confirm", Done: false},
								})
							},
						},
						{
							Label: "All complete",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StepsWithBoundary([]ui.StepProps{
									{Label: "Draft", Done: true},
									{Label: "Review", Done: true},
									{Label: "Published", Done: true},
								})
							},
						},
					},
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
					Name:        "Interactive",
					Description: "Accordion with multiple collapsible items.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.AccordionWithBoundary([]ui.AccordionItemProps{
							{Title: "What is go-daisy?", Content: templ.Raw("go-daisy is a Go UI component library for HTMX-driven web interfaces built with DaisyUI."), Open: true},
							{Title: "How do I install it?", Content: templ.Raw("<code>go get github.com/emergent-company/go-daisy</code>"), Open: false},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Multiple accordion items open/closed.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.AccordionWithBoundary([]ui.AccordionItemProps{
							{Title: "What is DaisyUI?", Content: rawHTML("DaisyUI is a plugin for Tailwind CSS that adds component classes."), Open: true},
							{Title: "Is it free?", Content: rawHTML("Yes, DaisyUI is free and open-source."), Open: false},
							{Title: "Does it support dark mode?", Content: rawHTML("Yes, DaisyUI supports light and dark themes out of the box."), Open: false},
						})
					},
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
						return ui.DropdownWithBoundary(align, ui.DropdownTrigger("Options", "btn-primary", nil), []ui.DropdownItemProps{
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
				{
					Name:        "Examples",
					Description: "Dropdown variants: left-aligned and right-aligned.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []ui.DropdownItemProps{
							{Label: "Edit"},
							{Label: "Duplicate"},
							{Label: "Delete", Danger: true},
						}
						return row(
							withText("Align left", ui.DropdownWithBoundary("", ui.Button(ui.ButtonProps{Variant: ui.ButtonPrimary, Size: ui.ButtonSM, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeDefault}), items)),
							withText("Align right", ui.DropdownWithBoundary(ui.DropdownEnd, ui.Button(ui.ButtonProps{Variant: ui.ButtonGhost, Size: ui.ButtonSM, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeDefault}), items)),
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
							ui.SimpleButton("ON", "btn-success", ui.ButtonSM, nil),
							ui.SimpleButton("OFF", ui.ButtonGhost, ui.ButtonSM, nil),
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
				{
					Name:        "Examples",
					Description: "Hero sections with different heights and copy.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="space-y-4">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-1 font-semibold uppercase px-4 pt-4">Compact</p>`); err != nil {
								return err
							}
							if err := ui.HeroWithBoundary("min-h-24", "Welcome", "Start building today.", "Get started").Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-1 font-semibold uppercase px-4">Full height</p>`); err != nil {
								return err
							}
							if err := ui.HeroWithBoundary("min-h-screen", "Build faster with go-daisy", "Type-safe HTMX components for Go.", "Explore components").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
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
					Name:        "Interactive",
					Description: "Single list row with live controls for every prop.",
					RenderFunc: func(params url.Values) templ.Component {
						name := params.Get("title")
						if name == "" {
							name = "Alice Johnson"
						}
						subtitle := params.Get("subtitle")
						if subtitle == "" {
							subtitle = "alice@example.com"
						}
						description := params.Get("description")
						status := params.Get("status")
						if status == "" {
							status = "active"
						}
						showLeading := params.Get("leading") != "no"
						showHeader := params.Get("header") != "no"

						var leading templ.Component
						if showLeading {
							leading = ui.PersonCellWithBoundary(ui.PersonCellProps{Name: name, Subtitle: subtitle})
						}
						var trailing []templ.Component
						if status != "none" {
							trailing = []templ.Component{ui.StatusBadgeWithBoundary(status)}
						}
						header := ""
						if showHeader {
							header = "Members"
						}
						// Use LeadingGrow so PersonCell (which already contains name+subtitle)
						// fills the available space; Title/Subtitle on the row stay empty.
						return ui.ListWithBoundary(ui.ListProps{Header: header}, []ui.ListRowProps{
							{
								Description: description,
								Leading:     leading,
								LeadingGrow: showLeading,
								Trailing:    trailing,
							},
						})
					},
					Tokens: ListTokens(),
				},
				{
					Name:        "Examples",
					Description: "All three list layout patterns: default, col-wrap description, and multiple trailing actions.",
					RenderFunc: func(_ url.Values) templ.Component {
					editBtn := ui.Button(ui.ButtonProps{Variant: ui.ButtonGhost, Size: ui.ButtonSM, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeDefault, Icon: "lucide--pencil"})
					deleteBtn := ui.Button(ui.ButtonProps{Variant: ui.ButtonGhost, Size: ui.ButtonSM, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeDefault, Icon: "lucide--trash-2"})
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}

							// ── Default: PersonCell + status badge ──────────────────────────
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Default — person cell + trailing badge</p>`); err != nil {
								return err
							}
							if err := ui.List(ui.ListProps{Header: "Members"}).Render(templ.WithChildren(ctx, seq(
								ui.ListRow(ui.ListRowProps{
									Leading:     ui.PersonCellWithBoundary(ui.PersonCellProps{Name: "Alice Johnson", Subtitle: "alice@example.com"}),
									LeadingGrow: true,
									Trailing:    []templ.Component{ui.StatusBadgeWithBoundary("active")},
								}),
								ui.ListRow(ui.ListRowProps{
									Leading:     ui.PersonCellWithBoundary(ui.PersonCellProps{Name: "Bob Smith", Subtitle: "bob@example.com"}),
									LeadingGrow: true,
									Trailing:    []templ.Component{ui.StatusBadgeWithBoundary("closed")},
								}),
								ui.ListRow(ui.ListRowProps{
									Leading:     ui.PersonCellWithBoundary(ui.PersonCellProps{Name: "Carol White", Subtitle: "carol@example.com"}),
									LeadingGrow: true,
									Trailing:    []templ.Component{ui.StatusBadgeWithBoundary("pending")},
								}),
							)), w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div>`); err != nil {
								return err
							}

							// ── Col-wrap: description wraps to second row ─────────────────────
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Col-wrap — description on second row</p>`); err != nil {
								return err
							}
							if err := ui.List(ui.ListProps{}).Render(templ.WithChildren(ctx, seq(
								ui.ListRow(ui.ListRowProps{
									Title:       "Design System Audit",
									Subtitle:    "Due Nov 15",
									Description: "Review all components for accessibility compliance and update token usage across the board.",
									Leading:     ui.AvatarWithBoundary("DS Audit", "", "lucide--clipboard-list", ui.AvatarSM),
									Trailing:    []templ.Component{ui.BadgeWithBoundary(ui.BadgeWarning, ui.BadgeStyleSoft, ui.BadgeSizeMD, false, "", "In Progress")},
								}),
								ui.ListRow(ui.ListRowProps{
									Title:       "Migrate API to v2",
									Subtitle:    "Due Dec 1",
									Description: "Refactor all client calls to use the new v2 endpoints. Coordinate with the backend team.",
									Leading:     ui.AvatarWithBoundary("API Migration", "", "lucide--code-2", ui.AvatarSM),
									Trailing:    []templ.Component{ui.BadgeWithBoundary(ui.BadgeSuccess, ui.BadgeStyleSoft, ui.BadgeSizeMD, false, "", "Done")},
								}),
							)), w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div>`); err != nil {
								return err
							}

							// ── Multiple trailing actions ─────────────────────────────────────
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Multiple trailing actions</p>`); err != nil {
								return err
							}
							if err := ui.List(ui.ListProps{Header: "Files"}).Render(templ.WithChildren(ctx, seq(
								ui.ListRow(ui.ListRowProps{
									Title:    "quarterly-report.pdf",
									Subtitle: "2.4 MB · Updated 3 days ago",
									Leading:  templ.Raw(`<span class="iconify lucide--file-text size-8 text-base-content/40"></span>`),
									Trailing: []templ.Component{editBtn, deleteBtn},
								}),
								ui.ListRow(ui.ListRowProps{
									Title:    "brand-assets.zip",
									Subtitle: "14.8 MB · Updated today",
									Leading:  templ.Raw(`<span class="iconify lucide--archive size-8 text-base-content/40"></span>`),
									Trailing: []templ.Component{editBtn, deleteBtn},
								}),
								ui.ListRow(ui.ListRowProps{
									Title:    "onboarding-deck.pptx",
									Subtitle: "5.1 MB · Updated 1 week ago",
									Leading:  templ.Raw(`<span class="iconify lucide--presentation size-8 text-base-content/40"></span>`),
									Trailing: []templ.Component{editBtn, deleteBtn},
								}),
							)), w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
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
							ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, Size: ui.ButtonMD, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeDefault, Icon: "lucide--bell"}),
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
						ui.Button(ui.ButtonProps{Variant: ui.ButtonGhost, Size: ui.ButtonSM, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeSquare, Icon: "lucide--bell"}),
						),
						ui.IndicatorWithBoundary("badge badge-primary badge-xs",
							templ.NopComponent,
							ui.Avatar("AJ", "", "", ui.AvatarMD, "", nil),
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Stacked cards in different color combinations.",
					RenderFunc: func(_ url.Values) templ.Component {
						return row(
							withText("Primary stack", ui.StackWithBoundary(
								ui.StackCard("Card A", "bg-primary text-primary-content"),
								ui.StackCard("Card B", "bg-primary/80 text-primary-content"),
								ui.StackCard("Card C", "bg-primary/60 text-primary-content"),
							)),
							withText("Neutral stack", ui.StackWithBoundary(
								ui.StackCard("Card 1", "bg-base-300"),
								ui.StackCard("Card 2", "bg-base-200"),
								ui.StackCard("Card 3", "bg-base-100"),
							)),
						)
					},
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
				{
					Name:        "Examples",
					Description: "Before/after comparison with text content.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Image comparison</p>`); err != nil {
								return err
							}
							if err := ui.DiffWithBoundary("Before: Original content", "After: Updated content").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
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
					Name:        "Interactive",
					Description: "Full-width one-at-a-time horizontal carousel (snap to start).",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.CarouselWithBoundary(ui.CarouselSnapStart, false, "w-full", []ui.CarouselItemProps{
							{ID: "slide1", ItemWidth: "w-full", Content: ui.CarouselSlide("Slide 1", "bg-primary text-primary-content")},
							{ID: "slide2", ItemWidth: "w-full", Content: ui.CarouselSlide("Slide 2", "bg-secondary text-secondary-content")},
							{ID: "slide3", ItemWidth: "w-full", Content: ui.CarouselSlide("Slide 3", "bg-accent text-accent-content")},
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Snap-start full-width, snap-center, snap-end, half-width, and vertical variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						fullSlides := []ui.CarouselItemProps{
							{ID: "s1a", ItemWidth: "w-full", Content: ui.CarouselSlide("Slide 1", "bg-primary text-primary-content")},
							{ID: "s2a", ItemWidth: "w-full", Content: ui.CarouselSlide("Slide 2", "bg-secondary text-secondary-content")},
							{ID: "s3a", ItemWidth: "w-full", Content: ui.CarouselSlide("Slide 3", "bg-accent text-accent-content")},
						}
						centerSlides := []ui.CarouselItemProps{
							{ID: "s1b", ItemWidth: "w-64", Content: ui.CarouselSlide("Slide 1", "bg-primary text-primary-content")},
							{ID: "s2b", ItemWidth: "w-64", Content: ui.CarouselSlide("Slide 2", "bg-secondary text-secondary-content")},
							{ID: "s3b", ItemWidth: "w-64", Content: ui.CarouselSlide("Slide 3", "bg-accent text-accent-content")},
						}
						halfSlides := []ui.CarouselItemProps{
							{ID: "s1c", ItemWidth: "w-1/2", Content: ui.CarouselSlide("Slide 1", "bg-primary text-primary-content")},
							{ID: "s2c", ItemWidth: "w-1/2", Content: ui.CarouselSlide("Slide 2", "bg-secondary text-secondary-content")},
							{ID: "s3c", ItemWidth: "w-1/2", Content: ui.CarouselSlide("Slide 3", "bg-accent text-accent-content")},
						}
						vertSlides := []ui.CarouselItemProps{
							{ID: "s1d", ItemWidth: "h-full", Content: ui.CarouselSlide("Slide 1", "bg-primary text-primary-content")},
							{ID: "s2d", ItemWidth: "h-full", Content: ui.CarouselSlide("Slide 2", "bg-secondary text-secondary-content")},
							{ID: "s3d", ItemWidth: "h-full", Content: ui.CarouselSlide("Slide 3", "bg-accent text-accent-content")},
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							sections := []struct {
								label string
								comp  templ.Component
							}{
								{"Snap Start — full width items", ui.CarouselWithBoundary(ui.CarouselSnapStart, false, "w-full", fullSlides)},
								{"Snap Center — fixed width items", ui.CarouselWithBoundary(ui.CarouselSnapCenter, false, "w-96", centerSlides)},
								{"Snap End — fixed width items", ui.CarouselWithBoundary(ui.CarouselSnapEnd, false, "w-96", centerSlides)},
								{"Half-width items", ui.CarouselWithBoundary(ui.CarouselSnapStart, false, "w-96", halfSlides)},
								{"Vertical", ui.CarouselWithBoundary(ui.CarouselSnapStart, true, "h-48", vertSlides)},
							}
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							for _, s := range sections {
								if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">`+s.label+`</p>`); err != nil {
									return err
								}
								if err := s.comp.Render(ctx, w); err != nil {
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
				{
					Name:        "Examples",
					Description: "Various countdown configurations.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Days remaining</p>`); err != nil {
								return err
							}
							if err := ui.CountdownWithBoundary(7, 0, 0, 0).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Full countdown</p>`); err != nil {
								return err
							}
							if err := ui.CountdownWithBoundary(2, 14, 30, 45).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
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
					Name:        "Interactive",
					Description: "Phone frame with an app screen placeholder.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupPhoneWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Phone mockup with placeholder content.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupPhoneWithBoundary()
					},
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
					Name:        "Interactive",
					Description: "Desktop window frame with content placeholder.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupWindowWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Window mockup with placeholder content.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.MockupWindowWithBoundary()
					},
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
					Description: "All status colors, each with a labeled dot.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Online (Success)",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(ui.StatusDotWithBoundary(ui.StatusSuccess, false))
							},
						},
						{
							Label: "Offline (Error)",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(ui.StatusDotWithBoundary(ui.StatusError, false))
							},
						},
						{
							Label: "Away (Warning)",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(ui.StatusDotWithBoundary(ui.StatusWarning, false))
							},
						},
						{
							Label: "Busy (Info)",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(ui.StatusDotWithBoundary(ui.StatusInfo, false))
							},
						},
						{
							Label: "Unknown (Neutral)",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(ui.StatusDotWithBoundary(ui.StatusNeutral, false))
							},
						},
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
						return ui.DropdownWithBoundary(align, ui.DropdownTrigger(label, "btn-primary", nil), []ui.DropdownItemProps{
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
						ui.DropdownWithBoundary("", ui.DropdownTrigger("Bottom ▼", "btn-outline", nil), items),
						ui.DropdownWithBoundary(ui.DropdownTop, ui.DropdownTrigger("Top ▲", "btn-outline", nil), items),
						ui.DropdownWithBoundary(ui.DropdownEnd, ui.DropdownTrigger("Options ⋮", "btn-primary", nil), []ui.DropdownItemProps{
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
				{
					Name:        "Examples",
					Description: "Short and long breadcrumb trails.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Two levels</p>`); err != nil {
								return err
							}
							if err := nav.BreadcrumbsWithBoundary([]nav.BreadcrumbItem{
								{Label: "Home", Href: "#"},
								{Label: "Cases"},
							}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Four levels</p>`); err != nil {
								return err
							}
							if err := nav.BreadcrumbsWithBoundary([]nav.BreadcrumbItem{
								{Label: "Home", Href: "#"},
								{Label: "Cases", Href: "#"},
								{Label: "Johnson v. Smith", Href: "#"},
								{Label: "Documents"},
							}).Render(ctx, w); err != nil {
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
				{
					Name:        "Examples",
					Description: "Mobile dock with 3 and 5 item configurations.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">3 items</p>`); err != nil {
								return err
							}
							if err := nav.DockWithBoundary([]nav.DockItem{
								{Label: "Home", Icon: "lucide--home", Href: "#", Active: true},
								{Label: "Search", Icon: "lucide--search", Href: "#"},
								{Label: "Profile", Icon: "lucide--user", Href: "#"},
							}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">5 items</p>`); err != nil {
								return err
							}
							if err := nav.DockWithBoundary([]nav.DockItem{
								{Label: "Home", Icon: "lucide--home", Href: "#", Active: true},
								{Label: "Cases", Icon: "lucide--folder", Href: "#"},
								{Label: "Search", Icon: "lucide--search", Href: "#"},
								{Label: "Alerts", Icon: "lucide--bell", Href: "#"},
								{Label: "Profile", Icon: "lucide--user", Href: "#"},
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
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
				{
					Name:        "Examples",
					Description: "File input with different accept types.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Any file</p>`); err != nil {
								return err
							}
							if err := form.FileInput("upload1", "Upload file", "", "").Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Images only</p>`); err != nil {
								return err
							}
							if err := form.FileInput("upload2", "Upload image", "image/*", "").Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">PDF only</p>`); err != nil {
								return err
							}
							if err := form.FileInput("upload3", "Upload PDF", ".pdf", "").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
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
								if err := withText(l.label, nav.Link("#", l.variant, "", nil)).Render(ctx, w); err != nil {
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
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Removable",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TagWithBoundary("Contract Law", "#")
							},
						},
						{
							Label: "Read-only",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TagWithBoundary("Family Law", "")
							},
						},
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
								if err := ui.Avatar("", "", "lucide--building-2", s.size, "", nil).Render(ctx, w); err != nil {
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
							if err := ui.Avatar("", "", "lucide--building-2", ui.AvatarXS, "", nil).Render(ctx, w); err != nil {
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
						return ui.PersonChipWithBoundary(name, "bg-primary", "text-primary-content", "from-primary/20", "to-primary/5", ui.PersonChipContact{
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
			Slug:        "person-cell",
			Name:        "Person Cell",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Compact avatar + name + subtitle identity block for use in lists, tables, and flex rows.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single person cell with live controls.",
					RenderFunc: func(params url.Values) templ.Component {
						name := params.Get("name")
						if name == "" {
							name = "Alice Johnson"
						}
						subtitle := params.Get("subtitle")
						if subtitle == "" {
							subtitle = "alice@example.com"
						}
						size := ui.AvatarSize(params.Get("size"))
						if size == "" {
							size = ui.AvatarSM
						}
						return ui.PersonCellWithBoundary(ui.PersonCellProps{
							Name:     name,
							Subtitle: subtitle,
							Size:     size,
						})
					},
					Tokens: PersonCellTokens(),
				},
				{
					Name:        "Examples",
					Description: "PersonCell in different sizes and contexts.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}

							// ── Sizes ────────────────────────────────────────────────────────
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Sizes</p><div class="flex flex-col gap-4">`); err != nil {
								return err
							}
							for _, size := range []ui.AvatarSize{ui.AvatarXS, ui.AvatarSM, ui.AvatarMD, ui.AvatarLG} {
								if err := ui.PersonCell(ui.PersonCellProps{
									Name:     "Alice Johnson",
									Subtitle: "alice@example.com",
									Size:     size,
								}).Render(ctx, w); err != nil {
									return err
								}
							}
							if _, err := io.WriteString(w, `</div></div>`); err != nil {
								return err
							}

							// ── Without subtitle ─────────────────────────────────────────────
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Name only</p><div class="flex flex-col gap-4">`); err != nil {
								return err
							}
							for _, person := range []string{"Alice Johnson", "Bob Smith", "Carol White"} {
								if err := ui.PersonCell(ui.PersonCellProps{Name: person}).Render(ctx, w); err != nil {
									return err
								}
							}
							if _, err := io.WriteString(w, `</div></div>`); err != nil {
								return err
							}

							// ── In a list ────────────────────────────────────────────────────
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Inside a list row</p>`); err != nil {
								return err
							}
							if err := ui.List(ui.ListProps{Header: "Team"}).Render(templ.WithChildren(ctx, seq(
								ui.ListRow(ui.ListRowProps{
									Leading:     ui.PersonCellWithBoundary(ui.PersonCellProps{Name: "Alice Johnson", Subtitle: "alice@example.com"}),
									LeadingGrow: true,
									Trailing:    []templ.Component{ui.StatusBadgeWithBoundary("active")},
								}),
								ui.ListRow(ui.ListRowProps{
									Leading:     ui.PersonCellWithBoundary(ui.PersonCellProps{Name: "Bob Smith", Subtitle: "bob@example.com"}),
									LeadingGrow: true,
									Trailing:    []templ.Component{ui.StatusBadgeWithBoundary("pending")},
								}),
							)), w); err != nil {
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
						return ui.ProgressCardWithBoundary(ui.ProgressCardProps{
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
					Description: "Vertical with stats and compact horizontal variant.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Vertical with stats",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressCardWithBoundary(ui.ProgressCardProps{
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
								})
							},
						},
						{
							Label: "Horizontal (compact)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ProgressCardWithBoundary(ui.ProgressCardProps{
									Title:         "Document Review",
									Subtitle:      "3 of 8 complete",
									ProgressValue: 38,
									Horizontal:    true,
								})
							},
						},
					},
					Tokens: []galleryruntime.DesignToken{},
				},
			},
		},
		{
			Slug:        "stat-card-minimal",
			Name:        "Stat Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "KPI stat card with trend indicator. Set Icon to show a floating icon-corner variant.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Single stat card — toggle Icon to switch between minimal and icon-corner layouts.",
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
						trendLabel := params.Get("trend_label")
						if trendLabel == "" {
							trendLabel = "12.3%"
						}
						trend := ui.StatTrend(params.Get("trend"))
						if trend == "" {
							trend = ui.StatTrendUp
						}
						return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{
							Label:      label,
							Value:      value,
							Icon:       icon,
							IconColor:  "text-primary",
							Trend:      trend,
							TrendLabel: trendLabel,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Value", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "142", QueryParam: "value"},
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Open Cases", QueryParam: "label"},
						{Label: "Icon", Group: "Style", Type: galleryruntime.TokenTypeSelect, Default: "", QueryParam: "icon", Options: []galleryruntime.TokenOption{
							{Value: "", Label: "None (minimal)"},
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
						{Label: "Trend label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "12.3%", QueryParam: "trend_label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Minimal style (no icon) and icon-corner style variants.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Minimal — Open Cases",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Label: "Open Cases", Value: "142", Trend: ui.StatTrendUp, TrendLabel: "12.3%"})
							},
						},
						{
							Label: "Minimal — Pending Tasks",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Label: "Pending Tasks", Value: "38", Trend: ui.StatTrendDown, TrendLabel: "4.1%"})
							},
						},
						{
							Label: "Minimal — Clients",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Label: "Clients", Value: "89", Trend: ui.StatTrendUp, TrendLabel: "7.8%"})
							},
						},
						{
							Label: "Minimal — Avg. Case Days",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Label: "Avg. Case Days", Value: "24", Trend: ui.StatTrendUp, TrendLabel: "2.5%"})
							},
						},
						{
							Label: "Icon corner — Open Cases",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Value: "142", Label: "Open Cases", Trend: ui.StatTrendUp, TrendLabel: "14.6%", Icon: "lucide--briefcase", IconColor: "text-primary"})
							},
						},
						{
							Label: "Icon corner — Pending Tasks",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Value: "38", Label: "Pending Tasks", Trend: ui.StatTrendDown, TrendLabel: "4.1%", Icon: "lucide--check-square", IconColor: "text-warning"})
							},
						},
						{
							Label: "Icon corner — Active Clients",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Value: "89", Label: "Active Clients", Trend: ui.StatTrendUp, TrendLabel: "7.8%", Icon: "lucide--users", IconColor: "text-success"})
							},
						},
						{
							Label: "Icon corner — Revenue (MTD)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardMinimalWithBoundary(ui.StatCardMinimalItem{Value: "$48K", Label: "Revenue (MTD)", Trend: ui.StatTrendUp, TrendLabel: "9.2%", Icon: "lucide--dollar-sign", IconColor: "text-secondary"})
							},
						},
					},
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Skeleton loaders for text, avatar, and card patterns.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Text line",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.SkeletonWithBoundary("h-4 w-48")
							},
						},
						{
							Label: "Avatar circle",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.SkeletonWithBoundary("size-16 rounded-full")
							},
						},
						{
							Label: "Card",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.SkeletonWithBoundary("h-32 w-full")
							},
						},
					},
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Section headers with different titles.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Account Settings",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.SectionHeaderWithBoundary("Account Settings")
							},
						},
						{
							Label: "Notifications",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.SectionHeaderWithBoundary("Notifications")
							},
						},
						{
							Label: "Billing & Payments",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.SectionHeaderWithBoundary("Billing & Payments")
							},
						},
					},
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
					Name:        "Interactive",
					Description: "A fixed permission-denied placeholder with no configurable props.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.NoPermissionsWithBoundary()
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "The no-permissions placeholder.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.NoPermissionsWithBoundary()
					},
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
					Name:        "Interactive",
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
						return ui.NotificationPanelWithBoundary(items, 2, "#")
					},
				},
		{
			Name:        "Examples",
			Description: "Panel with unread and read notifications.",
			SubExamples: []galleryruntime.GallerySubExample{
				{
					Label: "Data-driven",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []ui.NotificationItem{
							{Title: "Case assigned", Body: "Johnson v. Smith was assigned to you.", Time: "2 min ago", Unread: true},
							{Title: "Document uploaded", Body: "A new document was added to Case #142.", Time: "10 min ago", Unread: true},
							{Title: "Reminder", Body: "Court date in 3 days.", Time: "1 hour ago", Unread: false},
							{Title: "Workflow complete", Body: "Document review workflow finished.", Time: "2 hours ago", Unread: false},
						}
						return ui.NotificationPanelWithBoundary(items, 2, "#")
					},
				},
				{
					Label: "Composed via children",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							return ui.NotificationPanel(2, "#").Render(templ.WithChildren(ctx, seq(
								ui.NotificationRowWithBoundary(ui.NotificationItem{
									IconClass:     "bg-primary/10",
									IconTextClass: "text-primary",
									IconName:      "lucide--briefcase",
									Title:         "New case assigned",
									Body:          "Johnson v. Smith was assigned to you.",
									Time:          "2 min ago",
									Unread:        true,
								}),
								ui.NotificationRowWithBoundary(ui.NotificationItem{
									IconClass:     "bg-warning/10",
									IconTextClass: "text-warning",
									IconName:      "lucide--check-square",
									Title:         "Task deadline tomorrow",
									Body:          "File motion for Johnson v. Smith due soon.",
									Time:          "1 hour ago",
									Unread:        true,
								}),
								ui.NotificationRowWithBoundary(ui.NotificationItem{
									IconClass:     "bg-success/10",
									IconTextClass: "text-success",
									IconName:      "lucide--user",
									Title:         "Client signed in",
									Body:          "Alice Johnson accessed the client portal.",
									Time:          "Yesterday",
									Unread:        false,
								}),
							)), w)
						})
					},
				},
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
					Name:        "Interactive",
					Description: "FAB appears bottom-right. Click it to expand sub-actions.",
					FrameHeight: "300px",
					RenderFunc: func(_ url.Values) templ.Component {
						actions := []ui.FABAction{
							{Label: "New Case", Icon: "lucide--briefcase"},
							{Label: "Upload Doc", Icon: "lucide--file-up"},
							{Label: "Add Task", Icon: "lucide--check-square"},
						}
						return ui.FABWithBoundary("lucide--plus", actions)
					},
				},
				{
					Name:        "Examples",
					Description: "FAB with 2 and 4 action items.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">2 actions</p>`); err != nil {
								return err
							}
							if err := ui.FAB("lucide--plus", []ui.FABAction{
								{Icon: "lucide--file-plus", Label: "New case", Href: "#"},
								{Icon: "lucide--upload", Label: "Upload", Href: "#"},
							}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">4 actions</p>`); err != nil {
								return err
							}
							if err := ui.FAB("lucide--plus", []ui.FABAction{
								{Icon: "lucide--file-plus", Label: "New case", Href: "#"},
								{Icon: "lucide--users", Label: "Add contact", Href: "#"},
								{Icon: "lucide--upload", Label: "Upload doc", Href: "#"},
								{Icon: "lucide--calendar", Label: "Schedule", Href: "#"},
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
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
					Name:        "Interactive",
					Description: "Compact title bar with inline breadcrumb trail.",
					RenderFunc: func(_ url.Values) templ.Component {
						steps := []nav.PageTitleStep{
							{Label: "", Href: "#", Icon: "lucide--home"},
							{Label: "Cases", Href: "#", Icon: "lucide--briefcase"},
							{Label: "New"},
						}
						return nav.PageTitleMinimalWithBoundary("Create New Case", steps)
					},
				},
				{
					Name:        "Examples",
					Description: "Page title with different step counts.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="space-y-4">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-1 font-semibold uppercase px-4 pt-4">Single level</p>`); err != nil {
								return err
							}
							if err := nav.PageTitleMinimal("Dashboard", nil).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-1 font-semibold uppercase px-4">With breadcrumbs</p>`); err != nil {
								return err
							}
							if err := nav.PageTitleMinimal("Edit Record", []nav.PageTitleStep{
								{Label: "Cases", Href: "#"},
								{Label: "Johnson v. Smith", Href: "#"},
							}).Render(ctx, w); err != nil {
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
			Slug:        "page-title-editor",
			Name:        "Page Title — Editor",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Page Title",
			Description: "Full page title with DaisyUI breadcrumbs, subtitle meta, and action buttons.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
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
						return nav.PageTitleEditorWithBoundary(steps, "Johnson v. Smith", "Type: Civil Litigation", actions)
					},
				},
				{
					Name:        "Examples",
					Description: "Editor title bar with actions.",
					RenderFunc: func(_ url.Values) templ.Component {
						return nav.PageTitleEditor(
							[]nav.BreadcrumbStep{{Label: "Cases", URL: "#"}, {Label: "Johnson v. Smith", URL: "#"}},
							"Edit Document",
							"Last edited 2 minutes ago",
							[]nav.PageTitleEditorAction{
								{Label: "Preview", Icon: "lucide--eye", Href: "#"},
								{Label: "Save", Icon: "lucide--save", Href: "#", Class: "btn-primary"},
							},
						)
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Search dropdown with multiple result sections.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.SearchDropdown("Search cases, contacts...", []ui.SearchDropdownSection{
							{
								Title: "Cases",
								Items: []ui.SearchDropdownItem{
									{Label: "Johnson v. Smith", Href: "#", Icon: "lucide--folder"},
									{Label: "Garcia Estate", Href: "#", Icon: "lucide--folder"},
								},
							},
							{
								Title: "Contacts",
								Items: []ui.SearchDropdownItem{
									{Label: "Alice Johnson", Href: "#", Icon: "lucide--user"},
									{Label: "Bob Martinez", Href: "#", Icon: "lucide--user"},
								},
							},
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
			Description: "FilterCard wraps filter inputs in a card with Filter/Clear buttons. Set Inline=true for a bare horizontal bar (used above tables).",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "FilterCard (card style) and inline variant with sample search and status inputs.",
					RenderFunc: func(_ url.Values) templ.Component {
						filterInputs := seq(
							form.SearchInputWithBoundary("q", "", "Search cases…", "", "#", ""),
							form.SelectInputWithBoundary("status", "Status", "", [][2]string{
								{"", "All statuses"},
								{"active", "Active"},
								{"pending", "Pending"},
								{"closed", "Closed"},
							}, "", false),
						)
						compactInputs := seq(
							form.SearchInputWithBoundary("q", "", "Search…", "", "#", ""),
							form.SelectInputWithBoundary("status", "", "", [][2]string{
								{"", "All statuses"},
								{"active", "Active"},
								{"closed", "Closed"},
							}, "", false),
						)
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6"><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Card style</p>`); err != nil {
								return err
							}
							if err := withChildren(ui.FilterCard(ui.FilterCardProps{Action: "#"}), filterInputs).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Inline (Inline: true)</p>`); err != nil {
								return err
							}
							if err := withChildren(ui.FilterCard(ui.FilterCardProps{Action: "#", Inline: true}), compactInputs).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Filter card and compact bar variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						inputs := seq(
							form.SearchInput("q", "", "Search cases…", "", "#", ""),
							form.SelectInput("status", "Status", "", [][2]string{
								{"", "All statuses"},
								{"active", "Active"},
								{"closed", "Closed"},
							}, "", false),
						)
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Card style</p>`); err != nil {
								return err
							}
							if err := withChildren(ui.FilterCard(ui.FilterCardProps{Action: "#"}), inputs).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Inline bar</p>`); err != nil {
								return err
							}
							if err := withChildren(ui.FilterCard(ui.FilterCardProps{Action: "#", Inline: true}), inputs).Render(ctx, w); err != nil {
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Fieldsets grouping related form inputs.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if err := ui.FieldsetWithBoundary("Personal Info", seq(
								form.TextInput("name", "Full Name", "", "", true),
								form.TextInput("email", "Email", "", "", true),
							)).Render(ctx, w); err != nil {
								return err
							}
							if err := ui.FieldsetWithBoundary("Preferences", seq(
								form.CheckboxInput("newsletter", "Subscribe to newsletter", false, ""),
								form.CheckboxInput("updates", "Receive product updates", true, ""),
							)).Render(ctx, w); err != nil {
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
					Name:        "Interactive",
					Description: "Prompt bar with attach, image, voice, and token counter.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder:      "Describe what you want to generate or ask a question…",
							ShowTokenCounter: true,
							TokenCurrent:     88,
							TokenMax:         100,
							ShowAttach:       true,
							ShowImage:        true,
							ShowVoice:        true,
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Prompt bar with different placeholder text.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Default</p>`); err != nil {
								return err
							}
							if err := form.PromptBar(form.PromptBarProps{Placeholder: "Ask anything..."}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Legal context</p>`); err != nil {
								return err
							}
							if err := form.PromptBar(form.PromptBarProps{Placeholder: "Search case law, statutes, or documents..."}).Render(ctx, w); err != nil {
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
			Slug:        "prompt-bar-action",
			Name:        "Prompt Bar — Action",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "AI prompt input with quick-action buttons (Add File, Deep Thinking, Browsing).",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Prompt bar with labelled quick-action buttons.",
					RenderFunc: func(_ url.Values) templ.Component {
						actions := []form.PromptBarActionItem{
							{Label: "Add File", Icon: "lucide--circle-plus"},
							{Label: "Deep Think", Icon: "lucide--lightbulb"},
						}
						return form.PromptBarActionWithBoundary("Type your request or attach files to get started…", actions)
					},
				},
				{
					Name:        "Examples",
					Description: "Prompt bar with action toolbar variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">With 2 actions</p>`); err != nil {
								return err
							}
							if err := form.PromptBarAction("Ask a question...", []form.PromptBarActionItem{
								{Icon: "lucide--paperclip", Label: "Attach"},
								{Icon: "lucide--mic", Label: "Record"},
							}, false, nil).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">With 3 actions</p>`); err != nil {
								return err
							}
							if err := form.PromptBarAction("Type a message...", []form.PromptBarActionItem{
								{Icon: "lucide--image", Label: "Image"},
								{Icon: "lucide--paperclip", Label: "Attach"},
								{Icon: "lucide--smile", Label: "Emoji"},
							}, false, nil).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Loading state</p>`); err != nil {
								return err
							}
							if err := form.PromptBarAction("Waiting for response...", []form.PromptBarActionItem{
								{Icon: "lucide--circle-plus", Label: "Add File"},
								{Icon: "lucide--brain", Label: "Deep Thinking"},
								{Icon: "lucide--globe", Label: "Browsing"},
							}, true, nil).Render(ctx, w); err != nil {
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
			Slug:        "prompt-bar-model",
			Name:        "Prompt Bar — Model Selector",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "AI prompt bar with a model-selector dropdown and optional info banner.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Prompt bar with model selector and token counter.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarModelSelectorWithBoundary(form.PromptBarModelSelectorProps{
							Placeholder:      "Type your request and select a model to process it",
							SelectedModel:    "GPT-4o",
							ShowTokenCounter: true,
							TokenCurrent:     42,
							TokenMax:         128,
							Models: []form.PromptBarModelSelectorItem{
								{Label: "GPT-4o", Value: "gpt-4o", Icon: "lucide--cpu"},
								{Label: "GPT-4o Mini", Value: "gpt-4o-mini", Icon: "lucide--cpu"},
								{Label: "Claude 3.5 Sonnet", Value: "claude-3-5-sonnet", Icon: "lucide--cpu"},
								{Label: "Gemini 1.5 Pro", Value: "gemini-1.5-pro", Icon: "lucide--cpu"},
							},
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Model selector with and without info banner.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Plain</p>`); err != nil {
								return err
							}
							if err := form.PromptBarModelSelector(form.PromptBarModelSelectorProps{
								Placeholder:   "Type your request and select a model…",
								SelectedModel: "Claude 3.5 Sonnet",
								Models: []form.PromptBarModelSelectorItem{
									{Label: "Claude 3.5 Sonnet", Value: "claude-3-5-sonnet"},
									{Label: "GPT-4o", Value: "gpt-4o"},
								},
							}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">With info banner</p>`); err != nil {
								return err
							}
							if err := form.PromptBarModelSelector(form.PromptBarModelSelectorProps{
								Placeholder:   "Type your request and select a model…",
								SelectedModel: "GPT-4o",
								InfoBanner:    "This tool lets you choose and prompt different models.",
								Models: []form.PromptBarModelSelectorItem{
									{Label: "GPT-4o", Value: "gpt-4o"},
									{Label: "GPT-4o Mini", Value: "gpt-4o-mini"},
								},
							}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Loading</p>`); err != nil {
								return err
							}
							if err := form.PromptBarModelSelector(form.PromptBarModelSelectorProps{
								Placeholder:   "Waiting for response…",
								SelectedModel: "GPT-4o",
								Loading:       true,
							}).Render(ctx, w); err != nil {
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
			Slug:        "prompt-bar-ability",
			Name:        "Prompt Bar — Ability",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "AI prompt bar with an ability-selector dropdown and a token budget header strip.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Prompt bar with ability selector and token counter header.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarAbilityWithBoundary(form.PromptBarAbilityProps{
							Placeholder:     "Type your request and choose an ability to enhance it",
							SelectedAbility: "Ability",
							SavedLabel:      "Saved to Library",
							TokensLeft:      58,
							Abilities: []form.PromptBarAbilityItem{
								{Label: "Write", Icon: "lucide--pen-line"},
								{Label: "Summarise", Icon: "lucide--list"},
								{Label: "Translate", Icon: "lucide--languages"},
								{Label: "Code Review", Icon: "lucide--code"},
							},
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Ability selector with and without header strip.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">With header strip</p>`); err != nil {
								return err
							}
							if err := form.PromptBarAbility(form.PromptBarAbilityProps{
								Placeholder:     "Write your prompt or select one from My Prompts",
								SelectedAbility: "Write",
								SavedLabel:      "Saved to Library",
								TokensLeft:      58,
								Abilities: []form.PromptBarAbilityItem{
									{Label: "Write", Icon: "lucide--pen-line"},
									{Label: "Summarise", Icon: "lucide--list"},
								},
							}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Without header strip</p>`); err != nil {
								return err
							}
							if err := form.PromptBarAbility(form.PromptBarAbilityProps{
								Placeholder:     "Write your prompt…",
								SelectedAbility: "Ability",
								Abilities: []form.PromptBarAbilityItem{
									{Label: "Translate", Icon: "lucide--languages"},
									{Label: "Code Review", Icon: "lucide--code"},
								},
							}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
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
					Name:        "Interactive",
					Description: "Primary→secondary, success→info, and warning→error gradient examples.",
					RenderFunc: func(_ url.Values) templ.Component {
						return devmode.ComponentBoundary("GradientText", rawHTML(`<div class="p-6 space-y-6">
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
</div>`))
					},
				},
				{
					Name:        "Examples",
					Description: "Gradient text in multiple color directions.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 space-y-4">
	<p class="text-3xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">Primary to Secondary</p>
	<p class="text-3xl font-bold bg-gradient-to-r from-accent to-primary bg-clip-text text-transparent">Accent to Primary</p>
	<p class="text-3xl font-bold bg-gradient-to-br from-success to-info bg-clip-text text-transparent">Success to Info</p>
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
					Name:        "Interactive",
					Description: "Cards and buttons with colored drop shadows.",
					RenderFunc: func(_ url.Values) templ.Component {
						return devmode.ComponentBoundary("ColoredShadows", rawHTML(`<div class="p-8 space-y-6">
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
</div>`))
					},
				},
				{
					Name:        "Examples",
					Description: "Cards with colored drop shadows.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-8 flex flex-wrap gap-8">
	<div class="card bg-primary text-primary-content w-32 h-20 shadow-lg shadow-primary/50 flex items-center justify-center font-semibold">Primary</div>
	<div class="card bg-secondary text-secondary-content w-32 h-20 shadow-lg shadow-secondary/50 flex items-center justify-center font-semibold">Secondary</div>
	<div class="card bg-accent text-accent-content w-32 h-20 shadow-lg shadow-accent/50 flex items-center justify-center font-semibold">Accent</div>
	<div class="card bg-success text-success-content w-32 h-20 shadow-lg shadow-success/50 flex items-center justify-center font-semibold">Success</div>
	<div class="card bg-error text-error-content w-32 h-20 shadow-lg shadow-error/50 flex items-center justify-center font-semibold">Error</div>
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
					Name:        "Interactive",
					Description: "Heading levels, body, muted, overline, and link styles.",
					RenderFunc: func(_ url.Values) templ.Component {
						return devmode.ComponentBoundary("Typography", rawHTML(`<div class="p-6 space-y-3">
  <h1 class="text-3xl font-bold text-base-content">Heading 1</h1>
  <h2 class="text-2xl font-semibold text-base-content">Heading 2</h2>
  <h3 class="text-xl font-semibold text-base-content">Heading 3</h3>
  <h4 class="text-base font-semibold text-base-content">Heading 4</h4>
  <p class="text-base text-base-content/80">Body text — regular paragraph content used in cards and detail views.</p>
  <p class="text-sm text-base-content/60">Small / muted text — used for labels, hints, and secondary information.</p>
  <p class="text-xs text-base-content/50 uppercase tracking-wide font-semibold">Overline / label text</p>
  <a href="#" class="link link-primary text-sm">Link text</a>
</div>`))
					},
				},
				{
					Name:        "Examples",
					Description: "Typography hierarchy showcase.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 space-y-3">
	<h1 class="text-4xl font-bold">Heading 1</h1>
	<h2 class="text-3xl font-semibold">Heading 2</h2>
	<h3 class="text-2xl font-semibold">Heading 3</h3>
	<h4 class="text-xl font-medium">Heading 4</h4>
	<p class="text-base">Body text — the quick brown fox jumps over the lazy dog.</p>
	<p class="text-sm text-base-content/70">Small muted text for secondary information.</p>
	<p class="text-xs text-base-content/50 uppercase tracking-wider font-semibold">Label / Caption</p>
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
					Name:        "Interactive",
					Description: "Size scale from xs to 4xl and all font weights.",
					RenderFunc: func(_ url.Values) templ.Component {
						return devmode.ComponentBoundary("TypographyScale", rawHTML(`<div class="space-y-6 p-6">
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
</div>`))
					},
				},
				{
					Name:        "Examples",
					Description: "Full Tailwind font-size scale.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 space-y-2">
	<p class="text-xs">xs — 12px</p>
	<p class="text-sm">sm — 14px</p>
	<p class="text-base">base — 16px</p>
	<p class="text-lg">lg — 18px</p>
	<p class="text-xl">xl — 20px</p>
	<p class="text-2xl">2xl — 24px</p>
	<p class="text-3xl">3xl — 30px</p>
	<p class="text-4xl">4xl — 36px</p>
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
					Name:        "Interactive",
					Description: "Box shadow, inset shadow, and text shadow scales.",
					RenderFunc: func(_ url.Values) templ.Component {
						return devmode.ComponentBoundary("ShadowScale", rawHTML(`<div class="space-y-6 p-6">
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
</div>`))
					},
				},
				{
					Name:        "Examples",
					Description: "Tailwind shadow scale from sm to 2xl.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-8 flex flex-wrap gap-6 items-end">
	<div class="card bg-base-100 w-20 h-20 shadow-sm flex items-center justify-center text-xs font-medium">sm</div>
	<div class="card bg-base-100 w-20 h-20 shadow flex items-center justify-center text-xs font-medium">default</div>
	<div class="card bg-base-100 w-20 h-20 shadow-md flex items-center justify-center text-xs font-medium">md</div>
	<div class="card bg-base-100 w-20 h-20 shadow-lg flex items-center justify-center text-xs font-medium">lg</div>
	<div class="card bg-base-100 w-20 h-20 shadow-xl flex items-center justify-center text-xs font-medium">xl</div>
	<div class="card bg-base-100 w-20 h-20 shadow-2xl flex items-center justify-center text-xs font-medium">2xl</div>
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
					Name:        "Interactive",
					Description: "Image filter utility classes applied to sample images.",
					RenderFunc: func(_ url.Values) templ.Component {
						return devmode.ComponentBoundary("CSSFilters", rawHTML(`<div class="p-6">
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
</div>`))
					},
				},
				{
					Name:        "Examples",
					Description: "CSS filter effects: blur, brightness, contrast.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 flex flex-wrap gap-6">
	<div class="flex flex-col items-center gap-2"><div class="w-20 h-20 rounded bg-primary"></div><span class="text-xs">Normal</span></div>
	<div class="flex flex-col items-center gap-2"><div class="w-20 h-20 rounded bg-primary blur-sm"></div><span class="text-xs">blur-sm</span></div>
	<div class="flex flex-col items-center gap-2"><div class="w-20 h-20 rounded bg-primary brightness-75"></div><span class="text-xs">brightness-75</span></div>
	<div class="flex flex-col items-center gap-2"><div class="w-20 h-20 rounded bg-primary brightness-125"></div><span class="text-xs">brightness-125</span></div>
	<div class="flex flex-col items-center gap-2"><div class="w-20 h-20 rounded bg-primary contrast-50"></div><span class="text-xs">contrast-50</span></div>
	<div class="flex flex-col items-center gap-2"><div class="w-20 h-20 rounded bg-primary saturate-50"></div><span class="text-xs">saturate-50</span></div>
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
					Name:        "Interactive",
					Description: "Footer with copyright text only.",
					RenderFunc: func(_ url.Values) templ.Component {
						return nav.FooterMinimalWithBoundary("© 2025 LegalPlant. All rights reserved.", nil)
					},
				},
				{
					Name:        "Examples",
					Description: "Footer with various link configurations.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="space-y-4">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-1 font-semibold uppercase px-4 pt-4">Minimal (no links)</p>`); err != nil {
								return err
							}
							if err := nav.FooterMinimal("© 2026 Acme Corp.", nil).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-1 font-semibold uppercase px-4">With links</p>`); err != nil {
								return err
							}
							if err := nav.FooterMinimal("© 2026 Acme Corp.", []nav.FooterLink{
								{Label: "Privacy Policy", Href: "#"},
								{Label: "Terms of Service", Href: "#"},
								{Label: "Contact", Href: "#"},
							}).Render(ctx, w); err != nil {
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
					Name:        "Interactive",
					Description: "Profile dropdown with avatar, user info, menu items, and sign-out.",
					RenderFunc: func(_ url.Values) templ.Component {
						items := []nav.ProfileMenuItem{
							{Label: "Profile", Href: "#", Icon: "lucide--user"},
							{Label: "Settings", Href: "#", Icon: "lucide--pencil"},
							{Label: "Notifications", Href: "#", Icon: "lucide--bell", Badge: 3},
						}
						return nav.ProfileMenuWithBoundary("Jane Doe", "jane@example.com", "JD", items, "#")
					},
				},
				{
					Name:        "Examples",
					Description: "Variants: with badge count, and minimal (sign-out only).",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Full featured</p>`); err != nil {
								return err
							}
							if err := nav.ProfileMenu("Jane Doe", "jane@example.com", "JD", []nav.ProfileMenuItem{
								{Label: "Profile", Href: "#", Icon: "lucide--user"},
								{Label: "Settings", Href: "#", Icon: "lucide--settings"},
								{Label: "Notifications", Href: "#", Icon: "lucide--bell", Badge: 5},
							}, "#").Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div>`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Sign-out only</p>`); err != nil {
								return err
							}
							if err := nav.ProfileMenu("Bob Smith", "bob@example.com", "BS", []nav.ProfileMenuItem{}, "#").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
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
					Name:        "Interactive",
					Description: "Simple spinner with default styling.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.InputSpinnerWithBoundary("spin1", 0, 0, 99, true, "btn-outline", "w-24")
					},
				},
				{
					Name:        "Examples",
					Description: "Spinner variants: default, bounded, and no-display-bounds.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Default (0–99)</p>`); err != nil {
								return err
							}
							if err := form.InputSpinner("ex-spin1", 0, 0, 99, true, "btn-outline", "w-24").Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">Bounded (0–10)</p>`); err != nil {
								return err
							}
							if err := form.InputSpinner("ex-spin2", 5, 0, 10, true, "btn-primary btn-sm", "w-20 input-sm").Render(ctx, w); err != nil {
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
					Name:        "Interactive",
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
								Content: seq(
									form.FormField(form.FormFieldProps{
										Type:        form.FieldText,
										Name:        "case-title",
										Label:       "Case title",
										Value:       "",
										Placeholder: "e.g. Johnson v. Smith",
										Required:    true,
									}),
									form.FormField(form.FormFieldProps{
										Type: form.FieldSelect,
										Name: "case-type",
										Label: "Case type",
										Options: []form.SelectOption{
											{Value: "civil", Label: "Civil"},
											{Value: "criminal", Label: "Criminal"},
											{Value: "family", Label: "Family"},
										},
									}),
								),
							},
							{
								Title: "Step 2 — Details",
								Content: seq(
									form.FormField(form.FormFieldProps{
										Type:        form.FieldTextarea,
										Name:        "description",
										Label:       "Description",
										Placeholder: "Brief description of the case…",
										Rows:        3,
									}),
									form.FormField(form.FormFieldProps{
										Type: form.FieldSelect,
										Name: "priority",
										Label: "Priority",
										Options: []form.SelectOption{
											{Value: "normal", Label: "Normal"},
											{Value: "high", Label: "High"},
											{Value: "urgent", Label: "Urgent"},
										},
									}),
								),
							},
							{
								Title: "Step 3 — Team",
								Content: form.FormField(form.FormFieldProps{
									Type: form.FieldSelect,
									Name: "lead-attorney",
									Label: "Lead attorney",
									Options: []form.SelectOption{
										{Value: "alice", Label: "Alice Johnson"},
										{Value: "bob", Label: "Bob Smith"},
										{Value: "carol", Label: "Carol White"},
									},
								}),
							},
							{
								Title: "Step 4 — Review",
								Content: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
									_, err := io.WriteString(w, `<p class="text-sm text-base-content/60 mb-4">Review the case details before submitting.</p>
<div class="space-y-2 text-sm">
  <div class="flex gap-2"><span class="text-base-content/60 w-24">Title:</span><span class="font-medium">Johnson v. Smith</span></div>
  <div class="flex gap-2"><span class="text-base-content/60 w-24">Type:</span><span class="font-medium">Civil</span></div>
  <div class="flex gap-2"><span class="text-base-content/60 w-24">Attorney:</span><span class="font-medium">Alice Johnson</span></div>
</div>`)
									return err
								}),
							},
						}
						return form.WizardStepperWithBoundary("wizard-demo", steps, panels)
					},
				},
				{
					Name:        "Examples",
					Description: "Wizard with 2 and 4 steps.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-3 font-semibold uppercase">2-step wizard</p>`); err != nil {
								return err
							}
							if err := form.WizardStepper("wiz-2", []form.WizardStep{{Label: "Details"}, {Label: "Confirm"}}, []form.WizardStepPanel{
								{Title: "Step 1 — Details", Content: form.FormField(form.FormFieldProps{
									Type:        form.FieldText,
									Name:        "details",
									Placeholder: "Enter details",
								})},
								{Title: "Step 2 — Confirm", Content: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
									_, err := io.WriteString(w, `<p class="text-sm text-base-content/70">Review and submit.</p>`)
									return err
								})},
							}).Render(ctx, w); err != nil {
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
			Slug:        "clipboard-copy",
			Name:        "Clipboard Copy",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Click-to-copy buttons with visual feedback. Uses vanilla JS navigator.clipboard — no library needed.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "Clipboard copy with text, URL, and mono code variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 max-w-lg">`); err != nil {
								return err
							}
							if err := form.ClipboardCopy([]form.ClipboardCopyItem{
								{ID: "ex-copy1", Label: "Case ID", Value: "CASE-2026-00142", ButtonLabel: "Copy"},
								{ID: "ex-copy2", Label: "Share link", Value: "https://app.example.com/cases/CASE-2026-00142", ButtonLabel: "Copy Link", ButtonClass: "btn-primary"},
								{ID: "ex-copy3", Label: "API key", Value: "sk-live-abc123xyz789", Mono: true},
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
					Name:        "Interactive",
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
							return withText("Save changes", ui.ButtonWithBoundary(href, variant, size, ui.ButtonStyleDefault, typ, shape, icon, loading, false))
						}
						if typ == ui.ButtonTypeLink {
							return withText("Go to dashboard", ui.ButtonWithBoundary(href, variant, size, ui.ButtonStyleDefault, typ, shape, icon, loading, false))
						}
						return ui.ButtonWithBoundary(href, variant, size, ui.ButtonStyleDefault, typ, shape, icon, loading, false)
					},
					Tokens: ButtonTokens(),
				},
				{
					Name:        "Examples",
					Description: "All variants, sizes, and special states.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Primary",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Secondary",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonSecondary, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Accent",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonAccent, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Neutral",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonNeutral, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Ghost",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonGhost, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Outline",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonOutline, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Error",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Save changes", ui.ButtonWithBoundary("#", ui.ButtonError, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false))
							},
						},
						{
							Label: "Loading",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", true, false)
							},
						},
						{
							Label: "Icon + Label",
							RenderFunc: func(_ url.Values) templ.Component {
								return withText("Star", ui.ButtonWithBoundary("#", ui.ButtonSecondary, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "lucide--star", false, false))
							},
						},
						{
							Label: "Icon Only (Square)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ButtonWithBoundary("#", ui.ButtonAccent, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeSquare, "lucide--pencil", false, false)
							},
						},
						{
							Label: "Icon Only (Circle)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ButtonWithBoundary("#", ui.ButtonNeutral, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeCircle, "lucide--plus", false, false)
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "All intents × styles, dot variant, icon variant.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Primary",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgePrimary, ui.BadgeStyleDefault, ui.BadgeSizeMD, false, "", "Primary")
							},
						},
						{
							Label: "Success",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgeSuccess, ui.BadgeStyleDefault, ui.BadgeSizeMD, false, "", "Success")
							},
						},
						{
							Label: "Warning",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgeWarning, ui.BadgeStyleDefault, ui.BadgeSizeMD, false, "", "Warning")
							},
						},
						{
							Label: "Error",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgeError, ui.BadgeStyleDefault, ui.BadgeSizeMD, false, "", "Error")
							},
						},
						{
							Label: "Outline",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgePrimary, ui.BadgeStyleOutline, ui.BadgeSizeMD, false, "", "Outline")
							},
						},
						{
							Label: "Soft",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgePrimary, ui.BadgeStyleSoft, ui.BadgeSizeMD, false, "", "Soft")
							},
						},
						{
							Label: "Dot (animated)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgeWarning, ui.BadgeStyleDefault, ui.BadgeSizeMD, true, "", "Pending")
							},
						},
						{
							Label: "With icon",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.BadgeWithBoundary(ui.BadgeError, ui.BadgeStyleDefault, ui.BadgeSizeMD, false, "lucide--circle-x", "Error")
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "All supported status strings grouped by intent.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Positive states",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(
									withText("active", ui.StatusBadgeWithBoundary("active")),
									withText("open", ui.StatusBadgeWithBoundary("open")),
									withText("completed", ui.StatusBadgeWithBoundary("completed")),
									withText("approved", ui.StatusBadgeWithBoundary("approved")),
								)
							},
						},
						{
							Label: "Negative states",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(
									withText("closed", ui.StatusBadgeWithBoundary("closed")),
									withText("rejected", ui.StatusBadgeWithBoundary("rejected")),
									withText("cancelled", ui.StatusBadgeWithBoundary("cancelled")),
									withText("deleted", ui.StatusBadgeWithBoundary("deleted")),
								)
							},
						},
						{
							Label: "Neutral states",
							RenderFunc: func(_ url.Values) templ.Component {
								return row(
									withText("pending", ui.StatusBadgeWithBoundary("pending")),
									withText("in_progress", ui.StatusBadgeWithBoundary("in_progress")),
									withText("review", ui.StatusBadgeWithBoundary("review")),
									withText("draft", ui.StatusBadgeWithBoundary("draft")),
									withText("unknown", ui.StatusBadgeWithBoundary("unknown")),
								)
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Cards with different titles and body content.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "User Profile",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(ui.CardWithBoundary("User Profile"), rawHTML(`<p class="text-sm text-base-content/70">Name: Alice Johnson<br/>Role: Admin</p>`))
							},
						},
						{
							Label: "Statistics",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(ui.CardWithBoundary("Statistics"), rawHTML(`<p class="text-sm text-base-content/70">Active cases: 12<br/>Closed this month: 4</p>`))
							},
						},
						{
							Label: "Recent Activity",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(ui.CardWithBoundary("Recent Activity"), rawHTML(`<p class="text-sm text-base-content/70">Document uploaded 2m ago<br/>Comment added 5m ago</p>`))
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "All four toast types.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Success",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ToastWithBoundary(ui.ToastSuccess, "Record saved successfully.")
							},
						},
						{
							Label: "Error",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ToastWithBoundary(ui.ToastError, "Something went wrong. Please try again.")
							},
						},
						{
							Label: "Warning",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ToastWithBoundary(ui.ToastWarning, "Your session will expire in 5 minutes.")
							},
						},
						{
							Label: "Info",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ToastWithBoundary(ui.ToastInfo, "A new version is available.")
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Pagination at different pages within a 10-page set.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Page 1 of 10",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.PaginationWithBoundary(1, 10, "#", "content")
							},
						},
						{
							Label: "Page 5 of 10",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.PaginationWithBoundary(5, 10, "#", "content")
							},
						},
						{
							Label: "Page 10 of 10",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.PaginationWithBoundary(10, 10, "#", "content")
							},
						},
						{
							Label: "Circle style",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.PaginationCircleWithBoundary(3, 7, "#", "content")
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Empty states for different contexts.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "No results",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.EmptyWithBoundary("lucide--search", "No results found", "Try adjusting your search or filters.")
							},
						},
						{
							Label: "No cases",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.EmptyWithBoundary("lucide--folder-open", "No cases yet", "Create your first case to get started.")
							},
						},
						{
							Label: "No notifications",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.EmptyWithBoundary("lucide--bell-off", "No notifications", "You're all caught up!")
							},
						},
					},
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
					Name:        "Interactive",
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
				{
					Name:        "Examples",
					Description: "All three loader variants side by side.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Centered",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.LoaderWithBoundary(ui.LoaderCentered)
							},
						},
						{
							Label: "Inline",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.LoaderWithBoundary(ui.LoaderInline)
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Stat cards in different icon color schemes.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Active Cases",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardWithBoundary(ui.StatCardProps{Label: "Active Cases", Value: "42", Icon: "lucide--folder-open", IconColor: "bg-primary/10 text-primary"})
							},
						},
						{
							Label: "Contacts",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardWithBoundary(ui.StatCardProps{Label: "Contacts", Value: "128", Icon: "lucide--users", IconColor: "bg-secondary/10 text-secondary"})
							},
						},
						{
							Label: "Documents",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardWithBoundary(ui.StatCardProps{Label: "Documents", Value: "315", Icon: "lucide--file-text", IconColor: "bg-success/10 text-success"})
							},
						},
						{
							Label: "Overdue",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.StatCardWithBoundary(ui.StatCardProps{Label: "Overdue", Value: "7", Icon: "lucide--alert-circle", IconColor: "bg-error/10 text-error"})
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
				{
					Name:        "Examples",
					Description: "Various item configurations including dangerous actions and many items.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Basic 3 items",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ActionMenuWithBoundary([]ui.ActionMenuItem{
									{Label: "Edit", Icon: "lucide--pencil", HXGet: "#"},
									{Label: "Duplicate", Icon: "lucide--copy", HXGet: "#"},
									{Label: "Delete", Icon: "lucide--trash-2", HXGet: "#", Danger: true},
								})
							},
						},
						{
							Label: "View only",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ActionMenuWithBoundary([]ui.ActionMenuItem{
									{Label: "View details", Icon: "lucide--eye", HXGet: "#"},
									{Label: "Download", Icon: "lucide--download", HXGet: "#"},
									{Label: "Share", Icon: "lucide--share-2", HXGet: "#"},
								})
							},
						},
						{
							Label: "Many items",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.ActionMenuWithBoundary([]ui.ActionMenuItem{
									{Label: "Edit", Icon: "lucide--pencil", HXGet: "#"},
									{Label: "Rename", Icon: "lucide--text-cursor", HXGet: "#"},
									{Label: "Move", Icon: "lucide--folder-input", HXGet: "#"},
									{Label: "Copy link", Icon: "lucide--link", HXGet: "#"},
									{Label: "Archive", Icon: "lucide--archive", HXGet: "#"},
									{Label: "Delete", Icon: "lucide--trash-2", HXGet: "#", Danger: true},
								})
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "All sizes, initials fallback, icon placeholder, and image variants.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Initials (two-word)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AvatarWithBoundary("Bob Carter", "", "", ui.AvatarMD)
							},
						},
						{
							Label: "Icon placeholder",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AvatarWithBoundary("", "", "lucide--building-2", ui.AvatarMD)
							},
						},
						{
							Label: "With image",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AvatarWithBoundary("User", "https://i.pravatar.cc/150?img=3", "", ui.AvatarMD)
							},
						},
						{
							Label: "XS size",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AvatarWithBoundary("Jane Smith", "", "", ui.AvatarXS)
							},
						},
						{
							Label: "SM size",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AvatarWithBoundary("Jane Smith", "", "", ui.AvatarSM)
							},
						},
						{
							Label: "LG size",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.AvatarWithBoundary("Jane Smith", "", "", ui.AvatarLG)
							},
						},
					},
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
					Description: "Live label, value, required, and error controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Email address"
						}
						value := params.Get("value")
						errMsg := params.Get("errMsg")
						required := params.Get("required") == "true"
						return form.TextInputWithBoundary("email", label, value, errMsg, required)
					},
					Tokens: append(TextInputTokens(),
						galleryruntime.DesignToken{
							Label:      "Value",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "",
							QueryParam: "value",
						},
						galleryruntime.DesignToken{
							Label:      "Error Message",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "",
							QueryParam: "errMsg",
						},
					),
				},
				{
					Name:        "Examples",
					Description: "Default, pre-filled value, required, and error states.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Default",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.TextInputWithBoundary("name", "Full Name", "", "", false)
							},
						},
						{
							Label: "Pre-filled + required",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.TextInputWithBoundary("email", "Email", "jane@example.com", "", true)
							},
						},
						{
							Label: "Error state",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.TextInputWithBoundary("err-field", "Username", "taken_user", "Username is already taken.", false)
							},
						},
					},
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
					Description: "Live label, rows, required, and error controls.",
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
						errMsg := params.Get("errMsg")
						required := params.Get("required") == "true"
						return form.TextareaInputWithBoundary("description", label, "", errMsg, rows, required)
					},
					Tokens: append(TextareaInputTokens(),
						galleryruntime.DesignToken{
							Label:      "Error Message",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "",
							QueryParam: "errMsg",
						},
					),
				},
				{
					Name:        "Examples",
					Description: "Default, required, and error states.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Default",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.TextareaInputWithBoundary("bio", "Bio", "", "", 3, false)
							},
						},
						{
							Label: "Required",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.TextareaInputWithBoundary("notes", "Notes", "", "", 3, true)
							},
						},
						{
							Label: "Error state",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.TextareaInputWithBoundary("err-area", "Summary", "Too short", "Summary must be at least 50 characters.", 3, false)
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Checkboxes: unchecked, checked, and with error.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Unchecked",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CheckboxInputWithBoundary("opt1", "Enable notifications", false, "")
							},
						},
						{
							Label: "Checked",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CheckboxInputWithBoundary("opt2", "I agree to the terms", true, "")
							},
						},
						{
							Label: "Error state",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CheckboxInputWithBoundary("opt3", "Subscribe to newsletter", false, "This field is required.")
							},
						},
					},
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
					Description: "Live label, selected value, required, and error controls.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Country"
						}
						selected := params.Get("selected")
						errMsg := params.Get("errMsg")
						required := params.Get("required") == "true"
						return form.SelectInputWithBoundary("country", label, selected, [][2]string{
							{"us", "United States"},
							{"gb", "United Kingdom"},
							{"ca", "Canada"},
							{"au", "Australia"},
						}, errMsg, required)
					},
					Tokens: append(SelectInputTokens(),
						galleryruntime.DesignToken{
							Label:      "Error Message",
							Group:      "Component",
							Type:       galleryruntime.TokenTypeText,
							Default:    "",
							QueryParam: "errMsg",
						},
					),
				},
				{
					Name:        "Examples",
					Description: "Default, pre-selected, required, and error states.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Default",
							RenderFunc: func(_ url.Values) templ.Component {
								opts := [][2]string{{"us", "United States"}, {"gb", "United Kingdom"}, {"ca", "Canada"}, {"au", "Australia"}}
								return form.SelectInputWithBoundary("country1", "Country", "", opts, "", false)
							},
						},
						{
							Label: "Pre-selected + required",
							RenderFunc: func(_ url.Values) templ.Component {
								opts := [][2]string{{"us", "United States"}, {"gb", "United Kingdom"}, {"ca", "Canada"}, {"au", "Australia"}}
								return form.SelectInputWithBoundary("country2", "Country", "gb", opts, "", true)
							},
						},
						{
							Label: "Error state",
							RenderFunc: func(_ url.Values) templ.Component {
								opts := [][2]string{{"us", "United States"}, {"gb", "United Kingdom"}, {"ca", "Canada"}, {"au", "Australia"}}
								return form.SelectInputWithBoundary("country3", "Country", "", opts, "Please select a country.", false)
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Range sliders with different colors and values.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Primary",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.RangeInputWithBoundary("vol", "Volume", 70, 0, 100, 1, "range-primary")
							},
						},
						{
							Label: "Secondary",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.RangeInputWithBoundary("bright", "Brightness", 50, 0, 100, 10, "range-secondary")
							},
						},
						{
							Label: "Accent",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.RangeInputWithBoundary("speed", "Speed", 30, 0, 100, 5, "range-accent")
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Form fields of each type.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Text",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormFieldWithBoundary(form.FormFieldProps{Type: form.FieldText, Name: "name", Label: "Full Name", Placeholder: "Jane Smith", Required: true})
							},
						},
						{
							Label: "Email",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormFieldWithBoundary(form.FormFieldProps{Type: form.FieldEmail, Name: "email", Label: "Email", Placeholder: "jane@example.com"})
							},
						},
						{
							Label: "Textarea",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormFieldWithBoundary(form.FormFieldProps{Type: form.FieldTextarea, Name: "bio", Label: "Bio", Placeholder: "Tell us about yourself..."})
							},
						},
						{
							Label: "Select",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormFieldWithBoundary(form.FormFieldProps{Type: form.FieldSelect, Name: "role", Label: "Role", Options: []form.SelectOption{{Value: "admin", Label: "Admin"}, {Value: "member", Label: "Member"}, {Value: "viewer", Label: "Viewer"}}})
							},
						},
						{
							Label: "Checkbox",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormFieldWithBoundary(form.FormFieldProps{Type: form.FieldCheckbox, Name: "agree", Label: "I agree to the terms"})
							},
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
						return form.SearchInputWithBoundary("q", value, placeholder, "", "", "")
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
				{
					Name:        "Examples",
					Description: "Search inputs with different placeholders.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Default",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.SearchInputWithBoundary("q1", "", "Search...", "", "", "")
							},
						},
						{
							Label: "Pre-filled value",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.SearchInputWithBoundary("q2", "Johnson v. Smith", "Search cases...", "", "", "")
							},
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
				{
					Name:        "Examples",
					Description: "Top bars with different section titles.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Dashboard",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.TopBarWithBoundary("Dashboard")
							},
						},
						{
							Label: "Cases",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.TopBarWithBoundary("Cases")
							},
						},
						{
							Label: "Contacts",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.TopBarWithBoundary("Contacts")
							},
						},
						{
							Label: "Settings",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.TopBarWithBoundary("Settings")
							},
						},
					},
				},
			},
		},

		// nav.TabMenu / nav.SimpleTabs (unified)
		{
			Slug:        "tab-menu-real",
			Name:        "Tabs",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Tabs",
			Description: "Tab strip component. Full HTMX variant for page-level navigation; pass target=\"-\" for an in-panel lifted strip without HTMX.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "HTMX tab strip with configurable labels.",
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
				{
					Name:        "Examples",
					Description: "HTMX full-page strip and lifted in-panel strip.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "HTMX (TabMenu)",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.TabMenuWithBoundary([]nav.Tab{
									{Label: "Overview", Href: "#", Active: true},
									{Label: "Activity", Href: "#"},
									{Label: "Settings", Href: "#"},
								})
							},
						},
						{
							Label: "Lifted in-panel (SimpleTabs)",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.SimpleTabsWithBoundary([]nav.Tab{
									{Label: "All", Href: "#", Active: true},
									{Label: "Open", Href: "#"},
									{Label: "Closed", Href: "#"},
								})
							},
						},
					},
				},
			},
		},

		// ui.Tabs (general-purpose DaisyUI tabs with content switching)
		{
			Slug:        "tabs",
			Name:        "Tabs (DaisyUI)",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Tabs",
			Description: "General-purpose DaisyUI tabs with radio-based content switching. Supports lift/border/box styles, all sizes, bottom placement, icons, disabled tabs, and custom colors.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Tabs with configurable style, size, and placement.",
					RenderFunc: func(params url.Values) templ.Component {
						style := parseTabsStyle(params.Get("style"))
						size := parseTabsSize(params.Get("size"))
						bottom := params.Get("bottom") == "true"
						props := ui.TabsProps{
							Style:  style,
							Size:   size,
							Bottom: bottom,
							Name:   "demo-tabs",
							Items: []ui.TabItem{
								{Label: "Overview", Active: true, Content: tabsContent("Overview tab content — showing the overview details for this item.")},
								{Label: "Activity", Content: tabsContent("Activity tab content — recent activity and event history.")},
								{Label: "Settings", Content: tabsContent("Settings tab content — configuration and preferences.")},
							},
						}
						return ui.TabsWithBoundary(props)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Style", Group: "Appearance", Type: galleryruntime.TokenTypeSelect, Default: "lift", QueryParam: "style", Options: []galleryruntime.TokenOption{
							{Value: "lift", Label: "Lift"},
							{Value: "border", Label: "Border"},
							{Value: "box", Label: "Box"},
						}},
						{Label: "Size", Group: "Appearance", Type: galleryruntime.TokenTypeSelect, Default: "md", QueryParam: "size", Options: []galleryruntime.TokenOption{
							{Value: "xs", Label: "XS"},
							{Value: "sm", Label: "SM"},
							{Value: "md", Label: "MD"},
							{Value: "lg", Label: "LG"},
							{Value: "xl", Label: "XL"},
						}},
						{Label: "Bottom", Group: "Placement", Type: galleryruntime.TokenTypeSelect, Default: "false", QueryParam: "bottom", Options: []galleryruntime.TokenOption{
							{Value: "false", Label: "Top"},
							{Value: "true", Label: "Bottom"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "All DaisyUI tab variants with content.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Lift style",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsLift,
									Name:  "ex-lift",
									Items: []ui.TabItem{
										{Label: "Tab 1", Active: true, Content: tabsContent("Tab content 1")},
										{Label: "Tab 2", Content: tabsContent("Tab content 2")},
										{Label: "Tab 3", Content: tabsContent("Tab content 3")},
									},
								})
							},
						},
						{
							Label: "Border style",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsBorder,
									Name:  "ex-border",
									Items: []ui.TabItem{
										{Label: "Tab 1", Active: true, Content: tabsContent("Tab content 1")},
										{Label: "Tab 2", Content: tabsContent("Tab content 2")},
										{Label: "Tab 3", Content: tabsContent("Tab content 3")},
									},
								})
							},
						},
						{
							Label: "Box style",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsBox,
									Name:  "ex-box",
									Items: []ui.TabItem{
										{Label: "Tab 1", Active: true, Content: tabsContent("Tab content 1")},
										{Label: "Tab 2", Content: tabsContent("Tab content 2")},
										{Label: "Tab 3", Content: tabsContent("Tab content 3")},
									},
								})
							},
						},
						{
							Label: "Bottom placement",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style:  ui.TabsLift,
									Bottom: true,
									Name:   "ex-bottom",
									Items: []ui.TabItem{
										{Label: "Tab 1", Active: true, Content: tabsContent("Tab content 1")},
										{Label: "Tab 2", Content: tabsContent("Tab content 2")},
										{Label: "Tab 3", Content: tabsContent("Tab content 3")},
									},
								})
							},
						},
						{
							Label: "Custom colors",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsLift,
									Name:  "ex-color",
									Items: []ui.TabItem{
										{Label: "Tab 1", Active: true, ColorCSS: "text-primary [--tab-bg:var(--color-primary-100)] [--tab-border-color:var(--color-primary-500)]", Content: tabsContent("Custom colored tab content 1")},
										{Label: "Tab 2", Content: tabsContent("Tab content 2")},
										{Label: "Tab 3", Content: tabsContent("Tab content 3")},
									},
								})
							},
						},
						{
							Label: "With icons",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsLift,
									Name:  "ex-icons",
									Items: []ui.TabItem{
										{Label: "Live", Icon: "lucide--play", Active: true, Content: tabsContent("Live tab content")},
										{Label: "Comments", Icon: "lucide--message-circle", Content: tabsContent("Comments tab content")},
										{Label: "Likes", Icon: "lucide--heart", Content: tabsContent("Likes tab content")},
									},
								})
							},
						},
						{
							Label: "Disabled tabs",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsLift,
									Name:  "ex-disabled",
									Items: []ui.TabItem{
										{Label: "Available", Active: true, Content: tabsContent("Available tab content")},
										{Label: "Disabled", Disabled: true, Content: tabsContent("This tab is disabled")},
										{Label: "Also Here", Content: tabsContent("Another enabled tab")},
									},
								})
							},
						},
						{
							Label: "Small size (sm)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsLift,
									Size:  ui.TabsSM,
									Name:  "ex-sm",
									Items: []ui.TabItem{
										{Label: "Small 1", Active: true, Content: tabsContent("Small tabs content 1")},
										{Label: "Small 2", Content: tabsContent("Small tabs content 2")},
									},
								})
							},
						},
						{
							Label: "No content (visual only)",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.TabsWithBoundary(ui.TabsProps{
									Style: ui.TabsBox,
									Name:  "ex-no-content",
									Items: []ui.TabItem{
										{Label: "All", Active: true},
										{Label: "Open"},
										{Label: "Closed"},
									},
								})
							},
						},
					},
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
				{
					Name:        "Examples",
					Description: "Page headers with 2 and 3 breadcrumb levels.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "2-level breadcrumb",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.PageHeaderWithBoundary(nav.Crumbs("Home", "/", "Dashboard"))
							},
						},
						{
							Label: "3-level breadcrumb",
							RenderFunc: func(_ url.Values) templ.Component {
								return nav.PageHeaderWithBoundary(nav.Crumbs("Home", "/", "Cases", "/cases", "Johnson v. Smith"))
							},
						},
					},
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
			{
				Name:        "Examples",
				Description: "Menu in default and compact sizes.",
				SubExamples: []galleryruntime.GallerySubExample{
					{
						Label: "Default size",
						RenderFunc: func(_ url.Values) templ.Component {
							return nav.MenuWithBoundary("", []nav.MenuItem{
								{Label: "Dashboard", Icon: "lucide--layout-dashboard", Href: "#", Active: true},
								{Label: "Cases", Icon: "lucide--folder-open", Href: "#"},
								{Label: "Contacts", Icon: "lucide--users", Href: "#"},
								{Label: "Settings", Icon: "lucide--settings", Href: "#"},
							})
						},
					},
					{
						Label: "Compact (xs)",
						RenderFunc: func(_ url.Values) templ.Component {
							return nav.MenuWithBoundary("menu-xs", []nav.MenuItem{
								{Label: "Dashboard", Icon: "lucide--layout-dashboard", Href: "#", Active: true},
								{Label: "Cases", Icon: "lucide--folder-open", Href: "#"},
								{Label: "Contacts", Icon: "lucide--users", Href: "#"},
								{Label: "Settings", Icon: "lucide--settings", Href: "#"},
							})
						},
					},
				},
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
							withText("Cancel", ui.ButtonWithBoundary("", ui.ButtonGhost, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false)),
							withText("Confirm", ui.ButtonWithBoundary("", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false)),
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
			{
				Name:        "Examples",
				Description: "Modal in SM, MD, and LG sizes.",
				SubExamples: []galleryruntime.GallerySubExample{
					{
						Label: "Small (modal-sm)",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:240px;position:relative;">`); err != nil {
									return err
								}
								if err := withChildren(modal.ModalWithBoundary("Confirm — SM", modal.ModalSM), rawHTML(`<p class="text-sm text-base-content/70 mb-4">Modal body content.</p>`)).Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div>`)
								return err
							})
						},
					},
					{
						Label: "Default (modal-md)",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:280px;position:relative;">`); err != nil {
									return err
								}
								if err := withChildren(modal.ModalWithBoundary("Confirm — MD", modal.ModalMD), rawHTML(`<p class="text-sm text-base-content/70 mb-4">Modal body content.</p>`)).Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div>`)
								return err
							})
						},
					},
					{
						Label: "Large (modal-lg)",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:280px;position:relative;">`); err != nil {
									return err
								}
								if err := withChildren(modal.ModalWithBoundary("Confirm — LG", modal.ModalLG), rawHTML(`<p class="text-sm text-base-content/70 mb-4">Modal body content.</p>`)).Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div>`)
								return err
							})
						},
					},
				},
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
			{
				Name:        "Examples",
				Description: "Confirm dialogs for delete and archive actions.",
				SubExamples: []galleryruntime.GallerySubExample{
					{
						Label: "Delete confirmation",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:240px;position:relative;">`); err != nil {
									return err
								}
								if err := modal.ConfirmPopupWithBoundary("Delete record?", "This action cannot be undone.", "Delete", "#", "delete").Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div>`)
								return err
							})
						},
					},
					{
						Label: "Archive confirmation",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:240px;position:relative;">`); err != nil {
									return err
								}
								if err := modal.ConfirmPopupWithBoundary("Archive case?", "The case will be moved to your archive.", "Archive", "#", "patch").Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div>`)
								return err
							})
						},
					},
				},
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
						formFields := seq(
							form.FormField(form.FormFieldProps{Label: "Name", Name: "name", Type: "text", Placeholder: "Enter record name", Required: true}),
							form.FormField(form.FormFieldProps{Label: "Email", Name: "email", Type: "email", Placeholder: "user@example.com"}),
							form.FormField(form.FormFieldProps{Label: "Role", Name: "role", Type: "select", Options: []form.SelectOption{
								{Value: "admin", Label: "Admin"},
								{Value: "user", Label: "User"},
								{Value: "viewer", Label: "Viewer"},
							}}),
						)
						inner := devmode.ComponentBoundary("FormModal",
							templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								return modal.FormModal(modal.FormModalProps{
									ID: "gallery-form-modal", Title: title, Size: size,
									SubmitText: submitText, Action: "#", Method: "post",
								}).Render(templ.WithChildren(ctx, formFields), w)
							}),
							map[string]any{"id": "gallery-form-modal", "title": title, "size": string(size)},
						)
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
			{
				Name:        "Examples",
				Description: "Form modal at different sizes.",
				SubExamples: []galleryruntime.GallerySubExample{
					{
						Label: "Small form modal",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:300px;position:relative;">`); err != nil {
									return err
								}
								if err := modal.FormModalWithBoundary(modal.FormModalProps{
									ID: "ex-form-sm", Title: "New Case", Size: modal.ModalSM, SubmitText: "Create", Action: "#", Method: "post",
								}).Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div><script>document.addEventListener('DOMContentLoaded',function(){var d=document.getElementById('ex-form-sm');if(d&&d.showModal)d.showModal();});</script>`)
								return err
							})
						},
					},
					{
						Label: "Destructive form modal",
						RenderFunc: func(_ url.Values) templ.Component {
							return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div style="min-height:300px;position:relative;">`); err != nil {
									return err
								}
								if err := modal.FormModalWithBoundary(modal.FormModalProps{
									ID: "ex-form-err", Title: "Delete Case", Size: modal.ModalSM, SubmitText: "Delete", CancelText: "Keep", Action: "#", Method: "post", Variant: "error",
								}).Render(ctx, w); err != nil {
									return err
								}
								_, err := io.WriteString(w, `</div><script>document.addEventListener('DOMContentLoaded',function(){var d=document.getElementById('ex-form-err');if(d&&d.showModal)d.showModal();});</script>`)
								return err
							})
						},
					},
				},
			},
			},
		},

		// table — consolidated entry (covers text cells, badge cells, avatar cells, action menus, zebra, compact, empty)
		{
			Slug:        "table",
			Name:        "Table",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Composable data table. TableWithProps wraps TableHead/TableBody; content is composed from TableRow and TableCell primitives with any component as cell content.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Full-featured table: avatar + name, status badge, action menu, and pagination. Configurable row count, page, and all DaisyUI table modifiers.",
					RenderFunc: func(params url.Values) templ.Component {
						page := 1
						if v, err := parseInt(params.Get("page")); err == nil && v > 0 {
							page = v
						}
						type tableRow struct{ Name, Status, Role, Joined string }
						allRows := []tableRow{
							{"Alice Johnson", "active", "Admin", "2024-01-15"},
							{"Bob Smith", "pending", "Employee", "2024-03-02"},
							{"Carol White", "closed", "Employee", "2023-11-20"},
							{"David Kim", "active", "Viewer", "2024-06-10"},
							{"Eve Martinez", "pending", "Employee", "2024-08-22"},
						}
						rowCount := 3
						if v, err := parseInt(params.Get("rows")); err == nil && v >= 1 && v <= 5 {
							rowCount = v
						}
						props := table.TableProps{
							Size:     params.Get("size"),
							Striped:  params.Get("striped") != "false",
							PinRows:  params.Get("pin_rows") == "true",
							PinCols:  params.Get("pin_cols") == "true",
							Bordered: params.Get("bordered") == "true",
						}
						actionItems := []ui.ActionMenuItem{{Label: "View"}, {Label: "Edit"}, {Label: "Delete", Danger: true}}
						rowComponents := make([]templ.Component, rowCount)
						for i, r := range allRows[:rowCount] {
							r := r
							nameCell := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
								if _, err := io.WriteString(w, `<div class="flex items-center gap-3">`); err != nil {
									return err
								}
								if err := ui.AvatarWithBoundary(r.Name, "", "", ui.AvatarSM).Render(ctx, w); err != nil {
									return err
								}
								_, err := fmt.Fprintf(w, `<span class="font-medium">%s</span></div>`, r.Name)
								return err
							})
							rowComponents[i] = withChildren(
								table.TableRowWithBoundary("", true),
								seq(
									withChildren(table.TableCellWithBoundary(""), nameCell),
									withChildren(table.TableCellWithBoundary(""), ui.StatusBadgeWithBoundary(r.Status)),
									withChildren(table.TableCellWithBoundary("text-sm text-base-content/70"), rawHTML(r.Role)),
									withChildren(table.TableCellWithBoundary("text-sm text-base-content/60"), rawHTML(r.Joined)),
									withChildren(table.TableCellWithBoundary("text-right"), ui.ActionMenuWithBoundary(actionItems)),
								),
							)
						}
						tbl := withChildren(
							table.TableWithPropsWithBoundary(props),
							seq(
								withChildren(
									table.TableHeadWithBoundary(),
									withChildren(table.TableHeadRowWithBoundary(), seq(
										table.TableHeadCellWithBoundary("Name"),
										table.TableHeadCellWithBoundary("Status"),
										table.TableHeadCellWithBoundary("Role"),
										table.TableHeadCellWithBoundary("Joined"),
										table.TableHeadCellWithBoundary(""),
									)),
								),
								withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
							),
						)
						totalPages := 3
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6">`); err != nil {
								return err
							}
							if err := tbl.Render(ctx, w); err != nil {
								return err
							}
							if _, err := fmt.Fprintf(w, `<div class="flex items-center justify-between mt-4"><p class="text-sm text-base-content/60">Showing %d of 47 entries</p>`, rowCount); err != nil {
								return err
							}
							if err := ui.PaginationWithBoundary(page, totalPages, "#", "main-content").Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
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
						{Label: "Size", Group: "Style", Type: galleryruntime.TokenTypeSelect, Default: "", QueryParam: "size", Options: []galleryruntime.TokenOption{
							{Value: "xs", Label: "xs"},
							{Value: "sm", Label: "sm"},
							{Value: "", Label: "md (default)"},
							{Value: "lg", Label: "lg"},
							{Value: "xl", Label: "xl"},
						}},
						{Label: "Striped", Group: "Style", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "striped"},
						{Label: "Bordered", Group: "Style", Type: galleryruntime.TokenTypeBool, Default: "false", QueryParam: "bordered"},
						{Label: "Pin rows", Group: "Style", Type: galleryruntime.TokenTypeBool, Default: "false", QueryParam: "pin_rows"},
						{Label: "Pin cols", Group: "Style", Type: galleryruntime.TokenTypeBool, Default: "false", QueryParam: "pin_cols"},
					},
				},
				{
					Name:        "Examples",
					Description: "Cell type variations: text, status badge, avatar+name, action menu, empty state, and styling variants.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Text cells",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept, Joined string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering", "Jan 2024"},
									{"Bob Martinez", "Member", "Legal", "Mar 2024"},
									{"Carol White", "Viewer", "Finance", "Jun 2024"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Joined)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
											table.TableHeadCellWithBoundary("Joined"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Status badge cells",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Status, Role string }
								rows := []tableRow{
									{"Alice Johnson", "active", "Admin"},
									{"Bob Martinez", "pending", "Member"},
									{"Carol White", "closed", "Viewer"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), ui.StatusBadgeWithBoundary(r.Status)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Status"),
											table.TableHeadCellWithBoundary("Role"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Avatar + name cells",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Joined string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Jan 2024"},
									{"Bob Martinez", "Member", "Mar 2024"},
									{"Carol White", "Viewer", "Jun 2024"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									nameCell := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
										if _, err := io.WriteString(w, `<div class="flex items-center gap-3">`); err != nil {
											return err
										}
										if err := ui.AvatarWithBoundary(r.Name, "", "", ui.AvatarSM).Render(ctx, w); err != nil {
											return err
										}
										_, err := fmt.Fprintf(w, `<span class="font-medium">%s</span></div>`, r.Name)
										return err
									})
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), nameCell),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Joined)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Joined"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Action menu column",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role string }
								rows := []tableRow{
									{"Alice Johnson", "Admin"},
									{"Bob Martinez", "Member"},
									{"Carol White", "Viewer"},
								}
								actionItems := []ui.ActionMenuItem{
									{Label: "Edit", Icon: "lucide--pencil", HXGet: "#"},
									{Label: "Delete", Icon: "lucide--trash-2", HXGet: "#", Danger: true},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary("text-right"), ui.ActionMenuWithBoundary(actionItems)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary(""),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Zebra stripes",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
									{"David Kim", "Member", "Operations"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{Striped: true}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Compact (sm)",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
									{"David Kim", "Member", "Operations"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{Size: "sm"}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Bordered",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{Bordered: true}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Extra small (xs)",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
									{"David Kim", "Member", "Operations"},
									{"Eve Martinez", "Viewer", "HR"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{Size: "xs"}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Large (lg)",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{Size: "lg"}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Extra large (xl)",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{Size: "xl"}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Role"),
											table.TableHeadCellWithBoundary("Department"),
										))),
										withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
									),
								)
							},
						},
						{
							Label: "Pinned rows",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Role, Dept string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering"},
									{"Bob Martinez", "Member", "Legal"},
									{"Carol White", "Viewer", "Finance"},
									{"David Kim", "Member", "Operations"},
									{"Eve Martinez", "Viewer", "HR"},
									{"Frank Lee", "Admin", "Engineering"},
									{"Grace Chen", "Member", "Legal"},
									{"Henry Park", "Viewer", "Finance"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Role)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Dept)),
									))
								}
								return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
									if _, err := io.WriteString(w, `<div class="h-48 overflow-auto">`); err != nil {
										return err
									}
									if err := withChildren(
										table.TableWithPropsWithBoundary(table.TableProps{PinRows: true}),
										seq(
											withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
												table.TableHeadCellWithBoundary("Name"),
												table.TableHeadCellWithBoundary("Role"),
												table.TableHeadCellWithBoundary("Department"),
											))),
											withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
										),
									).Render(ctx, w); err != nil {
										return err
									}
									_, err := io.WriteString(w, `</div>`)
									return err
								})
							},
						},
						{
							Label: "Pinned cols",
							RenderFunc: func(_ url.Values) templ.Component {
								type tableRow struct{ Name, Col2, Col3, Col4, Col5 string }
								rows := []tableRow{
									{"Alice Johnson", "Admin", "Engineering", "New York", "Full-time"},
									{"Bob Martinez", "Member", "Legal", "Chicago", "Part-time"},
									{"Carol White", "Viewer", "Finance", "Austin", "Full-time"},
								}
								rowComponents := make([]templ.Component, len(rows))
								for i, r := range rows {
									r := r
									rowComponents[i] = withChildren(table.TableRowWithBoundary("", false), seq(
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Name)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Col2)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Col3)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Col4)),
										withChildren(table.TableCellWithBoundary(""), rawHTML(r.Col5)),
									))
								}
								return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
									if _, err := io.WriteString(w, `<div class="w-72 overflow-auto">`); err != nil {
										return err
									}
									if err := withChildren(
										table.TableWithPropsWithBoundary(table.TableProps{PinCols: true}),
										seq(
											withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
												table.TableHeadCellWithBoundary("Name"),
												table.TableHeadCellWithBoundary("Role"),
												table.TableHeadCellWithBoundary("Dept"),
												table.TableHeadCellWithBoundary("City"),
												table.TableHeadCellWithBoundary("Type"),
											))),
											withChildren(table.TableBodyWithBoundary(), seq(rowComponents...)),
										),
									).Render(ctx, w); err != nil {
										return err
									}
									_, err := io.WriteString(w, `</div>`)
									return err
								})
							},
						},
						{
							Label: "Empty state",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableWithPropsWithBoundary(table.TableProps{}),
									seq(
										withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
											table.TableHeadCellWithBoundary("Name"),
											table.TableHeadCellWithBoundary("Status"),
											table.TableHeadCellWithBoundary("Role"),
										))),
										withChildren(table.TableBodyWithBoundary(), table.TableEmptyWithBoundary(3, "No records found.")),
									),
								)
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
			{
				Name:        "Examples",
				Description: "Log table with all four log types.",
				SubExamples: []galleryruntime.GallerySubExample{
					{
						Label: "All log types",
						RenderFunc: func(_ url.Values) templ.Component {
							now := time.Now()
							return logs.LogsTableWithBoundary([]logs.LogEntry{
								{Type: "success", Message: "Case created successfully.", CreatedAt: now.Add(-1 * time.Minute)},
								{Type: "info", Message: "Workflow triggered: document-review.", CreatedAt: now.Add(-3 * time.Minute)},
								{Type: "warn", Message: "API rate limit at 80% of quota.", CreatedAt: now.Add(-10 * time.Minute)},
								{Type: "error", Message: "Integration sync failed: connection timeout.", CreatedAt: now.Add(-30 * time.Minute)},
								{Type: "info", Message: "User logged in from new device.", CreatedAt: now.Add(-1 * time.Hour)},
							})
						},
					},
					{
						Label: "Empty state",
						RenderFunc: func(_ url.Values) templ.Component {
							return logs.LogsTableWithBoundary([]logs.LogEntry{})
						},
					},
				},
			},
			},
		},

		// ── Feedback / Radial Progress ───────────────────────────────────────────
		{
			Slug:        "radial-progress",
			Name:        "Radial Progress",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Progress",
			Description: "Circular radial progress indicator with configurable colour, size and thickness.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Adjust value, size and thickness with design tokens.",
					RenderFunc: func(params url.Values) templ.Component {
						val := 70
						if v, err := parseInt(params.Get("value")); err == nil && v >= 0 && v <= 100 {
							val = v
						}
						return row(
							ui.RadialProgressWithBoundary(ui.ProgressPrimary, val, "6rem", "4px"),
							ui.RadialProgressWithBoundary(ui.ProgressSecondary, val, "6rem", "4px"),
							ui.RadialProgressWithBoundary(ui.ProgressSuccess, val, "6rem", "4px"),
						)
					},
					Tokens: []galleryruntime.DesignToken{
						{
							Label:      "Value",
							Group:      "Progress",
							Type:       galleryruntime.TokenTypeText,
							QueryParam: "value",
							Default:    "70",
						},
					},
				},
				{
					Name:        "Examples",
					Description: "Various sizes and colours.",
					Templ: row(
						ui.RadialProgressWithBoundary(ui.ProgressPrimary, 25, "4rem", ""),
						ui.RadialProgressWithBoundary(ui.ProgressSecondary, 50, "5rem", ""),
						ui.RadialProgressWithBoundary(ui.ProgressSuccess, 75, "6rem", ""),
						ui.RadialProgressWithBoundary(ui.ProgressWarning, 90, "7rem", "8px"),
						ui.RadialProgressWithBoundary(ui.ProgressError, 100, "8rem", "10px"),
					),
				},
			},
		},

		// ── Layout / Drawer ───────────────────────────────────────────────────────
		{
			Slug:        "drawer",
			Name:        "Drawer Sidebar",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Drawers",
			Description: "Slide-in sidebar overlay driven by a hidden checkbox toggle.",
			FrameHeight: "400px",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Left drawer",
					Description: "Sidebar slides in from the left.",
					FrameHeight: "400px",
					Templ: ui.DrawerWithBoundary(
						"drawer-left-demo",
						ui.DrawerLeft,
						rawHTML(`<div class="p-6 flex flex-col gap-4">
							<label for="drawer-left-demo" class="btn btn-primary drawer-button w-fit">Open sidebar</label>
							<p class="text-base-content/70">This is the main content area. Click the button to open the sidebar.</p>
						</div>`),
						rawHTML(`<ul class="menu gap-1"><li><a>Dashboard</a></li><li><a>Reports</a></li><li><a>Settings</a></li></ul>`),
						"",
					),
				},
				{
					Name:        "Right drawer",
					Description: "Sidebar slides in from the right.",
					FrameHeight: "400px",
					Templ: ui.DrawerWithBoundary(
						"drawer-right-demo",
						ui.DrawerRight,
						rawHTML(`<div class="p-6 flex flex-col gap-4">
							<label for="drawer-right-demo" class="btn btn-primary drawer-button w-fit">Open right sidebar</label>
							<p class="text-base-content/70">This is the main content area.</p>
						</div>`),
						rawHTML(`<ul class="menu gap-1"><li><a>Profile</a></li><li><a>Preferences</a></li><li><a>Log out</a></li></ul>`),
						"",
					),
				},
			},
		},

		// ── Forms / Label ─────────────────────────────────────────────────────────
		{
			Slug:        "label",
			Name:        "Label",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Layout",
			Description: "DaisyUI label wrapper with primary text and optional alt/hint text slots.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Examples",
					Description: "Label with top text, alt text and bottom hint row.",
					Templ: rawHTML(`<div class="p-6 flex flex-col gap-6 max-w-sm">
						<div class="form-control">
							<div class="label">
								<span class="label-text">Email address</span>
								<span class="label-text-alt">Required</span>
							</div>
							<input type="email" class="input input-bordered w-full" placeholder="you@example.com"/>
							<div class="label">
								<span class="label-text-alt">We'll never share your email.</span>
							</div>
						</div>
						<div class="form-control">
							<div class="label">
								<span class="label-text">Username</span>
							</div>
							<input type="text" class="input input-bordered w-full" placeholder="johndoe"/>
						</div>
					</div>`),
				},
			},
		},

		// ── Forms / Validator ─────────────────────────────────────────────────────
		{
			Slug:        "validator",
			Name:        "Validator",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "DaisyUI validator class adds green/red colour feedback on HTML5 validation states.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Examples",
					Description: "Inputs with built-in HTML5 validation feedback.",
					FrameHeight: "300px",
					Templ: rawHTML(`<div class="p-6 flex flex-col gap-6 max-w-sm">
						<div class="form-control">
							<div class="label"><span class="label-text">Email (required)</span></div>
							<div class="validator">
								<input type="email" required class="input input-bordered w-full" placeholder="you@example.com"/>
							</div>
							<p class="validator-hint">Please enter a valid email address.</p>
						</div>
						<div class="form-control">
							<div class="label"><span class="label-text">Password (min 8 chars)</span></div>
							<div class="validator">
								<input type="password" minlength="8" required class="input input-bordered w-full" placeholder="••••••••"/>
							</div>
							<p class="validator-hint">Password must be at least 8 characters.</p>
						</div>
					</div>`),
				},
			},
		},

		// ── Foundation / Theme Controller ─────────────────────────────────────────
		{
			Slug:        "theme-controller",
			Name:        "Theme Controller",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Checkbox or radio inputs with the theme-controller class change the page theme.",
			FrameHeight: "500px",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Toggle (checkbox)",
					Description: "A single checkbox that toggles between light and dark.",
					Templ: rawHTML(`<div class="p-6 flex flex-col gap-4 items-start">
						<p class="text-sm text-base-content/60">Check to switch theme:</p>
						<label class="flex items-center gap-3 cursor-pointer">
							<span class="label-text">Dark mode</span>
							<input type="checkbox" value="dark" class="theme-controller toggle toggle-primary"/>
						</label>
					</div>`),
				},
				{
					Name:        "Picker (radio)",
					Description: "Radio buttons to pick from a list of DaisyUI themes.",
					FrameHeight: "500px",
					Templ: rawHTML(`<div class="p-6 flex flex-col gap-2 max-w-xs">
						<p class="text-sm font-semibold mb-2">Choose a theme:</p>
						` + func() string {
						themes := []string{"light", "dark", "cupcake", "bumblebee", "emerald", "synthwave", "retro", "cyberpunk", "dracula", "nord"}
						out := ""
						for i, t := range themes {
							checked := ""
							if i == 0 {
								checked = " checked"
							}
							out += `<label class="flex items-center gap-3 cursor-pointer">
								<input type="radio" name="theme-controller" value="` + t + `" class="theme-controller radio radio-sm"` + checked + `/>
								<span class="label-text capitalize">` + t + `</span>
							</label>`
						}
						return out
					}() + `
					</div>`),
				},
			},
		},

		// ── Forms / Calendar ─────────────────────────────────────────────────────
		{
			Slug:        "calendar",
			Name:        "Calendar",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "DaisyUI calendar CSS classes applied to a static HTML month calendar demo.",
			FrameHeight: "350px",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Single month",
					Description: "A styled month calendar built with plain HTML and DaisyUI calendar classes.",
					Templ: form.CalendarDemoWithBoundary("June", "2025", 0, 30, 17, 23),
				},
				{
					// Two months side by side for a range-picker layout demo.
					// May 2025: starts Thursday (col 4), 31 days, 15th=today, 20th=selected.
					// June 2025: starts Sunday (col 0), 30 days, no today, 5th=selected.
					Name:        "Examples",
					Description: "Two calendar months side by side (range picker layout).",
					FrameHeight: "380px",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "May 2025",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CalendarDemoWithBoundary("May", "2025", 4, 31, 15, 20)
							},
						},
						{
							Label: "June 2025",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.CalendarDemoWithBoundary("June", "2025", 0, 30, 0, 5)
							},
						},
					},
				},
			},
		},

		// ── Data Display / Text Rotate ────────────────────────────────────────────
		{
			Slug:        "text-rotate",
			Name:        "Text Rotate",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Animated component that cycles through text lines with an infinite loop.",
			FrameHeight: "200px",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Cycling text with configurable duration.",
					Templ: rawHTML(`<div class="p-8 flex flex-col items-center gap-4">
						<h2 class="text-3xl font-bold flex gap-2">
							We build
							<span class="text-primary">
								<div class="text-rotate">
									<span>fast apps</span>
									<span>great UIs</span>
									<span>with Go</span>
									<span>with HTMX</span>
									<span>with DaisyUI</span>
								</div>
							</span>
						</h2>
					</div>`),
				},
				{
					Name:        "Standalone",
					Description: "TextRotate component on its own.",
					Templ: seq(
						rawHTML(`<div class="p-8 flex justify-center text-2xl font-semibold">`),
						ui.TextRotateWithBoundary([]string{"First item", "Second item", "Third item", "Fourth item"}, "8s"),
						rawHTML(`</div>`),
					),
				},
			},
		},

		// ── Data Display / Hover 3D Card ──────────────────────────────────────────
		{
			Slug:        "hover-3d",
			Name:        "Hover 3D Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Card that tilts and rotates in 3D space when hovered.",
			FrameHeight: "300px",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Default",
					Description: "Hover over the card to see the 3D tilt effect.",
					Templ: rawHTML(`<div class="p-10 flex justify-center">
						<div class="hover-3d rounded-2xl w-64 h-40 bg-primary text-primary-content shadow-xl">
							<div class="hover-3d-layer flex flex-col items-center justify-center h-full p-4">
								<h3 class="text-xl font-bold">3D Card</h3>
								<p class="text-sm opacity-80 mt-1">Hover me!</p>
							</div>
						</div>
					</div>`),
				},
				{
					Name:        "Examples",
					Description: "Multiple 3D cards side by side.",
					FrameHeight: "300px",
					Templ: rawHTML(`<div class="p-6 flex gap-4 justify-center flex-wrap">
						<div class="hover-3d rounded-2xl w-44 h-36 bg-secondary text-secondary-content shadow-lg">
							<div class="hover-3d-layer flex items-center justify-center h-full font-bold">Secondary</div>
						</div>
						<div class="hover-3d rounded-2xl w-44 h-36 bg-accent text-accent-content shadow-lg">
							<div class="hover-3d-layer flex items-center justify-center h-full font-bold">Accent</div>
						</div>
						<div class="hover-3d rounded-2xl w-44 h-36 bg-neutral text-neutral-content shadow-lg">
							<div class="hover-3d-layer flex items-center justify-center h-full font-bold">Neutral</div>
						</div>
					</div>`),
				},
			},
		},

		// ── Data Display / Hover Gallery ──────────────────────────────────────────
		{
			Slug:        "hover-gallery",
			Name:        "Hover Gallery",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Horizontal hover image gallery — first image visible by default, others revealed on hover.",
			FrameHeight: "280px",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Default",
					Description: "Hover horizontally across the gallery to reveal images.",
					Templ: ui.HoverGalleryWithBoundary([]ui.HoverGalleryImage{
						{Src: "https://picsum.photos/seed/g1/400/240", Alt: "Gallery image 1"},
						{Src: "https://picsum.photos/seed/g2/400/240", Alt: "Gallery image 2"},
						{Src: "https://picsum.photos/seed/g3/400/240", Alt: "Gallery image 3"},
						{Src: "https://picsum.photos/seed/g4/400/240", Alt: "Gallery image 4"},
					}),
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
			{
				Name:        "Examples",
				Description: "Navbar with different app names.",
				SubExamples: []galleryruntime.GallerySubExample{
					{
						Label: "LegalDesk",
						RenderFunc: func(_ url.Values) templ.Component {
							return layout.NavbarWithBoundary("LegalDesk")
						},
					},
					{
						Label: "CaseFlow",
						RenderFunc: func(_ url.Values) templ.Component {
							return layout.NavbarWithBoundary("CaseFlow")
						},
					},
					{
						Label: "DocVault",
						RenderFunc: func(_ url.Values) templ.Component {
							return layout.NavbarWithBoundary("DocVault")
						},
					},
				},
			},
			},
		},

		// ── Layout / Sidebar Variants ──────────────────────────────────────────────
		{
			Slug:        "sidebar-variant",
			Name:        "Sidebar Variants",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Sidebars",
			Description: "Themed sidebar variants for different app contexts: ecommerce, payment, project, chat, documentation.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Ecommerce",
					Description: "Full menu tree with collapsible sections, badges, trial promo card.",
					Templ: layout.SidebarVariantWithBoundary("ecommerce", layout.SidebarVariantOpts{
						AppName:   "ACME Inc",
						AvatarSrc: "./images/avatars/1.png",
					}),
					FrameHeight: "100vh",
				},
				{
					Name:        "Payment",
					Description: "Simplified payment dashboard with test mode toggle and logout.",
					Templ: layout.SidebarVariantWithBoundary("payment", layout.SidebarVariantOpts{
						AppName: "PayFlow",
					}),
					FrameHeight: "100vh",
				},
				{
					Name:        "Project",
					Description: "Org brand, project switcher, favorites, channels, in-meeting bar.",
					Templ: layout.SidebarVariantWithBoundary("project", layout.SidebarVariantOpts{
						OrgName: "Design Studio",
					}),
					FrameHeight: "100vh",
				},
				{
					Name:        "Chat",
					Description: "AI agent theme with search, chat threads, token usage meter.",
					Templ: layout.SidebarVariantWithBoundary("chat", layout.SidebarVariantOpts{
						AppName:           "Chat",
						SearchPlaceholder: "Search conversations...",
						TokenUsed:         6500,
						TokenMax:          10000,
					}),
					FrameHeight: "100vh",
				},
				{
					Name:        "Documentation",
					Description: "Learn theme with search, icon sections, tree menu, version switcher.",
					Templ: layout.SidebarVariantWithBoundary("documentation", layout.SidebarVariantOpts{
						AppName: "Docs",
						Version: "v3.0.0",
					}),
					FrameHeight: "100vh",
				},
			},
		},

		// ── Layout / Topbar Variants ───────────────────────────────────────────────
		{
			Slug:        "topbar-variant",
			Name:        "Topbar Variants",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Topbars",
			Description: "Topbar variants: classic, greeting, nav-menu, editor with different left and right zone content.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Classic",
					Description: "Hamburger + search bar, right: notifications, apps, language, avatar dropdown.",
					Templ: layout.TopbarVariantWithBoundary("classic", layout.TopbarVariantOpts{
						ShowSearch:         true,
						ShowNotifications:  true,
						ShowApps:           true,
						ShowLanguage:       true,
						AvatarSrc:          "./images/avatars/1.png",
					}),
					FrameHeight: "64px",
				},
				{
					Name:        "Greeting",
					Description: "Greeting text with subtitle, right: avatar with name + handle dropdown.",
					Templ: layout.TopbarVariantWithBoundary("greeting", layout.TopbarVariantOpts{
						Greeting:  "Good Morning",
						UserName:  "Denish",
						Subtitle:  "Welcome back, great to see you again!",
						AvatarSrc: "./images/avatars/1.png",
					}),
					FrameHeight: "64px",
				},
				{
					Name:        "Nav Menu 1",
					Description: "Horizontal nav links (Apps, Components, Pages) with active state.",
					Templ: layout.TopbarVariantWithBoundary("nav-menu-1", layout.TopbarVariantOpts{
						NavLinks: []layout.TopbarNavLink{
							{Label: "Apps", Href: "#"},
							{Label: "Components", Href: "#", Active: true},
							{Label: "Pages", Href: "#"},
						},
						ShowSearch:        true,
						ShowNotifications: true,
						ShowLanguage:      true,
						AvatarSrc:         "./images/avatars/1.png",
					}),
					FrameHeight: "64px",
				},
				{
					Name:        "Nav Menu 2",
					Description: "Extensive nav links (Dashboard, Analytics, Settings, etc.).",
					Templ: layout.TopbarVariantWithBoundary("nav-menu-2", layout.TopbarVariantOpts{
						ShowSearch:        true,
						ShowNotifications: true,
						ShowApps:          true,
						AvatarSrc:         "./images/avatars/1.png",
					}),
					FrameHeight: "64px",
				},
				{
					Name:        "Editor",
					Description: "Breadcrumb-style title + search, right: save status, avatar group, Edit button.",
					Templ: layout.TopbarVariantWithBoundary("editor", layout.TopbarVariantOpts{
						ShowSearch:      true,
						SearchPlaceholder: "Search files...",
						SaveStatus:      "Saved just now",
					}),
					FrameHeight: "64px",
				},
			},
		},

		// ── Layout / Footer Variants ──────────────────────────────────────────────
		{
			Slug:        "footer-variant",
			Name:        "Footer Variants",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Footers",
			Description: "9 footer styles: minimal, social, branding, legal, status, support, language pills, language+currency, custom background.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Minimal",
					Description: "Single centered copyright line.",
					Templ: nav.FooterVariantWithBoundary("minimal", nav.FooterVariantOpts{
						Copyright: "© 2025 Nexus. All rights reserved.",
					}),
				},
				{
					Name:        "Social",
					Description: "Brand left, copyright center, social icons right.",
					Templ: nav.FooterVariantWithBoundary("social", nav.FooterVariantOpts{
						Copyright:  "© 2025 Nexus. All rights reserved.",
						BrandName:  "Nexus",
						SocialLinks: []nav.SocialLink{
							{Icon: "lucide--github", Href: "#"},
							{Icon: "lucide--twitter", Href: "#"},
							{Icon: "lucide--linkedin", Href: "#"},
						},
					}),
				},
				{
					Name:        "Branding",
					Description: "Copyright left, 'Built with ❤️ daisyUI' right.",
					Templ: nav.FooterVariantWithBoundary("branding", nav.FooterVariantOpts{
						Copyright: "© 2025 Nexus. All rights reserved.",
					}),
				},
				{
					Name:        "Legal",
					Description: "Copyright left, legal links right.",
					Templ: nav.FooterVariantWithBoundary("legal", nav.FooterVariantOpts{
						Copyright: "© 2025 Nexus. All rights reserved.",
						LegalLinks: []nav.FooterLink{
							{Label: "Terms of Use", Href: "#"},
							{Label: "Privacy Policy", Href: "#"},
							{Label: "Legal & Compliance", Href: "#"},
						},
					}),
				},
				{
					Name:        "Status",
					Description: "System status badge left, copyright right.",
					Templ: nav.FooterVariantWithBoundary("status", nav.FooterVariantOpts{
						Copyright:   "© 2025 Nexus. All rights reserved.",
						StatusLabel: "System running smoothly",
						StatusOK:    true,
					}),
				},
				{
					Name:        "Support",
					Description: "Phone number + headset icon left, Follow + social icons right.",
					Templ: nav.FooterVariantWithBoundary("support", nav.FooterVariantOpts{
						PhoneNumber: "+1 (555) 123-4567",
						SocialLinks: []nav.SocialLink{
							{Icon: "lucide--github", Href: "#"},
							{Icon: "lucide--twitter", Href: "#"},
							{Icon: "lucide--linkedin", Href: "#"},
						},
					}),
				},
				{
					Name:        "Language Pills",
					Description: "Copyright left, language quick-switch pills right.",
					Templ: nav.FooterVariantWithBoundary("language-pills", nav.FooterVariantOpts{
						Copyright: "© 2025 Nexus. All rights reserved.",
						Languages: []nav.LangOption{
							{Label: "English", Href: "#"},
							{Label: "Spanish", Href: "#"},
							{Label: "Chinese", Href: "#"},
						},
					}),
				},
				{
					Name:        "Language + Currency",
					Description: "Copyright left, language dropdown + currency dropdown right.",
					Templ: nav.FooterVariantWithBoundary("language-currency", nav.FooterVariantOpts{
						Copyright: "© 2025 Nexus. All rights reserved.",
						Languages: []nav.LangOption{
							{Label: "English", Href: "#"},
							{Label: "Spanish", Href: "#"},
							{Label: "Chinese", Href: "#"},
						},
						Currencies: []nav.CurrencyOption{
							{Label: "USD", Symbol: "$", Href: "#"},
							{Label: "EUR", Symbol: "€", Href: "#"},
							{Label: "GBP", Symbol: "£", Href: "#"},
						},
					}),
				},
			},
		},

		// ── Layout / Page Title Variants ──────────────────────────────────────────
		{
			Slug:        "page-title-variant",
			Name:        "Page Title Variants",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Page Titles",
			Description: "Page title variants: minimal, ecommerce, editor, task, stepper, analytics.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Minimal",
					Description: "Title + breadcrumb on the right.",
					Templ: nav.PageTitleVariantWithBoundary("minimal", nav.PageTitleVariantOpts{
						Title: "Create New Tool",
						Steps: []nav.BreadcrumbStep{
							{Label: "Home", URL: "#"},
							{Label: "Tools", URL: "#"},
							{Label: "Create"},
						},
					}),
				},
				{
					Name:        "Ecommerce",
					Description: "Order summary with ID, date, status dropdown, Invoice + More buttons.",
					Templ: nav.PageTitleVariantWithBoundary("ecommerce", nav.PageTitleVariantOpts{
						Title:   "Order Summary",
						OrderID: "#12541",
						Date:    "Tue Jul 22 2025",
						Status:  "Paid",
					}),
				},
				{
					Name:        "Editor",
					Description: "Breadcrumbs + title + subtitle + Save/Preview/More actions.",
					Templ: nav.PageTitleVariantWithBoundary("editor", nav.PageTitleVariantOpts{
						Title:    "Smart Tool Builder",
						Subtitle: "Type: Custom Workflow",
						Date:     "Tue Jul 22 2025",
						Steps: []nav.BreadcrumbStep{
							{Label: "Dashboard", URL: "#"},
							{Label: "Tools", URL: "#"},
							{Label: "Builder"},
						},
						Actions: []nav.PageTitleEditorAction{
							{Label: "Save Changes", Class: "btn-primary btn-sm"},
							{Label: "Preview Tool", Class: "btn-outline btn-sm border-base-300"},
						},
					}),
				},
				{
					Name:        "Task",
					Description: "Schedule title with count badges + Sync/New Event buttons.",
					Templ: nav.PageTitleVariantWithBoundary("task", nav.PageTitleVariantOpts{
						Title:          "Today's Schedule",
						DueCount:       2,
						ProgressCount:  4,
						DoneCount:      7,
					}),
				},
				{
					Name:        "Stepper",
					Description: "Title + Preview button + 4-step progress indicator.",
					Templ: nav.PageTitleVariantWithBoundary("stepper", nav.PageTitleVariantOpts{
						Title:       "Build Your Smart Course",
						StepCurrent: 2,
						StepLabels:  []string{"Course Details", "Content Setup", "Appearance", "Launch Settings"},
					}),
				},
				{
					Name:        "Analytics",
					Description: "Analytics title with breadcrumbs + Filter/Export/New Report buttons.",
					Templ: nav.PageTitleVariantWithBoundary("analytics", nav.PageTitleVariantOpts{
						Title: "Analytics Overview",
						Steps: []nav.BreadcrumbStep{
							{Label: "Home", URL: "#"},
							{Label: "Analytics"},
						},
					}),
				},
			},
		},

		// ── Layout / Layout Builder ──────────────────────────────────────────────────
		{
			Slug:        "layout-builder",
			Name:        "Layout Builder",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Builder",
			Description: "Interactive mix-and-match builder for sidebar, topbar, and footer variants.",
			FrameHeight: "100vh",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Select sidebar, topbar, and footer variants to compose a live layout.",
					Templ: layout.LayoutBuilderWithBoundary(),
					FrameHeight: "100vh",
				},
			},
		},

		// ── Navigation / Notifications ─────────────────────────────────────────────
		{
			Slug:        "notification",
			Name:        "Notification",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Notifications",
			Description: "Notification dropdown/drawer variants: basic dropdown, tab dropdown, drawer side panel.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Basic",
					Description: "Dropdown with Today/Seen groups, mark read + view all.",
					Templ: nav.NotificationVariantWithBoundary("basic", nav.NotificationVariantOpts{
						Items: []nav.NotificationItem{
							{AvatarSrc: "./images/avatars/4.png", AvatarAlt: "User 4", Message: "Customer has requested a <span class=\"text-error\">return</span> for item", TimeAgo: "Now"},
							{AvatarSrc: "./images/avatars/5.png", AvatarAlt: "User 5", Message: "A new <span class=\"underline\">review</span> has been submitted for product", TimeAgo: "15 min ago"},
						},
					}),
				},
				{
					Name:        "Tab",
					Description: "Dropdown with All/Team/AI/@mention tabs, richer items with inline actions.",
					Templ: nav.NotificationVariantWithBoundary("tab", nav.NotificationVariantOpts{
						Tabs: []nav.NotificationTab{
							{Label: "All", Active: true, Badge: 4},
							{Label: "Team"},
							{Label: "AI"},
							{Label: "@mention"},
						},
						Items: []nav.NotificationItem{
							{AvatarSrc: "./images/avatars/2.png", AvatarAlt: "User 2", Message: "Lena submitted a draft for review.", TimeAgo: "15 min ago", Actions: []nav.NotificationAction{
								{Label: "Approve", Class: "btn-primary"},
								{Label: "Decline", Class: "btn-ghost"},
							}},
							{AvatarSrc: "./images/avatars/1.png", AvatarAlt: "User 1", Message: "Alex shared a new design file.", TimeAgo: "1 hour ago"},
						},
					}),
				},
				{
					Name:        "Drawer",
					Description: "Full-height side panel with tabs and notification list.",
					Templ: nav.NotificationVariantWithBoundary("drawer", nav.NotificationVariantOpts{
						Tabs: []nav.NotificationTab{
							{Label: "All", Active: true, Badge: 4},
							{Label: "Team"},
							{Label: "AI"},
						},
						Items: []nav.NotificationItem{
							{AvatarSrc: "./images/avatars/2.png", AvatarAlt: "User 2", Message: "New comment on your pull request.", TimeAgo: "5 min ago"},
							{AvatarSrc: "./images/avatars/3.png", AvatarAlt: "User 3", Message: "Deployment completed successfully.", TimeAgo: "1 hour ago"},
							{AvatarSrc: "./images/avatars/1.png", AvatarAlt: "User 1", Message: "You have a new follower.", TimeAgo: "2 hours ago"},
						},
					}),
				},
			},
		},

		// ── Navigation / Search Modal ─────────────────────────────────────────────
		{
			Slug:        "search-modal",
			Name:        "Search Modal",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Search",
			Description: "Command palette search dialog: minimal and split variants.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Minimal",
					Description: "Search dialog with Actions and Quick Links menus.",
					Templ: rawHTML(`<div class="p-6 flex justify-center">
						<button class="btn btn-circle btn-soft" onclick="document.getElementById('search-modal-minimal').showModal()" aria-label="Search">
							<span class="iconify lucide--search size-5"></span>
						</button>` + func() string {
						var buf strings.Builder
						nav.SearchModal("minimal", nav.SearchModalOpts{
							Actions: []nav.SearchAction{
								{Label: "Create a new folder", Icon: "lucide--folder-plus"},
								{Label: "Upload new document", Icon: "lucide--file-plus"},
								{Label: "Invite to project", Icon: "lucide--user-plus"},
							},
							QuickLinks: []nav.SearchQuickLink{
								{Label: "File Manager", Icon: "lucide--folders"},
								{Label: "Profile", Icon: "lucide--user"},
								{Label: "Dashboard", Icon: "lucide--layout-dashboard"},
								{Label: "Support", Icon: "lucide--help-circle"},
								{Label: "Keyboard Shortcuts", Icon: "lucide--keyboard"},
							},
						}).Render(context.Background(), &buf)
						return buf.String()
					}() + `</div>`),
				},
				{
					Name:        "Split",
					Description: "Larger search with keyboard hints, filter chips, and agent cards.",
					Templ: rawHTML(`<div class="p-6 flex justify-center">
						<button class="btn btn-circle btn-soft" onclick="document.getElementById('search-modal-split').showModal()" aria-label="Search">
							<span class="iconify lucide--search size-5"></span>
						</button>` + func() string {
						var buf strings.Builder
						nav.SearchModal("split", nav.SearchModalOpts{
							FilterChips: []nav.SearchFilterChip{
								{Label: "Writer", Active: true},
								{Label: "Editor"},
								{Label: "Explainer"},
							},
							AgentCards: []nav.SearchAgentCard{
								{Name: "Research Buddy", Description: "Helps with research tasks", Icon: "lucide--search"},
								{Name: "Task Planner", Description: "Organize your workflow", Icon: "lucide--list-checks"},
								{Name: "Spark Ideas", Description: "Brainstorming assistant", Icon: "lucide--lightbulb"},
							},
						}).Render(context.Background(), &buf)
						return buf.String()
					}() + `</div>`),
				},
			},
		},

		// ── Navigation / Profile Menu Variants ─────────────────────────────────────
		{
			Slug:        "profile-menu-variant",
			Name:        "Profile Menu Variants",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Profile Menus",
			Description: "Profile menu variants: default dropdown, switch account overlay, referral menu.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Default",
					Description: "Standard avatar dropdown with user info + menu items + sign-out.",
					Templ: nav.ProfileMenuVariantWithBoundary("default", nav.ProfileMenuVariantOpts{
						Name:       "Denish Navadiya",
						Email:      "denish@example.com",
						Initials:   "DN",
						SignOutHref: "#",
						Items: []nav.ProfileMenuItem{
							{Label: "My Profile", Href: "#", Icon: "lucide--user"},
							{Label: "Settings", Href: "#", Icon: "lucide--settings"},
							{Label: "Help", Href: "#", Icon: "lucide--help-circle"},
						},
					}),
				},
				{
					Name:        "Switch Account",
					Description: "Multi-account switcher with current user card + account list + add new.",
					Templ: nav.ProfileMenuVariantWithBoundary("switch-account", nav.ProfileMenuVariantOpts{
						Name:       "Denish Navadiya",
						Email:      "denish@example.com",
						Initials:   "DN",
						Role:       "Admin",
						Accounts: []nav.ProfileAccount{
							{Name: "Alex Johnson", Email: "alex@acme.com", Active: true},
							{Name: "Sarah Chen", Email: "sarah@acme.com"},
							{Name: "Mike Rivera", Email: "mike@acme.com"},
						},
					}),
				},
				{
					Name:        "Overlay",
					Description: "Profile overlay with avatar, verified badge, preferences/files/collections.",
					Templ: nav.ProfileMenuVariantWithBoundary("overlay", nav.ProfileMenuVariantOpts{
						Name:     "Denish Navadiya",
						Email:    "denish@example.com",
						Initials: "DN",
					}),
				},
				{
					Name:        "Referral",
					Description: "Referral-focused profile menu with rewards and invite CTA.",
					Templ: nav.ProfileMenuVariantWithBoundary("referral", nav.ProfileMenuVariantOpts{
						Name:     "Denish Navadiya",
						Email:    "denish@example.com",
						Initials: "DN",
					}),
				},
			},
		},

		// ── Forms / Form Control Wrappers ─────────────────────────────────────────
		{
			Slug:        "form-wrappers",
			Name:        "Form Control Wrappers",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Labeled form inputs with labelPosition (above / left), optional hint, error, and HTMX attrs.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "FormInput",
					Description: "Text input with label, labelPosition, hint, and error support.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Full Name"
						}
						pos := form.LabelPosition(params.Get("labelPosition"))
						if pos == "" {
							pos = form.LabelAbove
						}
						return form.FormInputWithBoundary("name", label, "", "Enter your name", pos, "", "", nil)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Full Name", QueryParam: "label"},
						{Label: "Label Position", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "above", QueryParam: "labelPosition", Options: []galleryruntime.TokenOption{
							{Value: "above", Label: "Above"},
							{Value: "left", Label: "Left"},
						}},
					},
				},
				{
					Name:        "FormSelect",
					Description: "Select dropdown with label, labelPosition, placeholder, and error support.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Country"
						}
						pos := form.LabelPosition(params.Get("labelPosition"))
						if pos == "" {
							pos = form.LabelAbove
						}
						return form.FormSelectWithBoundary("country", label, "", [][2]string{
							{"us", "United States"},
							{"ca", "Canada"},
							{"uk", "United Kingdom"},
						}, "Select a country", pos, "", "", nil)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Country", QueryParam: "label"},
						{Label: "Label Position", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "above", QueryParam: "labelPosition", Options: []galleryruntime.TokenOption{
							{Value: "above", Label: "Above"},
							{Value: "left", Label: "Left"},
						}},
					},
				},
				{
					Name:        "FormCheckbox",
					Description: "Checkbox with inline label and HTMX attrs.",
					RenderFunc: func(params url.Values) templ.Component {
						checked := params.Get("checked") == "true"
						label := params.Get("label")
						if label == "" {
							label = "Accept terms and conditions"
						}
						return form.FormCheckboxWithBoundary("terms", label, checked, "", nil)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Checked", Group: "State", Type: galleryruntime.TokenTypeSelect, Default: "false", QueryParam: "checked", Options: []galleryruntime.TokenOption{
							{Value: "false", Label: "Unchecked"},
							{Value: "true", Label: "Checked"},
						}},
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Accept terms and conditions", QueryParam: "label"},
					},
				},
				{
					Name:        "FormToggle",
					Description: "Toggle switch with inline label and HTMX attrs.",
					RenderFunc: func(params url.Values) templ.Component {
						checked := params.Get("checked") == "true"
						label := params.Get("label")
						if label == "" {
							label = "Enable notifications"
						}
						return form.FormToggleWithBoundary("notifications", label, checked, "", nil)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Checked", Group: "State", Type: galleryruntime.TokenTypeSelect, Default: "false", QueryParam: "checked", Options: []galleryruntime.TokenOption{
							{Value: "false", Label: "Unchecked"},
							{Value: "true", Label: "Checked"},
						}},
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Enable notifications", QueryParam: "label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Above and left label positions, with and without errors.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Label above (default)",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormInputWithBoundary("name-above", "Full Name", "", "Enter your name", form.LabelAbove, "", "", nil)
							},
						},
						{
							Label: "Label left",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormInputWithBoundary("name-left", "Full Name", "", "Enter your name", form.LabelLeft, "", "", nil)
							},
						},
						{
							Label: "With error",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormInputWithBoundary("email-err", "Email", "", "Enter your email", form.LabelAbove, "", "Invalid email address", nil)
							},
						},
						{
							Label: "With hint",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormInputWithBoundary("hint", "Username", "", "Choose a username", form.LabelAbove, "Must be 3-30 characters", "", nil)
							},
						},
						{
							Label: "Select above",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormSelectWithBoundary("sel-above", "Country", "us", [][2]string{{"us", "United States"}, {"ca", "Canada"}}, "Select", form.LabelAbove, "", "", nil)
							},
						},
						{
							Label: "Select left",
							RenderFunc: func(_ url.Values) templ.Component {
								return form.FormSelectWithBoundary("sel-left", "Country", "us", [][2]string{{"us", "United States"}, {"ca", "Canada"}}, "Select", form.LabelLeft, "", "", nil)
							},
						},
					},
				},
			},
		},

		// table-card — rich card-based table with header/footer compositions
		{
			Slug:        "table-card",
			Name:        "Table Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Card-based table composition: TableCardWrapper, TableCardHeader, TableCardFooter, TableSearch, TableSelect, and PaginationCircle — build full ecommerce-style table cards.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Full Nexus-style orders table card with togglable sections.",
					RenderFunc: func(params url.Values) templ.Component {
						showHeader := params.Get("header") != "false"
						showSearch := params.Get("search") != "false"
						showSelect := params.Get("select") != "false"
						showActions := params.Get("actions") != "false"
						showFooter := params.Get("footer") != "false"
						page := 1
						if v, err := parseInt(params.Get("page")); err == nil && v > 0 && v <= 5 {
							page = v
						}
						perPage := 10
						if v, err := parseInt(params.Get("per_page")); err == nil && (v == 10 || v == 20 || v == 50 || v == 100) {
							perPage = v
						}
						totalItems := int64(157)

						orderRows := makeOrdersTableBody(7)

						var pageComponents []templ.Component
						if showHeader {
							var headerChildren []templ.Component
							if showSearch {
								headerChildren = append(headerChildren, table.TableSearchWithBoundary("search", params.Get("search_q"), "Search orders", "#table-card-preview", "#"))
							}
							if showSelect {
								sel := table.TableSelectWithBoundary("category", "#table-card-preview", "#")
								sel = withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">Select Category</option>`),
									rawHTML(`<option value="fashion">Fashion</option>`),
									rawHTML(`<option value="daily">Daily Need</option>`),
									rawHTML(`<option value="cosmetic">Cosmetic</option>`),
									rawHTML(`<option value="electronics">Electronics</option>`),
									rawHTML(`<option value="food">Food</option>`),
								))
								headerChildren = append(headerChildren, sel)
							}
							if showActions {
								headerChildren = append(headerChildren, ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeSubmit, ui.ButtonShapeDefault, "lucide--monitor-dot", false, false))
							}
							if len(headerChildren) > 0 {
								pageComponents = append(pageComponents, withChildren(table.TableCardHeaderWithBoundary(), seq(headerChildren...)))
							}
						}
						tblProps := table.TableProps{Bordered: true, Size: "sm"}
						tbl := withChildren(
							table.TableWithPropsWithBoundary(tblProps),
							seq(
								withChildren(
									table.TableHeadWithBoundary(),
									withChildren(table.TableHeadRowWithBoundary(), seq(
										table.TableHeadCellWithBoundary("ID"),
										table.TableHeadCellWithBoundary("Customer"),
										table.TableHeadCellWithBoundary("Price"),
										table.TableHeadCellWithBoundary("Payment"),
										table.TableHeadCellWithBoundary("Status"),
										table.TableHeadCellWithBoundary("Action"),
									)),
								),
								withChildren(table.TableBodyWithBoundary(), orderRows),
							),
						)
						pageComponents = append(pageComponents, tbl)
						if showFooter {
							pageComponents = append(pageComponents, table.TableCardFooterWithBoundary(table.TableCardFooterProps{
								CurrentPage: page,
								TotalPages:  5,
								TotalItems:  totalItems,
								PageSize:    perPage,
								BaseURL:     "#",
								TargetID:    "table-card-preview",
							}))
						}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div id="table-card-preview">`); err != nil {
								return err
							}
							if err := withChildren(
								table.TableCardWrapperWithBoundary(),
								seq(pageComponents...),
							).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div>`)
							return err
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Header", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "header"},
						{Label: "Search", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "search"},
						{Label: "Category select", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "select"},
						{Label: "Action buttons", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "actions"},
						{Label: "Footer", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "footer"},
						{Label: "Page", Group: "Pagination", Type: galleryruntime.TokenTypeSelect, Default: "1", QueryParam: "page", Options: []galleryruntime.TokenOption{
							{Value: "1", Label: "1"}, {Value: "2", Label: "2"}, {Value: "3", Label: "3"}, {Value: "4", Label: "4"}, {Value: "5", Label: "5"},
						}},
						{Label: "Per page", Group: "Pagination", Type: galleryruntime.TokenTypeSelect, Default: "10", QueryParam: "per_page", Options: []galleryruntime.TokenOption{
							{Value: "10", Label: "10"}, {Value: "20", Label: "20"}, {Value: "50", Label: "50"}, {Value: "100", Label: "100"},
						}},
					},
				},
				{
					Name:        "Compositions",
					Description: "All composition patterns from minimal card to full ecommerce.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Card wrapper",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									withChildren(
										table.TableWithPropsWithBoundary(table.TableProps{Bordered: true}),
										seq(
											withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
												table.TableHeadCellWithBoundary("Name"),
												table.TableHeadCellWithBoundary("Role"),
												table.TableHeadCellWithBoundary("Department"),
											))),
											withChildren(table.TableBodyWithBoundary(), seq(
												withChildren(table.TableRowWithBoundary("", false), seq(
													withChildren(table.TableCellWithBoundary(""), rawHTML("Alice Johnson")),
													withChildren(table.TableCellWithBoundary(""), rawHTML("Admin")),
													withChildren(table.TableCellWithBoundary(""), rawHTML("Engineering")),
												)),
												withChildren(table.TableRowWithBoundary("", false), seq(
													withChildren(table.TableCellWithBoundary(""), rawHTML("Bob Martinez")),
													withChildren(table.TableCellWithBoundary(""), rawHTML("Member")),
													withChildren(table.TableCellWithBoundary(""), rawHTML("Legal")),
												)),
												withChildren(table.TableRowWithBoundary("", false), seq(
													withChildren(table.TableCellWithBoundary(""), rawHTML("Carol White")),
													withChildren(table.TableCellWithBoundary(""), rawHTML("Viewer")),
													withChildren(table.TableCellWithBoundary(""), rawHTML("Finance")),
												)),
											)),
										),
									),
								)
							},
						},
						{
							Label: "Header with search",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									seq(
										withChildren(table.TableCardHeaderWithBoundary(),
											table.TableSearchWithBoundary("search", "", "Search orders", "", ""),
										),
										withChildren(
											table.TableWithPropsWithBoundary(table.TableProps{Bordered: true}),
											makeOrdersTable(3),
										),
									),
								)
							},
						},
						{
							Label: "Header with filters",
							RenderFunc: func(_ url.Values) templ.Component {
								sel := table.TableSelectWithBoundary("category", "", "")
								sel = withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">All Categories</option>`),
									rawHTML(`<option value="fashion">Fashion</option>`),
									rawHTML(`<option value="electronics">Electronics</option>`),
									rawHTML(`<option value="food">Food</option>`),
								))
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									seq(
										withChildren(table.TableCardHeaderWithBoundary(), seq(
											table.TableSearchWithBoundary("search", "", "Search orders", "", ""),
											sel,
										)),
										withChildren(
											table.TableWithPropsWithBoundary(table.TableProps{Bordered: true}),
											makeOrdersTable(3),
										),
									),
								)
							},
						},
						{
							Label: "Header with actions",
							RenderFunc: func(_ url.Values) templ.Component {
								sel := table.TableSelectWithBoundary("category", "", "")
								sel = withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">All Categories</option>`),
									rawHTML(`<option value="fashion">Fashion</option>`),
								))
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									seq(
										withChildren(table.TableCardHeaderWithBoundary(), seq(
											table.TableSearchWithBoundary("search", "", "Search orders", "", ""),
											sel,
											ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeSubmit, ui.ButtonShapeDefault, "lucide--plus", false, false),
										)),
										withChildren(
											table.TableWithPropsWithBoundary(table.TableProps{Bordered: true}),
											makeOrdersTable(3),
										),
									),
								)
							},
						},
						{
							Label: "Circle pagination",
							RenderFunc: func(_ url.Values) templ.Component {
								return ui.PaginationCircleWithBoundary(3, 7, "#", "content")
							},
						},
						{
							Label: "Rich footer",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableCardFooterWithBoundary(table.TableCardFooterProps{
									CurrentPage: 2,
									TotalPages:  5,
									TotalItems:  157,
									PageSize:    20,
									BaseURL:     "#",
									TargetID:    "content",
								})
							},
						},
						{
							Label: "Card with footer",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									seq(
										withChildren(
											table.TableWithPropsWithBoundary(table.TableProps{Bordered: true}),
											makeOrdersTable(7),
										),
										table.TableCardFooterWithBoundary(table.TableCardFooterProps{
											CurrentPage: 1,
											TotalPages:  3,
											TotalItems:  157,
											PageSize:    10,
											BaseURL:     "#",
											TargetID:    "content",
										}),
									),
								)
							},
						},
						{
							Label: "Card full ecommerce",
							RenderFunc: func(_ url.Values) templ.Component {
								sel := table.TableSelectWithBoundary("category", "", "")
								sel = withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">Select Category</option>`),
									rawHTML(`<option value="fashion">Fashion</option>`),
									rawHTML(`<option value="daily">Daily Need</option>`),
									rawHTML(`<option value="cosmetic">Cosmetic</option>`),
									rawHTML(`<option value="electronics">Electronics</option>`),
									rawHTML(`<option value="food">Food</option>`),
								))
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									seq(
										withChildren(table.TableCardHeaderWithBoundary(), seq(
											table.TableSearchWithBoundary("search", "", "Search along orders", "", ""),
											sel,
											ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeSubmit, ui.ButtonShapeDefault, "lucide--monitor-dot", false, false),
										)),
										withChildren(
											table.TableWithPropsWithBoundary(table.TableProps{Bordered: true, Size: "sm"}),
											makeOrdersTable(10),
										),
										table.TableCardFooterWithBoundary(table.TableCardFooterProps{
											CurrentPage: 1,
											TotalPages:  5,
											TotalItems:  157,
											PageSize:    20,
											BaseURL:     "#",
											TargetID:    "content",
										}),
									),
								)
							},
						},
					},
				},
			},
		},

		// table-card-wrapper — standalone card container
		{
			Slug:        "table-card-wrapper",
			Name:        "Table Card Wrapper",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Card container for table compositions. Wraps content in a DaisyUI card with shadow and border.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Card wrapper with a basic 3-row table inside.",
					RenderFunc: func(params url.Values) templ.Component {
						return withChildren(
							table.TableCardWrapperWithBoundary(),
							withChildren(
								table.TableWithPropsWithBoundary(table.TableProps{}),
								makeOrdersTable(3),
							),
						)
					},
				},
				{
					Name:        "Examples",
					Description: "Wrapper with different content types.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Empty card",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableCardWrapperWithBoundary()
							},
						},
						{
							Label: "With table",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									withChildren(
										table.TableWithPropsWithBoundary(table.TableProps{}),
										makeOrdersTable(4),
									),
								)
							},
						},
						{
							Label: "With custom content",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardWrapperWithBoundary(),
									rawHTML(`<div class="p-6 text-center text-base-content/50">Custom content inside card wrapper</div>`),
								)
							},
						},
					},
				},
			},
		},

		// table-card-header — standalone header bar
		{
			Slug:        "table-card-header",
			Name:        "Table Card Header",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Header bar for table cards. Takes children for flexible composition: titles, search inputs, selects, and action buttons.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Header with togglable search input and action button.",
					RenderFunc: func(params url.Values) templ.Component {
						showSearch := params.Get("search") != "false"
						showActions := params.Get("actions") != "false"
						var children []templ.Component
						children = append(children, rawHTML(`<h3 class="text-base font-semibold">Orders</h3>`))
						if showSearch {
							children = append(children, table.TableSearchWithBoundary("q", "", "Search...", "", ""))
						}
						if showActions {
							children = append(children, ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeSubmit, ui.ButtonShapeDefault, "lucide--monitor-dot", false, false))
						}
						return withChildren(
							table.TableCardHeaderWithBoundary(),
							seq(children...),
						)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Search", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "search"},
						{Label: "Actions", Group: "Sections", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "actions"},
					},
				},
				{
					Name:        "Examples",
					Description: "Header with different child content.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Title only",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardHeaderWithBoundary(),
									rawHTML(`<h3 class="text-base font-semibold">Orders</h3>`),
								)
							},
						},
						{
							Label: "Title + search",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardHeaderWithBoundary(),
									seq(
										rawHTML(`<h3 class="text-base font-semibold">Orders</h3>`),
										table.TableSearchWithBoundary("q", "", "Search orders...", "", ""),
									),
								)
							},
						},
						{
							Label: "Title + search + actions",
							RenderFunc: func(_ url.Values) templ.Component {
								return withChildren(
									table.TableCardHeaderWithBoundary(),
									seq(
										rawHTML(`<h3 class="text-base font-semibold">Orders</h3>`),
										table.TableSearchWithBoundary("q", "", "Search orders...", "", ""),
										ui.ButtonWithBoundary("#", ui.ButtonPrimary, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeSubmit, ui.ButtonShapeDefault, "lucide--monitor-dot", false, false),
									),
								)
							},
						},
					},
				},
			},
		},

		// table-card-footer — standalone rich footer
		{
			Slug:        "table-card-footer",
			Name:        "Table Card Footer",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Rich footer for table cards: per-page selector, item range count, and circle-style pagination — matching the Nexus ecommerce orders pattern.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live page, total pages, page size, and total items controls.",
					RenderFunc: func(params url.Values) templ.Component {
						page := 1
						if v, err := parseInt(params.Get("page")); err == nil && v > 0 {
							page = v
						}
						totalPages := 5
						if v, err := parseInt(params.Get("totalPages")); err == nil && v > 0 {
							totalPages = v
						}
						pageSize := 10
						if v, err := parseInt(params.Get("pageSize")); err == nil && (v == 10 || v == 20 || v == 50 || v == 100) {
							pageSize = v
						}
						totalItems := int64(157)
						if v, err := parseInt(params.Get("totalItems")); err == nil && v > 0 {
							totalItems = int64(v)
						}
						return table.TableCardFooterWithBoundary(table.TableCardFooterProps{
							CurrentPage: page,
							TotalPages:  totalPages,
							TotalItems:  totalItems,
							PageSize:    pageSize,
							BaseURL:     "#",
							TargetID:    "content",
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Current Page", Group: "Pagination", Type: galleryruntime.TokenTypeRange, Default: "1", Min: 1, Max: 20, Step: 1, QueryParam: "page"},
						{Label: "Total Pages", Group: "Pagination", Type: galleryruntime.TokenTypeRange, Default: "5", Min: 1, Max: 50, Step: 1, QueryParam: "totalPages"},
						{Label: "Page Size", Group: "Pagination", Type: galleryruntime.TokenTypeSelect, Default: "10", QueryParam: "pageSize", Options: []galleryruntime.TokenOption{
							{Value: "10", Label: "10"},
							{Value: "20", Label: "20"},
							{Value: "50", Label: "50"},
							{Value: "100", Label: "100"},
						}},
						{Label: "Total Items", Group: "Pagination", Type: galleryruntime.TokenTypeRange, Default: "157", Min: 1, Max: 10000, Step: 1, QueryParam: "totalItems"},
					},
				},
				{
					Name:        "Examples",
					Description: "Footer at different page and dataset states.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Page 1 of 5",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableCardFooterWithBoundary(table.TableCardFooterProps{
									CurrentPage: 1,
									TotalPages:  5,
									TotalItems:  157,
									PageSize:    20,
									BaseURL:     "#",
									TargetID:    "content",
								})
							},
						},
						{
							Label: "Page 3 of 5",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableCardFooterWithBoundary(table.TableCardFooterProps{
									CurrentPage: 3,
									TotalPages:  5,
									TotalItems:  157,
									PageSize:    20,
									BaseURL:     "#",
									TargetID:    "content",
								})
							},
						},
						{
							Label: "Page 5 of 5",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableCardFooterWithBoundary(table.TableCardFooterProps{
									CurrentPage: 5,
									TotalPages:  5,
									TotalItems:  157,
									PageSize:    20,
									BaseURL:     "#",
									TargetID:    "content",
								})
							},
						},
						{
							Label: "Large dataset",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableCardFooterWithBoundary(table.TableCardFooterProps{
									CurrentPage: 42,
									TotalPages:  500,
									TotalItems:  9987,
									PageSize:    20,
									BaseURL:     "#",
									TargetID:    "content",
								})
							},
						},
					},
				},
			},
		},

		// table-search — standalone search input
		{
			Slug:        "table-search",
			Name:        "Table Search",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Search input with magnifying glass icon, designed for table card headers. Supports HTMX-powered live search.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Live value and placeholder controls.",
					RenderFunc: func(params url.Values) templ.Component {
						value := params.Get("value")
						placeholder := params.Get("placeholder")
						if placeholder == "" {
							placeholder = "Search orders..."
						}
						return table.TableSearchWithBoundary("q", value, placeholder, "", "")
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Value", Group: "Component", Type: galleryruntime.TokenTypeText, Default: "", QueryParam: "value"},
						{Label: "Placeholder", Group: "Component", Type: galleryruntime.TokenTypeText, Default: "Search orders...", QueryParam: "placeholder"},
					},
				},
				{
					Name:        "Examples",
					Description: "Search inputs with different placeholders and values.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Default placeholder",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableSearchWithBoundary("q", "", "Search orders...", "", "")
							},
						},
						{
							Label: "Pre-filled value",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableSearchWithBoundary("q", "Emily Johnson", "Search orders...", "", "")
							},
						},
						{
							Label: "Custom placeholder",
							RenderFunc: func(_ url.Values) templ.Component {
								return table.TableSearchWithBoundary("q", "", "Filter by customer name...", "", "")
							},
						},
					},
				},
			},
		},

		// table-select — standalone select filter
		{
			Slug:        "table-select",
			Name:        "Table Select",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Filter select for table card headers. Wraps a DaisyUI select-sm with HTMX change-triggered filtering.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Select filter with category options.",
					RenderFunc: func(params url.Values) templ.Component {
						sel := table.TableSelectWithBoundary("category", "", "")
						return withChildren(sel, seq(
							rawHTML(`<option value="" disabled="" selected="">All Categories</option>`),
							rawHTML(`<option value="fashion">Fashion</option>`),
							rawHTML(`<option value="electronics">Electronics</option>`),
							rawHTML(`<option value="food">Food</option>`),
						))
					},
				},
				{
					Name:        "Examples",
					Description: "Select with different option sets.",
					SubExamples: []galleryruntime.GallerySubExample{
						{
							Label: "Status filter",
							RenderFunc: func(_ url.Values) templ.Component {
								sel := table.TableSelectWithBoundary("status", "", "")
								return withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">All Status</option>`),
									rawHTML(`<option value="ordered">Ordered</option>`),
									rawHTML(`<option value="accepted">Accepted</option>`),
									rawHTML(`<option value="delivered">Delivered</option>`),
								))
							},
						},
						{
							Label: "Category filter",
							RenderFunc: func(_ url.Values) templ.Component {
								sel := table.TableSelectWithBoundary("category", "", "")
								return withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">Select Category</option>`),
									rawHTML(`<option value="fashion">Fashion</option>`),
									rawHTML(`<option value="daily">Daily Need</option>`),
									rawHTML(`<option value="cosmetic">Cosmetic</option>`),
									rawHTML(`<option value="electronics">Electronics</option>`),
									rawHTML(`<option value="food">Food</option>`),
								))
							},
						},
						{
							Label: "Sort by:",
							RenderFunc: func(_ url.Values) templ.Component {
								sel := table.TableSelectWithBoundary("sort", "", "")
								return withChildren(sel, seq(
									rawHTML(`<option value="" disabled="" selected="">Sort by</option>`),
									rawHTML(`<option value="date">Date</option>`),
									rawHTML(`<option value="name">Name</option>`),
									rawHTML(`<option value="price">Price</option>`),
								))
							},
						},
					},
				},
			},
		},

		// ── Basics / Buttons — Glow ─────────────────────────────────────────────
		{
			Slug:        "button-glow",
			Name:        "Button — Glow",
			Category:    galleryruntime.CategoryBasics,
			Subcategory: "Buttons",
			Description: "Button with a colored glow shadow effect for CTAs and hero sections.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Glow button with configurable variant.",
					RenderFunc: func(params url.Values) templ.Component {
						variantStr := params.Get("variant")
						variant := ui.ButtonPrimary
						switch variantStr {
						case "secondary":
							variant = ui.ButtonSecondary
						case "accent":
							variant = ui.ButtonAccent
						case "success":
							variant = ui.ButtonSuccess
						case "error":
							variant = ui.ButtonError
						case "neutral":
							variant = ui.ButtonNeutral
						case "outline":
							variant = ui.ButtonOutline
						}
						return ui.ButtonGlowWithBoundary("#", variant, ui.ButtonMD, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeDefault, "", false, false)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Variant", Group: "Style", Type: galleryruntime.TokenTypeSelect, Default: "primary", QueryParam: "variant", Options: []galleryruntime.TokenOption{
							{Value: "primary", Label: "Primary"},
							{Value: "secondary", Label: "Secondary"},
							{Value: "accent", Label: "Accent"},
							{Value: "success", Label: "Success"},
							{Value: "error", Label: "Error"},
							{Value: "neutral", Label: "Neutral"},
							{Value: "outline", Label: "Outline"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Glow buttons in primary, secondary, success, and accent variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="flex flex-wrap gap-4 p-6 items-center justify-center">`); err != nil {
								return err
							}
							for _, v := range []struct {
								variant ui.ButtonVariant
								label   string
							}{
								{ui.ButtonPrimary, "Primary"},
								{ui.ButtonSecondary, "Secondary"},
								{ui.ButtonSuccess, "Success"},
								{ui.ButtonError, "Error"},
								{ui.ButtonNeutral, "Neutral"},
							} {
								if err := ui.ButtonGlow(ui.ButtonProps{Variant: v.variant, Size: ui.ButtonMD, Type: ui.ButtonTypeButton, Shape: ui.ButtonShapeDefault}).Render(templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w2 io.Writer) error {
									_, err := io.WriteString(w2, v.label)
									return err
								})), w); err != nil {
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

		// ── Foundation / Effects — Animated Gradient ────────────────────────────
		{
			Slug:        "animated-gradient-text",
			Name:        "Animated Gradient Text",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "Animated gradient text that shifts colors using CSS animation. Also includes static gradient text examples.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive (animated)",
					Description: "Animated gradient text with configurable colors.",
					RenderFunc: func(params url.Values) templ.Component {
						text := params.Get("text")
						if text == "" {
							text = "go-daisy"
						}
						fromColor := params.Get("fromColor")
						if fromColor == "" {
							fromColor = "from-primary"
						}
						toColor := params.Get("toColor")
						if toColor == "" {
							toColor = "to-secondary"
						}
						return ui.AnimatedGradientTextWithBoundary(text, fromColor, toColor, "text-3xl", "font-black")
					},
					FrameHeight: "100px",
					Tokens: []galleryruntime.DesignToken{
						{Label: "Text", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "go-daisy", QueryParam: "text"},
						{Label: "From", Group: "Style", Type: galleryruntime.TokenTypeSelect, Default: "from-primary", QueryParam: "fromColor", Options: []galleryruntime.TokenOption{
							{Value: "from-primary", Label: "Primary"},
							{Value: "from-secondary", Label: "Secondary"},
							{Value: "from-accent", Label: "Accent"},
							{Value: "from-success", Label: "Success"},
							{Value: "from-warning", Label: "Warning"},
							{Value: "from-error", Label: "Error"},
						}},
						{Label: "To", Group: "Style", Type: galleryruntime.TokenTypeSelect, Default: "to-secondary", QueryParam: "toColor", Options: []galleryruntime.TokenOption{
							{Value: "to-secondary", Label: "Secondary"},
							{Value: "to-primary", Label: "Primary"},
							{Value: "to-accent", Label: "Accent"},
							{Value: "to-info", Label: "Info"},
							{Value: "to-success", Label: "Success"},
							{Value: "to-warning", Label: "Warning"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "Static and animated gradient text in multiple color combinations.",
					RenderFunc: func(_ url.Values) templ.Component {
						return rawHTML(`<div class="p-6 space-y-4">
	<p class="text-3xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">Static — Primary to Secondary</p>
	<p class="text-3xl font-bold bg-gradient-to-r from-success to-info bg-clip-text text-transparent">Static — Success to Info</p>
	<p class="bg-linear-to-r from-warning to-error bg-clip-text text-transparent text-3xl font-black animate-gradient-shift" style="background-size:200% 200%">Animated — Warning to Error</p>
	<p class="text-sm text-base-content/60 mt-4">Add <code class="bg-base-200 px-1 rounded text-xs">animate-gradient-shift</code> and <code class="bg-base-200 px-1 rounded text-xs">bg-linear-to-r from-X to-Y bg-clip-text text-transparent</code> for animated gradient text.</p>
</div>`)
					},
				},
			},
		},

		// ── Data Display / Cards — Testimonial ──────────────────────────────────
		{
			Slug:        "testimonial-card",
			Name:        "Testimonial Card",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "Quote card with star rating, avatar, name, and role. Composes existing Card, Avatar, and Rating primitives.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Testimonial with configurable quote, name, role, and rating.",
					RenderFunc: func(params url.Values) templ.Component {
						quote := params.Get("quote")
						if quote == "" {
							quote = "This product completely transformed our workflow. Highly recommended!"
						}
						name := params.Get("name")
						if name == "" {
							name = "Sarah Johnson"
						}
						role := params.Get("role")
						if role == "" {
							role = "CTO, TechCorp"
						}
						rating, _ := parseInt(params.Get("rating"))
						return ui.TestimonialCardWithBoundary(ui.TestimonialCardProps{
							Quote:       quote,
							Name:        name,
							Role:        role,
							Rating:      rating,
							RatingColor: "bg-orange-400",
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Quote", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "This product completely transformed our workflow. Highly recommended!", QueryParam: "quote"},
						{Label: "Name", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Sarah Johnson", QueryParam: "name"},
						{Label: "Role", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "CTO, TechCorp", QueryParam: "role"},
						{Label: "Rating", Group: "Content", Type: galleryruntime.TokenTypeRange, Default: "5", QueryParam: "rating", Min: 0, Max: 5, Step: 1},
					},
				},
				{
					Name:        "Examples",
					Description: "Testimonial cards with and without ratings, different roles.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 grid grid-cols-1 md:grid-cols-3 gap-4">`); err != nil {
								return err
							}
							cards := []ui.TestimonialCardProps{
								{Quote: "Amazing tool! Increased our productivity by 3x.", Name: "Alice M.", Role: "Product Manager", Rating: 5, RatingColor: "bg-orange-400"},
								{Quote: "Clean, well-documented, and a pleasure to use.", Name: "Bob K.", Role: "Lead Developer", Rating: 4, RatingColor: "bg-orange-400"},
								{Quote: "The support team is fantastic. They helped us get set up in no time.", Name: "Carol W.", Role: "Founder", Rating: 5, RatingColor: "bg-orange-400"},
							}
							for _, c := range cards {
								if err := ui.TestimonialCard(c).Render(ctx, w); err != nil {
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

		// ── Data Display / Cards — Sparkline Stat Card ──────────────────────────
		{
			Slug:        "stat-card-sparkline",
			Name:        "Stat Card — Sparkline",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Cards",
			Description: "KPI stat card with inline SVG sparkline chart showing trend data.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Stat card with sparkline and configurable label, value, and trend.",
					RenderFunc: func(params url.Values) templ.Component {
						value := params.Get("value")
						if value == "" {
							value = "$48,290"
						}
						label := params.Get("label")
						if label == "" {
							label = "Revenue"
						}
						trendLabel := params.Get("trendLabel")
						if trendLabel == "" {
							trendLabel = "12.5%"
						}
						trend := ui.StatTrend(params.Get("trend"))
						if trend == "" {
							trend = ui.StatTrendUp
						}
						return ui.StatCardSparklineWithBoundary(ui.StatCardSparklineProps{
							Label:      label,
							Value:      value,
							Trend:      trend,
							TrendLabel: trendLabel,
							Data:       []int{12, 19, 15, 22, 28, 24, 35, 42, 38, 48, 45, 52},
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Value", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "$48,290", QueryParam: "value"},
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Revenue", QueryParam: "label"},
						{Label: "Trend", Group: "Content", Type: galleryruntime.TokenTypeSelect, Default: "up", QueryParam: "trend", Options: []galleryruntime.TokenOption{
							{Value: "up", Label: "Up ↑"},
							{Value: "down", Label: "Down ↓"},
							{Value: "", Label: "Neutral"},
						}},
						{Label: "Trend label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "12.5%", QueryParam: "trendLabel"},
					},
				},
				{
					Name:        "Examples",
					Description: "Multiple sparkline stat cards showing different KPIs.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">`); err != nil {
								return err
							}
							cards := []ui.StatCardSparklineProps{
								{Label: "Revenue", Value: "$48,290", Trend: ui.StatTrendUp, TrendLabel: "12.5%", Data: []int{12, 19, 15, 22, 28, 24, 35, 42, 38, 48, 45, 52}},
								{Label: "Users", Value: "12,430", Trend: ui.StatTrendUp, TrendLabel: "8.3%", Data: []int{50, 55, 62, 58, 70, 68, 75, 82, 78, 88, 92, 95}},
								{Label: "Bounce Rate", Value: "24.1%", Trend: ui.StatTrendDown, TrendLabel: "3.2%", Data: []int{30, 28, 25, 22, 20, 18, 16, 15, 14, 13, 12, 11}},
								{Label: "Avg. Session", Value: "4m 32s", Trend: ui.StatTrendUp, TrendLabel: "5.7%", Data: []int{60, 65, 70, 68, 75, 72, 78, 82, 80, 85, 88, 90}},
							}
							for _, c := range cards {
								if err := ui.StatCardSparkline(c).Render(ctx, w); err != nil {
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

		// ── Forms / Inputs — Password Field ─────────────────────────────────────
		{
			Slug:        "password-field",
			Name:        "Password Field",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Password input with optional show/hide toggle button. Includes guard script for HTMX partial swaps.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Password input with configurable label and show toggle.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Password"
						}
						return form.PasswordFieldWithBoundary(form.PasswordFieldProps{
							Name:       "password-demo",
							Label:      label,
							ShowToggle: true,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Password", QueryParam: "label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Password field with and without toggle.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6 max-w-sm">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">With toggle</p>`); err != nil {
								return err
							}
							if err := form.PasswordField(form.PasswordFieldProps{Name: "pw1", Label: "Password", ShowToggle: true}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Without toggle</p>`); err != nil {
								return err
							}
							if err := form.PasswordField(form.PasswordFieldProps{Name: "pw2", Label: "Password"}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
			},
		},

		// ── Forms / Inputs — Password Meter ─────────────────────────────────────
		{
			Slug:        "password-meter",
			Name:        "Password Meter",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Password input with strength meter (progress bar + label) and optional show toggle.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Password meter with configurable label and min length.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "New Password"
						}
						return form.PasswordMeterWithBoundary(form.PasswordMeterProps{
							Name:       "pm-demo",
							Label:      label,
							ShowToggle: true,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "New Password", QueryParam: "label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Password meter with strength indicator in different states.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6 max-w-sm">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">With toggle & meter</p>`); err != nil {
								return err
							}
							if err := form.PasswordMeter(form.PasswordMeterProps{Name: "pm1", Label: "New Password", ShowToggle: true}).Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `</div><div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Without toggle</p>`); err != nil {
								return err
							}
							if err := form.PasswordMeter(form.PasswordMeterProps{Name: "pm2", Label: "Password"}).Render(ctx, w); err != nil {
								return err
							}
							_, err := io.WriteString(w, `</div></div>`)
							return err
						})
					},
				},
			},
		},

		// ── Foundation / Display — Theme Toggle ─────────────────────────────────
		{
			Slug:        "theme-toggle",
			Name:        "Theme Toggle",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Animated sun/moon theme toggle that persists the selection to localStorage. Uses the swap-rotate animation for smooth icon transition.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Theme toggle — click to switch between light and dark themes.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.ThemeToggleWithBoundary()
					},
					FrameHeight: "80px",
					Tokens:      []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Theme toggle rendered standalone.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 flex items-center gap-4">`); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<span class="text-sm text-base-content/60">Light</span>`); err != nil {
								return err
							}
							if err := ui.ThemeToggle().Render(ctx, w); err != nil {
								return err
							}
							if _, err := io.WriteString(w, `<span class="text-sm text-base-content/60">Dark</span>`); err != nil {
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

		// ── Forms / Inputs — Enhanced Select ────────────────────────────────────
		{
			Slug:        "enhanced-select",
			Name:        "Enhanced Select",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Choices.js-powered select with search, groups, tags, and removal. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Enhanced select with configurable search and remove.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Select country"
						}
						return form.EnhancedSelectWithBoundary(form.EnhancedSelectProps{
							Name:       "country",
							Label:      label,
							Searchable: true,
							Options: []form.EnhancedSelectOption{
								{Value: "us", Label: "United States"},
								{Value: "uk", Label: "United Kingdom"},
								{Value: "ca", Label: "Canada"},
								{Value: "de", Label: "Germany", GroupLabel: "Europe"},
								{Value: "fr", Label: "France", GroupLabel: "Europe"},
							},
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Single, searchable, and multi-remove enhanced selects.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6 max-w-md">`); err != nil {
								return err
							}
							sections := []struct {
								label string
								comp  templ.Component
							}{
								{"Single (searchable)", form.EnhancedSelect(form.EnhancedSelectProps{Name: "es1", Label: "Country", Searchable: true, Options: []form.EnhancedSelectOption{{Value: "us", Label: "US"}, {Value: "uk", Label: "UK"}, {Value: "ca", Label: "Canada"}}})},
								{"Multi (removable)", form.EnhancedSelect(form.EnhancedSelectProps{Name: "es2", Label: "Tags", Multiple: true, Searchable: true, RemoveItems: true, Options: []form.EnhancedSelectOption{{Value: "go", Label: "Go"}, {Value: "templ", Label: "Templ"}, {Value: "htmx", Label: "HTMX"}, {Value: "daisyui", Label: "DaisyUI"}}})},
								{"With groups", form.EnhancedSelect(form.EnhancedSelectProps{Name: "es3", Label: "Framework", Searchable: true, Options: []form.EnhancedSelectOption{{Value: "echo", Label: "Echo", GroupLabel: "Go"}, {Value: "gin", Label: "Gin", GroupLabel: "Go"}, {Value: "react", Label: "React", GroupLabel: "JS"}, {Value: "vue", Label: "Vue", GroupLabel: "JS"}}})},
							}
							for _, s := range sections {
								if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">`+s.label+`</p>`); err != nil {
									return err
								}
								if err := s.comp.Render(ctx, w); err != nil {
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

		// ── Forms / Inputs — Date Picker ────────────────────────────────────────
		{
			Slug:        "date-picker",
			Name:        "Date Picker",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Flatpickr date/time picker with single, range, multiple, time, datetime, month, and week modes. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Date picker with configurable mode.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Select date"
						}
						return form.DatePickerWithBoundary(form.DatePickerProps{
							Name:  "dp-demo",
							Label: label,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Select date", QueryParam: "label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Date, range, time, and datetime picker variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6 max-w-sm">`); err != nil {
								return err
							}
							pickers := []struct {
								label string
								comp  templ.Component
							}{
								{"Single date", form.DatePicker(form.DatePickerProps{Name: "dp1", Label: "Date", Placeholder: "Pick a date"})},
								{"Date range", form.DatePicker(form.DatePickerProps{Name: "dp2", Label: "Range", Mode: form.DatePickerRange, Placeholder: "Select range"})},
								{"Time picker", form.DatePicker(form.DatePickerProps{Name: "dp3", Label: "Time", Mode: form.DatePickerTime, Time24h: true, Placeholder: "Select time"})},
								{"Date + time", form.DatePicker(form.DatePickerProps{Name: "dp4", Label: "Date & time", Mode: form.DatePickerDateTime, Placeholder: "Select date & time"})},
							}
							for _, p := range pickers {
								if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">`+p.label+`</p>`); err != nil {
									return err
								}
								if err := p.comp.Render(ctx, w); err != nil {
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

		// ── Forms / Inputs — Rich Text Editor ───────────────────────────────────
		{
			Slug:        "rich-text-editor",
			Name:        "Rich Text Editor",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Quill rich text editor with snow and bubble themes. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Rich text editor with configurable label and theme.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Content"
						}
						return form.RichTextEditorWithBoundary(form.RichTextEditorProps{
							ID:    "rte-demo",
							Label: label,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Content", QueryParam: "label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Snow and bubble theme editors.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-8">`); err != nil {
								return err
							}
							editors := []struct {
								label string
								comp  templ.Component
							}{
								{"Snow theme", form.RichTextEditor(form.RichTextEditorProps{ID: "rte1", Label: "Editor (snow)", Theme: form.QuillThemeSnow, Height: "200px"})},
								{"Bubble theme", form.RichTextEditor(form.RichTextEditorProps{ID: "rte2", Label: "Editor (bubble)", Theme: form.QuillThemeBubble, Height: "200px"})},
							}
							for _, e := range editors {
								if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">`+e.label+`</p>`); err != nil {
									return err
								}
								if err := e.comp.Render(ctx, w); err != nil {
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

		// ── Forms / Inputs — File Upload ────────────────────────────────────────
		{
			Slug:        "file-upload",
			Name:        "File Upload",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "FilePond file upload with drag-and-drop, image preview, and avatar picker. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "File upload with configurable label and style.",
					RenderFunc: func(params url.Values) templ.Component {
						label := params.Get("label")
						if label == "" {
							label = "Upload files"
						}
						return form.FileUploadWithBoundary(form.FileUploadProps{
							Name:  "fu-demo",
							Label: label,
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Label", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Upload files", QueryParam: "label"},
					},
				},
				{
					Name:        "Examples",
					Description: "Default, image preview, and avatar picker upload styles.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 space-y-6 max-w-md">`); err != nil {
								return err
							}
							uploads := []struct {
								label string
								comp  templ.Component
							}{
								{"Default", form.FileUpload(form.FileUploadProps{Name: "fu1", Label: "Upload documents"})},
								{"Image preview", form.FileUpload(form.FileUploadProps{Name: "fu2", Label: "Upload images", Accept: "image/*", Style: form.FileUploadImagePreview})},
								{"Avatar picker", form.FileUpload(form.FileUploadProps{Name: "fu3", Label: "Profile photo", Accept: "image/*", Style: form.FileUploadAvatar, MaxFiles: 1})},
							}
							for _, u := range uploads {
								if _, err := io.WriteString(w, `<div><p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">`+u.label+`</p>`); err != nil {
									return err
								}
								if err := u.comp.Render(ctx, w); err != nil {
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

		// ── Forms / Inputs — Form Validation ────────────────────────────────────
		{
			Slug:        "form-validation",
			Name:        "Form Validation",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Client-side form validation with required, min/max length, email, pattern, and match rules. Uses vanilla JS.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Form with validation rules.",
					RenderFunc: func(_ url.Values) templ.Component {
						rules := []form.ValidationRule{
							{Field: "name", Label: "Name", Required: true},
							{Field: "email", Label: "Email", Required: true, Type: "email"},
							{Field: "password", Label: "Password", Required: true, MinLength: 8},
							{Field: "confirm", Label: "Confirm password", Required: true, Match: "password"},
						}
						children := seq(
							form.TextInput("name", "Name", "", "", true),
							form.TextInput("email", "Email", "", "", true),
							form.PasswordField(form.PasswordFieldProps{Name: "password", Label: "Password", Required: true, ShowToggle: true}),
							form.PasswordField(form.PasswordFieldProps{Name: "confirm", Label: "Confirm password", Required: true, ShowToggle: true}),
						)
						return form.FormValidationWithBoundary(form.FormValidationProps{ID: "fv-demo", Rules: rules, SubmitText: "Create Account"}, children)
					},
				},
			},
		},

		// ── Data Display / Display — Sortable List ──────────────────────────────
		{
			Slug:        "sortable-list",
			Name:        "Sortable List",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "SortableJS drag-and-drop list with optional handle. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Drag-to-reorder list items.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.SortableListWithBoundary("sort-demo", ui.SortableOptions{Animation: 150}, seq(
							ui.SortableItem("", nil),
							ui.SortableItem("", nil),
							ui.SortableItem("", nil),
						))
					},
				},
			},
		},

		// ── Data Display / Display — Swiper Carousel ────────────────────────────
		{
			Slug:        "swiper-carousel",
			Name:        "Swiper Carousel",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "SwiperJS-powered carousel with navigation, pagination, autoplay, and multiple effects. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Carousel with configurable effect.",
					RenderFunc: func(params url.Values) templ.Component {
						slides := []ui.SwiperSlide{
							{Content: rawHTML(`<div class="bg-primary text-primary-content text-xl font-bold p-12">Slide 1</div>`), Class: "bg-primary/5"},
							{Content: rawHTML(`<div class="bg-secondary text-secondary-content text-xl font-bold p-12">Slide 2</div>`), Class: "bg-secondary/5"},
							{Content: rawHTML(`<div class="bg-accent text-accent-content text-xl font-bold p-12">Slide 3</div>`), Class: "bg-accent/5"},
						}
						return ui.SwiperCarouselWithBoundary(ui.SwiperCarouselProps{
							ID:           "swiper-demo",
							Slides:       slides,
							Navigation:   true,
							Pagination:   true,
							Height:       "300px",
						})
					},
					Tokens: []galleryruntime.DesignToken{},
				},
				{
					Name:        "Examples",
					Description: "Carousel with navigation, pagination, autoplay, and multi-slide variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						slide := func(label string, bg string) ui.SwiperSlide {
							return ui.SwiperSlide{Content: rawHTML(`<div class="flex items-center justify-center h-full text-xl font-bold `+bg+`">`+label+`</div>`)}
						}
						slides := []ui.SwiperSlide{slide("Slide 1", "bg-primary text-primary-content"), slide("Slide 2", "bg-secondary text-secondary-content"), slide("Slide 3", "bg-accent text-accent-content")}
						return ui.SwiperCarousel(ui.SwiperCarouselProps{
							ID:           "swiper-example",
							Slides:       slides,
							Navigation:   true,
							Pagination:   true,
							Autoplay:     true,
							AutoplayDelay: 3000,
							Loop:         true,
							Height:       "300px",
						})
					},
				},
			},
		},

		// ── Table / Data Table ──────────────────────────────────────────────────
		{
			Slug:        "data-table",
			Name:        "Data Table",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Tables",
			Description: "Full-featured data table with search, sort, pagination, and column visibility toggle. Uses Alpine.js for client-side state.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Data table with search and sort enabled.",
					RenderFunc: func(_ url.Values) templ.Component {
						cols := []table.DataTableColumn{
							{ID: "id", Label: "ID", Sortable: true, Searchable: true},
							{ID: "customer", Label: "Customer", Sortable: true, Searchable: true},
							{ID: "price", Label: "Price", Sortable: true},
							{ID: "status", Label: "Status", Sortable: true},
						}
						rows := []table.DataTableRow{
							{"id": "#21001", "customer": "Emily Johnson", "price": "$342", "status": "Ordered"},
							{"id": "#21002", "customer": "Alex Thompson", "price": "$578", "status": "Accepted"},
							{"id": "#21003", "customer": "Sarah Davis", "price": "$215", "status": "On the Way"},
							{"id": "#21004", "customer": "Michael Wilson", "price": "$769", "status": "Delivered"},
							{"id": "#21005", "customer": "Jessica Miller", "price": "$431", "status": "Accepted"},
						}
						return table.DataTableWithBoundary(table.DataTableProps{
							ID:         "dt-demo",
							Columns:    cols,
							Rows:       rows,
							Searchable: true,
							Sortable:   true,
							Striped:    true,
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Data table with various features enabled.",
					RenderFunc: func(_ url.Values) templ.Component {
						cols := []table.DataTableColumn{
							{ID: "id", Label: "Order ID", Sortable: true},
							{ID: "customer", Label: "Customer", Sortable: true},
							{ID: "price", Label: "Amount", Sortable: true},
							{ID: "payment", Label: "Payment"},
							{ID: "status", Label: "Status", Sortable: true},
						}
						rows := sampleDataTableRows()
						return table.DataTable(table.DataTableProps{
							ID:         "dt-example",
							Columns:    cols,
							Rows:       rows,
							Searchable: true,
							Sortable:   true,
							Striped:    true,
							Compact:    true,
						})
					},
				},
			},
		},

		// ── Foundation / Display — Charts ───────────────────────────────────────
		{
			Slug:        "charts",
			Name:        "Charts",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "ApexCharts-powered chart components: area, bar, column, line, pie, donut. Loads from CDN dynamically.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Chart with configurable type.",
					RenderFunc: func(params url.Values) templ.Component {
						chartTypeStr := params.Get("type")
						chartType := ui.ChartArea
						switch chartTypeStr {
						case "bar":
							chartType = ui.ChartBar
						case "column":
							chartType = ui.ChartColumn
						case "line":
							chartType = ui.ChartLine
						case "pie":
							chartType = ui.ChartPie
						case "donut":
							chartType = ui.ChartDonut
						}
						title := params.Get("title")
						if title == "" {
							title = "Sales Overview"
						}
						return ui.ChartWithBoundary(ui.ChartProps{
							ID:     "chart-demo",
							Type:   chartType,
							Title:  title,
							Series: []ui.ChartSeries{{Name: "Sales", Data: []float64{30, 40, 35, 50, 49, 60, 70, 91, 125}}},
							Categories: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep"},
							Height: "300",
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Type", Group: "Style", Type: galleryruntime.TokenTypeSelect, Default: "area", QueryParam: "type", Options: []galleryruntime.TokenOption{
							{Value: "area", Label: "Area"},
							{Value: "bar", Label: "Bar"},
							{Value: "column", Label: "Column"},
							{Value: "line", Label: "Line"},
							{Value: "pie", Label: "Pie"},
							{Value: "donut", Label: "Donut"},
						}},
						{Label: "Title", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Sales Overview", QueryParam: "title"},
					},
				},
				{
					Name:        "Examples",
					Description: "Area, bar, column, line, and pie chart variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						series := []ui.ChartSeries{{Name: "Series 1", Data: []float64{30, 40, 35, 50, 49, 60, 70, 91, 125}}}
						cats := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep"}
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">`); err != nil {
								return err
							}
							charts := []struct {
								id    string
								title string
								ct    ui.ChartType
							}{
								{"ch1", "Area Chart", ui.ChartArea},
								{"ch2", "Column Chart", ui.ChartColumn},
								{"ch3", "Line Chart", ui.ChartLine},
								{"ch4", "Pie Chart", ui.ChartPie},
							}
							for _, c := range charts {
								if err := ui.Chart(ui.ChartProps{ID: c.id, Type: c.ct, Title: c.title, Series: series, Categories: cats, Height: "250"}).Render(ctx, w); err != nil {
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

		// ── Foundation / Display — Sparkline Chart ──────────────────────────────
		{
			Slug:        "sparkline-chart",
			Name:        "Sparkline Chart",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Display",
			Description: "Compact inline ApexCharts sparkline for embedding in stat cards and tables.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sparkline chart with configurable type.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.ChartWithBoundary(ui.ChartProps{
							ID:        "spark-demo",
							Type:      ui.ChartLine,
							Sparkline: true,
							Series:    []ui.ChartSeries{{Name: "", Data: []float64{12, 19, 15, 22, 28, 24, 35, 42, 38, 48, 45, 52}}},
							Height:    "100",
						})
					},
				},
				{
					Name:        "Examples",
					Description: "Multiple sparkline charts.",
					RenderFunc: func(_ url.Values) templ.Component {
						return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
							if _, err := io.WriteString(w, `<div class="p-6 grid grid-cols-1 md:grid-cols-4 gap-4">`); err != nil {
								return err
							}
							sparks := []struct {
								id   string
								data []float64
							}{
								{"sp1", []float64{12, 19, 15, 22, 28, 24, 35, 42, 38, 48, 45, 52}},
								{"sp2", []float64{50, 55, 62, 58, 70, 68, 75, 82, 78, 88, 92, 95}},
								{"sp3", []float64{30, 28, 25, 22, 20, 18, 16, 15, 14, 13, 12, 11}},
								{"sp4", []float64{60, 65, 70, 68, 75, 72, 78, 82, 80, 85, 88, 90}},
							}
							for _, s := range sparks {
								if err := ui.Chart(ui.ChartProps{ID: s.id, Type: ui.ChartArea, Sparkline: true, Series: []ui.ChartSeries{{Data: s.data}}, Height: "80"}).Render(ctx, w); err != nil {
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

		// ── Layout / Pages — Auth ──────────────────────────────────────────────
		{
			Slug:        "page-auth",
			Name:        "Auth Page",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Pages",
			Description: "Full authentication pages with split layout: login, register, forgot password, and reset password.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Auth page with configurable style.",
					RenderFunc: func(params url.Values) templ.Component {
						style := pages.AuthPageStyle(params.Get("style"))
						if style == "" {
							style = pages.AuthLogin
						}
						return pages.AuthPageWithBoundary(pages.AuthPageProps{
							Style:     style,
							BrandName: "go-daisy",
						})
					},
					FrameHeight: "600px",
					Tokens: []galleryruntime.DesignToken{
						{Label: "Style", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "login", QueryParam: "style", Options: []galleryruntime.TokenOption{
							{Value: "login", Label: "Login"},
							{Value: "register", Label: "Register"},
							{Value: "forgot", Label: "Forgot Password"},
							{Value: "reset", Label: "Reset Password"},
						}},
					},
				},
				{
					Name:        "Examples",
					Description: "All four auth page variants.",
					RenderFunc: func(_ url.Values) templ.Component {
						return seq(
							pages.AuthPage(pages.AuthPageProps{Style: pages.AuthLogin, BrandName: "go-daisy"}),
						)
					},
				},
			},
		},

		// ── Layout / Pages — Dashboard ──────────────────────────────────────────
		{
			Slug:        "page-dashboard",
			Name:        "Dashboard Page",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Pages",
			Description: "CRM dashboard with stat cards, sparklines, charts, and a recent orders table.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Full CRM dashboard.",
					RenderFunc: func(_ url.Values) templ.Component {
						return pages.DashboardPageWithBoundary(pages.DashboardPageProps{
							Style: pages.DashboardCRM,
						})
					},
					FrameHeight: "800px",
				},
			},
		},

		// ── Layout / Pages — Chat ──────────────────────────────────────────────
		{
			Slug:        "page-chat",
			Name:        "Chat Layout",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Pages",
			Description: "Full chat interface with conversation sidebar, message list, AI thinking indicator, and chat input.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Chat layout with conversation list and messages.",
					RenderFunc: func(_ url.Values) templ.Component {
						return pages.ChatLayoutWithBoundary(pages.ChatLayoutProps{
							ActiveConversation: "Alice Johnson",
						})
					},
					FrameHeight: "600px",
				},
			},
		},

		// ── Layout / Pages — Settings ───────────────────────────────────────────
		{
			Slug:        "page-settings",
			Name:        "Settings Page",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Pages",
			Description: "Tabbed settings page with Profile, Account, Appearance, and Notification sections.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Full settings page with tabs.",
					RenderFunc: func(_ url.Values) templ.Component {
						return pages.SettingsPageWithBoundary()
					},
					FrameHeight: "700px",
				},
			},
		},

		// ── Layout / Pages — Landing ───────────────────────────────────────────
		{
			Slug:        "page-landing",
			Name:        "Landing Page",
			Category:    galleryruntime.CategoryLayout,
			Subcategory: "Pages",
			Description: "Marketing landing page with animated gradient hero, feature cards, testimonials, and pricing tiers.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Landing page with configurable brand name.",
					RenderFunc: func(params url.Values) templ.Component {
						brand := params.Get("brandName")
						if brand == "" {
							brand = "go-daisy"
						}
						return pages.LandingPageWithBoundary(pages.LandingPageProps{
							BrandName: brand,
							Tagline:   "Type-safe Templ components styled with DaisyUI",
						})
					},
					FrameHeight: "800px",
					Tokens: []galleryruntime.DesignToken{
						{Label: "Brand", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "go-daisy", QueryParam: "brandName"},
					},
				},
			},
		},

		// ── Navigation / Headers — Scroll Topbar ───────────────────────────────
		{
			Slug:        "scroll-topbar",
			Name:        "Scroll-Aware Topbar",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Headers",
			Description: "Topbar that hides on scroll down and reveals on scroll up. Uses requestAnimationFrame for performant scroll detection.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Scroll-aware topbar.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" {
							title = "Dashboard"
						}
						return nav.ScrollTopbarWithBoundary(title)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Title", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Dashboard", QueryParam: "title"},
					},
				},
			},
		},

		// ── Layout / Sidebar — Dense Mode ─────────────────────────────────────
		{
			Slug:        "sidebar-dense",
			Name:        "Sidebar — Dense Mode",
			Category:    galleryruntime.CategoryLayout,
			Description: "Compact sidebar that collapses to icons and expands on hover or click toggle. State persists to localStorage.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Hover over or click the toggle button to expand/collapse.",
					RenderFunc: func(_ url.Values) templ.Component {
						return layout.SidebarDenseWithBoundary(layout.SidebarDenseProps{
							AppName: "go-daisy",
							Groups: []layout.SidebarGroup{
								{Label: "Main", Items: []layout.SidebarItem{
									{Label: "Dashboard", Icon: "lucide--layout-dashboard", Href: "#"},
									{Label: "Analytics", Icon: "lucide--bar-chart-2", Href: "#"},
								}},
								{Label: "Settings", Items: []layout.SidebarItem{
									{Label: "Profile", Icon: "lucide--user", Href: "#"},
									{Label: "Settings", Icon: "lucide--settings", Href: "#"},
								}},
							},
						})
					},
					FrameHeight: "400px",
				},
			},
		},

		// ── Overlays — Command Palette ──────────────────────────────────────────
		{
			Slug:        "command-palette",
			Name:        "Command Palette",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Dropdowns",
			Description: "⌘K command palette modal with search filtering. Opens with Cmd+K or Ctrl+K keyboard shortcut.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Press ⌘K or Ctrl+K to open.",
					RenderFunc: func(params url.Values) templ.Component {
						placeholder := params.Get("placeholder")
						if placeholder == "" {
							placeholder = "Type a command..."
						}
						return ui.CommandPaletteWithBoundary(ui.CommandPaletteProps{
							ID:          "cmd-palette-demo",
							Placeholder: placeholder,
							Open:        true,
							Items: []ui.CommandPaletteItem{
								{Label: "Go to Dashboard", Icon: "lucide--layout-dashboard", Shortcut: "⌘1", Group: "Navigation"},
								{Label: "Go to Settings", Icon: "lucide--settings", Shortcut: "⌘2", Group: "Navigation"},
								{Label: "New Project", Icon: "lucide--plus-circle", Shortcut: "⌘N", Group: "Actions"},
								{Label: "Search Files", Icon: "lucide--search", Shortcut: "⌘P", Group: "Actions"},
							},
							Groups: []string{"Navigation", "Actions"},
						})
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Placeholder", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Type a command...", QueryParam: "placeholder"},
					},
				},
			},
		},

		// ── Feedback / Notifications — Dropdown ────────────────────────────────
		{
			Slug:        "notification-dropdown",
			Name:        "Notification Dropdown",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Notifications",
			Description: "Bell icon with dropdown notification panel. Supports mark-as-read and mark-all-read with live badge counter updates.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Notification dropdown with unread count.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.NotificationDropdownWithBoundary(ui.NotificationDropdownProps{
							ID:          "notif-demo",
							UnreadCount: 3,
							ViewAllHref: "#",
							Items: []ui.NotificationItem{
								{IconClass: "bg-primary/10", IconTextClass: "text-primary", IconName: "lucide--briefcase", Title: "New case assigned", Body: "Johnson v. Smith was assigned to you.", Time: "2 min ago", Unread: true},
								{IconClass: "bg-warning/10", IconTextClass: "text-warning", IconName: "lucide--check-square", Title: "Task deadline tomorrow", Body: "File motion due soon.", Time: "1 hour ago", Unread: true},
								{IconClass: "bg-success/10", IconTextClass: "text-success", IconName: "lucide--user", Title: "Client signed in", Body: "Alice Johnson accessed the portal.", Time: "Yesterday", Unread: false},
							},
						})
					},
				},
			},
			},

		// ── Overlays / Popover ────────────────────────────────────────────────
		{
			Slug:        "popover",
			Name:        "Popover",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Popovers",
			Description: "Floating popover using the native Popover API with manual positioning. Supports click/hover triggers, arrow, and edge-aware placement.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Click",
					Description: "Click the button to open the popover.",
					RenderFunc: func(params url.Values) templ.Component {
						placement := ui.PopoverPlacement(params.Get("placement"))
						if placement == "" { placement = ui.PopoverBottom }
						showArrow := params.Get("showArrow") != "false"
						return ui.PopoverWithBoundary(placement, showArrow, ui.PopoverTriggerClick)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Placement", Group: "Position", Type: galleryruntime.TokenTypeSelect, Default: "bottom", QueryParam: "placement", Options: []galleryruntime.TokenOption{
							{Value: "top", Label: "Top"},
							{Value: "bottom", Label: "Bottom"},
							{Value: "left", Label: "Left"},
							{Value: "right", Label: "Right"},
						}},
						{Label: "Show Arrow", Group: "Style", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "showArrow"},
					},
				},
			},
		},

		// ── Foundation / Aspect Ratio ──────────────────────────────────────
		{
			Slug:        "aspect-ratio",
			Name:        "Aspect Ratio",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Container that maintains a fixed aspect ratio for its content.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Aspect ratio container with configurable ratio.",
					RenderFunc: func(params url.Values) templ.Component {
						ratio := params.Get("ratio")
						if ratio == "" { ratio = "16/9" }
						return ui.AspectRatioWithBoundary(ratio)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Ratio", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "16/9", QueryParam: "ratio", Options: []galleryruntime.TokenOption{
							{Value: "16/9", Label: "16:9"},
							{Value: "4/3", Label: "4:3"},
							{Value: "1/1", Label: "1:1 (Square)"},
							{Value: "3/2", Label: "3:2"},
							{Value: "21/9", Label: "21:9 (Ultrawide)"},
						}},
					},
				},
			},
		},

		// ── Foundation / Separator ─────────────────────────────────────────
		{
			Slug:        "separator",
			Name:        "Separator",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Visual divider between content sections.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Separator with orientation toggle.",
					RenderFunc: func(params url.Values) templ.Component {
						orientation := params.Get("orientation")
						if orientation == "" { orientation = "horizontal" }
						if orientation == "vertical" {
							return row(shared.StrComp("Left"), ui.SeparatorWithBoundary("vertical"), shared.StrComp("Right"))
						}
						return ui.SeparatorWithBoundary("horizontal")
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Orientation", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "horizontal", QueryParam: "orientation", Options: []galleryruntime.TokenOption{
							{Value: "horizontal", Label: "Horizontal"},
							{Value: "vertical", Label: "Vertical"},
						}},
					},
				},
			},
		},

		// ── Overlays / Tooltip ─────────────────────────────────────────────
		{
			Slug:        "tooltip",
			Name:        "Tooltip",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Popovers",
			Description: "Rich popover-based tooltip with positioning and arrow support.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Hover the button to see the tooltip.",
					RenderFunc: func(params url.Values) templ.Component {
						tip := params.Get("tip")
						if tip == "" { tip = "Helpful hint" }
						position := params.Get("position")
						if position == "" { position = "top" }
						return ui.TooltipWithBoundary(tip, position)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Tip", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Helpful hint", QueryParam: "tip"},
						{Label: "Position", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "top", QueryParam: "position", Options: []galleryruntime.TokenOption{
							{Value: "top", Label: "Top"},
							{Value: "bottom", Label: "Bottom"},
							{Value: "left", Label: "Left"},
							{Value: "right", Label: "Right"},
						}},
					},
				},
			},
		},

		// ── Overlays / Hover Card ──────────────────────────────────────────
		{
			Slug:        "hover-card",
			Name:        "Hover Card",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Popovers",
			Description: "Rich card displayed on hover using the Popover component.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Hover the button to see the card.",
					RenderFunc: func(params url.Values) templ.Component {
						side := params.Get("side")
						if side == "" { side = "bottom" }
						return ui.HoverCardWithBoundary(side)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Side", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "bottom", QueryParam: "side", Options: []galleryruntime.TokenOption{
							{Value: "top", Label: "Top"},
							{Value: "bottom", Label: "Bottom"},
							{Value: "left", Label: "Left"},
							{Value: "right", Label: "Right"},
						}},
					},
				},
			},
		},

		// ── Foundation / Collapsible ───────────────────────────────────────
		{
			Slug:        "collapsible",
			Name:        "Collapsible",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Single-item animated collapse/expand panel with optional icon.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Click the title to expand/collapse the content.",
					RenderFunc: func(params url.Values) templ.Component {
						title := params.Get("title")
						if title == "" { title = "Collapsible Section" }
						open := params.Get("open") == "true"
						return ui.CollapsibleWithBoundary(title, open)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Title", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "Collapsible Section", QueryParam: "title"},
						{Label: "Open", Group: "State", Type: galleryruntime.TokenTypeBool, Default: "false", QueryParam: "open"},
					},
				},
			},
		},

		// ── Foundation / Icon ──────────────────────────────────────────────
		{
			Slug:        "icon",
			Name:        "Icon",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Display",
			Description: "Accessible Iconify icon with colon-format names and size presets.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Configurable icon name and size.",
					RenderFunc: func(params url.Values) templ.Component {
						name := params.Get("name")
						if name == "" { name = "lucide:star" }
						size := params.Get("size")
						if size == "" { size = "md" }
						return ui.IconWithBoundary(name, size)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Name", Group: "Content", Type: galleryruntime.TokenTypeText, Default: "lucide:star", QueryParam: "name"},
						{Label: "Size", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "md", QueryParam: "size", Options: []galleryruntime.TokenOption{
							{Value: "xs", Label: "XS"},
							{Value: "sm", Label: "SM"},
							{Value: "md", Label: "MD"},
							{Value: "lg", Label: "LG"},
							{Value: "xl", Label: "XL"},
						}},
					},
				},
			},
		},

		// ── Overlays / Sheet ───────────────────────────────────────────────
		{
			Slug:        "sheet",
			Name:        "Sheet",
			Category:    galleryruntime.CategoryOverlays,
			Subcategory: "Panels",
			Description: "Floating slide-out panel using DaisyUI drawer overlay.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Sheet sliding in from the configurable side.",
					RenderFunc: func(params url.Values) templ.Component {
						side := params.Get("side")
						if side == "" { side = "left" }
						open := params.Get("open") != "false"
						return ui.SheetWithBoundary(side, open)
					},
					Tokens: []galleryruntime.DesignToken{
						{Label: "Side", Group: "Layout", Type: galleryruntime.TokenTypeSelect, Default: "left", QueryParam: "side", Options: []galleryruntime.TokenOption{
							{Value: "left", Label: "Left"},
							{Value: "right", Label: "Right"},
						}},
						{Label: "Open", Group: "State", Type: galleryruntime.TokenTypeBool, Default: "true", QueryParam: "open"},
					},
				},
			},
		},
	}

	return append(components, additionalComponents()...)
}

// ── helpers used by new real-component entries ────────────────────────────────

type orderRow struct {
	ID, Customer, Price, Payment, Status string
}

var sampleOrders = []orderRow{
	{"#21001", "Emily Johnson", "$342", "Paid", "Ordered"},
	{"#21002", "Alex Thompson", "$578", "Paid", "Accepted"},
	{"#21003", "Sarah Davis", "$215", "Pending", "On the Way"},
	{"#21004", "Michael Wilson", "$769", "Pending", "Delivered"},
	{"#21005", "Jessica Miller", "$431", "Paid", "Accepted"},
	{"#21006", "Brian Anderson", "$622", "Paid", "Ordered"},
	{"#21007", "Olivia Smith", "$894", "Pending", "On the Way"},
	{"#21008", "Daniel Robinson", "$156", "Paid", "Delivered"},
	{"#21009", "Emma Garcia", "$497", "Pending", "Ordered"},
	{"#21010", "Christopher Baker", "$783", "Paid", "Accepted"},
}

func makeOrdersTableBody(n int) templ.Component {
	if n > len(sampleOrders) {
		n = len(sampleOrders)
	}
	rows := make([]templ.Component, n)
	for i, o := range sampleOrders[:n] {
		o := o
		paymentBadge := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			cls := "badge badge-soft"
			if o.Payment == "Paid" {
				cls += " badge-success"
			} else {
				cls += " badge-error"
			}
			_, err := fmt.Fprintf(w, `<div aria-label="Badge" class="%s">%s</div>`, cls, o.Payment)
			return err
		})
		rows[i] = withChildren(table.TableRowWithBoundary("", true), seq(
			withChildren(table.TableCellWithBoundary("font-medium"), rawHTML(o.ID)),
			withChildren(table.TableCellWithBoundary(""), rawHTML(o.Customer)),
			withChildren(table.TableCellWithBoundary("text-sm font-medium"), rawHTML(o.Price)),
			withChildren(table.TableCellWithBoundary(""), paymentBadge),
			withChildren(table.TableCellWithBoundary("text-sm"), rawHTML(o.Status)),
			withChildren(table.TableCellWithBoundary("text-right"), ui.ButtonWithBoundary("#", ui.ButtonGhost, ui.ButtonSM, ui.ButtonStyleDefault, ui.ButtonTypeButton, ui.ButtonShapeSquare, "lucide--eye", false, false)),
		))
	}
	return seq(rows...)
}

func makeOrdersTable(n int) templ.Component {
	return seq(
		withChildren(table.TableHeadWithBoundary(), withChildren(table.TableHeadRowWithBoundary(), seq(
			table.TableHeadCellWithBoundary("ID"),
			table.TableHeadCellWithBoundary("Customer"),
			table.TableHeadCellWithBoundary("Price"),
			table.TableHeadCellWithBoundary("Payment"),
			table.TableHeadCellWithBoundary("Status"),
			table.TableHeadCellWithBoundary("Action"),
		))),
		withChildren(table.TableBodyWithBoundary(), makeOrdersTableBody(n)),
	)
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

func sampleDataTableRows() []table.DataTableRow {
	return []table.DataTableRow{
		{"id": "#21001", "customer": "Emily Johnson", "price": "$342", "payment": "Paid", "status": "Ordered"},
		{"id": "#21002", "customer": "Alex Thompson", "price": "$578", "payment": "Paid", "status": "Accepted"},
		{"id": "#21003", "customer": "Sarah Davis", "price": "$215", "payment": "Pending", "status": "On the Way"},
		{"id": "#21004", "customer": "Michael Wilson", "price": "$769", "payment": "Pending", "status": "Delivered"},
		{"id": "#21005", "customer": "Jessica Miller", "price": "$431", "payment": "Paid", "status": "Accepted"},
		{"id": "#21006", "customer": "Brian Anderson", "price": "$622", "payment": "Paid", "status": "Ordered"},
		{"id": "#21007", "customer": "Olivia Smith", "price": "$894", "payment": "Pending", "status": "On the Way"},
		{"id": "#21008", "customer": "Daniel Robinson", "price": "$156", "payment": "Paid", "status": "Delivered"},
	}
}

func parseTabsStyle(s string) ui.TabsStyle {
	switch s {
	case "border":
		return ui.TabsBorder
	case "box":
		return ui.TabsBox
	default:
		return ui.TabsLift
	}
}

func parseTabsSize(s string) ui.TabsSize {
	switch s {
	case "xs":
		return ui.TabsXS
	case "sm":
		return ui.TabsSM
	case "lg":
		return ui.TabsLG
	case "xl":
		return ui.TabsXL
	default:
		return ui.TabsMD
	}
}

func tabsContent(text string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	})
}
