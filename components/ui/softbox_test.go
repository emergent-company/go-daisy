package ui

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderSoftBox renders SoftBox to its HTML string with the given body.
func renderSoftBox(t *testing.T, body string, props SoftBoxProps) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if body != "" {
		ctx = templ.WithChildren(ctx, strComponent(body))
	}
	if err := SoftBox(props).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestSoftBoxOutputLocked locks SoftBox's pre-existing rendering. The Surface
// axis is additive: its zero value (Soft) must reproduce the old class list
// byte-for-byte — including the trailing space the old class expression left
// when Class was empty — so existing callers do not move.
func TestSoftBoxOutputLocked(t *testing.T) {
	cases := []struct {
		name  string
		body  string
		props SoftBoxProps
		want  string
	}{
		{
			name:  "default",
			body:  "Content",
			props: SoftBoxProps{},
			want:  `<div class="rounded-box border border-line-translucent bg-base-200/40 p-3 ">Content</div>`,
		},
		{
			name:  "padding",
			body:  "Content",
			props: SoftBoxProps{Padding: "p-5"},
			want:  `<div class="rounded-box border border-line-translucent bg-base-200/40 p-5 ">Content</div>`,
		},
		{
			name:  "class_appended",
			body:  "Content",
			props: SoftBoxProps{Class: "flex items-center gap-2"},
			want:  `<div class="rounded-box border border-line-translucent bg-base-200/40 p-3 flex items-center gap-2">Content</div>`,
		},
		{
			name:  "explicit_soft",
			body:  "Content",
			props: SoftBoxProps{Surface: SoftBoxDefault},
			want:  `<div class="rounded-box border border-line-translucent bg-base-200/40 p-3 ">Content</div>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderSoftBox(t, tc.body, tc.props)
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestSoftBoxSurfaceVariants pins the two new fills. Only the background
// utility moves; the border, radius and padding are shared.
func TestSoftBoxSurfaceVariants(t *testing.T) {
	cases := []struct {
		name    string
		surface SoftBoxSurface
		want    string
	}{
		{
			name:    "muted",
			surface: SoftBoxMuted,
			want:    `<div class="rounded-box border border-line-translucent bg-base-200/30 p-3 ">Content</div>`,
		},
		{
			name:    "plain",
			surface: SoftBoxPlain,
			want:    `<div class="rounded-box border border-line-translucent p-3 ">Content</div>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderSoftBox(t, "Content", SoftBoxProps{Surface: tc.surface})
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestSoftBoxAttrsEmptyIsByteIdentical proves the Attrs escape hatch is an
// addition, not a rewrite: an empty Attrs map (and a nil one) spread onto the
// root element emits nothing, so the markup stays byte-identical to the
// pre-Attrs output down to the trailing space in the class list. Existing
// callers that never set Attrs do not move.
func TestSoftBoxAttrsEmptyIsByteIdentical(t *testing.T) {
	const locked = `<div class="rounded-box border border-line-translucent bg-base-200/40 p-3 ">Content</div>`

	for _, attrs := range []templ.Attributes{nil, {}} {
		got := renderSoftBox(t, "Content", SoftBoxProps{Attrs: attrs})
		if got != locked {
			t.Errorf("empty Attrs %#v changed the markup:\n got: %q\nwant: %q", attrs, got, locked)
		}
	}
}

// TestSoftBoxAttrsRenders pins the new behaviour: Attrs values are spread onto
// the root element, including a data-* key carrying a dynamic Go value — the
// shape the gateway needs for data-testid, data-registry-name and
// data-mcp-configure-help selection.
func TestSoftBoxAttrsRenders(t *testing.T) {
	registryName := "mcp-entry"
	got := renderSoftBox(t, "Content", SoftBoxProps{
		Surface: SoftBoxMuted,
		Attrs: templ.Attributes{
			"data-testid":        "mcp-registry-entry",
			"data-registry-name": registryName,
		},
	})
	// templ.RenderAttributes emits attributes in sorted key order.
	want := `<div class="rounded-box border border-line-translucent bg-base-200/30 p-3 " data-registry-name="mcp-entry" data-testid="mcp-registry-entry">Content</div>`
	if got != want {
		t.Errorf("Attrs did not render on the root element:\n got: %q\nwant: %q", got, want)
	}
}
