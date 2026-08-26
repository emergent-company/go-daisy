package handler

import (
	"context"
	"embed"
	"io"
	"strings"

	"github.com/a-h/templ"
)

//go:embed nxhtml/*.html
var nxHTMLFiles embed.FS

func nxHTMLComponent(page string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		filename := page + ".html"
		data, err := nxHTMLFiles.ReadFile("nxhtml/" + filename)
		if err != nil {
			return err
		}
		html := string(data)
		// Extract #layout-content
		start := strings.Index(html, `id="layout-content"`)
		if start < 0 {
			_, err = io.WriteString(w, "<!-- no layout-content -->")
			return err
		}
		start += len(`id="layout-content">`)
		content := html[start:]
		// Find the closing </div> of the layout-content container
		depth := 1
		end := 0
		for i := 0; i < len(content); i++ {
			if strings.HasPrefix(content[i:], "<div") {
				depth++
			} else if strings.HasPrefix(content[i:], "</div>") {
				depth--
				if depth == 0 {
					end = i
					break
				}
			}
		}
		if end > 0 {
			content = content[:end]
		}
		_, err = io.WriteString(w, strings.TrimSpace(content))
		return err
	})
}
