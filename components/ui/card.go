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
// returns the default "border border-line"; CardBorderStrong uses the emphasised
// line-strong step and CardBorderNone removes the outline entirely.
//
// Border is a CLOSED enum (unlike Surface): the outline is structural and must
// map to one of the house line tokens. An unknown value panics in devmode (the
// loud-failure guarantee) and falls back to the default outline outside devmode,
// so a runtime-constructed value can never crash a production request.
func cardBorderClass(b CardBorder, ctx context.Context) string {
	switch b {
	case CardBorderLine:
		return "border border-line"
	case CardBorderStrong:
		return "border border-line-strong"
	case CardBorderNone:
		return ""
	default:
		if devmode.IsDevMode(ctx) {
			panic("ui.CardRaw: unsupported border " + strconv.Quote(string(b)))
		}
		return "border border-line"
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
