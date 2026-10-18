package ui

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// renderCardRawProps renders CardRaw with the given props (and no children) to
// its HTML string using a plain (non-devmode) context.
func renderCardRawProps(t *testing.T, props CardRawProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := CardRaw(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// renderCardRawPropsDevMode renders CardRaw with the given props using a devmode
// context, so the devmode loud-failure paths are active.
func renderCardRawPropsDevMode(t *testing.T, props CardRawProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := CardRaw(props).Render(devmode.WithDevMode(context.Background()), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestCardRawZeroValueOutput pins the full zero-value render. The surface stays
// exactly bg-base-100 (byte-identical to the pre-fix hard-coded value) and the
// default edge-line is emitted.
func TestCardRawZeroValueOutput(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{})

	want := `<div class="card bg-base-100 edge-line"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("zero-value CardRaw output:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawSurfaceReplacesDefault pins that a non-zero Surface REPLACES the
// bg-base-100 default instead of stacking a second bg-* class on top of it.
func TestCardRawSurfaceReplacesDefault(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Surface: CardSurfaceBase200})

	if !strings.Contains(got, `class="card bg-base-200 edge-line"`) {
		t.Errorf("Surface base-200 must replace bg-base-100:\n%s", got)
	}
	if strings.Contains(got, "bg-base-100") {
		t.Errorf("Surface base-200 must not leave bg-base-100 behind:\n%s", got)
	}
}

// TestCardRawBorderNone pins that BorderNone removes the default outline.
func TestCardRawBorderNone(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Border: CardBorderNone})

	want := `<div class="card bg-base-100"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("BorderNone output:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawBorderStrong pins that BorderStrong uses the emphasised step.
func TestCardRawBorderStrong(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Border: CardBorderStrong})

	if !strings.Contains(got, `class="card bg-base-100 edge-line-strong"`) {
		t.Errorf("BorderStrong must render border-line-strong:\n%s", got)
	}
}

// TestCardRawAsElement pins the full element render (opening and closing tag)
// for every documented CardElement value.
func TestCardRawAsElement(t *testing.T) {
	cases := []struct {
		name string
		as   CardElement
		want string
	}{
		{"div", CardElementDiv, `<div class="card bg-base-100 edge-line"><div class="card-body"></div></div>`},
		{"a", CardElementA, `<a class="card bg-base-100 edge-line"><div class="card-body"></div></a>`},
		{"button", CardElementButton, `<button type="button" class="card bg-base-100 edge-line"><div class="card-body"></div></button>`},
		{"section", CardElementSection, `<section class="card bg-base-100 edge-line"><div class="card-body"></div></section>`},
		{"article", CardElementArticle, `<article class="card bg-base-100 edge-line"><div class="card-body"></div></article>`},
		{"li", CardElementLi, `<li class="card bg-base-100 edge-line"><div class="card-body"></div></li>`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderCardRawProps(t, CardRawProps{As: tc.as})
			if got != tc.want {
				t.Errorf("As %q:\n got: %q\nwant: %q", tc.as, got, tc.want)
			}
		})
	}
}

// TestCardRawUnknownAsPanicsInDevMode pins that an unsupported as element fails
// loudly in devmode instead of silently degrading.
func TestCardRawUnknownAsPanicsInDevMode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic in devmode for unknown as element")
		}
	}()
	renderCardRawPropsDevMode(t, CardRawProps{As: CardElement("span")})
}

// TestCardRawUnknownAsFallsBackToDiv pins that outside devmode an unsupported as
// element degrades to the default <div> rather than crashing the request.
func TestCardRawUnknownAsFallsBackToDiv(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{As: CardElement("span")})

	want := `<div class="card bg-base-100 edge-line"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("unknown as must fall back to <div> outside devmode:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawUnknownBorderPanicsInDevMode pins that an unsupported border fails
// loudly in devmode.
func TestCardRawUnknownBorderPanicsInDevMode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic in devmode for unknown border")
		}
	}()
	renderCardRawPropsDevMode(t, CardRawProps{Border: CardBorder("dashed")})
}

// TestCardRawUnknownBorderFallsBackToDefault pins that outside devmode an
// unsupported border degrades to the default outline rather than crashing.
func TestCardRawUnknownBorderFallsBackToDefault(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Border: CardBorder("dashed")})

	want := `<div class="card bg-base-100 edge-line"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("unknown border must fall back to edge-line outside devmode:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawAttrsSpreadLast pins that caller Attrs land on the card root after
// the component's own class attribute.
func TestCardRawAttrsSpreadLast(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Attrs: templ.Attributes{"data-testid": "foo", "id": "bar"}})

	for _, attr := range []string{`data-testid="foo"`, `id="bar"`} {
		if !strings.Contains(got, attr) {
			t.Errorf("Attrs %s must be spread onto the card root:\n%s", attr, got)
		}
	}
	classIdx := strings.Index(got, `class="card`)
	attrIdx := strings.Index(got, `data-testid="foo"`)
	if classIdx < 0 || attrIdx < 0 || classIdx > attrIdx {
		t.Errorf("Attrs must be spread after the class attribute:\n%s", got)
	}
}

