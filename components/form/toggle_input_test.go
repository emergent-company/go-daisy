package form

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderForm renders a form component to its HTML string.
func renderForm(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestToggleInputRenders locks the bare toggle input's output for every
// consumer shape. Attribute order is canonical — type, name, value (when set),
// class, checked (when true), then caller Attrs in alphabetical key order — so
// the literal strings here are the contract.
func TestToggleInputRenders(t *testing.T) {
	cases := []struct {
		name  string
		props ToggleInputProps
		want  string
	}{
		{
			name:  "zero_value",
			props: ToggleInputProps{},
			want:  `<input type="checkbox" name="" class="toggle">`,
		},
		{
			name:  "name_only",
			props: ToggleInputProps{Name: "enabled"},
			want:  `<input type="checkbox" name="enabled" class="toggle">`,
		},
		{
			name: "class_only",
			props: ToggleInputProps{
				Name:    "tool",
				Checked: true,
				Class:   "toggle toggle-sm shrink-0",
			},
			want: `<input type="checkbox" name="tool" class="toggle toggle-sm shrink-0" checked>`,
		},
		{
			name: "value_attribute",
			props: ToggleInputProps{
				Name:    "prop_n",
				Value:   "true",
				Checked: true,
				Class:   "toggle",
			},
			want: `<input type="checkbox" name="prop_n" value="true" class="toggle" checked>`,
		},
		{
			name: "value_omitted_when_empty",
			props: ToggleInputProps{
				Name:  "enabled",
				Class: "toggle toggle-sm",
			},
			want: `<input type="checkbox" name="enabled" class="toggle toggle-sm">`,
		},
		{
			name: "unchecked_omits_checked",
			props: ToggleInputProps{
				Name:  "enabled",
				Class: "toggle toggle-sm",
			},
			want: `<input type="checkbox" name="enabled" class="toggle toggle-sm">`,
		},
		{
			name: "arbitrary_attrs",
			props: ToggleInputProps{
				Name:    "groupEnabled.1",
				Checked: true,
				Class:   "toggle toggle-sm",
				Attrs: templ.Attributes{
					"aria-label":  "Enable X tools",
					"data-testid": "tool-group-enabled-1",
					"onchange":    "this.form.submit()",
				},
			},
			want: `<input type="checkbox" name="groupEnabled.1" class="toggle toggle-sm" checked aria-label="Enable X tools" data-testid="tool-group-enabled-1" onchange="this.form.submit()">`,
		},
		{
			name: "data_hooks_sorted_alphabetically",
			props: ToggleInputProps{
				Name:    "srv",
				Checked: false,
				Class:   "toggle toggle-sm shrink-0",
				Attrs: templ.Attributes{
					"data-mcp-server": "s1",
					"data-mcp-tool":   "t1",
					"aria-label":      "Toggle tool",
				},
			},
			want: `<input type="checkbox" name="srv" class="toggle toggle-sm shrink-0" aria-label="Toggle tool" data-mcp-server="s1" data-mcp-tool="t1">`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderForm(t, ToggleInput(tc.props))
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestToggleInputMatchesConsumerSites proves the bare component reproduces the
// six hand-rolled consumer sites. Because the component renders a single
// canonical attribute order, each site maps to it via the Name/Class/Value/
// Checked/Attrs fields; the assertions here pin the exact element + attribute
// set for every site, including the optional value attribute and arbitrary
// aria-label / data-* / onchange hooks.
func TestToggleInputMatchesConsumerSites(t *testing.T) {
	cases := []struct {
		site  string
		props ToggleInputProps
		want  string
	}{
		{
			// agent.templ ~830: group enable toggle, bulk onchange + data-testid.
			site: "agent_group_toggle",
			props: ToggleInputProps{
				Name:    "groupEnabled.grp1",
				Value:   "on",
				Checked: true,
				Class:   "toggle toggle-sm",
				Attrs: templ.Attributes{
					"aria-label":  "Enable Group tools",
					"data-testid": "tool-group-enabled-grp1",
					"onchange":    "this.closest('[data-tool-group]').querySelectorAll('input[name=tool]').forEach(function(cb){cb.checked=this.checked},this)",
				},
			},
			want: `<input type="checkbox" name="groupEnabled.grp1" value="on" class="toggle toggle-sm" checked aria-label="Enable Group tools" data-testid="tool-group-enabled-grp1" onchange="this.closest(&#39;[data-tool-group]&#39;).querySelectorAll(&#39;input[name=tool]&#39;).forEach(function(cb){cb.checked=this.checked},this)">`,
		},
		{
			// agent.templ ~977: tool row toggle, optional onchange.
			site: "agent_tool_toggle",
			props: ToggleInputProps{
				Name:    "tool",
				Value:   "toolname",
				Checked: true,
				Class:   "toggle toggle-sm shrink-0",
			},
			want: `<input type="checkbox" name="tool" value="toolname" class="toggle toggle-sm shrink-0" checked>`,
		},
		{
			// mcp_servers.templ ~141: server enabled toggle inside caller <label>.
			site: "mcp_server_toggle",
			props: ToggleInputProps{
				Name:    "",
				Checked: true,
				Class:   "toggle toggle-sm",
				Attrs: templ.Attributes{
					"aria-label":              "Enable srv",
					"data-mcp-enabled-toggle": true,
					"data-mcp-server":         "s1",
				},
			},
			want: `<input type="checkbox" name="" class="toggle toggle-sm" checked aria-label="Enable srv" data-mcp-enabled-toggle data-mcp-server="s1">`,
		},
		{
			// mcp_servers.templ ~247: tool row toggle.
			site: "mcp_tool_toggle",
			props: ToggleInputProps{
				Checked: false,
				Class:   "toggle toggle-sm shrink-0",
				Attrs: templ.Attributes{
					"aria-label":           "Toggle tool",
					"data-mcp-tool-toggle": true,
					"data-mcp-server":      "s1",
					"data-mcp-tool":        "t1",
				},
			},
			want: `<input type="checkbox" name="" class="toggle toggle-sm shrink-0" aria-label="Toggle tool" data-mcp-server="s1" data-mcp-tool="t1" data-mcp-tool-toggle>`,
		},
		{
			// objects.templ ~223: form value toggle with sibling hidden=false.
			site: "object_prop_toggle",
			props: ToggleInputProps{
				Name:    "prop_active",
				Value:   "true",
				Checked: true,
				Class:   "toggle",
			},
			want: `<input type="checkbox" name="prop_active" value="true" class="toggle" checked>`,
		},
		{
			// schedules.templ ~128: enabled toggle that submits its form.
			site: "schedule_enabled_toggle",
			props: ToggleInputProps{
				Name:    "enabled",
				Checked: true,
				Class:   "toggle toggle-sm",
				Attrs: templ.Attributes{
					"aria-label": "Toggle Agent",
					"onchange":   "this.form.submit()",
				},
			},
			want: `<input type="checkbox" name="enabled" class="toggle toggle-sm" checked aria-label="Toggle Agent" onchange="this.form.submit()">`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.site, func(t *testing.T) {
			got := renderForm(t, ToggleInput(tc.props))
			if got != tc.want {
				t.Errorf("site %q output:\n got: %q\nwant: %q", tc.site, got, tc.want)
			}
		})
	}
}

// TestToggleCheckboxFormOutputLocked locks the pre-existing rendering of
// Toggle, Checkbox, FormToggle and FormCheckbox so the additive ToggleInput
// (and any future refactor of the internal checkableInput helper) cannot change
// their output byte-for-byte.
func TestToggleCheckboxFormOutputLocked(t *testing.T) {
	cases := []struct {
		name string
		c    templ.Component
		want string
	}{
		{
			name: "toggle_checked",
			c:    Toggle("n", true, "Label"),
			want: `<label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="toggle toggle-primary" checked> <span class="label-text">Label</span></label>`,
		},
		{
			name: "toggle_unchecked",
			c:    Toggle("n", false, "Label"),
			want: `<label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="toggle toggle-primary"> <span class="label-text">Label</span></label>`,
		},
		{
			name: "toggle_empty_label",
			c:    Toggle("n", true, ""),
			want: `<label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="toggle toggle-primary" checked> </label>`,
		},
		{
			name: "checkbox_checked",
			c:    Checkbox("n", true, "Label"),
			want: `<label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="checkbox checkbox-primary" checked> <span class="label-text">Label</span></label>`,
		},
		{
			name: "checkbox_empty_label",
			c:    Checkbox("n", false, ""),
			want: `<label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="checkbox checkbox-primary"> </label>`,
		},
		{
			name: "form_toggle_checked_err",
			c:    FormToggle("n", "Label", true, "err", nil),
			want: `<div class="form-control"><label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="toggle" checked> <span class="label-text">Label</span></label><div class="label pt-1"><span class="label-text-alt text-error">err</span></div></div>`,
		},
		{
			name: "form_toggle_unchecked_noerr",
			c:    FormToggle("n", "Label", false, "", nil),
			want: `<div class="form-control"><label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="toggle"> <span class="label-text">Label</span></label></div>`,
		},
		{
			name: "form_checkbox_checked_err",
			c:    FormCheckbox("n", "Label", true, "err", nil),
			want: `<div class="form-control"><label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="checkbox" checked> <span class="label-text">Label</span></label><div class="label pt-1"><span class="label-text-alt text-error">err</span></div></div>`,
		},
		{
			name: "form_checkbox_unchecked_noerr",
			c:    FormCheckbox("n", "Label", false, "", nil),
			want: `<div class="form-control"><label class="cursor-pointer label justify-start gap-3"><input type="checkbox" name="n" class="checkbox"> <span class="label-text">Label</span></label></div>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderForm(t, tc.c)
			if got != tc.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}
