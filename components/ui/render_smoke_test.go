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

func intPtr(n int) *int { return &n }

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

func TestButtonGroupContract(t *testing.T) {
	horizontal := renderUI(t, shared.RenderInto(ButtonGroup(ButtonGroupProps{}), uiStr("x")))
	for _, want := range []string{`role="group"`, `data-slot="button-group"`} {
		if !strings.Contains(horizontal, want) {
			t.Errorf("horizontal group missing %q:\n%s", want, horizontal)
		}
	}
	if strings.Contains(horizontal, "join-vertical") {
		t.Errorf("horizontal group must not add join-vertical:\n%s", horizontal)
	}

	vertical := renderUI(t, shared.RenderInto(ButtonGroup(ButtonGroupProps{Orientation: ButtonGroupVertical}), uiStr("x")))
	if !strings.Contains(vertical, "join-vertical") {
		t.Errorf("vertical group missing join-vertical:\n%s", vertical)
	}
}

func TestButtonGroupSeparatorOrientation(t *testing.T) {
	horizontal := renderUI(t, ButtonGroupSeparator(ButtonGroupSeparatorProps{Orientation: ButtonGroupHorizontal}))
	for _, want := range []string{`role="separator"`, `aria-orientation="vertical"`, "w-px self-stretch"} {
		if !strings.Contains(horizontal, want) {
			t.Errorf("horizontal separator missing %q:\n%s", want, horizontal)
		}
	}
	if strings.Contains(horizontal, "h-px w-full") {
		t.Errorf("horizontal separator must not render a horizontal line:\n%s", horizontal)
	}

	vertical := renderUI(t, ButtonGroupSeparator(ButtonGroupSeparatorProps{Orientation: ButtonGroupVertical}))
	for _, want := range []string{`aria-orientation="horizontal"`, "h-px w-full"} {
		if !strings.Contains(vertical, want) {
			t.Errorf("vertical separator missing %q:\n%s", want, vertical)
		}
	}
}

func TestMarkerVariants(t *testing.T) {
	def := renderUI(t, shared.RenderInto(Marker(MarkerProps{}), uiStr("content")))
	if !strings.Contains(def, `<div data-slot="marker"`) {
		t.Errorf("default marker should render a div:\n%s", def)
	}
	if !strings.Contains(def, `data-variant=""`) {
		t.Errorf("default marker missing data-variant:\n%s", def)
	}

	sep := renderUI(t, shared.RenderInto(Marker(MarkerProps{Variant: MarkerSeparator}), uiStr("content")))
	if !strings.Contains(sep, "before:h-px before:flex-1") || !strings.Contains(sep, "after:h-px after:flex-1") {
		t.Errorf("separator variant missing line classes:\n%s", sep)
	}

	border := renderUI(t, shared.RenderInto(Marker(MarkerProps{Variant: MarkerBorder}), uiStr("content")))
	if !strings.Contains(border, "rounded-box border border-base-300") {
		t.Errorf("border variant missing box classes:\n%s", border)
	}

	link := renderUI(t, shared.RenderInto(Marker(MarkerProps{Href: "#"}), uiStr("content")))
	if !strings.Contains(link, `<a href="#"`) {
		t.Errorf("marker with href should render an anchor:\n%s", link)
	}
	if strings.Contains(link, `<div data-slot="marker"`) {
		t.Errorf("marker with href must not render a div:\n%s", link)
	}

	// hrefs are caller-supplied, so unsafe schemes must be neutralised by templ.URL.
	unsafe := renderUI(t, shared.RenderInto(Marker(MarkerProps{Href: "javascript:alert(1)"}), uiStr("content")))
	if strings.Contains(unsafe, "javascript:") {
		t.Errorf("unsafe href was not sanitised:\n%s", unsafe)
	}
}

