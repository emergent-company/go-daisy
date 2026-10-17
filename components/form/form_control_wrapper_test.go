package form

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	ui "github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/devmode"
	"github.com/emergent-company/go-daisy/shared"
)

// renderDev renders a component with dev-mode annotation enabled so the
// data-component attributes are emitted (they are stripped in production).
func renderDev(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := devmode.WithDevMode(context.Background())
	if err := c.Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// firstTag returns the opening tag of the rendered document.
func firstTag(s string) string {
	end := strings.Index(s, ">")
	if end < 0 {
		return s
	}
	return s[:end+1]
}

// fieldGroupComponents lists every component that used to render the
// daisyUI-v4 `form-control` class on its root wrapper, paired with the root
// class the component is expected to keep and the dev-mode component name.
//
// `form-control` has no rules in daisyUI 5, so emitting it was inert markup that
// read as a layout contract that did not exist. Each root now carries only
// classes with live rules (see #71).
func fieldGroupComponents(t *testing.T) []struct {
	name      string
	component templ.Component
	wantRoot  string // exact opening tag prefix after `<div ` or `<label `
	devName   string
} {
	t.Helper()
	return []struct {
		name      string
		component templ.Component
		wantRoot  string
		devName   string
	}{
		{
			name:      "FormToggle",
			component: FormToggle("notify", "Notify", true, "", nil),
			wantRoot:  `<div class="flex flex-col items-start"`,
			devName:   "form/FormToggle",
		},
		{
			name:      "FormCheckbox",
			component: FormCheckbox("agree", "Agree", false, "", nil),
			wantRoot:  `<div class="flex flex-col items-start"`,
			devName:   "form/FormCheckbox",
		},
		{
			name:      "CheckboxInput",
			component: CheckboxInput("agree", "Agree", false, ""),
			wantRoot:  `<div class="flex flex-col items-start mb-3"`,
			devName:   "form/CheckboxInput",
		},
		{
			name:      "Field",
			component: fieldWithInput(),
			wantRoot:  `<div class="mb-3"`,
			devName:   "form/Field",
		},
		{
			name:      "FormField",
			component: FormField(FormFieldProps{Type: FieldText, Name: "email", Label: "Email"}),
			wantRoot:  `<div class="mb-3"`,
			devName:   "form/Field",
		},
		{
			name:      "FormFieldCheckbox",
			component: FormField(FormFieldProps{Type: FieldCheckbox, Name: "agree", Label: "Agree", Error: "Required"}),
			wantRoot:  `<div class="flex flex-col items-start mb-3"`,
			devName:   "form/FormField",
		},
		{
			name:      "ValidatedField",
			component: validatedFieldWithInput(),
			wantRoot:  `<div class="mb-3"`,
			devName:   "form/ValidatedField",
		},
		{
			name:      "DatePicker",
			component: DatePicker(DatePickerProps{Name: "due", Label: "Due"}),
			wantRoot:  `<div class="w-full"`,
			devName:   "form/DatePicker",
		},
		{
			name:      "EnhancedSelect",
			component: EnhancedSelect(EnhancedSelectProps{Name: "country", Label: "Country"}),
			wantRoot:  `<div class="w-full"`,
			devName:   "form/EnhancedSelect",
		},
		{
			name:      "FileUpload",
			component: FileUpload(FileUploadProps{Name: "avatar", Label: "Avatar"}),
			wantRoot:  `<div class="w-full"`,
			devName:   "form/FileUpload",
		},
		{
			name:      "RichTextEditor",
			component: RichTextEditor(RichTextEditorProps{ID: "body", Label: "Body"}),
			wantRoot:  `<div class="w-full"`,
			devName:   "form/RichTextEditor",
		},
		{
			name:      "Combobox",
			component: Combobox(ComboboxProps{ID: "pick", Name: "pick", Label: "Pick"}),
			wantRoot:  `<div class="w-full"`,
			devName:   "form/Combobox",
		},
		{
			name:      "StructuredInput",
			component: StructuredInput(StructuredInputProps{Name: "rows"}),
			wantRoot:  `<div class="w-full space-y-2"`,
			devName:   "form/StructuredInput",
		},
		{
			name:      "ThemeController",
			component: ui.ThemeController(ui.ThemeControllerProps{Theme: "dark", InputType: ui.ThemeCheckbox, Label: "Dark mode"}),
			wantRoot:  `<label class="label cursor-pointer gap-3"`,
			devName:   "ui/ThemeController",
		},
	}
}

func fieldWithInput() templ.Component {
	return shared.RenderInto(Field("Email", "email", ""), templ.Raw(`<input id="email" class="input w-full">`))
}

func validatedFieldWithInput() templ.Component {
	return shared.RenderInto(ValidatedField("Email", "Must be valid", "email"), templ.Raw(`<input id="email" type="email" class="input w-full">`))
}

// TestFieldGroupRootsDoNotEmitFormControl is the regression net for #71: no
// field-group component may reintroduce the daisyUI-v4 `form-control` class,
// which has no rules in the shipped daisyUI 5 stylesheet.
func TestFieldGroupRootsDoNotEmitFormControl(t *testing.T) {
	for _, tc := range fieldGroupComponents(t) {
		t.Run(tc.name, func(t *testing.T) {
			out := strings.TrimLeft(renderDev(t, tc.component), " \n\t")
			if strings.Contains(out, "form-control") {
				t.Errorf("rendered output still emits the inert daisyUI-4 class:\n%s", firstTag(out))
			}
			if !strings.HasPrefix(out, tc.wantRoot) {
				t.Errorf("root classes changed unexpectedly:\n got: %s\nwant prefix: %s", firstTag(out), tc.wantRoot)
			}
		})
	}
}

// TestFieldGroupRootsKeepDevModeAnnotation guards the gallery contract: the
// hover overlay and component tree resolve data-component on the component's
// root element. Every root here must still be annotated — including the two
// components whose wrapper div was removed.
func TestFieldGroupRootsKeepDevModeAnnotation(t *testing.T) {
	for _, tc := range fieldGroupComponents(t) {
		t.Run(tc.name, func(t *testing.T) {
			out := strings.TrimLeft(renderDev(t, tc.component), " \n\t")
			tag := firstTag(out)
			if !strings.Contains(tag, `data-component="`+tc.devName+`"`) {
				t.Errorf("root tag lost its dev-mode annotation:\n got: %s\nwant: data-component=%q", tag, tc.devName)
			}
		})
	}
}

// TestFieldGroupRootsKeepRealSpacingUtilities asserts the utilities that were
// load-bearing next to the removed class are still present: `mb-3` on stacked
// fields, `w-full` on full-width pickers, `space-y-2` on the structured input,
// and the flex column on the toggle/checkbox groups.
func TestFieldGroupRootsKeepRealSpacingUtilities(t *testing.T) {
	cases := []struct {
		name string
		c    templ.Component
		want []string
	}{
		{"field_spacing", fieldWithInput(), []string{`class="mb-3"`}},
		{"checked_group_stack", FormToggle("n", "L", false, "err", nil), []string{`class="flex flex-col items-start"`, `class="label pt-1"`}},
		{"datepicker_width", DatePicker(DatePickerProps{Name: "due"}), []string{`class="w-full"`}},
		{"structured_spacing", StructuredInput(StructuredInputProps{Name: "rows"}), []string{`class="w-full space-y-2"`}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := renderDev(t, tc.c)
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Errorf("missing %s in rendered output", want)
				}
			}
		})
	}
}

