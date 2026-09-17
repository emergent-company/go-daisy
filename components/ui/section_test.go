package ui

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// strComponent returns a component that writes a literal HTML string.
func strComponent(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	})
}

// renderSection renders Section to its HTML string, injecting child as the
// children slot. An empty child renders with no children.
func renderSection(t *testing.T, title string, props SectionProps, child string) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if child != "" {
		ctx = templ.WithChildren(ctx, strComponent(child))
	}
	if err := Section(title, props).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// renderCardRaw renders CardRaw with the given classes and a literal child.
func renderCardRaw(t *testing.T, cardClass, bodyClass, child string) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if child != "" {
		ctx = templ.WithChildren(ctx, strComponent(child))
	}
	if err := CardRaw(cardClass, bodyClass, nil).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestSectionDefaultOutputLocked locks the pre-existing (non-wrapperless)
// rendering for every prop combination that existed before the Wrapperless
// prop was added. Because Wrapperless defaults to false, these must reproduce
// today's output byte-for-byte — the additive change must not break callers.
func TestSectionDefaultOutputLocked(t *testing.T) {
	cases := []struct {
		name  string
		title string
		props SectionProps
		child string
		want  string
	}{
		{
			name:  "titled",
			title: "Object types",
			props: SectionProps{},
			child: `<div>r0</div>`,
			want:  `<section class=""><h2 class="mb-3 text-lg font-semibold tracking-tight">Object types</h2><div class="card bg-base-100 card-border"><div class="card-body "><div>r0</div></div></div></section>`,
		},
		{
			name:  "titleless",
			title: "Object types",
			props: SectionProps{Titleless: true},
			child: `<div>r0</div>`,
			want:  `<section class=""><div class="card bg-base-100 card-border"><div class="card-body "><div>r0</div></div></div></section>`,
		},
		{
			name:  "rows",
			title: "Object types",
			props: SectionProps{Rows: true},
			child: `<div>r0</div><div>r1</div>`,
			want:  `<section class=""><h2 class="mb-3 text-lg font-semibold tracking-tight">Object types</h2><div class="card bg-base-100 card-border"><div class="card-body "><div class="divide-y divide-base-200"><div>r0</div><div>r1</div></div></div></div></section>`,
		},
		{
			name:  "raw",
			title: "Object types",
			props: SectionProps{Raw: true},
			child: `<div>r0</div>`,
			want:  `<section class=""><h2 class="mb-3 text-lg font-semibold tracking-tight">Object types</h2><div>r0</div></section>`,
		},
		{
			name:  "empty",
			title: "Object types",
			props: SectionProps{Empty: true, EmptyIcon: "lucide--box", EmptyTitle: "No types", EmptyDesc: "None yet"},
			child: `<div>r0</div>`,
			want:  `<section class=""><h2 class="mb-3 text-lg font-semibold tracking-tight">Object types</h2><div class="card bg-base-100 card-border"><div class="card-body "><div class="flex flex-col items-center justify-center py-16 text-center px-4"><span class="iconify size-12 text-base-content/20 mb-4 lucide--box"></span><h3 class="text-base font-semibold text-base-content mb-1">No types</h3><p class="text-sm text-base-content/50 mb-4">None yet</p></div></div></div></section>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderSection(t, tc.title, tc.props, tc.child)
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestSectionWrapperlessMatchesCardRaw proves the new Wrapperless mode
// reproduces the consumer's bare card-list markup exactly: Section with
// Wrapperless + CardClass + BodyClass + Rows must equal CardRaw wrapping a
// divide-y row list. Both the literal string and the live CardRaw render are
// asserted so class names and nesting cannot silently drift.
func TestSectionWrapperlessMatchesCardRaw(t *testing.T) {
	const (
		cardClass = "card-border overflow-hidden"
		bodyClass = "p-0"
	)
	const child = `<div>r0</div><div>r1</div><div>r2</div>`

	got := renderSection(t, "ignored", SectionProps{
		Wrapperless: true,
		CardClass:   cardClass,
		BodyClass:   bodyClass,
		Rows:        true,
	}, child)

	literal := `<div class="card bg-base-100 card-border overflow-hidden"><div class="card-body p-0"><div class="divide-y divide-base-200"><div>r0</div><div>r1</div><div>r2</div></div></div></div>`

	if got != literal {
		t.Errorf("wrapperless output != literal target:\n got: %q\nwant: %q", got, literal)
	}

	want := renderCardRaw(t, cardClass, bodyClass, `<div class="divide-y divide-base-200">`+child+`</div>`)
	if got != want {
		t.Errorf("wrapperless output != CardRaw equivalent:\n got: %q\nwant: %q", got, want)
	}
}

// TestSectionWrapperlessOmitsSectionAndHeading confirms no <section> and no
// <h2> are emitted in Wrapperless mode, and that Titleless is implied.
func TestSectionWrapperlessOmitsSectionAndHeading(t *testing.T) {
	got := renderSection(t, "A heading that must not appear", SectionProps{
		Wrapperless: true,
		CardClass:   "card-border overflow-hidden",
		BodyClass:   "p-0",
		Rows:        true,
	}, `<div>r0</div>`)

	if strings.Contains(got, "<section") {
		t.Errorf("wrapperless emitted a <section>: %q", got)
	}
	if strings.Contains(got, "<h2") {
		t.Errorf("wrapperless emitted an <h2>: %q", got)
	}
	if strings.Contains(got, "A heading that must not appear") {
		t.Errorf("wrapperless emitted the title: %q", got)
	}
}

// TestSectionWrapperlessEmptySupported proves Empty is honoured in Wrapperless
// mode: the default empty-state card renders without a section wrapper.
func TestSectionWrapperlessEmptySupported(t *testing.T) {
	got := renderSection(t, "ignored", SectionProps{
		Wrapperless: true,
		Empty:       true,
		EmptyIcon:   "lucide--box",
		EmptyTitle:  "No types",
		EmptyDesc:   "None yet",
		// CardClass/BodyClass must be ignored, matching the non-wrapperless
		// Empty branch.
		CardClass: "card-border overflow-hidden",
		BodyClass: "p-0",
	}, `<div>r0</div>`)

	want := `<div class="card bg-base-100 card-border"><div class="card-body "><div class="flex flex-col items-center justify-center py-16 text-center px-4"><span class="iconify size-12 text-base-content/20 mb-4 lucide--box"></span><h3 class="text-base font-semibold text-base-content mb-1">No types</h3><p class="text-sm text-base-content/50 mb-4">None yet</p></div></div></div>`

	if got != want {
		t.Errorf("wrapperless empty output:\n got: %q\nwant: %q", got, want)
	}
}

// TestSectionWrapperlessRawIgnored proves Raw is ignored in Wrapperless mode:
// the card path always renders instead of the raw children.
func TestSectionWrapperlessRawIgnored(t *testing.T) {
	got := renderSection(t, "ignored", SectionProps{
		Wrapperless: true,
		Raw:         true,
		CardClass:   "card-border overflow-hidden",
		BodyClass:   "p-0",
	}, `<div>r0</div>`)

	want := `<div class="card bg-base-100 card-border overflow-hidden"><div class="card-body p-0"><div>r0</div></div></div>`

	if got != want {
		t.Errorf("wrapperless raw output:\n got: %q\nwant: %q", got, want)
	}
}

// TestSectionWrapperlessMarginOnCard proves Margin is applied to the card
// element (the new outermost element) when the section wrapper is absent.
func TestSectionWrapperlessMarginOnCard(t *testing.T) {
	got := renderSection(t, "ignored", SectionProps{
		Wrapperless: true,
		Margin:      "mt-6",
		CardClass:   "card-border overflow-hidden",
		BodyClass:   "p-0",
		Rows:        true,
	}, `<div>r0</div>`)

	want := `<div class="card bg-base-100 card-border overflow-hidden mt-6"><div class="card-body p-0"><div class="divide-y divide-base-200"><div>r0</div></div></div></div>`

	if got != want {
		t.Errorf("wrapperless margin output:\n got: %q\nwant: %q", got, want)
	}
}

// TestSectionWrapperlessDefaultCardClass proves an empty CardClass falls back
// to "card-border", matching the non-wrapperless card path default.
func TestSectionWrapperlessDefaultCardClass(t *testing.T) {
	got := renderSection(t, "ignored", SectionProps{
		Wrapperless: true,
		BodyClass:   "p-0",
	}, `<div>r0</div>`)

	want := `<div class="card bg-base-100 card-border"><div class="card-body p-0"><div>r0</div></div></div>`

	if got != want {
		t.Errorf("wrapperless default cardClass output:\n got: %q\nwant: %q", got, want)
	}
}
