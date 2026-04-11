package galleryruntime

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/emergent-company/go-daisy/staticfs"
)

// Options configures the gallery server.
type Options struct {
	// Title is shown in the browser tab and gallery header. Defaults to "Component Gallery".
	Title string
	// Logo is an optional templ component rendered in the sidebar header instead of the
	// plain title text. Useful for branded galleries that want a custom logo mark.
	// When nil, the Title string is rendered as a plain text link.
	Logo templ.Component
	// Components is the full registry of components to display.
	Components []GalleryComponent
	// Port is the TCP port to listen on. Defaults to 11000.
	Port int
	// StorePath is the path for the SQLite feedback database.
	// Pass an empty string to disable feedback persistence.
	StorePath string
	// GitHubCfg enables GitHub issue export. Nil disables it.
	GitHubCfg *GitHubConfig
	// ExtraStaticPrefixes lists additional URL prefixes under which the
	// embedded go-daisy static assets (CSS, JS) should also be served.
	// Use this when your component shell templates reference a custom static
	// path (e.g. "/dashboard/static/") instead of the default "/static/".
	// Example: []string{"/dashboard/static/"}
	ExtraStaticPrefixes []string
	// DevMode enables component boundary annotations in the gallery preview.
	// When true, templ components are wrapped in data-component/data-props
	// markers and the hover overlay + component tree panel are injected into
	// the preview iframe. Has no effect in production; safe to leave false.
	DevMode bool
}

// Serve starts the gallery HTTP server with the provided options.
// It blocks until the server exits or returns an error.
func Serve(opts Options) error {
	if opts.Title == "" {
		opts.Title = "Component Gallery"
	}
	if opts.Port == 0 {
		opts.Port = 11000
	}

	// Optional SQLite store.
	var store *Store
	if opts.StorePath != "" {
		var err error
		store, err = Open(opts.StorePath)
		if err != nil {
			log.Printf("warning: could not open gallery SQLite store at %s: %v (feedback disabled)", opts.StorePath, err)
		} else {
			log.Printf("gallery SQLite store opened at %s", opts.StorePath)
			defer store.Close()
		}
	}

	// Optional GitHub client.
	var gh *GitHubClient
	if opts.GitHubCfg != nil {
		var err error
		gh, err = NewGitHubClient(*opts.GitHubCfg)
		if err != nil {
			log.Printf("warning: could not create GitHub client: %v — GitHub integration disabled", err)
		} else {
			log.Printf("GitHub App integration enabled for repo %s", opts.GitHubCfg.Repo)
		}
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static assets (CSS, JS) embedded in go-daisy's staticfs.
	e.GET("/static/*", echo.WrapHandler(staticfs.Handler("/static/")))

	// Mount extra static prefixes declared by the caller (e.g. "/dashboard/static/").
	for _, prefix := range opts.ExtraStaticPrefixes {
		p := prefix
		// Ensure the prefix ends with / for the handler strip to work correctly.
		if len(p) > 0 && p[len(p)-1] != '/' {
			p += "/"
		}
		e.GET(p+"*", echo.WrapHandler(staticfs.Handler(p)))
	}

	// Redirect root to /gallery.
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/gallery")
	})

	// Build the full list of static prefixes: always include "/static/", plus any extras.
	staticPrefixes := append([]string{"/static/"}, opts.ExtraStaticPrefixes...)

	h := newGalleryHandler(opts.Title, opts.Logo, opts.Components, store, gh, staticPrefixes, opts.DevMode)
	h.register(e)

	addr := fmt.Sprintf("0.0.0.0:%d", opts.Port)
	log.Printf("gallery listening on %s", addr)
	return e.Start(addr)
}
