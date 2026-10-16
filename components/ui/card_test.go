package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// renderCardRawProps renders CardRaw with the given props (and no children) to
// its HTML string using a plain (non-devmode) context.
func renderCardRawProps(t *testing.T, props CardRawProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := CardRaw(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// renderCardRawPropsDevMode renders CardRaw with the given props using a devmode
// context, so the devmode loud-failure paths are active.
func renderCardRawPropsDevMode(t *testing.T, props CardRawProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := CardRaw(props).Render(devmode.WithDevMode(context.Background()), &buf); err != nil {
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

// TestCardRawAsElement pins the full element render (opening and closing tag)
// for every documented CardElement value.
func TestCardRawAsElement(t *testing.T) {
	cases := []struct {
		name string
		as   CardElement
		want string
	}{
		{"div", CardElementDiv, `<div class="card bg-base-100 border border-line"><div class="card-body"></div></div>`},
		{"a", CardElementA, `<a class="card bg-base-100 border border-line"><div class="card-body"></div></a>`},
		{"button", CardElementButton, `<button type="button" class="card bg-base-100 border border-line"><div class="card-body"></div></button>`},
		{"section", CardElementSection, `<section class="card bg-base-100 border border-line"><div class="card-body"></div></section>`},
		{"article", CardElementArticle, `<article class="card bg-base-100 border border-line"><div class="card-body"></div></article>`},
		{"li", CardElementLi, `<li class="card bg-base-100 border border-line"><div class="card-body"></div></li>`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderCardRawProps(t, CardRawProps{As: tc.as})
			if got != tc.want {
				t.Errorf("As %q:\n got: %q\nwant: %q", tc.as, got, tc.want)
			}
		})
	}
}

// TestCardRawUnknownAsPanicsInDevMode pins that an unsupported as element fails
// loudly in devmode instead of silently degrading.
func TestCardRawUnknownAsPanicsInDevMode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic in devmode for unknown as element")
		}
	}()
	renderCardRawPropsDevMode(t, CardRawProps{As: CardElement("span")})
}

// TestCardRawUnknownAsFallsBackToDiv pins that outside devmode an unsupported as
// element degrades to the default <div> rather than crashing the request.
func TestCardRawUnknownAsFallsBackToDiv(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{As: CardElement("span")})

	want := `<div class="card bg-base-100 border border-line"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("unknown as must fall back to <div> outside devmode:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawUnknownBorderPanicsInDevMode pins that an unsupported border fails
// loudly in devmode.
func TestCardRawUnknownBorderPanicsInDevMode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic in devmode for unknown border")
		}
	}()
	renderCardRawPropsDevMode(t, CardRawProps{Border: CardBorder("dashed")})
}

// TestCardRawUnknownBorderFallsBackToDefault pins that outside devmode an
// unsupported border degrades to the default outline rather than crashing.
func TestCardRawUnknownBorderFallsBackToDefault(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Border: CardBorder("dashed")})

	want := `<div class="card bg-base-100 border border-line"><div class="card-body"></div></div>`
	if got != want {
		t.Errorf("unknown border must fall back to border border-line outside devmode:\n got: %q\nwant: %q", got, want)
	}
}

// TestCardRawAttrsSpreadLast pins that caller Attrs land on the card root after
// the component's own class attribute.
func TestCardRawAttrsSpreadLast(t *testing.T) {
	got := renderCardRawProps(t, CardRawProps{Attrs: templ.Attributes{"data-testid": "foo", "id": "bar"}})

	for _, attr := range []string{`data-testid="foo"`, `id="bar"`} {
		if !strings.Contains(got, attr) {
			t.Errorf("Attrs %s must be spread onto the card root:\n%s", attr, got)
		}
	}
	classIdx := strings.Index(got, `class="card`)
	attrIdx := strings.Index(got, `data-testid="foo"`)
	if classIdx < 0 || attrIdx < 0 || classIdx > attrIdx {
		t.Errorf("Attrs must be spread after the class attribute:\n%s", got)
	}
}
