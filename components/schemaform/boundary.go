package schemaform

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// FieldsWithBoundary wraps Fields with a dev-mode component boundary annotation.
// gallery:token namePrefix
// gallery:hint namePrefix:default(values)
func FieldsWithBoundary(props FieldsProps) templ.Component {
	return devmode.ComponentBoundary("SchemaFields", Fields(props), map[string]any{
		"namePrefix":  props.NamePrefix,
		"fieldCount":  len(props.Fields),
	})
}