func TestListRowVariants(t *testing.T) {
	def := renderUI(t, ListRow(ListRowProps{Title: "A"}))
	if !strings.Contains(def, "list-row") {
		t.Errorf("default row missing list-row:\n%s", def)
	}
	if strings.Contains(def, "border-base-300") {
		t.Errorf("default row must not have border-base-300:\n%s", def)
	}

	outline := renderUI(t, ListRow(ListRowProps{Title: "A", Variant: ListRowOutline}))
	if !strings.Contains(outline, "rounded-box border border-base-300") {
		t.Errorf("outline row missing box classes:\n%s", outline)
	}

	muted := renderUI(t, ListRow(ListRowProps{Title: "A", Variant: ListRowMuted}))
	if !strings.Contains(muted, "bg-base-200") {
		t.Errorf("muted row missing bg-base-200:\n%s", muted)
	}

	sm := renderUI(t, ListRow(ListRowProps{Title: "A", Size: ListRowSizeSM}))
	if !strings.Contains(sm, "py-1") {
		t.Errorf("sm row missing py-1:\n%s", sm)
	}
	xs := renderUI(t, ListRow(ListRowProps{Title: "A", Size: ListRowSizeXS}))
	if !strings.Contains(xs, "py-0.5") {
		t.Errorf("xs row missing py-0.5:\n%s", xs)
	}
	if !strings.Contains(xs, "text-xs") {
		t.Errorf("xs row title should be text-xs:\n%s", xs)
	}

	icon := renderUI(t, ListRow(ListRowProps{Title: "A", Leading: uiStr("x"), MediaVariant: ListMediaIcon}))
	if !strings.Contains(icon, "size-9") || !strings.Contains(icon, "text-muted") {
		t.Errorf("icon media wrapper missing classes:\n%s", icon)
	}
	image := renderUI(t, ListRow(ListRowProps{Title: "A", Leading: uiStr("x"), MediaVariant: ListMediaImage}))
	if !strings.Contains(image, "overflow-hidden") || !strings.Contains(image, "[&_img]:object-cover") {
		t.Errorf("image media wrapper missing classes:\n%s", image)
	}

	hf := renderUI(t, ListRow(ListRowProps{Title: "A", Header: uiStr("HDR"), Footer: uiStr("FTR")}))
	if !strings.Contains(hf, "HDR") || !strings.Contains(hf, "FTR") {
		t.Errorf("header/footer text missing:\n%s", hf)
	}
	if !strings.Contains(hf, "text-label") || !strings.Contains(hf, "text-subtle") {
		t.Errorf("header/footer tone classes missing:\n%s", hf)
	}
	// Header lives inside the growing content column and Footer shares the
	// single wrapped cell, so neither may add a second list-col-wrap grid item.
	if strings.Count(hf, "list-col-wrap") != 1 {
		t.Errorf("header/footer must not add extra grid rows (want 1 list-col-wrap):\n%s", hf)
	}
	if strings.Contains(hf, `class="list-col-wrap text-xs font-semibold`) {
		t.Errorf("header must not be rendered as a wrapped row:\n%s", hf)
	}

	// MediaVariant is ignored when LeadingGrow is set, so the growing slot
	// keeps its list-col-grow wrapper.
	growMedia := renderUI(t, ListRow(ListRowProps{Title: "A", Leading: uiStr("x"), LeadingGrow: true, MediaVariant: ListMediaIcon}))
	if strings.Contains(growMedia, "size-9") {
		t.Errorf("LeadingGrow must win over MediaVariant:\n%s", growMedia)
	}
	if !strings.Contains(growMedia, `list-col-grow`) {
		t.Errorf("LeadingGrow media slot lost list-col-grow:\n%s", growMedia)
	}

	link := renderUI(t, ListRow(ListRowProps{Title: "A", Href: "https://example.com"}))
	if !strings.Contains(link, `<a class="list-row"`) || !strings.Contains(link, "href=") {
		t.Errorf("href row missing anchor:\n%s", link)
	}

	unsafe := renderUI(t, ListRow(ListRowProps{Title: "A", Href: "javascript:alert(1)"}))
	if strings.Contains(unsafe, "javascript:") {
		t.Errorf("unsafe href was not neutralised:\n%s", unsafe)
	}

	blank := renderUI(t, ListRow(ListRowProps{Title: "A", Href: "https://example.com", Target: "_blank"}))
	if !strings.Contains(blank, `rel="noopener noreferrer"`) {
		t.Errorf("_blank target missing rel:\n%s", blank)
	}

	sep := renderUI(t, ItemSeparator())
	if !strings.Contains(sep, `role="separator"`) {
		t.Errorf("separator missing role:\n%s", sep)
	}
}

