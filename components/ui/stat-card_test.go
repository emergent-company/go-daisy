package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderStatCard renders StatCard to its HTML string in a production context.
func renderStatCard(t *testing.T, props StatCardProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := StatCard(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestStatCardInlineLayoutUnchanged locks the default layout: adding Layout and
// Helper must not change what existing callers (zero values) already render.
func TestStatCardInlineLayoutUnchanged(t *testing.T) {
	html := renderStatCard(t, StatCardProps{
		Label:     "Active Sessions",
		Value:     "42",
		Icon:      "lucide--users",
		IconColor: "bg-primary/10 text-primary",
	})

	for _, want := range []string{
		`rounded-xl p-2.5`,
		`bg-primary/10 text-primary`,
		`text-xl font-bold text-base-content leading-none`,
		`text-xs text-muted mt-1 whitespace-nowrap`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("inline output missing %q\n%s", want, html)
		}
	}

	if strings.Contains(html, "eyebrow") {
		t.Errorf("inline layout must not render the metric label\n%s", html)
	}
	if strings.Index(html, "lucide--users") > strings.Index(html, "42") {
		t.Errorf("inline layout must put the icon tile before the value\n%s", html)
	}
}

// TestStatCardMetricLayout pins the metric arrangement: uppercase label and
// bare icon on the top row, then the value and helper line.
func TestStatCardMetricLayout(t *testing.T) {
	html := renderStatCard(t, StatCardProps{
		Label:     "Hosts",
		Value:     "32",
		Icon:      "lucide--server",
		IconColor: "text-primary",
		Layout:    StatCardLayoutMetric,
		Helper:    "28 reporting · 0 down",
	})

	for _, want := range []string{
		`eyebrow text-xs`,
		`size-4 shrink-0`,
		`text-primary`,
		`text-2xl font-semibold text-base-content tabular-nums`,
		`text-xs break-words text-subtle`,
		`28 reporting · 0 down`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("metric output missing %q\n%s", want, html)
		}
	}

	if strings.Contains(html, "rounded-xl p-2.5") {
		t.Errorf("metric layout must not render the inline icon tile\n%s", html)
	}

	label := strings.Index(html, "Hosts")
	icon := strings.Index(html, "lucide--server")
	value := strings.Index(html, "32")
	helper := strings.Index(html, "28 reporting")
	if label < 0 || icon < 0 || value < 0 || helper < 0 || label >= icon || icon >= value || value >= helper {
		t.Errorf("metric order must be label, icon, value, helper (got %d, %d, %d, %d)\n%s", label, icon, value, helper, html)
	}
}

// TestStatCardMetricOptionalParts guards the empty-value paths.
func TestStatCardMetricOptionalParts(t *testing.T) {
	html := renderStatCard(t, StatCardProps{
		Label:  "Uptime",
		Value:  "99.98%",
		Layout: StatCardLayoutMetric,
	})

	if strings.Contains(html, "text-xs break-words text-subtle") {
		t.Errorf("metric layout must omit an empty helper line\n%s", html)
	}
	if !strings.Contains(html, "iconify") {
		return
	}
	t.Errorf("metric layout must omit the icon when Icon is empty\n%s", html)
}
