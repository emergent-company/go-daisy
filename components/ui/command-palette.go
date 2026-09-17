package ui

import "strings"

// commandPaletteMobileCSS returns the CSS for the FullScreenMobile command
// palette behaviour, with selectors derived from the palette's ID (the same
// ID+"-toggle" / ID+"-input" / ID+"-results" suffixes used by the markup).
// The placeholder "ID" is replaced by the component's ID.
//
// The ID must be a valid CSS identifier (kebab-case in practice, e.g.
// "spotlight"); it is interpolated verbatim into the selectors.
func commandPaletteMobileCSS(id string) string {
	css := `@media (max-width: 767px) {
  /* the toggle input and .modal are siblings, so the checked state drives the
     layout. Fill the viewport edge-to-edge; desktop keeps centered max-w-lg. */
  #ID-toggle:checked + .modal { padding: 0; }
  #ID-toggle:checked + .modal .modal-box {
    width: 100%;
    max-width: none;
    height: 100vh;
    height: 100dvh;
    max-height: 100dvh;
    border-radius: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  /* input row pinned at the top, clear of notches/status bars */
  #ID-toggle:checked + .modal .modal-box > div:first-child {
    flex-shrink: 0;
    padding-top: calc(0.75rem + env(safe-area-inset-top));
  }
  /* results fill the remaining height and scroll inside the sheet */
  #ID-toggle:checked + .modal #ID-results {
    flex: 1 1 auto;
    min-height: 0;
    max-height: none;
    padding-bottom: calc(0.5rem + env(safe-area-inset-bottom));
  }
  /* close affordance: replace the "Esc" label text with a tap-sized X.
     font-size:0 keeps the real "Esc" text in the DOM for screen readers. */
  #ID-toggle:checked + .modal .modal-box label[for="ID-toggle"] {
    width: 2rem;
    height: 2rem;
    min-height: 2rem;
    padding: 0;
    font-size: 0;
  }
  #ID-toggle:checked + .modal .modal-box label[for="ID-toggle"]::after {
    content: "✕";
    font-size: 0.9375rem;
    line-height: 1;
  }
}`
	return strings.ReplaceAll(css, "ID", id)
}