// attrValue returns the value of name="..." inside an opening tag, or "".
func attrValue(tag, name string) string {
	key := name + `="`
	i := strings.Index(tag, key)
	if i < 0 {
		return ""
	}
	rest := tag[i+len(key):]
	j := strings.Index(rest, `"`)
	if j < 0 {
		return ""
	}
	return rest[:j]
}

// TestCheckableFieldRootsShareOneShape pins the invariant the per-component
// assertions above cannot: every root that renders checkableInput + FieldError
// is the same widget (a checkable input with an error slot), so it must be the
// same vertical stack — label first, error on its own line after it, never
// inline beside the control, and never a bare block wrapper with an
// inline-flex label (which reintroduces the line-box strut measured in #71).
//
// The follow-up bug this catches: FormCheckbox stacked its error while
// FormField's checkbox branch kept a plain `mb-3` block, giving callers two
// different error layouts and row heights for one widget.
func TestCheckableFieldRootsShareOneShape(t *testing.T) {
	cases := []struct {
		name string
		c    templ.Component
	}{
		{"FormToggle", FormToggle("n", "Label", true, "Required", nil)},
		{"FormCheckbox", FormCheckbox("n", "Label", true, "Required", nil)},
		{"FormFieldCheckbox", FormField(FormFieldProps{Type: FieldCheckbox, Name: "n", Label: "Label", Checked: true, Error: "Required"})},
		{"CheckboxInput", CheckboxInput("n", "Label", true, "Required")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := strings.TrimLeft(renderDev(t, tc.c), " \n\t")
			tag := firstTag(out)
			class := attrValue(tag, "class")

			// `mb-3` is field rhythm, not widget shape, so it may accompany the
			// stack but nothing else may.
			var shape []string
			for _, tok := range strings.Fields(class) {
				if tok == "mb-3" {
					continue
				}
				shape = append(shape, tok)
			}
			if got, want := strings.Join(shape, " "), "flex flex-col items-start"; got != want {
				t.Fatalf("checkable field root shape = %q, want %q (class=%q)", got, want, class)
			}

			labelEnd := strings.Index(out, "</label>")
			errStart := strings.Index(out, `class="label pt-1"`)
			if labelEnd < 0 || errStart < 0 {
				t.Fatalf("expected a <label> and an error slot, got: %s", out)
			}
			if errStart < labelEnd {
				t.Fatalf("error slot is not after the label (it would share the control's line):\n%s", out)
			}
			if n := strings.Count(out, `class="label pt-1"`); n != 1 {
				t.Fatalf("error slot count = %d, want 1:\n%s", n, out)
			}

			tagEnd := strings.Index(out, ">")
			if body := out[tagEnd+1:]; !strings.HasPrefix(body, `<label `) {
				t.Fatalf("label is not the first child of the group root: %s", body)
			}
		})
	}
}

// TestComboboxAppendsCallerClassToRoot locks the Class passthrough that shared
// the wrapper's class attribute with the removed class name.
func TestComboboxAppendsCallerClassToRoot(t *testing.T) {
	out := strings.TrimLeft(renderDev(t, Combobox(ComboboxProps{ID: "pick", Name: "pick", Class: "max-w-xs"})), " \n\t")
	if !strings.HasPrefix(out, `<div class="w-full max-w-xs"`) {
		t.Errorf("caller class not appended to root:\n got: %s", firstTag(out))
	}
}
