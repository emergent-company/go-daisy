package form

import (
	"strings"
	"testing"
)

// TestComboboxPopoverSearchEnableSearchLocks the search input rendered by the
// popover's search slot across the EnableSearch axis.
//
// It is a regression test for a live bug on master: templ compiled master's
// bare `return` statement in comboboxPopoverSearch as a literal WriteString
// rather than control flow, so with EnableSearch: false the component emitted
// the text `return ` into the HTML *and* fell through to render the search
// input anyway — the EnableSearch prop was effectively ignored. The nested
// `if props.EnableSearch { ... }` form fixes both symptoms: no literal `return`
// text and no `<input` element at all when search is disabled.
func TestComboboxPopoverSearchEnableSearch(t *testing.T) {
	cases := []struct {
		name    string
		props   ComboboxProps
		want    string
		exclude []string
	}{
		{
			name:    "search_disabled_renders_nothing",
			props:   ComboboxProps{EnableSearch: false},
			want:    "",
			exclude: []string{"return", "<input"},
		},
		{
			name:  "search_enabled_lazy_endpoint_set",
			props: ComboboxProps{EnableSearch: true, Source: ComboboxSource{LazyEndpoint: "/lazy"}},
			want:  `<input type="text" name="q" class="input input-sm w-full mb-1" placeholder="Search..." x-model="search" hx-get="/lazy" hx-target="next .combobox-options-inner" hx-trigger="keyup changed delay:200ms" hx-include="closest .combobox-popover">`,
		},
		{
			name:  "search_enabled_lazy_endpoint_unset",
			props: ComboboxProps{EnableSearch: true},
			want:  `<input type="text" class="input input-sm w-full mb-1" placeholder="Search..." x-model="search">`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderForm(t, comboboxPopoverSearch(tc.props))
			if tc.want != "" && got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
			if tc.want == "" && got != "" {
				t.Errorf("expected empty output, got: %q", got)
			}
			for _, sub := range tc.exclude {
				if strings.Contains(got, sub) {
					t.Errorf("output must not contain %q, got: %q", sub, got)
				}
			}
		})
	}
}
