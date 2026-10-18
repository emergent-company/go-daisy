package ui

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/shared"
)

// ChatBubbleLegacy renders a ChatBubble from the legacy flat parameter list.
//
// Deprecated: use ChatBubble(ChatBubbleProps{...}) instead. This shim is a plain
// Go function, so .templ block-children call sites must move to the props form.
func ChatBubbleLegacy(sent bool, author, timestamp, avatarSrc string, botIcon bool, bubbleClass string, showActions bool, attrs templ.Attributes, children ...templ.Component) templ.Component {
	props := ChatBubbleProps{
		Sent:        sent,
		Author:      author,
		Timestamp:   timestamp,
		AvatarSrc:   avatarSrc,
		BotIcon:     botIcon,
		BubbleClass: bubbleClass,
		ShowActions: showActions,
		Attrs:       attrs,
	}

	// Compose is nil-safe for zero children but panics on a nil child, so drop
	// any nil children the caller may have passed.
	filtered := make([]templ.Component, 0, len(children))
	for _, c := range children {
		if c != nil {
			filtered = append(filtered, c)
		}
	}
	content := shared.Compose(filtered...)
	return shared.RenderInto(ChatBubble(props), content)
}