// TestToggleButtonContract locks the aria-pressed toggle contract: the zero value
// renders an unpressed button with no variant/active/label classes, and each
// prop maps to exactly one emitted attribute or class.
func TestToggleButtonContract(t *testing.T) {
	def := renderUI(t, ToggleButton(ToggleButtonProps{}))
	for _, want := range []string{`type="button"`, `data-slot="toggle"`, `aria-pressed="false"`} {
		if !strings.Contains(def, want) {
			t.Errorf("ToggleButton missing %q:\n%s", want, def)
		}
	}
	for _, unwanted := range []string{"btn-outline", "btn-active", "btn-sm", "btn-lg", "aria-label="} {
		if strings.Contains(def, unwanted) {
			t.Errorf("ToggleButton zero value must not emit %q:\n%s", unwanted, def)
		}
	}

	pressed := renderUI(t, ToggleButton(ToggleButtonProps{Pressed: true}))
	if !strings.Contains(pressed, `aria-pressed="true"`) || !strings.Contains(pressed, "btn-active") {
		t.Errorf("pressed ToggleButton missing pressed state:\n%s", pressed)
	}

	outline := renderUI(t, ToggleButton(ToggleButtonProps{Variant: ToggleOutline}))
	if !strings.Contains(outline, "btn-outline") {
		t.Errorf("outline ToggleButton missing btn-outline:\n%s", outline)
	}

	sm := renderUI(t, ToggleButton(ToggleButtonProps{Size: ToggleSizeSM}))
	if !strings.Contains(sm, "btn-sm") {
		t.Errorf("sm ToggleButton missing btn-sm:\n%s", sm)
	}
	lg := renderUI(t, ToggleButton(ToggleButtonProps{Size: ToggleSizeLG}))
	if !strings.Contains(lg, "btn-lg") {
		t.Errorf("lg ToggleButton missing btn-lg:\n%s", lg)
	}

	disabled := renderUI(t, ToggleButton(ToggleButtonProps{Disabled: true}))
	if !strings.Contains(disabled, " disabled") {
		t.Errorf("disabled ToggleButton missing disabled attribute:\n%s", disabled)
	}

	valued := renderUI(t, ToggleButton(ToggleButtonProps{Value: "center"}))
	if !strings.Contains(valued, `data-gd-toggle-value="center"`) {
		t.Errorf("ToggleButton missing data-gd-toggle-value:\n%s", valued)
	}
	if strings.Contains(def, "data-gd-toggle-value") {
		t.Errorf("ToggleButton without Value must not emit an empty data-gd-toggle-value:\n%s", def)
	}

	labelled := renderUI(t, ToggleButton(ToggleButtonProps{AriaLabel: "Bold", Icon: "lucide--bold"}))
	if !strings.Contains(labelled, `aria-label="Bold"`) {
		t.Errorf("ToggleButton missing aria-label:\n%s", labelled)
	}
	if !strings.Contains(labelled, "lucide--bold") {
		t.Errorf("ToggleButton missing icon:\n%s", labelled)
	}
}

// TestToggleGroupContract locks the container classes and conditional data
// attributes the toggle runtime depends on.
func TestToggleGroupContract(t *testing.T) {
	joined := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{Spacing: intPtr(0)}), uiStr("x")))
	for _, want := range []string{`role="group"`, `data-gd-toggle-group`, "join"} {
		if !strings.Contains(joined, want) {
			t.Errorf("joined group missing %q:\n%s", want, joined)
		}
	}
	if strings.Contains(joined, "gap-") || strings.Contains(joined, "join-vertical") {
		t.Errorf("horizontal joined group must not add gap or join-vertical:\n%s", joined)
	}

	vJoined := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{Orientation: ToggleGroupVertical, Spacing: intPtr(0)}), uiStr("x")))
	if !strings.Contains(vJoined, "join-vertical") {
		t.Errorf("vertical joined group missing join-vertical:\n%s", vJoined)
	}

	spaced := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{}), uiStr("x")))
	if !strings.Contains(spaced, "gap-2") {
		t.Errorf("default-spaced group missing gap-2:\n%s", spaced)
	}
	if strings.Contains(spaced, "join") {
		t.Errorf("spaced group must not add join:\n%s", spaced)
	}

	plain := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{}), uiStr("x")))
	if strings.Contains(plain, "data-multiple") || strings.Contains(plain, "data-controlled") {
		t.Errorf("plain group must not emit data-multiple/data-controlled:\n%s", plain)
	}

	multi := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{Multiple: true}), uiStr("x")))
	if !strings.Contains(multi, "data-multiple") || strings.Contains(multi, "data-controlled") {
		t.Errorf("multiple group must emit only data-multiple:\n%s", multi)
	}

	controlled := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{Value: []string{"a"}}), uiStr("x")))
	if !strings.Contains(controlled, "data-controlled") {
		t.Errorf("controlled group missing data-controlled:\n%s", controlled)
	}
}

