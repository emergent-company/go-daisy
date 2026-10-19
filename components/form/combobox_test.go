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

// TestComboboxTriggerAccessibility locks the control semantics that give the
// visible Label a real accessible name (the trigger is a div, so aria-labelledby
// rather than for): the trigger is focusable, exposes aria-expanded /
// aria-controls / aria-haspopup, and the option container is the listbox it
// controls.
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
		`@keydown.space.prevent="open || toggle()"`,
		`@keydown.arrow-down.prevent="open || toggle()"`,
		`id="pick-listbox"`,
		`role="listbox"`,
		`:aria-selected=`,
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

// TestComboboxClearAllOutsideTrigger locks the structural reason the clear-all
// control is not a trigger child: role=button makes descendants presentational
// in the accessibility tree, so a nested <button> would disappear for AT.
func TestComboboxClearAllOutsideTrigger(t *testing.T) {
	props := ComboboxProps{
		ID:             "pick",
		Name:           "pick",
		EnableClearAll: true,
		Source: ComboboxSource{Static: []ComboboxOption{
			{Value: "a", Label: "A"},
		}},
	}

	trigger := renderForm(t, comboboxTrigger(props))
	if strings.Contains(trigger, "<button") {
		t.Errorf("trigger must not contain interactive descendants:\n%s", trigger)
	}
	if strings.Contains(trigger, "Clear all") {
		t.Errorf("trigger must not contain the clear-all control:\n%s", trigger)
	}

	full := renderForm(t, Combobox(props))
	if !strings.Contains(full, `aria-label="Clear all"`) {
		t.Errorf("clear-all control missing when enabled:\n%s", full)
	}
	if !strings.Contains(full, "pe-16") {
		t.Errorf("trigger must reserve space for the adornment layer:\n%s", full)
	}

	// Disabled comboboxes never render the clear-all control.
	props.Disabled = true
	disabled := renderForm(t, Combobox(props))
	if strings.Contains(disabled, "Clear all") {
		t.Errorf("disabled combobox must not render clear-all:\n%s", disabled)
	}
}

// TestComboboxOptionARIA locks option-level listbox semantics: selected state is
// exposed, disabled options are announced and cannot be clicked.
func TestComboboxOptionARIA(t *testing.T) {
	props := ComboboxProps{
		ID:   "pick",
		Name: "pick",
		Source: ComboboxSource{Static: []ComboboxOption{
			{Value: "a", Label: "A"},
			{Value: "b", Label: "B", Disabled: true},
		}},
	}
	got := renderForm(t, Combobox(props))

	if n := strings.Count(got, ":aria-selected="); n != 2 {
		t.Errorf("want aria-selected on both options, got %d:\n%s", n, got)
	}
	if !strings.Contains(got, `aria-disabled="true"`) {
		t.Errorf("disabled option missing aria-disabled:\n%s", got)
	}
	if !strings.Contains(got, `@click="selectItem(&#39;a&#39;)"`) {
		t.Errorf("enabled option missing click handler:\n%s", got)
	}
	if strings.Contains(got, `@click="selectItem(&#39;b&#39;)"`) {
		t.Errorf("disabled option must not be clickable:\n%s", got)
	}
}
