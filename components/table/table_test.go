package table

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// renderWithText renders a component that expects `children...` with a single
// text child.
func renderWithText(t *testing.T, c templ.Component, text string) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := templ.WithChildren(context.Background(), templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	}))
	if err := c.Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestTableHeadCellLegacyOutputIsUnchanged guards the 78 call sites of the
// positional API: only the a11y scope attribute may differ from v0.18.0.
func TestTableHeadCellLegacyOutputIsUnchanged(t *testing.T) {
	got := render(t, TableHeadCell("Name"))
	want := `<th scope="col" class="text-sm font-medium text-muted">Name</th>`
	if got != want {
		t.Fatalf("legacy header markup changed:\n got: %s\nwant: %s", got, want)
	}
}

func TestTableHeadCellWithPropsZeroValueMatchesLegacy(t *testing.T) {
	propsForm := render(t, TableHeadCellWithProps(TableHeadCellProps{Label: "Name"}))
	legacy := render(t, TableHeadCell("Name"))
	if propsForm != legacy {
		t.Fatalf("zero-value props form must match the positional form:\n props: %s\nlegacy: %s", propsForm, legacy)
	}
}

func TestTableHeadCellAlignAndNumeric(t *testing.T) {
	tests := []struct {
		name     string
		props    TableHeadCellProps
		want     string
		wantNot  string
		includes string
	}{
		{
			name:     "align end",
			props:    TableHeadCellProps{Label: "Score", Align: AlignEnd},
			want:     `class="text-sm font-medium text-muted text-right"`,
			wantNot:  "tabular-nums",
			includes: ">Score</th>",
		},
		{
			name:     "align start",
			props:    TableHeadCellProps{Label: "Name", Align: AlignStart},
			want:     `class="text-sm font-medium text-muted text-left"`,
			wantNot:  "text-right",
			includes: ">Name</th>",
		},
		{
			name:     "numeric implies right alignment",
			props:    TableHeadCellProps{Label: "Score", Numeric: true},
			want:     `class="text-sm font-medium text-muted text-right tabular-nums"`,
			includes: ">Score</th>",
		},
		{
			name:     "numeric wins over align start",
			props:    TableHeadCellProps{Label: "Score", Align: AlignStart, Numeric: true},
			want:     `class="text-sm font-medium text-muted text-right tabular-nums"`,
			wantNot:  "text-left",
			includes: ">Score</th>",
		},
		{
			name:     "caller class is appended",
			props:    TableHeadCellProps{Label: "Score", Align: AlignEnd, Class: "whitespace-nowrap"},
			want:     `class="text-sm font-medium text-muted text-right whitespace-nowrap"`,
			includes: ">Score</th>",
		},
		{
			name:     "attrs are spread",
			props:    TableHeadCellProps{Label: "Score", Attrs: templ.Attributes{"data-testid": "score-head"}},
			want:     `data-testid="score-head"`,
			includes: ">Score</th>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := render(t, TableHeadCellWithProps(tc.props))
			if !strings.Contains(got, tc.want) {
				t.Fatalf("missing %q in %s", tc.want, got)
			}
			if tc.wantNot != "" && strings.Contains(got, tc.wantNot) {
				t.Fatalf("unexpected %q in %s", tc.wantNot, got)
			}
			if !strings.Contains(got, `scope="col"`) {
				t.Fatalf("missing scope: %s", got)
			}
			if !strings.Contains(got, tc.includes) {
				t.Fatalf("missing %q in %s", tc.includes, got)
			}
			if strings.Contains(got, "data-component") {
				t.Fatalf("plain render must not emit dev-mode attributes: %s", got)
			}
		})
	}
}

func TestTableCellLegacyOutputIsUnchanged(t *testing.T) {
	tests := []struct {
		class string
		attrs templ.Attributes
		want  string
	}{
		{class: "", attrs: nil, want: `<td class="">Value</td>`},
		{class: "text-sm", attrs: nil, want: `<td class="text-sm">Value</td>`},
		{class: "", attrs: templ.Attributes{"id": "cell"}, want: `<td class="" id="cell">Value</td>`},
	}
	for _, tc := range tests {
		got := renderWithText(t, TableCell(tc.class, tc.attrs), "Value")
		if !strings.HasPrefix(got, tc.want) {
			t.Fatalf("legacy cell markup changed:\n got: %s\nwant prefix: %s", got, tc.want)
		}
	}
}

func TestTableCellWithPropsAlignAndNumeric(t *testing.T) {
	numeric := renderWithText(t, TableCellWithProps(TableCellProps{Numeric: true}), "92")
	if !strings.Contains(numeric, `class="text-right tabular-nums"`) {
		t.Fatalf("numeric cell missing alignment: %s", numeric)
	}

	aligned := renderWithText(t, TableCellWithProps(TableCellProps{Align: AlignEnd, Class: "text-sm"}), "92")
	if !strings.Contains(aligned, `class="text-right text-sm"`) {
		t.Fatalf("aligned cell missing classes: %s", aligned)
	}

	attrs := renderWithText(t, TableCellWithProps(TableCellProps{Attrs: templ.Attributes{"colspan": "2"}}), "92")
	if !strings.Contains(attrs, `colspan="2"`) {
		t.Fatalf("cell attrs not spread: %s", attrs)
	}
}

