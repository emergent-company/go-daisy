package ui

import "strconv"

// cardSurfaceClass maps a CardSurface to its background class. The zero value
// returns the default "bg-base-100"; any other surface REPLACES the default
// class rather than stacking on top of it, so a caller never has to fight the
// CSS cascade to change the surface.
func cardSurfaceClass(s CardSurface) string {
	if s == "" {
		return "bg-base-100"
	}
	return string(s)
}

// cardBorderClass maps a CardBorder to its outline class string. The zero value
// returns the default "border border-line"; CardBorderStrong uses the emphasised
// line-strong step and CardBorderNone removes the outline entirely.
func cardBorderClass(b CardBorder) string {
	switch b {
	case CardBorderLine:
		return "border border-line"
	case CardBorderStrong:
		return "border border-line-strong"
	case CardBorderNone:
		return ""
	default:
		panic("ui.CardRaw: unsupported border " + strconv.Quote(string(b)))
	}
}

// validateCardElement panics on an unsupported CardElement value so CardRaw
// fails loudly instead of silently degrading to a <div> (the pre-fix behaviour
// for any unrecognised "as" value, e.g. a typo).
func validateCardElement(e CardElement) {
	switch e {
	case CardElementDiv, CardElementA, CardElementButton, CardElementSection, CardElementArticle, CardElementLi:
		return
	default:
		panic("ui.CardRaw: unsupported as element " + strconv.Quote(string(e)) +
			" (supported: div, a, button, section, article, li)")
	}
}
