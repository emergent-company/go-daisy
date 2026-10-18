package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/shared"
)

func renderUI(t *testing.T, c templ.Component, opts ...func(context.Context) context.Context) string {
	t.Helper()
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	var buf bytes.Buffer
	if err := c.Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func uiStr(s string) templ.Component { return shared.StrComp(s) }

// TestDropdownEmitsMenuContract locks the data attributes the a11y runtime
// depends on. If these disappear the CSS-only dropdown silently loses keyboard
// support, so they are part of the component contract, not styling.
func TestDropdownEmitsMenuContract(t *testing.T) {
	inner := shared.Compose(
		DropdownTrigger("Actions", "btn-primary", nil),
		DropdownMenu(nil),
	)
	got := renderUI(t, shared.RenderInto(Dropdown(DropdownEnd, nil), inner))

	for _, want := range []string{
		"data-gd-menu-root",
		"data-gd-menu-trigger",
		`role="menu"`,
		"data-gd-menu",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("dropdown missing %q in:\n%s", want, got)
		}
	}
}

// TestDropdownHasNoStaticAriaExpanded guards the honest-ARIA rule: before JS
// runs there is no open state, so a hardcoded "false" would be a lie for
// assistive tech that reads the markup before hydration.
func TestDropdownHasNoStaticAriaExpanded(t *testing.T) {
	got := renderUI(t, DropdownTrigger("Actions", "", nil))
	if strings.Contains(got, "aria-expanded") {
		t.Errorf("trigger must not hardcode aria-expanded before the runtime runs:\n%s", got)
	}
}

func TestContextMenuItems(t *testing.T) {
	props := ContextMenuProps{
		ID: "cm",
		Items: []ContextMenuItem{
			{Label: "Open", Icon: "lucide--external-link", Href: "#"},
			{Divider: true},
			{Label: "Delete", Danger: true, Href: "#"},
			{Label: "Server action", HXGet: "/x", HXTarget: "#y"},
		},
	}
	got := renderUI(t, ContextMenu(props))

	checks := map[string]string{
		"region marker":  "data-gd-context-root",
		"menu marker":    "data-gd-context-menu",
		"menu id":        `id="cm"`,
		"menu role":      `role="menu"`,
		"divider":        `role="separator"`,
		"danger styling": "text-error",
		"htmx attr":      `hx-get="/x"`,
		"htmx target":    `hx-target="#y"`,
	}
	for name, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("%s: missing %q", name, want)
		}
	}
	if n := strings.Count(got, "data-gd-menu-item"); n != 3 {
		t.Errorf("expected 3 menu items, got %d", n)
	}
	if strings.Contains(got, `href=""`) {
		t.Error("items without Href must not render an empty href")
	}
}

