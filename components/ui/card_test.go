package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderCardRawProps renders CardRaw with the given props (and no children) to
// its HTML string.
func renderCardRawProps(t *testing.T, props CardRawProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := CardRaw(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestCardRawZeroValueOutput pins the full zero-value render. The surface stays
// exactly bg-base-100 (byte-identical to the pre-fix hard-coded value) and the
// default border border-line is emitted.
func TestCardRawZeroValueOutput(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{})

	want := `<div class="card bg-base-100 border border-line"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("zero-value CardRaw output:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawSurfaceReplacesDefault pins that a non-zero Surface REPLACES the
// bg-base-100 default instead of stacking a second bg-* class on top of it.
func TestCardRawSurfaceReplacesDefault(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Surface: CardSurfaceBase200})

	if !strings.Contains(got, `class="card bg-base-200 border border-line"`) {
		t.Errorf("Surface base-200 must replace bg-base-100:\n%s", got)
	}
	if strings.Contains(got, "bg-base-100") {
		t.Errorf("Surface base-200 must not leave bg-base-100 behind:\n%s", got)
	}
}

// TestCardRawBorderNone pins that BorderNone removes the default outline.
func TestCardRawBorderNone(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Border: CardBorderNone})

	want := `<div class="card bg-base-100"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("BorderNone output:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawBorderStrong pins that BorderStrong uses the emphasised step.
func TestCardRawBorderStrong(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Border: CardBorderStrong})

	if !strings.Contains(got, `class="card bg-base-100 border border-line-strong"`) {
		t.Errorf("BorderStrong must render border-line-strong:\n%s", got)
	}
}

// TestCardRawAsElement pins that As selects the rendered element.
func TestCardRawAsElement(t *testing.T) {
	cases := []struct {
		name    string
		as      CardElement
		openTag string
	}{
		{"div", CardElementDiv, `<div class="card bg-base-100 border border-line">`},
		{"a", CardElementA, `<a class="card bg-base-100 border border-line">`},
		{"button", CardElementButton, `<button type="button" class="card bg-base-100 border border-line">`},
		{"section", CardElementSection, `<section class="card bg-base-100 border border-line">`},
		{"article", CardElementArticle, `<article class="card bg-base-100 border border-line">`},
		{"li", CardElementLi, `<li class="card bg-base-100 border border-line">`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderCardRawProps(t, CardRawProps{As: tc.as})
			if !strings.HasPrefix(got, tc.openTag) {
				t.Errorf("As %q must render %s:\n%s", tc.as, tc.openTag, got)
			}
		})
	}
}

// TestCardRawUnknownAsPanics pins that an unsupported as element fails loudly
// instead of silently degrading to a <div>.
func TestCardRawUnknownAsPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown as element")
		}
	}()
	renderCardRawProps(t, CardRawProps{As: CardElement("span")})
}

// TestCardRawUnknownBorderPanics pins that an unsupported border fails loudly.
func TestCardRawUnknownBorderPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown border")
		}
	}()
	renderCardRawProps(t, CardRawProps{Border: CardBorder("dashed")})
}