// TestToggleGroupChildrenInheritDefaults proves the group's variant/size/disabled
// propagate to child ToggleButtons via the render context, and that a child only
// carries data-gd-toggle-value when it has an explicit Value.
func TestToggleGroupChildrenInheritDefaults(t *testing.T) {
	got := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{
		Variant:  ToggleOutline,
		Size:     ToggleSizeSM,
		Disabled: true,
		Value:    []string{"a"},
	}), ToggleButton(ToggleButtonProps{Value: "a"})))

	for _, want := range []string{"btn-outline", "btn-sm", " disabled", `data-gd-toggle-value="a"`, `data-gd-toggle-values="[&#34;a&#34;]"`, "data-controlled"} {
		if !strings.Contains(got, want) {
			t.Errorf("child must inherit/carry %q from group:\n%s", want, got)
		}
	}

	// An inherited group must not invent an empty value for a child that has
	// none: every item in a group needs a distinct value.
	valueless := renderUI(t, shared.RenderInto(ToggleGroup(ToggleGroupProps{}), ToggleButton(ToggleButtonProps{})))
	if strings.Contains(valueless, `data-gd-toggle-value=""`) {
		t.Errorf("valueless child must not emit an empty data-gd-toggle-value:\n%s", valueless)
	}
	if !strings.Contains(valueless, `data-gd-toggle-values="[]"`) {
		t.Errorf("group with no values must emit an empty JSON array:\n%s", valueless)
	}
}

// TestMessagePartsContract locks the data-slot/data-align markers and the
// group-data alignment variants the message stack relies on. The alignment of
// the content/header/footer slots is derived from the parent group/message
// marker, so the slots must not carry their own Go alignment props.
func TestMessagePartsContract(t *testing.T) {
	start := renderUI(t, Message(MessageProps{}))
	for _, want := range []string{
		`data-slot="message"`,
		`data-align=""`,
		"group/message",
	} {
		if !strings.Contains(start, want) {
			t.Errorf("start Message missing %q:\n%s", want, start)
		}
	}
	if strings.Contains(start, "flex-row-reverse") {
		t.Errorf("start Message must not add flex-row-reverse:\n%s", start)
	}

	end := renderUI(t, Message(MessageProps{Align: MessageEnd}))
	if !strings.Contains(end, `data-align="end"`) {
		t.Errorf("end Message missing data-align=end:\n%s", end)
	}
	if !strings.Contains(end, "flex-row-reverse") {
		t.Errorf("end Message missing flex-row-reverse:\n%s", end)
	}

	content := renderUI(t, MessageContent(MessageContentProps{}))
	for _, want := range []string{
		`data-slot="message-content"`,
		"group-data-[align=end]/message:items-end",
		"group-data-[align=end]/message:text-end",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("MessageContent missing %q:\n%s", want, content)
		}
	}

	header := renderUI(t, MessageHeader(MessageHeaderProps{}))
	if !strings.Contains(header, `data-slot="message-header"`) || !strings.Contains(header, "group-data-[align=end]/message:flex-row-reverse") {
		t.Errorf("MessageHeader missing header contract:\n%s", header)
	}

	footer := renderUI(t, MessageFooter(MessageFooterProps{}))
	if !strings.Contains(footer, `data-slot="message-footer"`) || !strings.Contains(footer, "text-muted") {
		t.Errorf("MessageFooter missing footer contract:\n%s", footer)
	}

	def := renderUI(t, MessageBubble(MessageBubbleProps{}))
	if !strings.Contains(def, "bg-base-100") || !strings.Contains(def, "border-base-300") {
		t.Errorf("default bubble missing default surface:\n%s", def)
	}
	primary := renderUI(t, MessageBubble(MessageBubbleProps{Variant: MessageBubblePrimary}))
	if !strings.Contains(primary, "bg-primary/10") || !strings.Contains(primary, "text-primary") {
		t.Errorf("primary bubble missing primary surface:\n%s", primary)
	}
	muted := renderUI(t, MessageBubble(MessageBubbleProps{Variant: MessageBubbleMuted}))
	if !strings.Contains(muted, "bg-base-200") {
		t.Errorf("muted bubble missing muted surface:\n%s", muted)
	}
	for _, out := range []string{def, primary, muted} {
		if !strings.Contains(out, "group-data-[align=end]/message:bg-primary") {
			t.Errorf("bubble missing end alignment group variant:\n%s", out)
		}
	}
}

