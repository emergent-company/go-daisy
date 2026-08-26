package handler

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
)

func initials(name string) string {
	if len(name) == 0 { return "?" }
	r := []rune(name)
	if len(r) >= 2 { return string(r[:2]) }
	return string(r[0])
}

func capFirst(s string) string {
	if len(s) == 0 { return s }
	return strings.ToUpper(s[:1]) + s[1:]
}

func productNameOnly(full string) string {
	if idx := strings.LastIndex(full, " #"); idx > 0 {
		return full[:idx]
	}
	return full
}

func productSKU(full string) string {
	if idx := strings.LastIndex(full, " #"); idx > 0 {
		return full[idx+1:]
	}
	return ""
}

func fileTableRowSharedWith(sharedWithType, count string) templ.Component {
	switch sharedWithType {
	case "SharedWithMembers":
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<span class="text-base-content/60">`+count+` members</span>`)
			return err
		})
	case "SharedWithPrivate":
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<span class="text-error flex items-center gap-2"><span class="iconify lucide--shield size-4"></span> Private</span>`)
			return err
		})
	case "SharedWithPublic":
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<span class="text-success flex items-center gap-2"><span class="iconify lucide--globe size-4"></span> Public</span>`)
			return err
		})
	default:
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, fmt.Sprintf(`<span class="text-base-content/60">%s</span>`, count))
			return err
		})
	}
}

func strComponent(s string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	})
}

func boolPtr(b bool) *bool { return &b }
