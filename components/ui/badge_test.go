package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderBadge renders Badge to its HTML string.
func renderBadge(t *testing.T, props BadgeProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Badge(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestBadgeRendersSpanByDefault pins the default span output, so the Href
// addition does not change existing non-link callers.
func TestBadgeRendersSpanByDefault(t *testing.T) {
	got := renderBadge(t, BadgeProps{Label: "Active"})
	if !strings.Contains(got, `<span class="badge gap-1`) {
		t.Errorf("Badge without Href must render a <span>:\n%s", got)
	}
	if strings.Contains(got, `<a `) {
		t.Errorf("Badge without Href must not render an <a>:\n%s", got)
	}
	if !strings.Contains(got, `>Active</span>`) {
		t.Errorf("Badge label missing:\n%s", got)
	}
}

// TestBadgeRendersLinkWhenHrefSet pins the link-styled badge (<a> instead of <span>).
func TestBadgeRendersLinkWhenHrefSet(t *testing.T) {
	got := renderBadge(t, BadgeProps{Label: "Docs", Href: "#"})
	if !strings.Contains(got, `<a href="#" class="badge gap-1`) {
		t.Errorf("Badge with Href must render an <a href>:\n%s", got)
	}
	if strings.Contains(got, `<span class="badge gap-1`) {
		t.Errorf("Badge with Href must not render a <span>:\n%s", got)
	}
}

// TestBadgeDirectLabel pins the DirectLabel flag: the label must be the
// element's direct text content (no nested label <span>), so a consumer that
// reads the badge's text node directly sees the count.
func TestBadgeDirectLabel(t *testing.T) {
	got := renderBadge(t, BadgeProps{Label: "5", DirectLabel: true})
	if !strings.Contains(got, `>5</span>`) {
		t.Errorf("DirectLabel must render the label as direct text content:\n%s", got)
	}
	if strings.Contains(got, `<span class="">5</span>`) {
		t.Errorf("DirectLabel must not wrap the label in a nested <span>:\n%s", got)
	}

	// Non-direct keeps the nested label span for backward compatibility.
	nested := renderBadge(t, BadgeProps{Label: "5"})
	if !strings.Contains(nested, `<span class="">5</span>`) {
		t.Errorf("Default badge must keep the nested label <span>:\n%s", nested)
	}
}

// TestBadgePositionClass pins the Position field is appended to the badge's class.
func TestBadgePositionClass(t *testing.T) {
	got := renderBadge(t, BadgeProps{Label: "5", DirectLabel: true, Position: "indicator-item indicator-top indicator-end"})
	if !strings.Contains(got, `indicator-item indicator-top indicator-end`) {
		t.Errorf("Position classes not emitted:\n%s", got)
	}
}
