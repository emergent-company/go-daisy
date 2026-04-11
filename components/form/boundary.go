package form

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// TextInputWithBoundary wraps TextInput with a dev-mode component boundary annotation.
// gallery:token label,required
// gallery:hint label:default(Full Name)
func TextInputWithBoundary(name string, label string, value string, errMsg string, required bool) templ.Component {
	return devmode.ComponentBoundary("TextInput", map[string]any{
		"name":     name,
		"label":    label,
		"required": required,
	}, TextInput(name, label, value, errMsg, required))
}

// TextareaInputWithBoundary wraps TextareaInput with a dev-mode component boundary annotation.
// gallery:token label,rows,required
// gallery:hint rows:range(2,10,1)
// gallery:hint label:default(Description)
func TextareaInputWithBoundary(name string, label string, value string, errMsg string, rows int, required bool) templ.Component {
	return devmode.ComponentBoundary("TextareaInput", map[string]any{
		"name":     name,
		"label":    label,
		"rows":     rows,
		"required": required,
	}, TextareaInput(name, label, value, errMsg, rows, required))
}

// CheckboxInputWithBoundary wraps CheckboxInput with a dev-mode component boundary annotation.
// gallery:token label,checked
// gallery:hint label:default(Accept terms and conditions)
func CheckboxInputWithBoundary(name string, label string, checked bool, errMsg string) templ.Component {
	return devmode.ComponentBoundary("CheckboxInput", map[string]any{
		"name":    name,
		"label":   label,
		"checked": checked,
	}, CheckboxInput(name, label, checked, errMsg))
}

// SelectInputWithBoundary wraps SelectInput with a dev-mode component boundary annotation.
// gallery:token label,required
// gallery:hint label:default(Country)
func SelectInputWithBoundary(name string, label string, selected string, options [][2]string, errMsg string, required bool) templ.Component {
	return devmode.ComponentBoundary("SelectInput", map[string]any{
		"name":     name,
		"label":    label,
		"selected": selected,
		"required": required,
	}, SelectInput(name, label, selected, options, errMsg, required))
}

// SearchInputWithBoundary wraps SearchInput with a dev-mode component boundary annotation.
func SearchInputWithBoundary(name string, value string, placeholder string, hxTarget string, hxGet string) templ.Component {
	return devmode.ComponentBoundary("SearchInput", map[string]any{
		"name":        name,
		"placeholder": placeholder,
		"hxTarget":    hxTarget,
	}, SearchInput(name, value, placeholder, hxTarget, hxGet))
}

// FormFieldWithBoundary wraps FormField with a dev-mode component boundary annotation.
func FormFieldWithBoundary(props FormFieldProps) templ.Component {
	return devmode.ComponentBoundary("FormField", props, FormField(props))
}

// RangeInputWithBoundary wraps RangeInput with a dev-mode component boundary annotation.
// gallery:token label,value,color
// gallery:hint value:range(0,100,1)
// gallery:hint label:default(Volume)
// gallery:hint color:default(range-primary)
// gallery:hint value:default(50)
func RangeInputWithBoundary(name string, label string, value int, min int, max int, step int, color string) templ.Component {
	return devmode.ComponentBoundary("RangeInput", map[string]any{
		"name":  name,
		"label": label,
		"min":   min,
		"max":   max,
		"step":  step,
	}, RangeInput(name, label, value, min, max, step, color))
}
