package ui

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderDialogTitle renders DialogTitle to its HTML string. An empty title
// renders with no children.
func renderDialogTitle(t *testing.T, title string, props DialogTitleProps) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if title != "" {
		ctx = templ.WithChildren(ctx, strComponent(title))
	}
	if err := DialogTitle(props).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestDialogTitleOutputLocked pins DialogTitle's rendering to the exact markup
// the in-tree titles already emit, so adopting the primitive at a call site is
// a pure structural swap. The dominant shape — an <h3> at text-lg
// font-semibold — must stay byte-identical, and each axis must move exactly one
// utility.
func TestDialogTitleOutputLocked(t *testing.T) {
	cases := []struct {
		name  string
		title string
		props DialogTitleProps
		want  string
	}{
		{
			name:  "default",
			title: "Connect object",
			props: DialogTitleProps{},
			want:  `<h3 class="text-lg font-semibold">Connect object</h3>`,
		},
		{
			name:  "id",
			title: "New agent",
			props: DialogTitleProps{ID: "agent-form-title"},
			want:  `<h3 id="agent-form-title" class="text-lg font-semibold">New agent</h3>`,
		},
		{
			name:  "tracking",
			title: "BPL-1",
			props: DialogTitleProps{Tracking: true, Class: "mt-3"},
			want:  `<h3 class="text-lg font-semibold tracking-tight mt-3">BPL-1</h3>`,
		},
		{
			name:  "text-base",
			title: "Test result",
			props: DialogTitleProps{ID: "memory-alert-title", Size: DialogTitleSizeMD},
			want:  `<h3 id="memory-alert-title" class="text-base font-semibold">Test result</h3>`,
		},
		{
			name:  "h2",
			title: "Appearance",
			props: DialogTitleProps{Level: DialogTitleLevelH2},
			want:  `<h2 class="text-lg font-semibold">Appearance</h2>`,
		},
		{
			name:  "h2_text-base_tracking",
			title: "Migrate",
			props: DialogTitleProps{Level: DialogTitleLevelH2, Size: DialogTitleSizeMD, Tracking: true},
			want:  `<h2 class="text-base font-semibold tracking-tight">Migrate</h2>`,
		},
		{
			name:  "attrs",
			title: "T",
			props: DialogTitleProps{Attrs: templ.Attributes{"data-testid": "t"}},
			want:  `<h3 class="text-lg font-semibold" data-testid="t">T</h3>`,
		},
		{
			name:  "inline_markup",
			title: `Transfer <span id="transfer-project-name" class="text-base-content"></span>`,
			props: DialogTitleProps{ID: "transfer-title"},
			want:  `<h3 id="transfer-title" class="text-lg font-semibold">Transfer <span id="transfer-project-name" class="text-base-content"></span></h3>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderDialogTitle(t, tc.title, tc.props)
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestDialogTitleConfirmDialogParity proves a dogfooded caller renders exactly
// the same heading as the hand-written markup it replaced. ConfirmDialog's
// heading is the canonical "Delete <noun>?" title text-lg font-semibold, so the
// two must agree byte-for-byte.
func TestDialogTitleConfirmDialogParity(t *testing.T) {
	got := renderDialogTitle(t, "Delete agent?", DialogTitleProps{})
	const want = `<h3 class="text-lg font-semibold">Delete agent?</h3>`
	if got != want {
		t.Errorf("ConfirmDialog heading drifted:\n got: %q\nwant: %q", got, want)
	}
}
