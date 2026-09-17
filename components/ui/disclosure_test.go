package ui

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderDisclosure renders Disclosure to its HTML string, injecting header as
// the summary leading content, an optional trailing component before the
// chevron, and body as the children slot.
func renderDisclosure(t *testing.T, props DisclosureProps, header, trailing, body string) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if body != "" {
		ctx = templ.WithChildren(ctx, strComponent(body))
	}
	var tr templ.Component
	if trailing != "" {
		tr = strComponent(trailing)
	}
	if err := Disclosure(props, strComponent(header), tr).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestDisclosureDefaultOutputLocked locks the pre-existing rendering for every
// prop combination that existed before DetailsBase/SummaryBase/BodyBase and
// SummaryAttrs were added. Because the new fields default to their zero value,
// these must reproduce today's output byte-for-byte — the additive change must
// not break callers.
func TestDisclosureDefaultOutputLocked(t *testing.T) {
	cases := []struct {
		name     string
		props    DisclosureProps
		header   string
		trailing string
		want     string
	}{
		{
			name:   "zero",
			props:  DisclosureProps{},
			header: "HDR",
			want:   `<details class="group rounded-box border border-base-content/10 bg-base-200/40 "><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-center ">HDR<div class="flex shrink-0 items-center gap-2"><span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 "></div></details>`,
		},
		{
			name:   "open",
			props:  DisclosureProps{Open: true},
			header: "HDR",
			want:   `<details class="group rounded-box border border-base-content/10 bg-base-200/40 " open><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-center ">HDR<div class="flex shrink-0 items-center gap-2"><span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 "></div></details>`,
		},
		{
			name:   "items_start",
			props:  DisclosureProps{ItemsStart: true},
			header: "HDR",
			want:   `<details class="group rounded-box border border-base-content/10 bg-base-200/40 "><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-start ">HDR<div class="flex shrink-0 items-center gap-2"><span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 "></div></details>`,
		},
		{
			name:   "appended_classes",
			props:  DisclosureProps{Class: "my-details", SummaryClass: "my-sum", BodyClass: "my-body"},
			header: "HDR",
			want:   `<details class="group rounded-box border border-base-content/10 bg-base-200/40 my-details"><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-center my-sum">HDR<div class="flex shrink-0 items-center gap-2"><span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 my-body"></div></details>`,
		},
		{
			name:   "named_group",
			props:  DisclosureProps{GroupClass: "group/cap"},
			header: "HDR",
			want:   `<details class="group/cap rounded-box border border-base-content/10 bg-base-200/40 "><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-center ">HDR<div class="flex shrink-0 items-center gap-2"><span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open/cap:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 "></div></details>`,
		},
		{
			name:     "trailing",
			props:    DisclosureProps{},
			header:   "HDR",
			trailing: "TRAIL",
			want:     `<details class="group rounded-box border border-base-content/10 bg-base-200/40 "><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-center ">HDR<div class="flex shrink-0 items-center gap-2">TRAIL<span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 "></div></details>`,
		},
		{
			name:   "attrs",
			props:  DisclosureProps{Attrs: templ.Attributes{"data-x": "1"}},
			header: "HDR",
			want:   `<details class="group rounded-box border border-base-content/10 bg-base-200/40 " data-x="1"><summary class="flex cursor-pointer list-none justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden items-center ">HDR<div class="flex shrink-0 items-center gap-2"><span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 "></div></details>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderDisclosure(t, tc.props, tc.header, tc.trailing, "")
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestDisclosureTargetVariants proves the new Base/SummaryAttrs API reproduces
// the consumer's two tool-group disclosure variants (source and relay) exactly.
// The two differ only in the details border/background and the body border
// tone; the summary classes and the data-testid/data-tool-group hooks are
// identical, so a caller must not have to restate everything to change one
// utility.
func TestDisclosureTargetVariants(t *testing.T) {
	const summaryBase = "flex cursor-pointer list-none items-center justify-between gap-3 px-3 py-2.5 select-none [&::-webkit-details-marker]:hidden"

	source := renderDisclosure(t, DisclosureProps{
		DetailsBase: "rounded-box border border-base-content/10 bg-base-200/40",
		SummaryBase: summaryBase,
		BodyBase:    "flex flex-col border-t border-base-content/10",
	}, "HDR", "TRAIL", "BODY")

	sourceWant := `<details class="group rounded-box border border-base-content/10 bg-base-200/40 "><summary class="flex cursor-pointer list-none items-center justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden ">HDR<div class="flex shrink-0 items-center gap-2">TRAIL<span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col border-t border-base-content/10 ">BODY</div></details>`

	if source != sourceWant {
		t.Errorf("source variant:\n got: %q\nwant: %q", source, sourceWant)
	}

	relay := renderDisclosure(t, DisclosureProps{
		DetailsBase: "rounded-box border border-primary/25 bg-primary/[0.04]",
		SummaryBase: summaryBase,
		BodyBase:    "flex flex-col border-t border-primary/15",
	}, "HDR", "TRAIL", "BODY")

	relayWant := `<details class="group rounded-box border border-primary/25 bg-primary/[0.04] "><summary class="flex cursor-pointer list-none items-center justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden ">HDR<div class="flex shrink-0 items-center gap-2">TRAIL<span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col border-t border-primary/15 ">BODY</div></details>`

	if relay != relayWant {
		t.Errorf("relay variant:\n got: %q\nwant: %q", relay, relayWant)
	}
}

// TestDisclosureSummaryAttrsAndNamedGroup proves SummaryAttrs lands on the
// <summary> element (not the <details>) and that a named group (group/cap)
// still drives the chevron rotation while ItemsStart supplies items-start when
// no custom SummaryBase is given.
func TestDisclosureSummaryAttrsAndNamedGroup(t *testing.T) {
	got := renderDisclosure(t, DisclosureProps{
		GroupClass:  "group/cap",
		DetailsBase: "rounded-box border border-base-content/10 bg-base-200/40",
		SummaryBase: "flex cursor-pointer list-none items-start justify-between gap-3 px-3 py-2.5 select-none [&::-webkit-details-marker]:hidden",
		SummaryAttrs: templ.Attributes{"data-testid": "tool-group-header-grp1"},
		Attrs:        templ.Attributes{"data-testid": "tool-group", "data-tool-group": "grp1"},
	}, "HDR", "TRAIL", "BODY")

	want := `<details class="group/cap rounded-box border border-base-content/10 bg-base-200/40 " data-testid="tool-group" data-tool-group="grp1"><summary class="flex cursor-pointer list-none items-start justify-between gap-3 px-3 py-2.5 select-none [&amp;::-webkit-details-marker]:hidden " data-testid="tool-group-header-grp1">HDR<div class="flex shrink-0 items-center gap-2">TRAIL<span class="iconify lucide--chevron-down size-4 shrink-0 text-base-content/40 transition-transform group-open/cap:rotate-180" aria-hidden="true"></span></div></summary><div class="flex flex-col gap-2 border-t border-base-content/10 p-3 ">BODY</div></details>`

	if got != want {
		t.Errorf("capability variant:\n got: %q\nwant: %q", got, want)
	}
}
