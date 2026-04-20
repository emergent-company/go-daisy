package form

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// TextInputWithBoundary wraps TextInput with a dev-mode component boundary annotation.
// gallery:token label,required
// gallery:hint label:default(Full Name)
func TextInputWithBoundary(name string, label string, value string, errMsg string, required bool) templ.Component {
	return devmode.ComponentBoundary("TextInput", TextInput(name, label, value, errMsg, required), map[string]any{
		"name":     name,
		"label":    label,
		"required": required,
	})
}

// TextareaInputWithBoundary wraps TextareaInput with a dev-mode component boundary annotation.
// gallery:token label,rows,required
// gallery:hint rows:range(2,10,1)
// gallery:hint label:default(Description)
func TextareaInputWithBoundary(name string, label string, value string, errMsg string, rows int, required bool) templ.Component {
	return devmode.ComponentBoundary("TextareaInput", TextareaInput(name, label, value, errMsg, rows, required), map[string]any{
		"name":     name,
		"label":    label,
		"rows":     rows,
		"required": required,
	})
}

// CheckboxInputWithBoundary wraps CheckboxInput with a dev-mode component boundary annotation.
// gallery:token label,checked
// gallery:hint label:default(Accept terms and conditions)
func CheckboxInputWithBoundary(name string, label string, checked bool, errMsg string) templ.Component {
	return devmode.ComponentBoundary("CheckboxInput", CheckboxInput(name, label, checked, errMsg), map[string]any{
		"name":    name,
		"label":   label,
		"checked": checked,
	})
}

// SelectInputWithBoundary wraps SelectInput with a dev-mode component boundary annotation.
// gallery:token label,required
// gallery:hint label:default(Country)
func SelectInputWithBoundary(name string, label string, selected string, options [][2]string, errMsg string, required bool) templ.Component {
	return devmode.ComponentBoundary("SelectInput", SelectInput(name, label, selected, options, errMsg, required), map[string]any{
		"name":     name,
		"label":    label,
		"selected": selected,
		"required": required,
	})
}

// SearchInputWithBoundary wraps SearchInput with a dev-mode component boundary annotation.
func SearchInputWithBoundary(name string, value string, placeholder string, hxTarget string, hxGet string) templ.Component {
	return devmode.ComponentBoundary("SearchInput", SearchInput(name, value, placeholder, hxTarget, hxGet), map[string]any{
		"name":        name,
		"placeholder": placeholder,
		"hxTarget":    hxTarget,
	})
}

// FormFieldWithBoundary wraps FormField with a dev-mode component boundary annotation.
func FormFieldWithBoundary(props FormFieldProps) templ.Component {
	return devmode.ComponentBoundary("FormField", FormField(props), props)
}

// RangeInputWithBoundary wraps RangeInput with a dev-mode component boundary annotation.
// gallery:token label,value,color
// gallery:hint value:range(0,100,1)
// gallery:hint label:default(Volume)
// gallery:hint color:default(range-primary)
// gallery:hint value:default(50)
func RangeInputWithBoundary(name string, label string, value int, min int, max int, step int, color string) templ.Component {
	return devmode.ComponentBoundary("RangeInput", RangeInput(name, label, value, min, max, step, color), map[string]any{
		"name":  name,
		"label": label,
		"min":   min,
		"max":   max,
		"step":  step,
	})
}

// RadioGroupWithBoundary wraps RadioGroup with a dev-mode component boundary annotation.
// gallery:token color
// gallery:hint color:default(radio-primary)
func RadioGroupWithBoundary(name string, selected string, options [][2]string, color string) templ.Component {
	return devmode.ComponentBoundary("RadioGroup", RadioGroup(name, selected, options, color), map[string]any{
		"name":     name,
		"selected": selected,
		"color":    color,
	})
}

// RatingWithBoundary wraps Rating with a dev-mode component boundary annotation.
// gallery:token value,max,shape,color,size
// gallery:hint value:range(1,10,1)
// gallery:hint value:default(3)
// gallery:hint max:range(1,10,1)
// gallery:hint max:default(5)
// gallery:hint color:default(bg-orange-400)
func RatingWithBoundary(name string, value int, max int, shape RatingShape, color string, size string) templ.Component {
	return devmode.ComponentBoundary("Rating", Rating(name, value, max, shape, color, size), map[string]any{
		"name":  name,
		"value": value,
		"max":   max,
		"shape": string(shape),
		"color": color,
		"size":  size,
	})
}

// FileInputWithBoundary wraps FileInput with a dev-mode component boundary annotation.
// gallery:token label,accept
// gallery:hint label:default(Upload file)
func FileInputWithBoundary(name string, label string, accept string) templ.Component {
	return devmode.ComponentBoundary("FileInput", FileInput(name, label, accept, ""), map[string]any{
		"name":   name,
		"label":  label,
		"accept": accept,
	})
}

// CheckboxWithBoundary wraps Checkbox with a dev-mode component boundary annotation.
// gallery:token label,checked
// gallery:hint label:default(Accept terms and conditions)
func CheckboxWithBoundary(name string, checked bool, label string) templ.Component {
	return devmode.ComponentBoundary("Checkbox", Checkbox(name, checked, label), map[string]any{
		"name":    name,
		"checked": checked,
		"label":   label,
	})
}

// ToggleWithBoundary wraps Toggle with a dev-mode component boundary annotation.
// gallery:token label,checked
// gallery:hint label:default(Enable notifications)
func ToggleWithBoundary(name string, checked bool, label string) templ.Component {
	return devmode.ComponentBoundary("Toggle", Toggle(name, checked, label), map[string]any{
		"name":    name,
		"checked": checked,
		"label":   label,
	})
}

// PromptBarWithBoundary wraps PromptBar with a dev-mode component boundary annotation.
func PromptBarWithBoundary(props PromptBarProps) templ.Component {
	return devmode.ComponentBoundary("PromptBar", PromptBar(props), props)
}

// PromptBarActionWithBoundary wraps PromptBarAction with a dev-mode component boundary annotation.
func PromptBarActionWithBoundary(placeholder string, actions []PromptBarActionItem) templ.Component {
	return devmode.ComponentBoundary("PromptBarAction", PromptBarAction(placeholder, actions), map[string]any{
		"placeholder": placeholder,
		"actionCount": len(actions),
	})
}

// InputSpinnerWithBoundary wraps InputSpinner with a dev-mode component boundary annotation.
func InputSpinnerWithBoundary(id string, value, min, max int, hasMinMax bool, btnClass, inputClass string) templ.Component {
	return devmode.ComponentBoundary("InputSpinner", InputSpinner(id, value, min, max, hasMinMax, btnClass, inputClass), map[string]any{
		"id":         id,
		"value":      value,
		"min":        min,
		"max":        max,
		"hasMinMax":  hasMinMax,
		"btnClass":   btnClass,
		"inputClass": inputClass,
	})
}

// WizardStepperWithBoundary wraps WizardStepper with a dev-mode component boundary annotation.
func WizardStepperWithBoundary(id string, steps []WizardStep, panels []WizardStepPanel) templ.Component {
	return devmode.ComponentBoundary("WizardStepper", WizardStepper(id, steps, panels), map[string]any{
		"id":         id,
		"stepCount":  len(steps),
		"panelCount": len(panels),
	})
}
