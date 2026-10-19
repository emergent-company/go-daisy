package galleryruntime

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// mimeLLMsText is the Content-Type for the machine-readable /llms.txt index.
const mimeLLMsText = "text/plain; charset=utf-8"

// mimeMarkdownText is the Content-Type for the per-component markdown endpoint.
const mimeMarkdownText = "text/markdown; charset=utf-8"

// handleLLMs renders the machine-readable gallery index at /llms.txt.
//
// The body is a Markdown-formatted plain-text document: an H1 title, a
// blockquote describing the library, one "## <Category>" section per gallery
// category (ordered by BuildCategoryGroups), a bulleted component list, and a
// final "## Links" section. Ordering is deterministic so agents can diff it.
func (h *galleryHandler) handleLLMs(c echo.Context) error {
	var b strings.Builder
	b.WriteString("# " + h.title + "\n\n")
	b.WriteString("> The component reference for the github.com/emergent-company/go-daisy Templ component library.\n\n")

	for _, cg := range BuildCategoryGroups(h.components) {
		b.WriteString("## " + string(cg.Name) + "\n")
		for _, comp := range cg.Components {
			fmt.Fprintf(&b, "- [%s](/gallery/docs/%s.md): %s\n", comp.Name, comp.Slug, comp.Description)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Links\n")
	b.WriteString("- [Gallery documentation](/gallery/docs)\n")
	b.WriteString("- [Component detail pages](/gallery/<slug>)\n")

	return c.Blob(http.StatusOK, mimeLLMsText, []byte(b.String()))
}

// handleComponentMarkdown renders a single component's reference document as
// Markdown at /gallery/docs/:slug.md.
func (h *galleryHandler) handleComponentMarkdown(c echo.Context) error {
	// Echo matches a ":name" param up to the next slash, so the ".md" suffix
	// becomes part of the param name: the registered route param is literally
	// "slug.md" and its value includes the extension. Read it and strip the suffix.
	slug := strings.TrimSuffix(c.Param("slug.md"), ".md")
	comp, ok := ComponentBySlug(h.components, slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "component not found")
	}
	return c.Blob(http.StatusOK, mimeMarkdownText, []byte(componentMarkdown(comp)))
}

// componentMarkdown renders the reference document for a single component.
// Layout: H1 name, Category/Subcategory lines, description paragraph, Variants
// bullet list, a Design tokens table (omitted when there are none; the Parameter
// column shows the CSSVar, or the QueryParam prefixed with "?" when the token
// drives the render iframe's URL query string instead of a CSS custom property),
// an Example code block (omitted when no story has Templ or HTML content), and a
// Gallery link back to the interactive detail page.
func componentMarkdown(comp GalleryComponent) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", comp.Name)
	fmt.Fprintf(&b, "Category: %s\n", comp.Category)
	if comp.Subcategory != "" {
		fmt.Fprintf(&b, "Subcategory: %s\n", comp.Subcategory)
	}
	b.WriteString("\n")
	if comp.Description != "" {
		b.WriteString(comp.Description + "\n\n")
	}

	variants := comp.EffectiveVariants()
	if len(variants) > 0 {
		b.WriteString("## Variants\n")
		for _, v := range variants {
			b.WriteString("- **" + v.Name + "**")
			if v.Description != "" {
				b.WriteString(" — " + v.Description)
			}
			b.WriteString("\n")
			for _, se := range v.SubExamples {
				b.WriteString("  - " + se.Label + "\n")
			}
		}
		b.WriteString("\n")
	}

	if tokens := collectTokens(comp); len(tokens) > 0 {
		b.WriteString("## Design tokens\n\n")
		b.WriteString("| Label | Parameter | Type | Default | Options |\n")
		b.WriteString("|---|---|---|---|---|\n")
		for _, t := range tokens {
			label := t.Label
			if label == "" {
				label = t.CSSVar
			}
			parameter := t.CSSVar
			if parameter == "" && t.QueryParam != "" {
				parameter = "?" + t.QueryParam
			}
			options := ""
			if len(t.Options) > 0 {
				values := make([]string, 0, len(t.Options))
				for _, o := range t.Options {
					values = append(values, o.Value)
				}
				options = strings.Join(values, ", ")
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
				escapeCell(label), escapeCell(parameter), string(t.Type), escapeCell(t.Default), escapeCell(options))
		}
		b.WriteString("\n")
	}

	if html, ok := renderExample(variants); ok {
		b.WriteString("## Example\n\n````\n")
		b.WriteString(html)
		b.WriteString("\n````\n\n")
	}

	fmt.Fprintf(&b, "## Gallery\n\n[Open in gallery](/gallery/%s)\n", comp.Slug)

	return b.String()
}

// collectTokens aggregates a component's design tokens with those of its
// variants, in order, so the markdown table covers every controllable parameter.
func collectTokens(comp GalleryComponent) []DesignToken {
	tokens := make([]DesignToken, 0, len(comp.Tokens))
	tokens = append(tokens, comp.Tokens...)
	for _, v := range comp.EffectiveVariants() {
		tokens = append(tokens, v.Tokens...)
	}
	return tokens
}

// renderExample renders the first story that has static content (Templ or HTML)
// as a standalone snippet. Stories whose only content is a RenderFunc are
// skipped because they require request query params to produce output.
func renderExample(variants []GalleryStory) (string, bool) {
	for _, v := range variants {
		if v.RenderFunc != nil {
			continue
		}
		if v.Templ != nil {
			var buf bytes.Buffer
			if err := v.Templ.Render(context.Background(), &buf); err != nil {
				continue
			}
			return buf.String(), true
		}
		if v.HTML != "" {
			return v.HTML, true
		}
	}
	return "", false
}

// escapeCell escapes markdown-table-sensitive characters in a single cell so a
// value containing "|" or a newline cannot break the table layout.
func escapeCell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// markdown is the shared renderer for GalleryComponent.Docs. GFM enables
// tables, strikethrough, task lists, and autolinks. Raw HTML is intentionally
// NOT enabled (no html.WithUnsafe), so inline HTML in Markdown sources is
// escaped instead of being passed through.
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// renderMarkdown converts Markdown source to an HTML string using GitHub
// Flavored Markdown. Raw HTML is escaped.
//
// Empty (or whitespace-only) input returns "". On render error it returns an
// escaped <pre> fallback, which is safe to embed with templ.Raw.
func renderMarkdown(src string) string {
	if strings.TrimSpace(src) == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(src), &buf); err != nil {
		return "<pre>" + html.EscapeString(src) + "</pre>"
	}
	return buf.String()
}
