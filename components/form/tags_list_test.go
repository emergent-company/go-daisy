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
	// The badge may shrink and may not grow past its field, but it must NOT clip
	// (overflow-hidden hides the remove button's focus-visible outline).
	chipClasses := []string{"badge", "badge-sm", "min-w-0", "max-w-full"}
	labelClasses := []string{"line-clamp-1", "overflow-hidden", "wrap-anywhere"}
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
			chip := spanWith(t, html, `:style="tagColor(tag)"`)
			for _, want := range chipClasses {
				if !hasToken(classOf(chip), want) {
					t.Errorf("chip class %q missing from %q", want, classOf(chip))
				}
			}
			// clipping the chip would also clip the remove button's
			// focus-visible outline, so it must not be used
			if hasToken(classOf(chip), "overflow-hidden") {
				t.Errorf("chip must not clip its content: %q", classOf(chip))
			}
			// the clamp applies to the label, so the remove action and the
			// hidden input stay outside the clipped span
			label := spanWith(t, html, `x-text="tag"`)
			for _, want := range labelClasses {
				if !hasToken(classOf(label), want) {
					t.Errorf("label class %q missing from %q", want, classOf(label))
				}
			}
			// the full value stays reachable when the label is truncated
			if !strings.Contains(html, `:title="tag"`) {
				t.Errorf("chip missing the title binding in:\n%s", html)
			}
		})
	}
}

// spanWith returns the opening-tag attributes of the <span> whose attributes
// contain marker.
func spanWith(t *testing.T, html, marker string) string {
	t.Helper()
	for _, part := range strings.Split(html, "<span ") {
		if i := strings.Index(part, ">"); i >= 0 && strings.Contains(part[:i], marker) {
			return part[:i]
		}
	}
	t.Fatalf("no <span> with %q in:\n%s", marker, html)
	return ""
}

// classOf returns the class attribute value of an opening tag.
func classOf(tag string) string {
	const attr = `class="`
	i := strings.Index(tag, attr)
	if i < 0 {
		return ""
	}
	rest := tag[i+len(attr):]
	if j := strings.Index(rest, `"`); j >= 0 {
		return rest[:j]
	}
	return rest
}

// hasToken reports whether a space-separated class list contains want.
func hasToken(classes, want string) bool {
	for _, c := range strings.Fields(classes) {
		if c == want {
			return true
		}
	}
	return false
}
