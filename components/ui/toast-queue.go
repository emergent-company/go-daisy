package ui

// ToastQueuePosition selects the fixed corner/edge the toast queue anchors to.
// The zero value (empty string) is the default top-right behaviour.
const (
	// ToastQueueTopRight is the zero-value default: fixed to the top-right with
	// newest toasts on top (flex-col-reverse).
	ToastQueueTopRight string = ""
	// ToastQueueBottomCenter anchors to the bottom-center on small screens and
	// the bottom-right on sm+, with newest toasts on the bottom (flex-col).
	ToastQueueBottomCenter string = "bottom-center"
)

// toastQueuePositionClass returns the container's fixed-position and flex
// classes for a position. Anything other than ToastQueueBottomCenter falls back
// to the default top-right layout (identical to ui.ToastQueue).
func toastQueuePositionClass(position string) string {
	if position == ToastQueueBottomCenter {
		return "fixed bottom-4 left-1/2 z-70 flex w-80 max-w-[calc(100vw-2rem)] -translate-x-1/2 flex-col gap-2 sm:left-auto sm:right-4 sm:translate-x-0"
	}
	return "fixed top-4 right-4 z-70 flex flex-col-reverse gap-2 w-80 max-w-[calc(100vw-2rem)]"
}

// toastQueueSlideClass returns the enter/leave translate class: top-right slides
// down from above (-translate-y-2); bottom-center slides up from below
// (translate-y-2). Mirrors the translate direction of the existing ToastQueue.
func toastQueueSlideClass(position string) string {
	if position == ToastQueueBottomCenter {
		return "translate-y-2"
	}
	return "-translate-y-2"
}
