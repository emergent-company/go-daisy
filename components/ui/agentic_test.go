package ui

import (
	"strings"
	"testing"
)

func TestHumanizeToolName(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"":                       "Tool",
		"search-knowledge":       "Search knowledge",
		"run_cypher_query":       "Run cypher query",
		"schema.list-installed":  "Schema list installed",
		"  spaced__tool--name  ": "Spaced tool name",
	}
	for input, want := range cases {
		if got := humanizeToolName(input); got != want {
			t.Errorf("humanizeToolName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestRenderMarkdownGFM(t *testing.T) {
	t.Parallel()

	got := renderMarkdown("**bold**\n\n| A | B |\n| - | - |\n| 1 | 2 |")
	if !strings.Contains(got, "<strong>bold</strong>") {
		t.Fatalf("expected strong tag in %q", got)
	}
	if !strings.Contains(got, "<table>") {
		t.Fatalf("expected GFM table in %q", got)
	}
}

func TestRenderMarkdownDoesNotInjectRawHTML(t *testing.T) {
	t.Parallel()

	got := renderMarkdown("<script>alert('x')</script>\n\nhello")
	if strings.Contains(got, "<script>") {
		t.Fatalf("raw script tag was injected: %q", got)
	}
	if !strings.Contains(got, "hello") {
		t.Fatalf("expected surrounding text to render: %q", got)
	}
}
