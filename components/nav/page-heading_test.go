package nav

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

// renderPageHeading renders PageHeading to its HTML string. actions, when
// non-empty, is injected into the Actions slot.
func renderPageHeading(t *testing.T, props PageHeadingProps, actions string) string {
	t.Helper()
	if actions != "" {
		props.Actions = strComponent(actions)
	}
	var buf bytes.Buffer
	if err := PageHeading(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestPageHeadingDefaultOutputLocked locks the pre-existing rendering for every
// prop combination that existed before Leading and HideBreadcrumbs were added.
// Because the new fields default to their zero value, these must reproduce
// today's output byte-for-byte.
func TestPageHeadingDefaultOutputLocked(t *testing.T) {
	cases := []struct {
		name    string
		props   PageHeadingProps
		actions string
		want    string
	}{
		{
			name:  "zero",
			props: PageHeadingProps{Title: "T"},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><h1 class="text-2xl font-bold tracking-tight">T</h1></div></div></div>`,
		},
		{
			name:  "crumbs",
			props: PageHeadingProps{Title: "T", Breadcrumbs: []BreadcrumbItem{{Label: "Home", Href: "/"}, {Label: "Cases"}}},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul><li class=""><a href="/" class="">Home</a></li><li class=""><span class="">Cases</span></li></ul></div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><h1 class="text-2xl font-bold tracking-tight">T</h1></div></div></div>`,
		},
		{
			name:  "kicker_subtitle",
			props: PageHeadingProps{Title: "T", Kicker: "Eyebrow", Subtitle: "Sub"},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><p class="text-base-content/45 font-semibold tracking-[0.16em] uppercase text-[11px] mb-1">Eyebrow</p><h1 class="text-2xl font-bold tracking-tight">T</h1><p class="text-base-content/55 mt-1 max-w-2xl text-sm">Sub</p></div></div></div>`,
		},
		{
			name:  "dashboard",
			props: PageHeadingProps{Title: "T", Dashboard: true, Subtitle: "Sub"},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2 flex flex-wrap items-end justify-between gap-4"><div><h1 class="text-2xl font-bold tracking-tight lg:text-3xl">T</h1><p class="text-base-content/55 mt-1 max-w-2xl text-sm">Sub</p></div></div></div>`,
		},
		{
			name:  "bare",
			props: PageHeadingProps{Title: "T", Bare: true, Subtitle: "Sub"},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2"><h1 class="text-2xl font-bold tracking-tight">T</h1><p class="text-base-content/55 mt-1 text-sm">Sub</p></div></div>`,
		},
		{
			name:  "subtitle_full",
			props: PageHeadingProps{Title: "T", Subtitle: "Sub", SubtitleFull: true},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><h1 class="text-2xl font-bold tracking-tight">T</h1><p class="text-base-content/55 mt-1 text-sm">Sub</p></div></div></div>`,
		},
		{
			name:    "actions",
			props:   PageHeadingProps{Title: "T"},
			actions: "ACT",
			want:    `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><h1 class="text-2xl font-bold tracking-tight">T</h1></div>ACT</div></div>`,
		},
		{
			name:  "adornment",
			props: PageHeadingProps{Title: "T", TitleAdornment: strComponent("ADORN")},
			want:  `<div class="mb-4"><div class="breadcrumbs text-sm" style="width: 100%;"><ul></ul></div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><div class="flex flex-wrap items-center gap-2"><h1 class="text-2xl font-bold tracking-tight">T</h1>ADORN</div></div></div></div>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderPageHeading(t, tc.props, tc.actions)
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestPageHeadingHideBreadcrumbs omits the breadcrumbs region entirely when
// HideBreadcrumbs is set, even though Breadcrumbs is non-empty.
func TestPageHeadingHideBreadcrumbs(t *testing.T) {
	got := renderPageHeading(t, PageHeadingProps{
		Title:           "T",
		Breadcrumbs:     []BreadcrumbItem{{Label: "Home", Href: "/"}},
		HideBreadcrumbs: true,
	}, "")

	want := `<div class="mb-4"><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><h1 class="text-2xl font-bold tracking-tight">T</h1></div></div></div>`

	if got != want {
		t.Errorf("hideBreadcrumbs output:\n got: %q\nwant: %q", got, want)
	}
	if strings.Contains(got, "breadcrumbs") {
		t.Errorf("hideBreadcrumbs still emitted a breadcrumbs region: %q", got)
	}
}

// TestPageHeadingLeadingSlot proves Leading renders before the title region
// (before breadcrumbs and the h1) and composes with HideBreadcrumbs, kicker,
// subtitle, and actions — the no-breadcrumbs leading-slot header shape.
func TestPageHeadingLeadingSlot(t *testing.T) {
	got := renderPageHeading(t, PageHeadingProps{
		Title:           "Agents",
		Kicker:          "Workspace",
		Subtitle:        "Manage your agents",
		Leading:         strComponent(`<div class="shrink-0">ICON</div>`),
		Actions:         strComponent(`<button>New</button>`),
		HideBreadcrumbs: true,
		Margin:          "mb-6",
	}, "")

	want := `<div class="mb-6"><div class="shrink-0">ICON</div><div class="mt-2 flex flex-wrap items-center justify-between gap-4"><div class="min-w-0"><p class="text-base-content/45 font-semibold tracking-[0.16em] uppercase text-[11px] mb-1">Workspace</p><h1 class="text-2xl font-bold tracking-tight">Agents</h1><p class="text-base-content/55 mt-1 max-w-2xl text-sm">Manage your agents</p></div><button>New</button></div></div>`

	if got != want {
		t.Errorf("leading slot output:\n got: %q\nwant: %q", got, want)
	}
}
