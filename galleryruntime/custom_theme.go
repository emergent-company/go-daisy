package galleryruntime

import "encoding/json"

// CustomThemeStorageKey is the localStorage key holding the user-generated
// theme. The value is a JSON CustomTheme. The gallery theme studio writes it;
// every gallery page (and every preview iframe) reads it on load and injects a
// scoped `[data-theme="custom"]` rule, so a generated theme can be previewed
// across the whole gallery without rebuilding CSS.
const CustomThemeStorageKey = "gallery-custom-theme"

// CustomThemeName is the DaisyUI theme name used for the generated theme.
const CustomThemeName = "custom"

// CustomTheme is the persisted shape of a user-generated DaisyUI theme.
//
// Vars maps DaisyUI custom properties (e.g. "--color-primary") to CSS values
// (DaisyUI themes use OKLCH). The studio exports a `@plugin "daisyui/theme"`
// block from the same data.
type CustomTheme struct {
	Name   string            `json:"name"`
	Scheme string            `json:"scheme"` // "light" or "dark"
	Vars   map[string]string `json:"vars"`
}

// CSS returns the scoped rule for the theme, targeting the given selector.
func (t CustomTheme) CSS(selector string) string {
	if selector == "" {
		selector = `[data-theme="` + CustomThemeName + `"]`
	}
	scheme := t.Scheme
	if scheme == "" {
		scheme = "light"
	}
	out := selector + "{color-scheme:" + scheme + ";"
	for k, v := range t.Vars {
		out += k + ":" + v + ";"
	}
	return out + "}"
}

// JSON serialises the theme for storage.
func (t CustomTheme) JSON() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

// customThemeStyleScript injects the stored custom theme's scoped rule. It
// runs early in the gallery body so the rule exists before components paint.
func customThemeStyleScript() string {
	return customThemeBootScript(false)
}

// customThemeApplyScript activates the custom theme (when the user selected it)
// and runs after the gallery theme switcher, which otherwise resets data-theme.
func customThemeApplyScript() string {
	return customThemeBootScript(true)
}

func customThemeBootScript(apply bool) string {
	applyJS := ""
	if apply {
		applyJS = `
    if (localStorage.getItem('gallery-ui-theme') === '` + CustomThemeName + `') {
      document.documentElement.setAttribute('data-theme', '` + CustomThemeName + `');
    }`
	}
	return `<script>
(function() {
  try {
    var raw = localStorage.getItem('` + CustomThemeStorageKey + `');
    if (!raw) return;
    var theme = JSON.parse(raw);
    if (!theme || !theme.vars) return;
    if (!document.getElementById('gallery-custom-theme')) {
      var css = '[data-theme="` + CustomThemeName + `"]{color-scheme:' + (theme.scheme || 'light') + ';';
      for (var k in theme.vars) { if (Object.prototype.hasOwnProperty.call(theme.vars, k)) css += k + ':' + theme.vars[k] + ';'; }
      css += '}';
      var style = document.createElement('style');
      style.id = 'gallery-custom-theme';
      style.textContent = css;
      document.head.appendChild(style);
    }` + applyJS + `
  } catch (e) {}
})();
</script>`
}
