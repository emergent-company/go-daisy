package ui

import (
	"context"
	"strconv"

	"github.com/emergent-company/go-daisy/devmode"
)

// cardSurfaceClass maps a CardSurface to its background class. The zero value
// returns the default "bg-base-100"; any other surface REPLACES the default
// class rather than stacking on top of it.
//
// Surface is deliberately an OPEN class passthrough (see CardSurface): a caller
// may legitimately want an arbitrary background utility, so the value is emitted
// verbatim with no validation or sentinel mapping.
func cardSurfaceClass(s CardSurface) string {
	if s == "" {
		return "bg-base-100"
	}
	return string(s)
}

// cardBorderClass maps a CardBorder to its outline class string. The zero value
// returns the default house edge "edge-line"; CardBorderStrong uses the
// emphasised line-strong step and CardBorderNone removes the outline entirely.
//
// The edge-* classes come from components/css/custom.css, not from Tailwind's
// `border` utility. `border` hard-codes `border-width: 1px`, which silently
// dropped the width daisyUI's `.card-border` took from the theme's `--border`
// token — so cards stopped thickening for `prefers-contrast: more` users
// (issue #79). `edge-line` keeps the same colours as `border-line` /
// `border-line-strong` and takes its width from `--border`, so it is 1px
// normally and 2px under high contrast. Do not "simplify" it back to
// `border border-line`: that is the regression this replaced.
//
// Border is a CLOSED enum (unlike Surface): the outline is structural and must
// map to one of the house edge tokens. An unknown value panics in devmode (the
// loud-failure guarantee) and falls back to the default outline outside devmode,
// so a runtime-constructed value can never crash a production request.
func cardBorderClass(b CardBorder, ctx context.Context) string {
	switch b {
	case CardBorderLine:
		return "edge-line"
	case CardBorderStrong:
		return "edge-line-strong"
	case CardBorderNone:
		return ""
	default:
		if devmode.IsDevMode(ctx) {
			panic("ui.CardRaw: unsupported border " + strconv.Quote(string(b)))
		}
		return "edge-line"
	}
}

// cardElement validates and returns the effective HTML element CardRaw renders.
// Element is a CLOSED enum (unlike Surface): it selects a structural HTML tag
// from a finite documented set. An unknown value panics in devmode so a typo is
// caught loudly, and falls back to the default <div> outside devmode so a
// runtime-constructed value can never crash a request.
func cardElement(e CardElement, ctx context.Context) CardElement {
	switch e {
	case CardElementDiv, CardElementA, CardElementButton, CardElementSection, CardElementArticle, CardElementLi:
		return e
	default:
		if devmode.IsDevMode(ctx) {
			panic("ui.CardRaw: unsupported as element " + strconv.Quote(string(e)) +
				" (supported: div, a, button, section, article, li)")
		}
		return CardElementDiv
	}
}
