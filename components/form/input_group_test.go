package form

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/shared"
)

func renderFormComponent(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func TestInputGroupPrefixSuffix(t *testing.T) {
	got := renderFormComponent(t, shared.RenderInto(
		InputGroup(InputGroupProps{Prefix: "https://", Suffix: ".com", Class: "max-w-sm"}),
		shared.StrComp(`<input class="input join-item w-full" aria-label="Host"/>`),
	))

	for _, want := range []string{
		"join w-full",
		"max-w-sm",
		`>https://</span>`,
		`>.com</span>`,
		"join-item",
		`aria-label="Host"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestInputGroupIconsWhenTextAbsent(t *testing.T) {
	got := renderFormComponent(t, shared.RenderInto(
		InputGroup(InputGroupProps{PrefixIcon: "lucide--link", SuffixIcon: "lucide--search"}),
		shared.StrComp(`<input class="input join-item w-full" aria-label="Query"/>`),
	))
	if !strings.Contains(got, "lucide--link") || !strings.Contains(got, "lucide--search") {
		t.Errorf("icon addons missing:\n%s", got)
	}
}

func TestInputGroupTextWinsOverIcon(t *testing.T) {
	got := renderFormComponent(t, shared.RenderInto(
		InputGroup(InputGroupProps{Prefix: "USD", PrefixIcon: "lucide--dollar-sign"}),
		shared.StrComp(`<input class="input join-item w-full" aria-label="Amount"/>`),
	))
	if strings.Contains(got, "lucide--dollar-sign") {
		t.Errorf("icon should be ignored when text is set:\n%s", got)
	}
	if !strings.Contains(got, ">USD</span>") {
		t.Errorf("prefix text missing:\n%s", got)
	}
}

func TestInputGroupOmitsEmptyAddons(t *testing.T) {
	got := renderFormComponent(t, shared.RenderInto(
		InputGroup(InputGroupProps{}),
		shared.StrComp(`<input class="input join-item w-full"/>`),
	))
	if n := strings.Count(got, "join-item flex"); n != 0 {
		t.Errorf("expected no addon spans, got %d:\n%s", n, got)
	}
}

func TestLabelRendersForAndLabelText(t *testing.T) {
	got := renderFormComponent(t, shared.RenderInto(
		Label(LabelProps{For: "email"}),
		shared.StrComp("Email address"),
	))
	for _, want := range []string{
		`for="email"`,
		"label pb-1",
		"label-text",
		">Email address</span>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}
