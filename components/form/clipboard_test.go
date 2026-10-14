package form

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderClipboard renders ClipboardCopy to its HTML string.
func renderClipboard(t *testing.T, items []ClipboardCopyItem) string {
	t.Helper()
	var buf bytes.Buffer
	if err := ClipboardCopy(items).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestClipboardCopyRendersMarkup locks the component's rendered markup for the
// input and mono shapes, including the onclick contract: the target button is
// passed explicitly as `this`, and the copy argument (the input id, or the
// inlined value in Mono mode) is the first argument.
func TestClipboardCopyRendersMarkup(t *testing.T) {
	cases := []struct {
		name       string
		items      []ClipboardCopyItem
		contains   []string
		notContain []string
	}{
		{
			name: "input_variant",
			items: []ClipboardCopyItem{
				{ID: "copy-input-1", Label: "Case ID", Value: "CASE-2026-00142", ButtonLabel: "Copy"},
			},
			contains: []string{
				`<p class="text-xs text-base-content/60 mb-2 font-semibold uppercase">Case ID</p>`,
				`id="copy-input-1" type="text" value="CASE-2026-00142" readonly class="input input-sm flex-1 font-mono text-sm"`,
				`onclick="clipboardCopyInput(&#34;copy-input-1&#34;,this)"`,
				// The copy button must be a plain button, not an implicit submit.
				`<button type="button" class="btn btn-sm btn-outline"`,
				`Copy</button>`,
			},
			notContain: []string{`<code class="bg-base-200`, `type=""`},
		},
		{
			name: "mono_variant",
			items: []ClipboardCopyItem{
				{ID: "copy-input-3", Label: "Inline copy badge", Value: "CASE-2026-00142", Mono: true},
			},
			contains: []string{
				`<code class="bg-base-200 text-base-content px-3 py-1.5 rounded text-sm font-mono flex-1">CASE-2026-00142</code>`,
				`onclick="clipboardCopyText(&#34;CASE-2026-00142&#34;,this)"`,
				// The copy button must be a plain button, not an implicit submit.
				`<button type="button" class="btn btn-sm btn-outline"`,
				// Mono omits the button label, falling back to the default "Copy".
				`Copy</button>`,
			},
			notContain: []string{`type="text"`, `type=""`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderClipboard(t, tc.items)
			for _, want := range tc.contains {
				if !strings.Contains(got, want) {
					t.Errorf("missing %q in output:\n%s", want, got)
				}
			}
			for _, unwanted := range tc.notContain {
				if strings.Contains(got, unwanted) {
					t.Errorf("unexpected %q in output:\n%s", unwanted, got)
				}
			}
		})
	}
}

// TestClipboardCopyButtonExplicitType pins the copy button to a non-submit
// type. The Button component renders type="" by default when Type is unset,
// and an empty type attribute defaults to submit — so inside a <form> a copy
// click would submit the form. This fails if the explicit type regresses.
func TestClipboardCopyButtonExplicitType(t *testing.T) {
	got := renderClipboard(t, []ClipboardCopyItem{{ID: "copy-input-1", Value: "CASE-2026-00142"}})

	if !strings.Contains(got, `type="button"`) {
		t.Errorf("copy button must carry an explicit type=\"button\":\n%s", got)
	}
	for _, unwanted := range []string{`type=""`, `type="submit"`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("copy button renders %s — it would submit an enclosing form:\n%s", unwanted, got)
		}
	}
}

// TestClipboardCopyScriptWiresHandlers asserts the inline script exposes the
// success, failure and fallback paths that a render test can observe. A
// templ-render test cannot execute the JavaScript or drive a non-secure
// context, a clipboard permission denial, or a write contention — those are
// browser behaviours — so this pins the *wiring* instead: the fallback exists,
// the rejection path is attached, the failure state is distinct from success,
// and the target button is a parameter rather than the implicit global `event`.
func TestClipboardCopyScriptWiresHandlers(t *testing.T) {
	got := renderClipboard(t, []ClipboardCopyItem{
		{ID: "copy-input-1", Value: "CASE-2026-00142"},
		{ID: "copy-input-3", Value: "CASE-2026-00142", Mono: true},
	})

	// Fallback for non-secure origins where navigator.clipboard is undefined.
	for _, want := range []string{
		`document.execCommand('copy')`,
		`document.createElement('textarea')`,
		`navigator.clipboard && navigator.clipboard.writeText`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("fallback not wired: missing %q", want)
		}
	}

	// Rejection handling: a failed write reaches a distinct failure state.
	for _, want := range []string{
		`.catch(function()`,
		`clipboardFail(btn)`,
		`'✗ Copy failed'`,
		`btn.classList.add('btn-error')`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("failure path not wired: missing %q", want)
		}
	}

	// Success state is distinct from failure.
	for _, want := range []string{
		`'✓ Copied!'`,
		`'✓'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("success state not wired: missing %q", want)
		}
	}

	// The handlers take the target button as an explicit parameter and never
	// read the implicit global `event`.
	for _, want := range []string{
		`window.clipboardCopyInput = function(inputId, btn)`,
		`window.clipboardCopyText = function(text, btn)`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("handler signature changed: missing %q", want)
		}
	}
	if strings.Contains(got, "event.currentTarget") {
		t.Errorf("handler still relies on the implicit global `event`:\n%s", got)
	}
}

// TestClipboardCopyOnclickPassesThis pins the exact onclick the Button renders,
// so a future regression that drops the explicit `this` (and falls back to the
// implicit global event) fails loudly.
func TestClipboardCopyOnclickPassesThis(t *testing.T) {
	got := renderClipboard(t, []ClipboardCopyItem{{ID: "copy-input-1", Value: "CASE-2026-00142"}})
	if !strings.Contains(got, `clipboardCopyInput(&#34;copy-input-1&#34;,this)`) {
		t.Errorf("onclick does not pass `this` explicitly:\n%s", got)
	}
}

// TestClipboardCopyScriptGuardsMissingButton pins the `!btn` guards in the
// emitted script. `this` is passed explicitly, so a stray call with no button
// (or a missing input element in the input variant) would otherwise dereference
// a nil/undefined button inside clipboardWrite. Reverting the guard must fail
// this test; a render test cannot exercise the JS, so presence is the contract.
func TestClipboardCopyScriptGuardsMissingButton(t *testing.T) {
	got := renderClipboard(t, []ClipboardCopyItem{
		{ID: "copy-input-1", Value: "CASE-2026-00142"},
		{ID: "copy-input-3", Value: "CASE-2026-00142", Mono: true},
	})

	// Input variant: bail out if the element or the button is missing.
	if !strings.Contains(got, `if (!el || !btn) return;`) {
		t.Errorf("input handler is missing the `!el || !btn` guard:\n%s", got)
	}
	// Mono variant: bail out if the button is missing.
	if !strings.Contains(got, `if (!btn) return;`) {
		t.Errorf("text handler is missing the `!btn` guard:\n%s", got)
	}
}
