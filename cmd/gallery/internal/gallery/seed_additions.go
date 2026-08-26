package gallery

import (
	"net/url"

	"github.com/emergent-company/go-daisy/components/form"
	"github.com/emergent-company/go-daisy/components/head"
	"github.com/emergent-company/go-daisy/components/nav"
	"github.com/emergent-company/go-daisy/components/schemaform"
	"github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/galleryruntime"
	"github.com/a-h/templ"
)

func additionalComponents() []galleryruntime.GalleryComponent {
	return []galleryruntime.GalleryComponent{
		// ── Basics / Aura ──────────────────────────────────────────────────────
		{
			Slug:        "aura",
			Name:        "Aura",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Effects",
			Description: "DaisyUI 5 aura effect — animated border glow that responds to cursor position. Wraps any child element.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Aura wrapping a card.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.AuraWithBoundary(ui.DashboardCardWithBoundary("Aura Effect", "Hover over this card to see the glow"))
					},
					FrameHeight: "180px",
				},
			},
		},

		// ── Navigation / Megamenu ─────────────────────────────────────────────
		{
			Slug:        "megamenu",
			Name:        "Megamenu",
			Category:    galleryruntime.CategoryNavigation,
			Subcategory: "Headers",
			Description: "DaisyUI 5 megamenu — horizontal nav bar with multi-column popover dropdowns.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Megamenu with products, resources, and company columns.",
					RenderFunc: func(_ url.Values) templ.Component {
						return nav.MegamenuWithBoundary([]nav.MegamenuItem{
							{Label: "Home", Href: "#"},
							{
								Label: "Products",
								Columns: []nav.MegamenuColumn{
									{
										Title: "Software",
										Links: []nav.MegamenuLink{
											{Label: "Dashboard", Href: "#", Icon: "lucide--layout-dashboard"},
											{Label: "Analytics", Href: "#", Icon: "lucide--bar-chart-2"},
											{Label: "Automation", Href: "#", Icon: "lucide--zap"},
										},
									},
									{
										Title: "Services",
										Links: []nav.MegamenuLink{
											{Label: "Consulting", Href: "#", Icon: "lucide--briefcase"},
											{Label: "Support", Href: "#", Icon: "lucide--headphones"},
											{Label: "Training", Href: "#", Icon: "lucide--book-open"},
										},
									},
								},
							},
							{
								Label: "Resources",
								Columns: []nav.MegamenuColumn{
									{
										Title: "Learn",
										Links: []nav.MegamenuLink{
											{Label: "Documentation", Href: "#", Icon: "lucide--file-text"},
											{Label: "API Reference", Href: "#", Icon: "lucide--code-2"},
											{Label: "Guides", Href: "#", Icon: "lucide--book"},
										},
									},
								},
							},
							{Label: "Contact", Href: "#"},
						})
					},
					FrameHeight: "280px",
				},
			},
		},

		// ── Forms / OTP Input ──────────────────────────────────────────────────
		{
			Slug:        "otp-input",
			Name:        "OTP Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Multi-digit one-time password input with auto-advance, paste support, and keyboard navigation.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "6-digit OTP input.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.OTPInputWithBoundary("otp-demo", 6)
					},
					FrameHeight: "120px",
				},
			},
		},

		// ── Forms / Color Input ────────────────────────────────────────────────
		{
			Slug:        "color-input",
			Name:        "Color Picker",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Native color picker input with DaisyUI styling.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Color picker with default teal value.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.ColorInputWithBoundary("color-demo", "#0d9488")
					},
					FrameHeight: "100px",
				},
			},
		},

		// ── Forms / Datalist Input ─────────────────────────────────────────────
		{
			Slug:        "datalist-input",
			Name:        "Datalist Autocomplete",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Text input with native browser autocomplete via <datalist> element.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Country name autocomplete.",
					RenderFunc: func(_ url.Values) templ.Component {
						opts := []form.DatalistOption{
							{Value: "US", Label: "United States"},
							{Value: "GB", Label: "United Kingdom"},
							{Value: "DE", Label: "Germany"},
							{Value: "FR", Label: "France"},
							{Value: "JP", Label: "Japan"},
							{Value: "BR", Label: "Brazil"},
							{Value: "IN", Label: "India"},
							{Value: "CA", Label: "Canada"},
						}
						return form.DatalistInputWithBoundary("country-demo", "Search countries...", opts)
					},
					FrameHeight: "100px",
				},
			},
		},

		// ── Forms / TagList ────────────────────────────────────────────────────
		{
			Slug:        "tags-list",
			Name:        "Tag List",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Alpine-managed tag input — add tags with Enter/comma, remove with X, submit as hidden inputs.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Tag list with default values.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.TagListWithBoundary(form.TagListProps{
							ID: "demo-tags", Name: "tags",
							Values:      []string{"production", "gpu"},
							Placeholder: "Add a tag...",
						})
					},
					FrameHeight: "120px",
				},
			},
		},

		// ── Forms / Combobox ────────────────────────────────────────────────────
		{
			Slug:        "combobox",
			Name:        "Combobox",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Alpine searchable select — client-mode filter or server-mode lazy search, single/multi, keyboard nav.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Client Mode",
					Description: "Combobox with pre-rendered options.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.ComboboxWithBoundary(form.ComboboxProps{
							ID:   "demo-combo",
							Name: "industry",
							Mode: form.ComboboxSingle,
							Source: form.ComboboxSource{Static: []form.ComboboxOption{
								{Value: "tech", Label: "Technology"},
								{Value: "health", Label: "Healthcare"},
								{Value: "finance", Label: "Finance"},
								{Value: "retail", Label: "Retail"},
							}},
							Placeholder: "Select industry...",
						})
					},
					FrameHeight: "250px",
				},
			},
		},

		// ── Forms / StructuredInput ─────────────────────────────────────────────
		{
			Slug:        "structured-input",
			Name:        "Structured Input",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Alpine repeatable key-value rows with text+select columns. Add/remove rows. Submits as name[0][key].",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Labels with key+value columns.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.StructuredInputWithBoundary(form.StructuredInputProps{
							ID: "demo-structured", Name: "labels",
							Columns: []form.StructuredColumn{
								{Key: "key", Label: "Key", Placeholder: "key", Separator: "="},
								{Key: "value", Label: "Value", Placeholder: "value"},
							},
							Entries: []form.StructuredEntry{
								{"key": "app", "value": "web"},
								{"key": "env", "value": "prod"},
							},
						})
					},
					FrameHeight: "200px",
				},
			},
		},

		// ── Forms / Palette ─────────────────────────────────────────────────────
		{
			Slug:        "palette",
			Name:        "Palette",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Alpine color swatch picker — Tailwind hues × shades, neutral colors, hex display, reset button.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Full palette with hex display.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PaletteWithBoundary(form.PaletteProps{
							ID: "demo-palette", ShowHex: true,
						})
					},
					FrameHeight: "400px",
				},
			},
		},

		// ── Forms / SelectShell ─────────────────────────────────────────────────
		{
			Slug:        "select-shell",
			Name:        "Select Shell",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Composable Alpine dropdown trigger+panel. Chevron rotates, click-outside closes, custom children inside.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Select shell with stat options.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.SelectShellWithBoundary(form.SelectShellProps{
							ID: "demo-shell", Label: "Choose", Placeholder: "Select...",
							ValueExpr: "picked",
						})
					},
					FrameHeight: "200px",
				},
			},
		},

		// ── Feedback / Banner ───────────────────────────────────────────────────
		{
			Slug:        "banner",
			Name:        "Banner",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Alerts",
			Description: "Alpine dismissable announcement bar — 5 color variants, CTA button, cookie-consent mode.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Default",
					Description: "Info banner with dismiss.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.BannerWithBoundary(ui.BannerProps{
							Description: "Limited Time Offer! Explore exclusive deals & savings",
						})
					},
					FrameHeight: "80px",
				},
				{
					Name:        "With CTA",
					Description: "Success banner with call-to-action button.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.BannerWithBoundary(ui.BannerProps{
							Description: "Get Fit Anywhere, Anytime!",
							Variant:     ui.BannerSuccess,
							CTA:         &ui.BannerCTAProps{Label: "Start free trial", Href: "#"},
						})
					},
					FrameHeight: "80px",
				},
			},
		},

		// ── Feedback / Toast Queue ──────────────────────────────────────────────
		{
			Slug:        "toast-queue",
			Name:        "Toast Queue",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Notifications",
			Description: "Alpine toast notification queue — stack multiple toasts, auto-dismiss, transitions, click-to-dismiss.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Toast queue container.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.ToastQueueWithBoundary()
					},
					FrameHeight: "300px",
				},
			},
		},

		// ── Data Display / CodeBlock ────────────────────────────────────────────
		{
			Slug:        "code-block",
			Name:        "Code Block",
			Category:    galleryruntime.CategoryDataDisplay,
			Subcategory: "Content",
			Description: "Code display with language label, copy button, scrollable height. Plain <pre> fallback or Chroma highlighter.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Go Code",
					Description: "Go source code example.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.CodeBlockWithBoundary(ui.CodeBlockProps{
							Language: "go",
							Label:    "main.go",
							Code:     `package main\n\nimport "fmt"\n\nfunc main() {\n    fmt.Println("hello, world")\n}`,
						})
					},
					FrameHeight: "200px",
				},
			},
		},

		// ── Data Display / Skeleton Loaders ────────────────────────────────────
		{
			Slug:        "skeleton",
			Name:        "Skeleton Loaders",
			Category:    galleryruntime.CategoryFeedback,
			Subcategory: "Loading",
			Description: "DaisyUI skeleton shimmer placeholders — card, table row, text lines, avatar, button, stats.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Card",
					Description: "Skeleton card placeholder.",
					RenderFunc:  func(_ url.Values) templ.Component { return ui.SkeletonCardWithBoundary() },
					FrameHeight: "250px",
				},
				{
					Name:        "Table Row",
					Description: "Skeleton table row.",
					RenderFunc:  func(_ url.Values) templ.Component { return ui.SkeletonTableRowWithBoundary(4) },
					FrameHeight: "80px",
				},
				{
					Name:        "Text Lines",
					Description: "Skeleton text lines.",
					RenderFunc:  func(_ url.Values) templ.Component { return ui.SkeletonTextWithBoundary("w-full", "w-3/4", "w-1/2") },
					FrameHeight: "100px",
				},
				{
					Name:        "Avatar",
					Description: "Skeleton avatar circle.",
					RenderFunc:  func(_ url.Values) templ.Component { return ui.SkeletonAvatarWithBoundary("size-16") },
					FrameHeight: "100px",
				},
				{
					Name:        "Button",
					Description: "Skeleton button.",
					RenderFunc:  func(_ url.Values) templ.Component { return ui.SkeletonButtonWithBoundary() },
					FrameHeight: "80px",
				},
				{
					Name:        "Stats",
					Description: "Skeleton stat card.",
					RenderFunc:  func(_ url.Values) templ.Component { return ui.SkeletonStatsWithBoundary() },
					FrameHeight: "80px",
				},
			},
		},

		// ── Layout / Head Dependencies ─────────────────────────────────────────
		{
			Slug:        "dependencies",
			Name:        "Head Dependencies",
			Category:    galleryruntime.CategoryFoundation,
			Subcategory: "Setup",
			Description: "Unified <head> asset tag emitter — CSS, HTMX, Alpine, Morph, Stimulus, SSE, WS with cache-busting.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Full",
					Description: "All dependencies enabled.",
					RenderFunc: func(_ url.Values) templ.Component {
						return head.DependenciesWithBoundary(head.DepsProps{
							Alpine:  true,
							Morph:   true,
							Stimulus: true,
							SSE:     true,
							WS:      true,
						})
					},
					FrameHeight: "100px",
				},
			},
		},

		// ── Forms / Schema Form Fields ─────────────────────────────────────────
		{
			Slug:        "schema-fields",
			Name:        "Schema Form",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Layout",
			Description: "Generate form controls from JSON Schema subset. Walk(schema, defaults, values, allowList) → go-daisy form fields.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Example",
					Description: "Generated form from pod config schema.",
					RenderFunc: func(_ url.Values) templ.Component {
						fields := schemaform.Walk(
							map[string]any{"properties": map[string]any{
								"replicaCount": map[string]any{"type": "integer", "title": "Replicas"},
								"imageTag":     map[string]any{"type": "string", "title": "Image tag"},
								"tlsEnabled":   map[string]any{"type": "boolean", "title": "TLS enabled"},
								"serviceType":  map[string]any{"enum": []any{"ClusterIP", "LoadBalancer", "NodePort"}},
							}},
							map[string]any{"replicaCount": float64(3), "serviceType": "ClusterIP"},
							map[string]any{"serviceType": "LoadBalancer"},
							map[string]schemaform.AllowMode{},
						)
						return schemaform.FieldsWithBoundary(schemaform.FieldsProps{
							Fields:     fields,
							NamePrefix: "config",
						})
					},
					FrameHeight: "350px",
				},
			},
		},
	}
}
