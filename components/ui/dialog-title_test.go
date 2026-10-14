package ui

import (
	"bytes"
	"context"
	"strings"
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
// heading is the canonical "Delete <noun>?" title, so the two must agree
// byte-for-byte.
//
// It renders the real ConfirmDialog rather than DialogTitle in isolation: the
// adoption is @DialogTitle() { { title } }, and a children block holding a bare
// identifier renders that identifier as literal text, so an isolated render
// would stay green through exactly the regression this test exists to catch.
func TestDialogTitleConfirmDialogParity(t *testing.T) {
	var buf bytes.Buffer
	if err := ConfirmDialog(ConfirmDialogProps{Noun: "agent"}, nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	got := buf.String()

	const want = `<h3 class="text-lg font-semibold">Delete agent?</h3>`
	if !strings.Contains(got, want) {
		t.Errorf("ConfirmDialog heading drifted — want to contain:\n  %s\ngot:\n  %s", want, got)
	}
	// The literal identifier must not leak into the heading as text.
	if strings.Contains(got, ">title<") || strings.Contains(got, ">title ") {
		t.Errorf("title rendered as the literal word \"title\":\n  %s", got)
	}
	if !strings.Contains(got, "Delete agent?") {
		t.Errorf("ConfirmDialog did not render its title text:\n  %s", got)
	}
}

// TestDialogTitleConfirmDialogOverrideParity pins the Title override path — the
// other branch that reaches the same heading through DialogTitle's children.
func TestDialogTitleConfirmDialogOverrideParity(t *testing.T) {
	var buf bytes.Buffer
	if err := ConfirmDialog(ConfirmDialogProps{Title: "Remove integration?"}, nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, `<h3 class="text-lg font-semibold">Remove integration?</h3>`) {
		t.Errorf("Title override did not reach the heading:\n  %s", got)
	}
}
