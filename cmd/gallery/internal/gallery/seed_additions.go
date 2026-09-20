package gallery

import (
	"net/url"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/components/alpine"
	"github.com/emergent-company/go-daisy/components/form"
	"github.com/emergent-company/go-daisy/components/head"
	"github.com/emergent-company/go-daisy/components/nav"
	"github.com/emergent-company/go-daisy/components/schemaform"
	"github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/galleryruntime"
	"github.com/emergent-company/go-daisy/shared"
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
						return ui.ToastQueueWithBoundary(
							alpine.ToastItem{Type: "success", Message: "Changes saved", Action: "Undo"},
							alpine.ToastItem{Type: "info", Message: "New update available", Action: "View"},
							alpine.ToastItem{Type: "error", Message: "Upload failed — file too large"},
						)
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
							Alpine:   true,
							Morph:    true,
							Stimulus: true,
							SSE:      true,
							WS:       true,
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

		// ── Forms / Icon Pick ──────────────────────────────────────────────────
		{
			Slug:        "icon-picker",
			Name:        "Icon Picker",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Searchable Popover icon listbox with roving-tabindex keyboard navigation, a default/reset tile, and a hidden form input. Catalog-agnostic: the caller injects the options.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "A small Lucide catalog with a default tile.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.IconPickerWithBoundary(ui.IconPickerProps{
							ID:           "demo-icon-picker",
							Name:         "icon",
							Value:        "lucide--file-text",
							DefaultValue: "",
							DefaultLabel: "Default icon",
							DefaultClass: "lucide--box",
							Options: []ui.IconPickerOption{
								{Name: "file-text", Label: "File Text", Class: "lucide--file-text"},
								{Name: "database", Label: "Database", Class: "lucide--database"},
								{Name: "braces", Label: "Braces", Class: "lucide--braces"},
								{Name: "user", Label: "User", Class: "lucide--user"},
								{Name: "calendar", Label: "Calendar", Class: "lucide--calendar"},
								{Name: "mail", Label: "Mail", Class: "lucide--mail"},
								{Name: "tag", Label: "Tag", Class: "lucide--tag"},
								{Name: "settings", Label: "Settings", Class: "lucide--settings"},
								{Name: "git-branch", Label: "Git Branch", Class: "lucide--git-branch"},
								{Name: "package", Label: "Package", Class: "lucide--package"},
								{Name: "image", Label: "Image", Class: "lucide--image"},
								{Name: "link", Label: "Link", Class: "lucide--link"},
							},
						})
					},
					FrameHeight: "300px",
				},
			},
		},

		// ── Forms / Color Picker ───────────────────────────────────────────────
		{
			Slug:        "color-picker",
			Name:        "Color Picker",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Inputs",
			Description: "Text colour field paired with a native colour swatch and optional preset palette. Submits as a normal named input and works without JS.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Presets, native swatch, and a clear affordance.",
					RenderFunc: func(_ url.Values) templ.Component {
						return ui.ColorPickerWithBoundary(ui.ColorPickerProps{
							ID:         "demo-color-picker",
							Name:       "color",
							Value:      "#4F46E5",
							AllowEmpty: true,
							Presets: []string{
								"#4F46E5", "#0EA5E9", "#10B981",
								"#F59E0B", "#EF4444", "#8B5CF6",
							},
						})
					},
					FrameHeight: "160px",
				},
			},
		},

		// ── Forms / Prompt Bar — Agent Picker ──────────────────────────────────
		{
			Slug:        "prompt-bar-agent",
			Name:        "Prompt Bar — Agent Picker",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Flagship prompt bar composing an agent picker with model, mode and feature pickers via slots. Pickers submit hidden inputs; loading and running states swap the send action.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Agent, model, mode and feature pickers with token counter and toolbar icons.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Ask the agent anything…",
							Agents: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Description: "Deep research with citations", Icon: "lucide--bot", Group: "Assistants", Meta: "128k", Selected: true},
								{Value: "coder", Label: "Code Companion", Description: "Writes and reviews code", Icon: "lucide--code-2", Group: "Assistants"},
								{Value: "analyst", Label: "Data Analyst", Description: "Charts, SQL and summaries", Icon: "lucide--bar-chart-2", Group: "Assistants"},
								{Value: "writer", Label: "Copywriter", Description: "Product and marketing copy", Icon: "lucide--pen-line", Group: "Specialists"},
								{Value: "tutor", Label: "Tutor", Description: "Explains concepts step by step", Icon: "lucide--graduation-cap", Group: "Specialists"},
							},
							SelectedAgent: "Research Buddy",
							AgentName:     "agent",
							Models: []form.PromptBarPickerItem{
								{Value: "gpt-4o", Label: "GPT-4o", Icon: "lucide--cpu", Meta: "128k", Selected: true},
								{Value: "claude-3-7", Label: "Claude 3.7 Sonnet", Icon: "lucide--cpu", Meta: "200k"},
								{Value: "gemini-2-5", Label: "Gemini 2.5 Pro", Icon: "lucide--cpu", Meta: "1M"},
							},
							SelectedModel: "GPT-4o",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "chat", Label: "Chat", Icon: "lucide--message-square", Selected: true},
								{Value: "plan", Label: "Plan", Icon: "lucide--map"},
								{Value: "build", Label: "Build", Icon: "lucide--hammer"},
							},
							SelectedMode: "Chat",
							ModeName:     "mode",
							Features: []form.PromptBarFeatureItem{
								{Value: "web", Label: "Web search", Description: "Browse the web for fresh sources", Icon: "lucide--globe", On: true},
								{Value: "thinking", Label: "Deep thinking", Description: "Extended reasoning pass", Icon: "lucide--lightbulb"},
							},
							FeatureName:      "feature",
							ShowTokenCounter: true,
							TokenCurrent:     42,
							TokenMax:         128,
							ShowAttach:       true,
							ShowImage:        true,
							ShowVoice:        true,
						})
					},
					FrameHeight: "240px",
				},
				{
					Name:        "Running",
					Description: "Active run — the send action is replaced by a stop button.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Working…",
							Agents: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Description: "Deep research with citations", Icon: "lucide--bot", Group: "Assistants", Selected: true},
								{Value: "coder", Label: "Code Companion", Description: "Writes and reviews code", Icon: "lucide--code-2", Group: "Assistants"},
								{Value: "analyst", Label: "Data Analyst", Description: "Charts, SQL and summaries", Icon: "lucide--bar-chart-2", Group: "Assistants"},
								{Value: "writer", Label: "Copywriter", Description: "Product and marketing copy", Icon: "lucide--pen-line", Group: "Specialists"},
							},
							SelectedAgent: "Research Buddy",
							AgentName:     "agent",
							Models: []form.PromptBarPickerItem{
								{Value: "gpt-4o", Label: "GPT-4o", Icon: "lucide--cpu", Meta: "128k", Selected: true},
								{Value: "claude-3-7", Label: "Claude 3.7 Sonnet", Icon: "lucide--cpu", Meta: "200k"},
								{Value: "gemini-2-5", Label: "Gemini 2.5 Pro", Icon: "lucide--cpu", Meta: "1M"},
							},
							SelectedModel: "GPT-4o",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "chat", Label: "Chat", Icon: "lucide--message-square", Selected: true},
								{Value: "plan", Label: "Plan", Icon: "lucide--map"},
								{Value: "build", Label: "Build", Icon: "lucide--hammer"},
							},
							SelectedMode: "Chat",
							ModeName:     "mode",
							Features: []form.PromptBarFeatureItem{
								{Value: "web", Label: "Web search", Description: "Browse the web for fresh sources", Icon: "lucide--globe", On: true},
								{Value: "thinking", Label: "Deep thinking", Description: "Extended reasoning pass", Icon: "lucide--lightbulb"},
							},
							FeatureName:      "feature",
							ShowTokenCounter: true,
							TokenCurrent:     42,
							TokenMax:         128,
							ShowAttach:       true,
							ShowImage:        true,
							ShowVoice:        true,
							Running:          true,
						})
					},
					FrameHeight: "240px",
				},
			},
		},

		// ── Forms / Prompt Bar Picker ──────────────────────────────────────────
		{
			Slug:        "prompt-bar-picker",
			Name:        "Prompt Bar Picker",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Standalone picker dropdown for the prompt bar toolbar — grouped rows with descriptions and metadata, an optional hidden input, and a disabled row.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Interactive",
					Description: "Grouped agent picker with one selected and one disabled row.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarPickerWithBoundary(form.PromptBarPickerProps{
							Label:     "Agent",
							Icon:      "lucide--bot",
							Selected:  "Research Buddy",
							Name:      "agent",
							MenuClass: "w-80",
							Items: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Description: "Deep research with citations", Icon: "lucide--bot", Group: "Assistants", Meta: "128k", Selected: true},
								{Value: "coder", Label: "Code Companion", Description: "Writes and reviews code", Icon: "lucide--code-2", Group: "Assistants"},
								{Value: "analyst", Label: "Data Analyst", Description: "Charts, SQL and summaries", Icon: "lucide--bar-chart-2", Group: "Assistants"},
								{Value: "writer", Label: "Copywriter", Description: "Product and marketing copy", Icon: "lucide--pen-line", Group: "Specialists"},
								{Value: "legacy", Label: "Legacy Agent", Description: "Retired — read only", Icon: "lucide--archive", Group: "Specialists", Disabled: true},
							},
						})
					},
					FrameHeight: "240px",
				},
			},
		},

		// ── Forms / Prompt Bar — Paseo Parity ──────────────────────────────────
		{
			Slug:        "prompt-bar-paseo",
			Name:        "Prompt Bar — Paseo Parity",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Paseo's real composer: one combined model selector (provider glyph + grouped models), a thinking picker in LeftSlot, a bot mode selector, a green icon-only auto-accept feature toggle, and dictation + voice + send on the right. No token counter and no context meter.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Composer",
					Description: "Model, thinking, mode and feature controls with dictation, voice and send.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Message the agent, tag @files, or use /commands and /skills",
							Models: []form.PromptBarPickerItem{
								{Value: "deepseek:v4-flash", Label: "DeepSeek V4 Flash", Icon: "lucide--brain-circuit", Group: "DeepSeek", Meta: "1M", Selected: true},
								{Value: "deepseek:v4", Label: "DeepSeek V4", Icon: "lucide--brain-circuit", Group: "DeepSeek", Meta: "1M"},
								{Value: "anthropic:claude-sonnet-4.5", Label: "Claude Sonnet 4.5", Icon: "lucide--sparkles", Group: "Anthropic", Meta: "200k"},
								{Value: "anthropic:claude-haiku", Label: "Claude Haiku", Icon: "lucide--sparkles", Group: "Anthropic", Meta: "200k"},
								{Value: "openai:gpt-5.4", Label: "GPT-5.4", Icon: "lucide--cpu", Group: "OpenAI", Meta: "400k"},
								{Value: "openai:gpt-5.4-mini", Label: "GPT-5.4 Mini", Icon: "lucide--cpu", Group: "OpenAI", Meta: "400k"},
							},
							SelectedModel: "DeepSeek V4 Flash",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "orchestrator", Label: "Orchestrator", Description: "Plans and delegates work across subagents", Icon: "lucide--bot", Selected: true},
								{Value: "plan", Label: "Plan Mode", Description: "Analyze the codebase without executing tools or edits", Icon: "lucide--shield-ellipsis", Meta: "planning"},
								{Value: "ask", Label: "Always Ask", Description: "Prompts for permission the first time a tool is used", Icon: "lucide--shield", Meta: "safe"},
								{Value: "auto", Label: "Auto mode", Description: "Uses a model classifier to review permission prompts automatically", Icon: "lucide--shield-check", Meta: "moderate"},
								{Value: "full", Label: "Full Access", Description: "Edit files, run commands, and access the network without additional prompts.", Icon: "lucide--shield-off", Meta: "dangerous"},
							},
							SelectedMode: "Orchestrator",
							ModeName:     "mode",
							Features: []form.PromptBarFeatureItem{
								{Value: "auto_accept", Label: "Auto accept", AriaLabel: "Auto accept permission prompts", Icon: "lucide--shield-check", On: true, Toggle: true},
							},
							LeftSlot: form.PromptBarPicker(form.PromptBarPickerProps{
								Label:    "Select thinking option (High)",
								Icon:     "lucide--brain",
								Selected: "High",
								Name:     "thinking",
								Items: []form.PromptBarPickerItem{
									{Value: "low", Label: "Low"},
									{Value: "medium", Label: "Medium"},
									{Value: "high", Label: "High", Selected: true},
									{Value: "think", Label: "Think"},
									{Value: "think-hard", Label: "Think hard"},
									{Value: "xhigh", Label: "Extra high"},
								},
							}),
							RightSlot: shared.Compose(
								ui.Button(ui.ButtonProps{Type: ui.ButtonTypeButton, Variant: ui.ButtonGhost, Size: ui.ButtonSM, Shape: ui.ButtonShapeSquare, Icon: "lucide--mic", AriaLabel: "Start dictation"}),
								ui.Button(ui.ButtonProps{Type: ui.ButtonTypeButton, Variant: ui.ButtonGhost, Size: ui.ButtonSM, Shape: ui.ButtonShapeSquare, Icon: "lucide--audio-lines", AriaLabel: "Enable Voice mode"}),
							),
							Class:            "max-w-5xl",
							Name:             "message",
							ShowAttach:       true,
							ShowImage:        false,
							ShowVoice:        false,
							ShowTokenCounter: false,
						})
					},
					FrameHeight: "240px",
				},
				{
					Name:        "Running",
					Description: "Active run — the send action is replaced by the stop button.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Message the agent, tag @files, or use /commands and /skills",
							Models: []form.PromptBarPickerItem{
								{Value: "deepseek:v4-flash", Label: "DeepSeek V4 Flash", Icon: "lucide--brain-circuit", Group: "DeepSeek", Meta: "1M", Selected: true},
								{Value: "deepseek:v4", Label: "DeepSeek V4", Icon: "lucide--brain-circuit", Group: "DeepSeek", Meta: "1M"},
								{Value: "anthropic:claude-sonnet-4.5", Label: "Claude Sonnet 4.5", Icon: "lucide--sparkles", Group: "Anthropic", Meta: "200k"},
								{Value: "anthropic:claude-haiku", Label: "Claude Haiku", Icon: "lucide--sparkles", Group: "Anthropic", Meta: "200k"},
								{Value: "openai:gpt-5.4", Label: "GPT-5.4", Icon: "lucide--cpu", Group: "OpenAI", Meta: "400k"},
								{Value: "openai:gpt-5.4-mini", Label: "GPT-5.4 Mini", Icon: "lucide--cpu", Group: "OpenAI", Meta: "400k"},
							},
							SelectedModel: "DeepSeek V4 Flash",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "orchestrator", Label: "Orchestrator", Description: "Plans and delegates work across subagents", Icon: "lucide--bot", Selected: true},
								{Value: "plan", Label: "Plan Mode", Description: "Analyze the codebase without executing tools or edits", Icon: "lucide--shield-ellipsis", Meta: "planning"},
								{Value: "ask", Label: "Always Ask", Description: "Prompts for permission the first time a tool is used", Icon: "lucide--shield", Meta: "safe"},
								{Value: "auto", Label: "Auto mode", Description: "Uses a model classifier to review permission prompts automatically", Icon: "lucide--shield-check", Meta: "moderate"},
								{Value: "full", Label: "Full Access", Description: "Edit files, run commands, and access the network without additional prompts.", Icon: "lucide--shield-off", Meta: "dangerous"},
							},
							SelectedMode: "Orchestrator",
							ModeName:     "mode",
							Features: []form.PromptBarFeatureItem{
								{Value: "auto_accept", Label: "Auto accept", AriaLabel: "Auto accept permission prompts", Icon: "lucide--shield-check", On: true, Toggle: true},
							},
							LeftSlot: form.PromptBarPicker(form.PromptBarPickerProps{
								Label:    "Select thinking option (High)",
								Icon:     "lucide--brain",
								Selected: "High",
								Name:     "thinking",
								Items: []form.PromptBarPickerItem{
									{Value: "low", Label: "Low"},
									{Value: "medium", Label: "Medium"},
									{Value: "high", Label: "High", Selected: true},
									{Value: "think", Label: "Think"},
									{Value: "think-hard", Label: "Think hard"},
									{Value: "xhigh", Label: "Extra high"},
								},
							}),
							RightSlot: shared.Compose(
								ui.Button(ui.ButtonProps{Type: ui.ButtonTypeButton, Variant: ui.ButtonGhost, Size: ui.ButtonSM, Shape: ui.ButtonShapeSquare, Icon: "lucide--mic", AriaLabel: "Start dictation"}),
								ui.Button(ui.ButtonProps{Type: ui.ButtonTypeButton, Variant: ui.ButtonGhost, Size: ui.ButtonSM, Shape: ui.ButtonShapeSquare, Icon: "lucide--audio-lines", AriaLabel: "Enable Voice mode"}),
							),
							Class:            "max-w-5xl",
							Name:             "message",
							ShowAttach:       true,
							ShowImage:        false,
							ShowVoice:        false,
							ShowTokenCounter: false,
							Running:          true,
						})
					},
					FrameHeight: "240px",
				},
			},
		},

		// ── Forms / Prompt Bar — Composer ──────────────────────────────────────
		{
			Slug:        "prompt-bar-composer",
			Name:        "Prompt Bar — Composer",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "A working composer: real form submit with Enter-to-send, Shift+Enter newline, an empty-input guard, and an HTMX-wired stop action while running.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Submit",
					Description: "GET form submit — Enter sends, Shift+Enter inserts a newline.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Type a message and press Enter…",
							Action:      "/gallery/render/prompt-bar-composer",
							Method:      "get",
							SubmitLabel: "Send",
							ShowAttach:  true,
							Agents: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Icon: "lucide--bot", Group: "Assistants", Selected: true},
								{Value: "coder", Label: "Code Companion", Icon: "lucide--code-2", Group: "Assistants"},
							},
							SelectedAgent: "Research Buddy",
							AgentName:     "agent",
							Models: []form.PromptBarPickerItem{
								{Value: "gpt-4o", Label: "GPT-4o", Icon: "lucide--cpu", Selected: true},
								{Value: "claude-3-7", Label: "Claude 3.7 Sonnet", Icon: "lucide--cpu"},
							},
							SelectedModel: "GPT-4o",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "chat", Label: "Chat", Icon: "lucide--message-square", Selected: true},
								{Value: "plan", Label: "Plan", Icon: "lucide--map"},
							},
							SelectedMode: "Chat",
							ModeName:     "mode",
							Features: []form.PromptBarFeatureItem{
								{Value: "web", Label: "Web search", Icon: "lucide--globe", On: true},
								{Value: "thinking", Label: "Deep thinking", Icon: "lucide--lightbulb"},
							},
							FeatureName: "feature",
							Name:        "message",
						})
					},
					FrameHeight: "240px",
				},
				{
					Name:        "Running",
					Description: "Active run — the send action is replaced by an HTMX-wired stop button.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Working…",
							Action:      "/gallery/render/prompt-bar-composer",
							Method:      "get",
							SubmitLabel: "Send",
							ShowAttach:  true,
							Running:     true,
							StopAction:  "/gallery/render/prompt-bar-composer",
							StopMethod:  "get",
							Agents: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Icon: "lucide--bot", Group: "Assistants", Selected: true},
								{Value: "coder", Label: "Code Companion", Icon: "lucide--code-2", Group: "Assistants"},
							},
							SelectedAgent: "Research Buddy",
							AgentName:     "agent",
							Models: []form.PromptBarPickerItem{
								{Value: "gpt-4o", Label: "GPT-4o", Icon: "lucide--cpu", Selected: true},
								{Value: "claude-3-7", Label: "Claude 3.7 Sonnet", Icon: "lucide--cpu"},
							},
							SelectedModel: "GPT-4o",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "chat", Label: "Chat", Icon: "lucide--message-square", Selected: true},
								{Value: "plan", Label: "Plan", Icon: "lucide--map"},
							},
							SelectedMode: "Chat",
							ModeName:     "mode",
							Features: []form.PromptBarFeatureItem{
								{Value: "web", Label: "Web search", Icon: "lucide--globe", On: true},
								{Value: "thinking", Label: "Deep thinking", Icon: "lucide--lightbulb"},
							},
							FeatureName: "feature",
							Name:        "message",
						})
					},
					FrameHeight: "240px",
				},
				{
					Name:        "Full",
					Description: "Attachments, tracks and pickers composed in one dense prompt bar.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Message the agent…",
							Action:      "/gallery/render/prompt-bar-composer",
							Method:      "get",
							SubmitLabel: "Send",
							ShowAttach:  true,
							Attachments: []form.PromptBarAttachmentItem{
								{Label: "design.png", Subtitle: "2.1 MB", Thumbnail: "https://picsum.photos/96/96", Kind: form.PromptBarAttachmentImage},
								{Label: "notes.md", Subtitle: "12 KB", Icon: "lucide--file-text", Kind: form.PromptBarAttachmentFile},
								{Label: "Figma board", Subtitle: "figma.com/board/x1", Icon: "lucide--link", Kind: form.PromptBarAttachmentLink, Href: "#"},
							},
							Tracks: []form.PromptBarTrackItem{
								{
									Label: "3/5 tasks", Icon: "lucide--list-checks", Title: "Tasks",
									Rows: []form.PromptBarTrackRow{
										{Label: "Research", Subtitle: "Complete", Icon: "lucide--circle-check", Done: true},
										{Label: "Write", Subtitle: "In progress", Icon: "lucide--circle-dot"},
									},
								},
								{
									Label: "Diff", Icon: "lucide--git-branch", Title: "Changes",
									Segments: []form.PromptBarTrackSegment{
										{Text: "+128", Bucket: form.PromptBarTrackDone},
										{Text: "−42", Bucket: form.PromptBarTrackFailed},
									},
								},
							},
							Agents: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Icon: "lucide--bot", Group: "Assistants", Selected: true},
								{Value: "coder", Label: "Code Companion", Icon: "lucide--code-2", Group: "Assistants"},
							},
							SelectedAgent: "Research Buddy",
							AgentName:     "agent",
							Models: []form.PromptBarPickerItem{
								{Value: "gpt-4o", Label: "GPT-4o", Icon: "lucide--cpu", Selected: true},
								{Value: "claude-3-7", Label: "Claude 3.7 Sonnet", Icon: "lucide--cpu"},
							},
							SelectedModel: "GPT-4o",
							ModelName:     "model",
							Modes: []form.PromptBarPickerItem{
								{Value: "chat", Label: "Chat", Icon: "lucide--message-square", Selected: true},
								{Value: "plan", Label: "Plan", Icon: "lucide--map"},
							},
							SelectedMode: "Chat",
							ModeName:     "mode",
							Name:         "message",
						})
					},
					FrameHeight: "380px",
				},
			},
		},

		// ── Forms / Prompt Bar — Attachments ──────────────────────────────────
		{
			Slug:        "prompt-bar-attachments",
			Name:        "Prompt Bar — Attachments",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Paseo-style attachment chips with thumbnail/file/link leading visuals, an optional subtitle, and a hover-reveal remove action.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Tray",
					Description: "Standalone attachment tray.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarAttachmentsWithBoundary(form.PromptBarAttachmentsProps{
							Items: []form.PromptBarAttachmentItem{
								{Label: "hero-shot.png", Subtitle: "2.1 MB", Thumbnail: "https://picsum.photos/96/96", Kind: form.PromptBarAttachmentImage},
								{Label: "research-notes.md", Subtitle: "12 KB · Markdown", Icon: "lucide--file-text", Kind: form.PromptBarAttachmentFile},
								{Label: "Figma board", Subtitle: "figma.com/board/x1", Icon: "lucide--link", Kind: form.PromptBarAttachmentLink, Href: "#"},
								{Label: "Wireframes", Subtitle: "Final revisions after the design review on Tuesday afternoon", Icon: "lucide--file", Kind: form.PromptBarAttachmentFile},
							},
						})
					},
					FrameHeight: "160px",
				},
				{
					Name:        "In Composer",
					Description: "The same chips rendered above the PromptBar textarea via Attachments.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Add context…",
							Attachments: []form.PromptBarAttachmentItem{
								{Label: "hero-shot.png", Subtitle: "2.1 MB", Thumbnail: "https://picsum.photos/96/96", Kind: form.PromptBarAttachmentImage},
								{Label: "research-notes.md", Subtitle: "12 KB · Markdown", Icon: "lucide--file-text", Kind: form.PromptBarAttachmentFile},
								{Label: "Figma board", Subtitle: "figma.com/board/x1", Icon: "lucide--link", Kind: form.PromptBarAttachmentLink, Href: "#"},
							},
							ShowAttach: true,
							Name:       "message",
						})
					},
					FrameHeight: "280px",
				},
			},
		},

		// ── Forms / Prompt Bar — Tracks ───────────────────────────────────────
		{
			Slug:        "prompt-bar-tracks",
			Name:        "Prompt Bar — Tracks",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Ambient status pills floated above the composer — task progress, subagent states, and a diff stat with colour-coded segments.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Rail",
					Description: "Standalone track rail.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarTracksWithBoundary(form.PromptBarTracksProps{
							Items: []form.PromptBarTrackItem{
								{
									Label: "3/5 tasks", Icon: "lucide--list-checks", Title: "Tasks",
									Rows: []form.PromptBarTrackRow{
										{Label: "Research the API", Subtitle: "Complete", Icon: "lucide--circle-check", Done: true},
										{Label: "Draft the proposal", Subtitle: "In progress", Icon: "lucide--circle-dot"},
										{Label: "Review changes", Subtitle: "Queued", Icon: "lucide--circle"},
									},
								},
								{
									Label: "2 subagents", Icon: "lucide--bot", Title: "Subagents",
									Rows: []form.PromptBarTrackRow{
										{Label: "explorer", Subtitle: "Scanning the repo", Icon: "lucide--search"},
										{Label: "reviewer", Subtitle: "Checking the diff", Icon: "lucide--shield-check"},
									},
								},
								{
									Label: "Diff", Icon: "lucide--git-branch", Title: "Changes",
									Segments: []form.PromptBarTrackSegment{
										{Text: "+128", Bucket: form.PromptBarTrackDone},
										{Text: "−42", Bucket: form.PromptBarTrackFailed},
									},
								},
							},
						})
					},
					FrameHeight: "120px",
				},
				{
					Name:        "Above Composer",
					Description: "The same pills rendered above a PromptBar via Tracks.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Message the agent…",
							Tracks: []form.PromptBarTrackItem{
								{
									Label: "3/5 tasks", Icon: "lucide--list-checks", Title: "Tasks",
									Rows: []form.PromptBarTrackRow{
										{Label: "Research", Icon: "lucide--circle-check", Done: true},
										{Label: "Write", Icon: "lucide--circle-dot"},
									},
								},
								{
									Label: "Diff", Icon: "lucide--git-branch",
									Segments: []form.PromptBarTrackSegment{
										{Text: "+128", Bucket: form.PromptBarTrackDone},
										{Text: "−42", Bucket: form.PromptBarTrackFailed},
									},
								},
							},
							Name: "message",
						})
					},
					FrameHeight: "240px",
				},
			},
		},

		// ── Forms / Prompt Bar — Autocomplete ─────────────────────────────────
		{
			Slug:        "prompt-bar-autocomplete",
			Name:        "Prompt Bar — Autocomplete",
			Category:    galleryruntime.CategoryForms,
			Subcategory: "Prompt Bar",
			Description: "Slash-command and @mention autocomplete — a popover above the input, filtered as you type, keyboard-navigable, replacing the trigger query on accept.",
			Variants: []galleryruntime.GalleryStory{
				{
					Name:        "Commands",
					Description: "Slash command options.",
					RenderFunc: func(_ url.Values) templ.Component {
						return shared.Compose(
							templ.Raw(`<div class="relative mt-52">`),
							form.PromptBarAutocomplete(form.PromptBarAutocompleteProps{
								ID: "autocomplete-commands",
								Options: []form.PromptBarAutocompleteOption{
									{Label: "/help", Value: "/help", Description: "Show available commands", Detail: "⌘ /", Kind: form.PromptBarAutocompleteCommand},
									{Label: "/clear", Value: "/clear", Description: "Clear the conversation", Kind: form.PromptBarAutocompleteCommand},
									{Label: "/plan", Value: "/plan", Description: "Switch to planning mode", Kind: form.PromptBarAutocompleteCommand},
									{Label: "/model", Value: "/model", Description: "Change the model", Detail: "gpt-4o", Kind: form.PromptBarAutocompleteCommand},
									{Label: "/review", Value: "/review", Description: "Review the current diff", Kind: form.PromptBarAutocompleteCommand},
								},
							}),
							templ.Raw(`</div>`),
						)
					},
					FrameHeight: "280px",
				},
				{
					Name:        "Mentions",
					Description: "@ mention options — a directory, files, and an agent.",
					RenderFunc: func(_ url.Values) templ.Component {
						return shared.Compose(
							templ.Raw(`<div class="relative mt-40">`),
							form.PromptBarAutocomplete(form.PromptBarAutocompleteProps{
								ID: "autocomplete-mentions",
								Options: []form.PromptBarAutocompleteOption{
									{Label: "src/components/", Value: "src/components/", Description: "Directory", Kind: form.PromptBarAutocompleteDirectory},
									{Label: "docs/roadmap.md", Value: "docs/roadmap.md", Description: "File", Detail: "12 KB", Kind: form.PromptBarAutocompleteFile},
									{Label: "prompt-bar.templ", Value: "prompt-bar.templ", Description: "File", Detail: "21 KB", Kind: form.PromptBarAutocompleteFile},
									{Label: "Research Buddy", Value: "@research", Description: "Deep research with citations", Kind: form.PromptBarAutocompleteAgent},
								},
							}),
							templ.Raw(`</div>`),
						)
					},
					FrameHeight: "240px",
				},
				{
					Name:        "In composer",
					Description: "Autocomplete anchored above the PromptBar textarea — type / or @ to open.",
					RenderFunc: func(_ url.Values) templ.Component {
						return form.PromptBarWithBoundary(form.PromptBarProps{
							Placeholder: "Type / for commands or @ to mention…",
							Action:      "/gallery/render/prompt-bar-autocomplete",
							Method:      "get",
							SubmitLabel: "Send",
							ShowAttach:  true,
							Autocomplete: &form.PromptBarAutocompleteProps{
								ID: "autocomplete-composer",
								Options: []form.PromptBarAutocompleteOption{
									{Label: "/help", Value: "/help", Description: "Show available commands", Kind: form.PromptBarAutocompleteCommand},
									{Label: "/clear", Value: "/clear", Description: "Clear the conversation", Kind: form.PromptBarAutocompleteCommand},
									{Label: "/plan", Value: "/plan", Description: "Switch to planning mode", Kind: form.PromptBarAutocompleteCommand},
									{Label: "src/main.go", Value: "src/main.go", Kind: form.PromptBarAutocompleteFile},
									{Label: "Research Buddy", Value: "@research", Description: "Deep research with citations", Kind: form.PromptBarAutocompleteAgent},
								},
							},
							Agents: []form.PromptBarPickerItem{
								{Value: "research", Label: "Research Buddy", Icon: "lucide--bot", Group: "Assistants", Selected: true},
								{Value: "coder", Label: "Code Companion", Icon: "lucide--code-2", Group: "Assistants"},
							},
							SelectedAgent: "Research Buddy",
							AgentName:     "agent",
							Models: []form.PromptBarPickerItem{
								{Value: "gpt-4o", Label: "GPT-4o", Icon: "lucide--cpu", Selected: true},
								{Value: "claude-3-7", Label: "Claude 3.7 Sonnet", Icon: "lucide--cpu"},
							},
							SelectedModel: "GPT-4o",
							ModelName:     "model",
							Name:          "message",
						})
					},
					FrameHeight: "240px",
				},
			},
		},
	}
}
