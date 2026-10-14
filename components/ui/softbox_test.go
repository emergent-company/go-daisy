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
