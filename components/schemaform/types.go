// Package schemaform generates form fields from a JSON Schema object subset
// and renders them as go-daisy form controls.
package schemaform

// Kind identifies the generated control type.
type Kind string

const (
	KindString  Kind = "string"
	KindNumber  Kind = "number"
	KindInteger Kind = "integer"
	KindBoolean Kind = "boolean"
	KindEnum    Kind = "enum"
	KindArray   Kind = "array"
	KindObject  Kind = "object"
	KindUnknown Kind = "unknown"
)

// AllowMode controls field visibility in the form.
type AllowMode string

const (
	AllowManaged  AllowMode = "managed"  // visible, disabled, marked managed
	AllowDisabled AllowMode = "disabled" // hidden entirely
)

// Field is a generated form control.
type Field struct {
	Path       string  // dotted config path, e.g. "resources.cpu"
	Label      string  // display label
	HelperText string  // supporting copy
	Kind       Kind    // control type
	Managed    bool    // true when allow-listed as managed
	Default    string  // serialized default value
	Value      string  // current value (submitted wins over default)
	Options    []string // enum values
	Children   []Field // nested object fields
}
