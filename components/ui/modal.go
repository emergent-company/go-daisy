package ui

import "github.com/a-h/templ"

// modalBoxAttrs returns the optional id attribute for a modal-box element
// (empty string omits it, keeping the markup clean).
func modalBoxAttrs(boxID string) templ.Attributes {
	if boxID == "" {
		return nil
	}
	return templ.Attributes{"id": boxID}
}