// TestMessageScrollerContract locks the data attributes the scroller runtime
// depends on, the static marker, the jump button, and the height default.
func TestMessageScrollerContract(t *testing.T) {
	got := renderUI(t, MessageScroller(MessageScrollerProps{}))
	for _, want := range []string{
		"data-gd-messagescroller",
		"data-gd-messagescroller-viewport",
		"data-gd-messagescroller-jump",
		"overflow-y-auto",
		"h-96",
		"lucide--arrow-down",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("MessageScroller missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "data-gd-messagescroller-static") {
		t.Errorf("default MessageScroller must not emit static marker:\n%s", got)
	}

	static := renderUI(t, MessageScroller(MessageScrollerProps{Static: true}))
	if !strings.Contains(static, "data-gd-messagescroller-static") {
		t.Errorf("static MessageScroller missing static marker:\n%s", static)
	}

	hidden := renderUI(t, MessageScroller(MessageScrollerProps{HideJumpButton: true}))
	if strings.Contains(hidden, "data-gd-messagescroller-jump") {
		t.Errorf("HideJumpButton must omit the jump button:\n%s", hidden)
	}

	custom := renderUI(t, MessageScroller(MessageScrollerProps{HeightClass: "h-64"}))
	if !strings.Contains(custom, "h-64") {
		t.Errorf("custom height class missing:\n%s", custom)
	}
}

// TestChatBubbleLegacyShim proves the deprecated flat-parameter shim renders the
// same output as the props-struct API for the same inputs.
func TestChatBubbleLegacyShim(t *testing.T) {
	legacy := renderUI(t, ChatBubbleLegacy(true, "You", "10:33 AM", "", false, "bg-base-200", false, nil, uiStr("hi")))
	props := renderUI(t, shared.RenderInto(ChatBubble(ChatBubbleProps{Sent: true, Author: "You", Timestamp: "10:33 AM", BubbleClass: "bg-base-200"}), uiStr("hi")))

	for _, want := range []string{"chat group", "chat-end", "chat-bubble", "bg-base-200"} {
		if !strings.Contains(legacy, want) {
			t.Errorf("legacy shim missing %q:\n%s", want, legacy)
		}
	}
	if legacy != props {
		t.Errorf("legacy shim output differs from props API:\nlegacy: %s\nprops:  %s", legacy, props)
	}
}

// TestAlertDialogContract locks the native <dialog> alert contract: the dialog
// id drives the derived title/description ids (via the render context), the
// trigger carries the open marker, the content is a method="dialog" form, and
// the footer/cancel/action emit their native submit values.
func TestAlertDialogContract(t *testing.T) {
	inner := shared.Compose(
		shared.RenderInto(AlertDialogTrigger(AlertDialogTriggerProps{DialogID: "x"}), uiStr("Delete")),
		shared.RenderInto(AlertDialog(AlertDialogProps{ID: "x"}), shared.Compose(
			shared.RenderInto(AlertDialogContent(AlertDialogContentProps{}), shared.Compose(
				shared.RenderInto(AlertDialogHeader(AlertDialogHeaderProps{}), shared.Compose(
					shared.RenderInto(AlertDialogTitle(AlertDialogTitleProps{}), uiStr("Delete account")),
					shared.RenderInto(AlertDialogDescription(AlertDialogDescriptionProps{}), uiStr("This cannot be undone.")),
				)),
				shared.RenderInto(AlertDialogFooter(AlertDialogFooterProps{}), shared.Compose(
					shared.RenderInto(AlertDialogCancel(AlertDialogCancelProps{}), uiStr("Cancel")),
					shared.RenderInto(AlertDialogAction(AlertDialogActionProps{Variant: ButtonError}), uiStr("Delete")),
				)),
			)),
		)),
	)
	got := renderUI(t, inner)

	for _, want := range []string{
		`<dialog id="x"`,
		"data-gd-alert-dialog",
		`aria-labelledby="x-title"`,
		`aria-describedby="x-desc"`,
		`class="modal "`,
		`data-gd-alert-dialog-open="x"`,
		`<form method="dialog"`,
		`data-slot="alert-dialog-content"`,
		"modal-box",
		`data-slot="alert-dialog-header"`,
		`<h2 id="x-title"`,
		`<p id="x-desc"`,
		"text-muted",
		`data-slot="alert-dialog-footer"`,
		"modal-action",
		`type="submit" value="cancel"`,
		`type="submit" value="confirm"`,
		"btn-error",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("alert dialog missing %q in:\n%s", want, got)
		}
	}
}

