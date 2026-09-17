package ui

import "strings"

// accentStyle returns the inline style declarations that tint an element with an
// arbitrary CSS color: the color's text/border plus a translucent background,
// mirroring the house "tone/10 bg + tone/15 border" tile look.
//
// Colors are arbitrary CSS values (hex in practice), so they are applied via a
// style attribute — never a class. Returns "" when no color is declared (the
// caller then keeps its default tone classes). An invalid color value is
// dropped by the browser declaration by declaration, leaving the default look.
//
// Used by IconTile (props.Color) and Badge (props.Color).
func accentStyle(color string) string {
	c := strings.TrimSpace(color)
	if c == "" {
		return ""
	}
	return "color:" + c +
		";background-color:color-mix(in oklch," + c + " 10%,transparent)" +
		";border-color:color-mix(in oklch," + c + " 15%,transparent)"
}
