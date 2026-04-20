package modal

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// ModalWithBoundary wraps Modal with a dev-mode component boundary annotation.
// gallery:token title,size
// gallery:hint title:default(Modal Title)
func ModalWithBoundary(title string, size ModalSize) templ.Component {
	return devmode.ComponentBoundary("Modal", Modal(title, size), map[string]any{
		"title": title,
		"size":  string(size),
	})
}

// FormModalWithBoundary wraps FormModal with a dev-mode component boundary annotation.
func FormModalWithBoundary(props FormModalProps) templ.Component {
	return devmode.ComponentBoundary("FormModal", FormModal(props), map[string]any{
		"id":    props.ID,
		"title": props.Title,
		"size":  string(props.Size),
	})
}

// ConfirmPopupWithBoundary wraps ConfirmPopup with a dev-mode component boundary annotation.
// gallery:token title,message
// gallery:hint title:default(Are you sure?)
// gallery:hint message:default(This action cannot be undone.)
func ConfirmPopupWithBoundary(title string, message string, confirmLabel string, confirmURL string, confirmHXMethod string) templ.Component {
	return devmode.ComponentBoundary("ConfirmPopup", ConfirmPopup(title, message, confirmLabel, confirmURL, confirmHXMethod), map[string]any{
		"title":   title,
		"message": message,
	})
}