// TestAlertDialogTitleFallsBackWithoutContext proves a title/description
// rendered outside an AlertDialog emits no empty id attribute.
func TestAlertDialogTitleFallsBackWithoutContext(t *testing.T) {
	title := renderUI(t, shared.RenderInto(AlertDialogTitle(AlertDialogTitleProps{}), uiStr("Loose")))
	if strings.Contains(title, `id=""`) {
		t.Errorf("title outside a dialog must not emit an empty id:\n%s", title)
	}
	if !strings.Contains(title, "text-lg font-semibold") {
		t.Errorf("title missing base classes:\n%s", title)
	}

	desc := renderUI(t, shared.RenderInto(AlertDialogDescription(AlertDialogDescriptionProps{}), uiStr("Loose")))
	if strings.Contains(desc, `id=""`) {
		t.Errorf("description outside a dialog must not emit an empty id:\n%s", desc)
	}
}

// TestMenubarContract locks the data attributes the a11y + menubar runtimes
// depend on, plus the honest-ARIA rule (no static aria-expanded) and the
// disabled/shortcut/tone conventions.
func TestMenubarContract(t *testing.T) {
	file := shared.RenderInto(MenubarMenu(MenubarMenuProps{}), shared.Compose(
		shared.RenderInto(MenubarTrigger(MenubarTriggerProps{}), uiStr("File")),
		shared.RenderInto(MenubarContent(MenubarContentProps{}), shared.Compose(
			shared.RenderInto(MenubarItem(MenubarItemProps{Href: "#"}), uiStr("New")),
			MenubarSeparator(MenubarSeparatorProps{}),
			shared.RenderInto(MenubarLabel(MenubarLabelProps{}), uiStr("Zoom")),
			shared.RenderInto(MenubarItem(MenubarItemProps{Disabled: true}), uiStr("Paste")),
		)),
	))
	got := renderUI(t, shared.RenderInto(Menubar(MenubarProps{}), file))

	for _, want := range []string{
		`role="menubar"`,
		"data-gd-menubar",
		"data-gd-menubar-menu",
		"data-gd-menu-root",
		"data-gd-menu-trigger",
		"data-gd-menubar-trigger",
		`aria-haspopup="menu"`,
		`role="menu"`,
		"data-gd-menubar-content",
		`role="menuitem"`,
		"data-gd-menu-item",
		`role="separator"`,
		"menu-label",
		"go-daisy-menubar.js",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("menubar missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "aria-expanded") {
		t.Errorf("trigger must not hardcode aria-expanded before the runtime runs:\n%s", got)
	}

	// Href items neutralise unsafe schemes.
	unsafe := renderUI(t, shared.RenderInto(MenubarItem(MenubarItemProps{Href: "javascript:alert(1)"}), uiStr("x")))
	if strings.Contains(unsafe, "javascript:") {
		t.Errorf("unsafe href not neutralised:\n%s", unsafe)
	}

	// Disabled button items use the native attribute.
	btn := renderUI(t, shared.RenderInto(MenubarItem(MenubarItemProps{Disabled: true}), uiStr("x")))
	if !strings.Contains(btn, " disabled") {
		t.Errorf("disabled button item missing disabled:\n%s", btn)
	}

	// Disabled link items use aria-disabled plus pointer-events off.
	link := renderUI(t, shared.RenderInto(MenubarItem(MenubarItemProps{Href: "#", Disabled: true}), uiStr("x")))
	for _, want := range []string{`aria-disabled="true"`, "pointer-events-none", "opacity-50"} {
		if !strings.Contains(link, want) {
			t.Errorf("disabled link item missing %q:\n%s", want, link)
		}
	}

	// Danger items tint text-error.
	danger := renderUI(t, shared.RenderInto(MenubarItem(MenubarItemProps{Danger: true}), uiStr("x")))
	if !strings.Contains(danger, "text-error") {
		t.Errorf("danger item missing text-error:\n%s", danger)
	}

	// Shortcuts use the faint tone.
	short := renderUI(t, shared.RenderInto(MenubarShortcut(MenubarShortcutProps{}), uiStr("⌘N")))
	if !strings.Contains(short, "text-faint") {
		t.Errorf("shortcut missing faint tone:\n%s", short)
	}
}
