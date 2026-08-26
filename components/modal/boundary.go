package modal

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// ModalWithBoundary wraps Modal with a dev-mode component boundary annotation.
// gallery:token title,size,driver
// gallery:hint title:default(Modal Title)
func ModalWithBoundary(title string, size ModalSize, driver ModalDriver) templ.Component {
	props := ModalProps{Title: title, Size: size, Driver: driver}
	return devmode.ComponentBoundary("Modal", Modal(props), map[string]any{
		"title":  title,
		"size":   string(size),
		"driver": string(driver),
	})
}

// FormModalWithBoundary wraps FormModal with a dev-mode component boundary annotation.
// gallery:token id,title,size,driver
// gallery:hint id:default(form-modal-1)
// gallery:hint title:default(Edit Item)
// gallery:hint size:default(md)
func FormModalWithBoundary(props FormModalProps, driver ModalDriver) templ.Component {
	props.Driver = driver
	return devmode.ComponentBoundary("FormModal", FormModal(props), map[string]any{
		"id":     props.ID,
		"title":  props.Title,
		"size":   string(props.Size),
		"driver": string(props.Driver),
	})
}

// ConfirmPopupWithBoundary wraps ConfirmPopup with a dev-mode component boundary annotation.
// gallery:token title,message
// gallery:hint title:default(Are you sure?)
// gallery:hint message:default(This action cannot be undone.)
func ConfirmPopupWithBoundary(title string, message string, confirmLabel string, confirmURL string, confirmHXMethod string) templ.Component {
	return devmode.ComponentBoundary("ConfirmPopup", ConfirmPopup(title, message, confirmLabel, confirmURL, confirmHXMethod, nil), map[string]any{
		"title":   title,
		"message": message,
	})
}

// LoaderModalWithBoundary wraps LoaderModal with a dev-mode component boundary annotation.
func LoaderModalWithBoundary() templ.Component {
	return devmode.ComponentBoundary("LoaderModal", LoaderModal(), nil)
}

// OpenModalButtonWithBoundary wraps OpenModalButton with a dev-mode component boundary annotation.
// gallery:token modalID,label,driver
// gallery:hint modalID:default(my-modal)
// gallery:hint label:default(Open Modal)
func OpenModalButtonWithBoundary(modalID string, label string, driver ModalDriver) templ.Component {
	return devmode.ComponentBoundary("OpenModalButton", OpenModalButton(modalID, label, nil, driver), map[string]any{
		"modalID": modalID,
		"label":   label,
		"driver":  string(driver),
	})
}

// DeleteButtonWithBoundary wraps DeleteButton with a dev-mode component boundary annotation.
// gallery:token url,confirm,label
// gallery:hint url:default(#)
// gallery:hint confirm:default(Are you sure?)
// gallery:hint label:default(Delete)
func DeleteButtonWithBoundary(url string, confirm string, label string) templ.Component {
	return devmode.ComponentBoundary("DeleteButton", DeleteButton(url, confirm, "", "", label, nil), map[string]any{
		"url":     url,
		"confirm": confirm,
		"label":   label,
	})
}
