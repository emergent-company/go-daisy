package form

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderTagList renders TagList to its HTML string.
func renderTagList(t *testing.T, props TagListProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := TagList(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestTagListChipsClampTheirLabel locks the chip contract: a chip is a
// fixed-height, single-line badge, so its label must be clamped to one line and
// be allowed to shrink instead of widening the field that holds it. Without the
// clamp a long value — an array element written as a JSON document, a long
// label — paints over everything below the input.
func TestTagListChipsClampTheirLabel(t *testing.T) {
	chip := []string{
		// the chip itself must not grow past its container, and must clip
		`badge badge-sm gap-1 pr-0.5 min-w-0 max-w-full overflow-hidden`,
		// the label is one line with an ellipsis, and long unbroken tokens may
		// break so they never contribute an intrinsic width to the field
		`class="line-clamp-1 overflow-hidden wrap-anywhere"`,
		// the full value stays reachable when the label is truncated
		`:title="tag"`,
	}
	cases := []struct {
		name  string
		props TagListProps
	}{
		{name: "enabled", props: TagListProps{ID: "tags", Name: "tags", Values: []string{"a"}}},
		{name: "disabled", props: TagListProps{ID: "tags", Name: "tags", Values: []string{"a"}, Disabled: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			html := renderTagList(t, tc.props)
			for _, want := range chip {
				if !strings.Contains(html, want) {
					t.Errorf("chip markup missing %q in:\n%s", want, html)
				}
			}
			// the clamp must apply to the label, not to the whole chip: the
			// hidden input and the remove action stay outside the clamped span
			if !strings.Contains(html, `<span class="line-clamp-1 overflow-hidden wrap-anywhere" x-text="tag"></span>`) {
				t.Errorf("label span is not the clamped element in:\n%s", html)
			}
		})
	}
}
