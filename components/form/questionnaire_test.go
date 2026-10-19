package form

import (
	"strings"
	"testing"
)

// questionnaireFixture returns a three-step props struct covering the contract
// shapes: a required single-choice step, a multi-choice step with a disabled
// option, and a plain single-choice step. Letter shortcuts are enabled.
func questionnaireFixture() QuestionnaireProps {
	return QuestionnaireProps{
		ID:        "q",
		Method:    "post",
		Shortcuts: QuestionnaireShortcutsLetters,
		Steps: []QuestionnaireStep{
			{
				Name:     "role",
				Title:    "What's your role?",
				Required: true,
				Choices: []QuestionnaireOption{
					{Value: "developer", Label: "Developer"},
					{Value: "designer", Label: "Designer"},
				},
			},
			{
				Name:     "interests",
				Title:    "Which technologies interest you?",
				Multiple: true,
				Choices: []QuestionnaireOption{
					{Value: "htmx", Label: "HTMX"},
					{Value: "alpine", Label: "Alpine.js", Disabled: true},
				},
			},
			{
				Name:  "experience",
				Title: "Rate your experience",
				Choices: []QuestionnaireOption{
					{Value: "beginner", Label: "Beginner"},
				},
			},
		},
	}
}

// TestQuestionnaireContract locks the progressive-enhancement contract of the
// wizard: one native form with the data-gd-questionnaire root hook, per-step
// fieldsets carrying data-name/data-required/data-multiple, radios vs
// checkboxes, id/for pairing, shortcut kbd badges, a hidden alert error slot,
// server-hidden nav buttons with a visible Submit, the "Question X of N" label,
// and a single runtime script tag.
func TestQuestionnaireContract(t *testing.T) {
	got := renderForm(t, Questionnaire(questionnaireFixture()))

	mustContain(t, got, []string{
		// Root form hook + native method.
		`<form id="q"`,
		`method="post"`,
		`data-gd-questionnaire`,
		`data-gd-questionnaire-mode="letters"`,
		// Progress.
		`Question 1 of 3`,
		`<progress`,
		`max="3"`,
		// Step fieldsets + data hooks.
		`id="q-role"`,
		`data-name="role"`,
		`data-required="true"`,
		`id="q-interests"`,
		`data-multiple="true"`,
		// Radios vs checkboxes.
		`<input type="radio" id="q-role-developer" name="role" value="developer" class="radio radio-sm">`,
		`<input type="checkbox" id="q-interests-htmx" name="interests" value="htmx" class="checkbox checkbox-sm">`,
		// Disabled choice label + input.
		`id="q-interests-alpine"`,
		`disabled`,
		`opacity-50 pointer-events-none`,
		// id/for pairing.
		`<label for="q-role-developer"`,
		// Shortcut badge.
		`<kbd data-gd-questionnaire-shortcut`,
		`>a</kbd>`,
		// Error slot.
		`role="alert"`,
		`hidden`,
		// Nav buttons hidden server-side.
		`data-gd-questionnaire-previous`,
		`data-gd-questionnaire-skip`,
		`data-gd-questionnaire-next`,
		// Submit visible.
		`type="submit"`,
		// Runtime script emitted exactly once.
		`go-daisy-questionnaire.js`,
	})

	if n := strings.Count(got, "go-daisy-questionnaire.js"); n != 1 {
		t.Errorf("runtime script emitted %d times, want exactly 1", n)
	}

	// No-JS contract: the Submit button is the only usable server-side control,
	// so it must be visible (no `hidden` attribute) and carry the default label.
	if !strings.Contains(got, `<button type="submit" data-gd-questionnaire-submit`) || !strings.Contains(got, `>Submit</button>`) {
		t.Errorf("Submit button not rendered visible with default label")
	}
	submitStart := strings.Index(got, `<button type="submit"`)
	submitEnd := strings.Index(got[submitStart:], `</button>`) + submitStart
	if strings.Contains(got[submitStart:submitEnd], "hidden") {
		t.Errorf("Submit button must not be hidden server-side")
	}
}

// TestQuestionnaireMethodDispatch pins the native method attribute mapping and
// the htmx dispatch for verbs a browser cannot express natively.
func TestQuestionnaireMethodDispatch(t *testing.T) {
	cases := []struct {
		method   string
		wantAttr string // e.g. `method="post"`; always present
		wantHX   string // e.g. `hx-put=`; "" when no htmx verb is emitted
	}{
		{"", `method="post"`, ""},
		{"get", `method="get"`, ""},
		{"put", `method="post"`, "hx-put="},
		{"patch", `method="post"`, "hx-patch="},
		{"delete", `method="post"`, "hx-delete="},
	}

	for _, tc := range cases {
		t.Run(tc.method, func(t *testing.T) {
			props := questionnaireFixture()
			props.Method = tc.method
			got := renderForm(t, Questionnaire(props))
			if !strings.Contains(got, tc.wantAttr) {
				t.Errorf("want form to contain %q, got %q", tc.wantAttr, got)
			}
			if strings.Contains(got, "data-method=") {
				t.Errorf("data-method is dead weight and must not be emitted, got %q", got)
			}
			if tc.wantHX == "" {
				if strings.Contains(got, "hx-put=") || strings.Contains(got, "hx-patch=") || strings.Contains(got, "hx-delete=") {
					t.Errorf("want no htmx verb for method %q, got %q", tc.method, got)
				}
			} else if !strings.Contains(got, tc.wantHX) {
				t.Errorf("want %q for method %q, got %q", tc.wantHX, tc.method, got)
			}
		})
	}
}

// TestQuestionnaireSubmitLabelDefault pins the default Submit label and custom
// label override.
func TestQuestionnaireSubmitLabelDefault(t *testing.T) {
	props := questionnaireFixture()
	props.SubmitLabel = ""
	got := renderForm(t, Questionnaire(props))
	if !strings.Contains(got, ">Submit</button>") {
		t.Errorf("default submit label: got %q", got)
	}

	props.SubmitLabel = "Finish"
	got = renderForm(t, Questionnaire(props))
	if !strings.Contains(got, ">Finish</button>") {
		t.Errorf("custom submit label: got %q", got)
	}
}

// mustContain fails t if got lacks any of the wanted substrings.
func mustContain(t *testing.T, got string, wants []string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(got, w) {
			t.Errorf("output missing %q", w)
		}
	}
}