// TestTableCellMatchesHeaderAlignment exercises the issue's core promise: one
// Align/Numeric value drives both the header and the body cells.
func TestTableCellMatchesHeaderAlignment(t *testing.T) {
	head := render(t, TableHeadCellWithProps(TableHeadCellProps{Label: "Score", Numeric: true}))
	body := renderWithText(t, TableCellWithProps(TableCellProps{Numeric: true}), "92")
	if !strings.Contains(head, "text-right tabular-nums") || !strings.Contains(body, "text-right tabular-nums") {
		t.Fatalf("header and cell alignment drifted:\n head: %s\n body: %s", head, body)
	}
}

func TestSortableTableHeaderEndAligned(t *testing.T) {
	got := render(t, SortableTableHeader(SortableTableHeaderProps{
		Label:    "Score",
		Column:   "score",
		Sort:     SortContext{Params: SortParams{By: "score", Dir: SortAsc}, BasePath: "/items", Query: url.Values{}},
		FirstDir: SortDesc,
		Numeric:  true,
		Class:    "text-xs font-semibold",
	}))

	if !strings.Contains(got, `class="text-right tabular-nums text-xs font-semibold"`) {
		t.Fatalf("header cell classes: %s", got)
	}
	if !strings.Contains(got, `aria-sort="ascending"`) {
		t.Fatalf("aria-sort lost: %s", got)
	}
	if !strings.Contains(got, "w-full justify-end") {
		t.Fatalf("end-aligned link must span the cell and pack right: %s", got)
	}

	icon := strings.Index(got, "lucide--chevron-up")
	label := strings.Index(got, "<span>Score</span>")
	if icon < 0 || label < 0 {
		t.Fatalf("missing icon or label: %s", got)
	}
	if icon > label {
		t.Fatalf("end-aligned icon must lead the label:\n icon@%d label@%d\n%s", icon, label, got)
	}
}

func TestSortableTableHeaderStartAlignedKeepsIconTrailing(t *testing.T) {
	got := render(t, SortableTableHeader(SortableTableHeaderProps{
		Label:    "Name",
		Column:   "name",
		Sort:     SortContext{BasePath: "/items", Query: url.Values{}},
		FirstDir: SortAsc,
	}))

	if strings.Contains(got, "justify-end") {
		t.Fatalf("start-aligned link must not justify-end: %s", got)
	}
	icon := strings.Index(got, "lucide--chevrons-up-down")
	label := strings.Index(got, "<span>Name</span>")
	if icon < 0 || label < 0 {
		t.Fatalf("missing icon or label: %s", got)
	}
	if icon < label {
		t.Fatalf("start-aligned icon must follow the label:\n icon@%d label@%d\n%s", icon, label, got)
	}
}

func TestSortableTableHeaderPlainColumnAligns(t *testing.T) {
	got := render(t, SortableTableHeader(SortableTableHeaderProps{
		Label:  "Actions",
		Column: "",
		Sort:   SortContext{BasePath: "/items"},
		Align:  AlignEnd,
	}))
	if !strings.Contains(got, `class="text-right"`) {
		t.Fatalf("plain end-aligned header missing class: %s", got)
	}
	if strings.Contains(got, "<a ") {
		t.Fatalf("plain header must not render a link: %s", got)
	}

	numeric := render(t, SortableTableHeader(SortableTableHeaderProps{
		Label:   "Score",
		Column:  "",
		Sort:    SortContext{BasePath: "/items"},
		Numeric: true,
	}))
	if !strings.Contains(numeric, `class="text-right tabular-nums"`) {
		t.Fatalf("plain numeric header missing alignment: %s", numeric)
	}
}

// TestSortableTableHeaderStartAlignedLink covers the opt-out path: AlignStart on
// a linked header keeps the icon trailing and does not stretch the link.
func TestSortableTableHeaderStartAlignedLink(t *testing.T) {
	got := render(t, SortableTableHeader(SortableTableHeaderProps{
		Label:    "Name",
		Column:   "name",
		Sort:     SortContext{BasePath: "/items", Query: url.Values{}},
		FirstDir: SortAsc,
		Align:    AlignStart,
		Class:    "text-xs",
	}))
	if !strings.Contains(got, `class="text-left text-xs"`) {
		t.Fatalf("align must precede caller class: %s", got)
	}
	if strings.Contains(got, "justify-end") || strings.Contains(got, "w-full") {
		t.Fatalf("start-aligned link must not span the cell: %s", got)
	}
}

func TestCellClasses(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{nil, ""},
		{[]string{"", ""}, ""},
		{[]string{"a"}, "a"},
		{[]string{"a", "", "b"}, "a b"},
		{[]string{" a ", " b"}, "a b"},
	}
	for _, tc := range tests {
		if got := cellClasses(tc.in...); got != tc.want {
			t.Fatalf("cellClasses(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
