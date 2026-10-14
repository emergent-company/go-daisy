package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderEmptyNote renders EmptyNote to its HTML string with the given body.
func renderEmptyNote(t *testing.T, body string, props EmptyNoteProps) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if body != "" {
		ctx = templ.WithChildren(ctx, strComponent(body))
	}
	if err := EmptyNote(props).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestEmptyNoteOutputLocked pins EmptyNote's rendering. The class order is the
// contract: the muted foreground first, then hidden, then the icon row, then
// the type scale, then the caller's Class — matching the shapes already in the
// tree (e.g. "… hidden text-sm", "… flex items-center gap-1.5 text-xs").
func TestEmptyNoteOutputLocked(t *testing.T) {
	cases := []struct {
		name  string
		body  string
		props EmptyNoteProps
		want  string
	}{
		{
			name:  "default",
			body:  "No pending invitations.",
			props: EmptyNoteProps{},
			want:  `<p class="text-base-content/40 text-sm">No pending invitations.</p>`,
		},
		{
			name:  "xs",
			body:  "No projects in this organization yet.",
			props: EmptyNoteProps{Size: EmptyNoteSizeXS},
			want:  `<p class="text-base-content/40 text-xs">No projects in this organization yet.</p>`,
		},
		{
			name:  "2xs_micro",
			body:  "Enter to send · Shift+Enter for a new line",
			props: EmptyNoteProps{Size: EmptyNoteSize2XS},
			want:  `<p class="text-base-content/40 text-[11px]">Enter to send · Shift+Enter for a new line</p>`,
		},
		{
			name:  "hidden_id",
			body:  "Nothing matches.",
			props: EmptyNoteProps{ID: "picker-empty", Hidden: true, Size: EmptyNoteSizeXS},
			want:  `<p id="picker-empty" class="text-base-content/40 hidden text-xs">Nothing matches.</p>`,
		},
		{
			name:  "class_appended",
			body:  "No emails sent yet.",
			props: EmptyNoteProps{Size: EmptyNoteSizeXS, Class: "mt-1"},
			want:  `<p class="text-base-content/40 text-xs mt-1">No emails sent yet.</p>`,
		},
		{
			name:  "icon",
			body:  `Restore is not available yet — it ships in a later release.`,
			props: EmptyNoteProps{Size: EmptyNoteSizeXS, Icon: "lucide--info"},
			want:  `<p class="text-base-content/40 flex items-center gap-1.5 text-xs"><span class="iconify lucide--info size-3.5 shrink-0" aria-hidden="true"></span>Restore is not available yet — it ships in a later release.</p>`,
		},
		{
			name:  "icon_class_override",
			body:  "Runs on the memory host.",
			props: EmptyNoteProps{Size: EmptyNoteSizeXS, Icon: "lucide--shield-check", IconClass: "size-3.5 shrink-0 text-warning"},
			want:  `<p class="text-base-content/40 flex items-center gap-1.5 text-xs"><span class="iconify lucide--shield-check size-3.5 shrink-0 text-warning" aria-hidden="true"></span>Runs on the memory host.</p>`,
		},
		{
			name:  "attrs",
			body:  "Truncated.",
			props: EmptyNoteProps{Size: EmptyNoteSizeXS, Attrs: templ.Attributes{"data-testid": "preview-truncated"}},
			want:  `<p class="text-base-content/40 text-xs" data-testid="preview-truncated">Truncated.</p>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderEmptyNote(t, tc.body, tc.props)
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestEmptyNoteSizeStepsAreExclusive guards the reason 2XS exists at all: each
// size step must emit exactly one font-size utility. A "text-[11px]" passed
// through Class instead would emit two competing font-size utilities, whose
// winner is decided by CSS cascade order rather than class order.
func TestEmptyNoteSizeStepsAreExclusive(t *testing.T) {
	steps := []struct {
		name string
		size EmptyNoteSize
		want string
	}{
		{"sm", EmptyNoteSizeSM, "text-sm"},
		{"xs", EmptyNoteSizeXS, "text-xs"},
		{"2xs", EmptyNoteSize2XS, "text-[11px]"},
	}
	for _, tc := range steps {
		t.Run(tc.name, func(t *testing.T) {
			got := renderEmptyNote(t, "note", EmptyNoteProps{Size: tc.size})
			if !strings.Contains(got, tc.want) {
				t.Fatalf("size %q did not emit %q: %q", tc.name, tc.want, got)
			}
			for _, other := range steps {
				if other.want != tc.want && strings.Contains(got, other.want) {
					t.Errorf("size %q also emitted %q — two competing font-size utilities: %q", tc.name, other.want, got)
				}
			}
		})
	}
}

// TestEmptyNoteIsNotEmptyState guards the boundary between the two components:
// a note renders a single <p>, never EmptyState's flex column or its title.
func TestEmptyNoteIsNotEmptyState(t *testing.T) {
	got := renderEmptyNote(t, "None yet.", EmptyNoteProps{})
	if got != `<p class="text-base-content/40 text-sm">None yet.</p>` {
		t.Errorf("a note must be a single <p>:\n got: %q", got)
	}
}
