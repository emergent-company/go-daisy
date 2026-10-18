package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestBannerBindsOpen locks Banner's x-data contract. Banner seeds
// alpine.Toggle(true), whose only property is `open`, but master's markup bound
// x-show="show" and dismissed with `show = false`. Alpine evaluated an
// undeclared `show` on both variants: the banner logged
// "ReferenceError: show is not defined" at init and the dismiss button did
// nothing. Same defect class as form/Palette's `picked`.
func TestBannerBindsOpen(t *testing.T) {
	cases := []struct {
		name  string
		props BannerProps
	}{
		{name: "alert_bar", props: BannerProps{Description: "Limited time offer"}},
		{name: "cookie_banner", props: BannerProps{CookieBanner: true, Description: "We use cookies"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Banner(tc.props).Render(context.Background(), &buf); err != nil {
				t.Fatalf("render: %v", err)
			}
			html := buf.String()

			if strings.Contains(html, `x-show="show"`) {
				t.Errorf("Banner binds x-show=\"show\" but x-data only declares `open`:\n%s", html)
			}
			if !strings.Contains(html, `x-show="open"`) {
				t.Errorf("Banner lost its x-show visibility binding:\n%s", html)
			}
			if strings.Contains(html, "show = false") {
				t.Errorf("Banner dismiss handler assigns undeclared `show`:\n%s", html)
			}
		})
	}
}
