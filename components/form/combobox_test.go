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

// TestComboboxTriggerAccessibility locks the control semantics that make the
// visible Label's for association real: the trigger is focusable, exposes
// aria-expanded / aria-controls / aria-haspopup, and the option container is the
// listbox it controls.
func TestComboboxTriggerAccessibility(t *testing.T) {
	props := ComboboxProps{
		ID:    "pick",
		Name:  "pick",
		Label: "Pick",
		Mode:  ComboboxSingle,
		Source: ComboboxSource{Static: []ComboboxOption{
			{Value: "a", Label: "A"},
		}},
	}
	got := renderForm(t, Combobox(props))

	for _, want := range []string{
		`id="pick-label"`,
		`aria-labelledby="pick-label"`,
		`id="pick"`,
		`role="button"`,
		`aria-haspopup="listbox"`,
		`aria-controls="pick-listbox"`,
		`:aria-expanded=`,
		`tabindex="0"`,
		`@keydown.enter.prevent="open || toggle()"`,
		`@keydown.space.prevent="toggle()"`,
		`@keydown.arrow-down.prevent="open || toggle()"`,
		`id="pick-listbox"`,
		`role="listbox"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("combobox missing %q in:\n%s", want, got)
		}
	}

	if strings.Contains(got, `aria-disabled="true"`) {
		t.Errorf("enabled combobox must not be aria-disabled:\n%s", got)
	}

	// Disabled: removed from the tab order and announced as disabled.
	props.Disabled = true
	disabled := renderForm(t, Combobox(props))
	for _, want := range []string{`tabindex="-1"`, `aria-disabled="true"`} {
		if !strings.Contains(disabled, want) {
			t.Errorf("disabled combobox missing %q in:\n%s", want, disabled)
		}
	}

	// Multiple mode: the listbox announces multi-selectability.
	props.Disabled = false
	props.Mode = ComboboxMultiple
	multiple := renderForm(t, Combobox(props))
	if !strings.Contains(multiple, `aria-multiselectable="true"`) {
		t.Errorf("multiple combobox missing aria-multiselectable:\n%s", multiple)
	}

	// ID without a label: no dangling aria-labelledby.
	props.ID = "pick"
	props.Label = ""
	unlabelled := renderForm(t, Combobox(props))
	if strings.Contains(unlabelled, "aria-labelledby") {
		t.Errorf("combobox without a label must not emit aria-labelledby:\n%s", unlabelled)
	}

	// No ID: no orphan aria-controls/listbox id is emitted.
	props.ID = ""
	props.Label = ""
	bare := renderForm(t, Combobox(props))
	if strings.Contains(bare, "aria-controls") || strings.Contains(bare, "-listbox") || strings.Contains(bare, "aria-labelledby") {
		t.Errorf("combobox without ID must not emit association attributes:\n%s", bare)
	}
}