func TestResizableSeparatorAria(t *testing.T) {
	props := ResizableProps{
		ID:          "split",
		Initial:     "40%",
		Min:         "20%",
		Max:         "75%",
		PersistKey:  "split-key",
		Orientation: "",
		Start:       uiStr("start"),
		End:         uiStr("end"),
	}
	got := renderUI(t, Resizable(props))
	for _, want := range []string{
		`id="split"`,
		"data-gd-resizable",
		`data-orientation="horizontal"`,
		`data-initial="40"`,
		`data-min="20"`,
		`data-max="75"`,
		`data-persist="split-key"`,
		`role="separator"`,
		`aria-orientation="vertical"`,
		`aria-valuenow="40"`,
		`aria-valuemin="20"`,
		`aria-valuemax="75"`,
		"--gd-split:40%",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestResizableVerticalOrientation(t *testing.T) {
	got := renderUI(t, Resizable(ResizableProps{Orientation: "vertical", Start: uiStr("a"), End: uiStr("b")}))
	for _, want := range []string{`data-orientation="vertical"`, "flex-col", `aria-orientation="horizontal"`, "cursor-row-resize"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestResizableSanitisesBadPercentages(t *testing.T) {
	got := renderUI(t, Resizable(ResizableProps{Initial: "nonsense", Min: "-3", Max: "150", Start: uiStr("a"), End: uiStr("b")}))
	for _, want := range []string{`data-initial="50"`, `data-min="10"`, `data-max="90"`} {
		if !strings.Contains(got, want) {
			t.Errorf("fallback missing %q in:\n%s", want, got)
		}
	}
}

func TestScrollAreaVariants(t *testing.T) {
	vertical := renderUI(t, ScrollArea(ScrollAreaProps{MaxHeight: "12rem", Class: "border"}))
	if !strings.Contains(vertical, "gd-scrollarea") || strings.Contains(vertical, "gd-scrollarea-x") {
		t.Errorf("vertical variant wrong:\n%s", vertical)
	}
	if !strings.Contains(vertical, "max-height:12rem") {
		t.Errorf("max-height missing:\n%s", vertical)
	}

	horizontal := renderUI(t, ScrollArea(ScrollAreaProps{Orientation: "x"}))
	if !strings.Contains(horizontal, "gd-scrollarea-x") {
		t.Errorf("horizontal variant missing gd-scrollarea-x:\n%s", horizontal)
	}

	both := renderUI(t, ScrollArea(ScrollAreaProps{Orientation: "both"}))
	if !strings.Contains(both, "gd-scrollarea-both") {
		t.Errorf("both variant missing gd-scrollarea-both:\n%s", both)
	}
}

func TestDirection(t *testing.T) {
	if got := (DirectionProps{}).DirValue(); got != "ltr" {
		t.Errorf("zero DirValue = %q, want ltr", got)
	}
	if got := (DirectionProps{Dir: "rtl"}).IsRTL(); !got {
		t.Error("rtl should report IsRTL")
	}
	html := renderUI(t, Direction(DirectionProps{Dir: "rtl", Class: "gap-2"}))
	if !strings.Contains(html, `dir="rtl"`) || !strings.Contains(html, "gap-2") {
		t.Errorf("direction markup wrong:\n%s", html)
	}
}

func TestAttachmentKinds(t *testing.T) {
	cases := map[string]string{
		"image":   "lucide--image",
		"audio":   "lucide--music",
		"video":   "lucide--video",
		"archive": "lucide--file-archive",
		"code":    "lucide--file-code",
		"":        "lucide--file",
	}
	for kind, icon := range cases {
		got := renderUI(t, Attachment(AttachmentProps{Name: "file.bin", Kind: kind, Size: "1 MB"}))
		if !strings.Contains(got, icon) {
			t.Errorf("kind %q: want icon %q in:\n%s", kind, icon, got)
		}
	}
}

func TestAttachmentActions(t *testing.T) {
	got := renderUI(t, Attachment(AttachmentProps{
		Name:         "report.pdf",
		Href:         "/download/report.pdf",
		RemoveAction: "/files/1",
		RemoveTarget: "closest .attachment",
	}))
	for _, want := range []string{
		`href="/download/report.pdf"`,
		`aria-label="Download report.pdf"`,
		`hx-delete="/files/1"`,
		`hx-target="closest .attachment"`,
		`aria-label="Remove report.pdf"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

// TestComponentScriptsCarryNonce proves a strict CSP works end to end: with a
// nonce in the context, every script tag the component emits carries it.
func TestComponentScriptsCarryNonce(t *testing.T) {
	got := renderUI(t,
		shared.RenderInto(Dropdown("", nil), shared.Compose(
			DropdownTrigger("Actions", "", nil),
			DropdownMenu(nil),
		)),
		func(ctx context.Context) context.Context { return templ.WithNonce(ctx, "abc123") },
	)
	if !strings.Contains(got, `nonce="abc123"`) {
		t.Errorf("dropdown script tag missing nonce:\n%s", got)
	}
}

func TestThemeRuntimeTagNonce(t *testing.T) {
	plain := ThemeRuntimeTag(context.Background())
	if !strings.HasPrefix(plain, "<script>") {
		t.Errorf("no-nonce output should stay a plain script tag, got: %.40s", plain)
	}
	nonced := ThemeRuntimeTag(templ.WithNonce(context.Background(), "n0nce"))
	if !strings.HasPrefix(nonced, `<script nonce="n0nce">`) {
		t.Errorf("nonce missing, got: %.40s", nonced)
	}
}

// TestContextMenuStartsHidden pins the visibility mechanism. The menu is a
// fixed-position list that must stay out of the layout until opened; the runtime
// toggles the `hidden` attribute, so a static Tailwind `hidden` class (which the
// JS cannot remove) would make the menu impossible to open.
func TestContextMenuStartsHidden(t *testing.T) {
	got := renderUI(t, ContextMenu(ContextMenuProps{
		Items: []ContextMenuItem{{Label: "Open", Href: "#"}},
	}))
	start := strings.Index(got, "<ul")
	if start < 0 {
		t.Fatalf("no menu list in output:\n%s", got)
	}
	end := strings.Index(got[start:], ">")
	if end < 0 {
		t.Fatalf("unterminated menu tag:\n%s", got)
	}
	tag := got[start : start+end]

	if !strings.Contains(tag, " hidden") && !strings.Contains(tag, "\thidden") {
		t.Errorf("menu must start with the hidden attribute: %s", tag)
	}
	if classStart := strings.Index(tag, "class=\""); classStart >= 0 {
		classEnd := strings.Index(tag[classStart:], "\"")
		class := tag[classStart : classStart+classEnd]
		for _, token := range strings.Fields(class) {
			if token == "hidden" {
				t.Errorf("menu must not use the static Tailwind hidden class: %s", class)
			}
		}
	}
}

// TestDocumentDirectionScriptInterpolatesDir guards against templ's rule that
// <script> bodies are not interpolated: building the body in the template once
// shipped the literal text "{ dir }" and silently never set the direction.
func TestDocumentDirectionScriptInterpolatesDir(t *testing.T) {
	got := renderUI(t, DocumentDirectionScript("rtl"))
	if !strings.Contains(got, `var target="rtl"`) {
		t.Errorf("direction value was not interpolated:\n%s", got)
	}
	if strings.Contains(got, "{ dir }") {
		t.Errorf("template expression leaked into the script body:\n%s", got)
	}
	if strings.Contains(got, `nonce=""`) {
		t.Errorf("empty nonce must be omitted, got:\n%s", got)
	}
}