// TestCardBorderClassPinsHouseEdge pins the class each CardBorder maps to. The
// edge-* names are the point of #79: they carry the theme's --border width, so a
// later "tidy up" back to `border border-line` fails here before it reaches a
// high-contrast user.
func TestCardBorderClassPinsHouseEdge(t *testing.T) {
	cases := []struct {
		name string
		in   CardBorder
		want string
	}{
		{"line", CardBorderLine, "edge-line"},
		{"strong", CardBorderStrong, "edge-line-strong"},
		{"none", CardBorderNone, ""},
		{"unknown_falls_back", CardBorder("dashed"), "edge-line"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := cardBorderClass(tc.in, context.Background()); got != tc.want {
				t.Errorf("cardBorderClass(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestCardOutlineWidthComesFromThemeToken is the width-parity net for #79. The
// class-string tests pin names; this one pins the mechanism, so the outline
// cannot silently regress to a fixed 1px:
//
//   - the house edge must take its width from var(--border …), not a length;
//   - --border must still be raised for `prefers-contrast: more`;
//   - no rendered card outline may fall back to Tailwind's fixed-width `border`
//     utility or to a per-call-site arbitrary width.
func TestCardOutlineWidthComesFromThemeToken(t *testing.T) {
	src, err := os.ReadFile("../../components/css/custom.css")
	if err != nil {
		t.Fatalf("read house stylesheet: %v", err)
	}
	// Collapse whitespace so the assertions pin the declarations, not the
	// stylesheet's line wrapping.
	css := strings.Join(strings.Fields(string(src)), " ")

	for _, want := range []string{
		".edge-line { border: var(--border, 1px) solid var(--color-line); }",
		".edge-line-strong { border: var(--border, 1px) solid var(--color-line-strong); }",
		"@media (prefers-contrast: more) { :root { --border: 2px; } }",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("components/css/custom.css is missing %q — the card outline would no longer follow the theme border width", want)
		}
	}

	// A literal width in the house edge is the exact defect this replaced.
	for _, pinned := range []string{
		"border: 1px solid var(--color-line)",
		"border: 2px solid var(--color-line)",
	} {
		if strings.Contains(css, pinned) {
			t.Errorf("house edge pins a literal width (%q); it must read var(--border, 1px)", pinned)
		}
	}

	for name, props := range map[string]CardRawProps{
		"default": {},
		"strong":  {Border: CardBorderStrong},
	} {
		for _, tok := range strings.Fields(cardRootClass(t, props)) {
			if tok == "border" || strings.HasPrefix(tok, "border-[length:") || strings.HasPrefix(tok, "border-(length:") {
				t.Errorf("%s outline uses the fixed-width %q; want an edge-* class", name, tok)
			}
		}
	}
}

// cardRootClass returns the class attribute of the rendered card root.
func cardRootClass(t *testing.T, props CardRawProps) string {
	t.Helper()
	out := renderCardRawProps(t, props)
	const key = `class="`
	i := strings.Index(out, key)
	if i < 0 {
		t.Fatalf("no class attribute in card root: %q", out)
	}
	rest := out[i+len(key):]
	j := strings.Index(rest, `"`)
	if j < 0 {
		t.Fatalf("unterminated class attribute in card root: %q", out)
	}
	return rest[:j]
}

// TestSurfaceRootsEmitHouseEdge covers the components that were switched from
// `border border-line` to the house edge alongside the card family. Rendering
// is the point: the class strings here are literals in .templ files, so nothing
// else would catch a partial revert.
func TestSurfaceRootsEmitHouseEdge(t *testing.T) {
	cases := []struct {
		name string
		c    templ.Component
	}{
		{"Card", Card("Title", nil)},
		{"CardPlaceholder", CardPlaceholder("Card")},
		{"StatCardMinimal", StatCardMinimal(StatCardMinimalItem{Label: "Open", Value: "1", TrendLabel: "2%"})},
		{"StatCardMinimalWithIcon", StatCardMinimal(StatCardMinimalItem{Label: "Open", Value: "1", TrendLabel: "2%", Icon: "lucide--briefcase", IconColor: "text-primary"})},
		{"Banner", Banner(BannerProps{CookieBanner: true, CookieTitle: "Cookies", Description: "We use them."})},
		{"AccordionItem", AccordionItem("Details", false)},
		{"AlpineAccordionItem", AlpineAccordionItem("Details", false)},
		{"ActionMenu", ActionMenu([]ActionMenuItem{{Label: "Edit", Href: "#"}})},

		// The remaining surface roots converted in #82: the outline that says
		// "this is a card / panel / floating surface". Cheap to instantiate, so
		// they are pinned here rather than by a snapshot of a whole tree.
		{"DashboardCard", DashboardCard("Revenue", "")},
		{"FilterCard", FilterCard(FilterCardProps{})},
		{"SkeletonCard", SkeletonCard()},
		{"StatsGroup", StatsGroup([]StatsGroupItem{{Title: "Open", Value: "1"}})},
		{"List", List(ListProps{})},
		{"CodeBlock", CodeBlock(CodeBlockProps{Language: "go", Code: "x"})},
		{"DiffContainer", DiffContainer()},
		{"NotificationPanel", NotificationPanel(0, "#")},
		{"Collapsible", Collapsible()},
		{"PopoverContent", PopoverContent()},
		{"SearchDropdown", SearchDropdown("Search", nil)},
		{"NotificationDropdown", NotificationDropdown(NotificationDropdownProps{})},
		{"DropdownMenu", DropdownMenu(nil)},
		{"AlpineDropdownMenu", AlpineDropdownMenu()},
		{"PersonCellCompact", PersonCell(PersonCellProps{Name: "Ada", Compact: true})},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tc.c.Render(context.Background(), &buf); err != nil {
				t.Fatalf("render: %v", err)
			}
			out := buf.String()
			if !strings.Contains(out, "edge-line") {
				t.Errorf("%s does not emit the house edge; outline width would fall back to a fixed 1px:\n%s", tc.name, out)
			}
			if strings.Contains(out, "border border-line") {
				t.Errorf("%s still emits the fixed-width `border border-line` pair:\n%s", tc.name, out)
			}
		})
	}
}
