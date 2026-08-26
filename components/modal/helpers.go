package modal

// modalDefaults returns the resolved submit text, cancel text, and size for modal form components.
func modalDefaults(props FormModalProps) (submitText, cancelText string, size ModalSize) {
	submitText = props.SubmitText
	if submitText == "" {
		submitText = "Save"
	}
	cancelText = props.CancelText
	if cancelText == "" {
		cancelText = "Cancel"
	}
	size = props.Size
	if size == "" {
		size = ModalSM
	}
	return
}
