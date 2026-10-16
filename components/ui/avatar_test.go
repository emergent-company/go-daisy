package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderAvatarFull renders AvatarFull to its HTML string.
func renderAvatarFull(t *testing.T, props AvatarProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := AvatarFull(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestAvatarInitialsOverride pins the Initials field: it must take precedence
// over initials computed from Name. This is the migration path for existing
// initials helpers (e.g. first-two-runes) when hand-rolled avatars move to
// AvatarFull.
func TestAvatarInitialsOverride(t *testing.T) {
	got := renderAvatarFull(t, AvatarProps{Name: "Ada Lovelace", Initials: "XX"})
	if !strings.Contains(got, ">XX</div>") {
		t.Errorf("Initials must override Name-derived initials:\n%s", got)
	}
	if strings.Contains(got, ">AL</div>") {
		t.Errorf("Name-derived initials leaked when Initials was set:\n%s", got)
	}
}

// TestAvatarNameDerivesInitials pins the default initials path: when Initials
// is empty, AvatarFull computes up to two initials from Name.
func TestAvatarNameDerivesInitials(t *testing.T) {
	got := renderAvatarFull(t, AvatarProps{Name: "Ada Lovelace"})
	if !strings.Contains(got, ">AL</div>") {
		t.Errorf("expected Name-derived initials AL:\n%s", got)
	}
}

// TestAvatarBadgeClassReplacesTone pins BadgeClass: when set it must replace the
// built-in Tone colors, not append after them.
func TestAvatarBadgeClassReplacesTone(t *testing.T) {
	got := renderAvatarFull(t, AvatarProps{Name: "Ada", BadgeClass: "bg-secondary text-secondary-content font-semibold"})
	if !strings.Contains(got, "bg-secondary text-secondary-content") {
		t.Errorf("BadgeClass classes not emitted:\n%s", got)
	}
	if strings.Contains(got, "bg-primary") {
		t.Errorf("BadgeClass must replace Tone colors, not render bg-primary alongside:\n%s", got)
	}
	// Structural flex centering must survive the color replacement.
	if !strings.Contains(got, "flex items-center justify-center") {
		t.Errorf("flex centering lost when BadgeClass replaced Tone:\n%s", got)
	}
}

// TestAvatarInnerClassAppended pins InnerClass: it is appended to the initials
// badge element for extra treatment (e.g. ring/border).
func TestAvatarInnerClassAppended(t *testing.T) {
	got := renderAvatarFull(t, AvatarProps{Name: "Ada", InnerClass: "font-bold ring-2 ring-base-100"})
	if !strings.Contains(got, "font-bold ring-2 ring-base-100") {
		t.Errorf("InnerClass not appended to badge:\n%s", got)
	}
}

// TestAvatarNumericSizes pins the explicit numeric size classes added for
// migrating the app's hand-rolled avatar diameters.
func TestAvatarNumericSizes(t *testing.T) {
	for _, tc := range []struct {
		size AvatarSize
		want string
	}{
		{AvatarSize6, "size-6"},
		{AvatarSize7, "size-7"},
		{AvatarSize12, "size-12"},
		{AvatarSize14, "size-14"},
		{AvatarSize16, "size-16"},
	} {
		got := renderAvatarFull(t, AvatarProps{Name: "Ada", Size: tc.size})
		if !strings.Contains(got, tc.want) {
			t.Errorf("AvatarSize %q must emit %q:\n%s", tc.size, tc.want, got)
		}
	}
}

// renderPersonCell renders PersonCell to its HTML string.
func renderPersonCell(t *testing.T, p PersonCellProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := PersonCell(p).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestPersonCellCompactUsesAvatarFull pins that compact mode (PersonChip) renders
// its initials via AvatarFull rather than the old hand-rolled
// "size-6 rounded-full ... bg-primary text-primary-content" div.
func TestPersonCellCompactUsesAvatarFull(t *testing.T) {
	got := renderPersonCell(t, PersonCellProps{Name: "Jane Doe", Compact: true})
	if !strings.Contains(got, `class="avatar"`) {
		t.Errorf("compact PersonCell must render an AvatarFull wrapper:\n%s", got)
	}
	if !strings.Contains(got, "size-6") {
		t.Errorf("compact avatar must keep size-6:\n%s", got)
	}
	if !strings.Contains(got, ">JD</div>") {
		t.Errorf("compact avatar initials missing:\n%s", got)
	}
	// Regression guards for the bare-`return` footgun: on master, PersonCell's
	// compact branch was `if p.Compact { @PersonCellCompact(p); return }`. Templ
	// compiled that `return` as a literal " return" write and fell through to
	// render the non-compact markup too. The if/else form fixes both symptoms.
	if strings.Contains(got, "return") {
		t.Errorf("compact PersonCell must not leak a bare-return literal:\n%s", got)
	}
	if strings.Contains(got, `class="flex items-center gap-3"`) {
		t.Errorf("compact PersonCell must not fall through to non-compact markup:\n%s", got)
	}
	// The hand-rolled avatars carried shrink-0 (they are flex items in an
	// inline-flex chip and a flex card); the AvatarFull migration must keep it.
	if !strings.Contains(got, "shrink-0") {
		t.Errorf("compact PersonCell avatars must keep shrink-0:\n%s", got)
	}
}

// TestPersonCellCompactCustomColors pins that the compact avatar still honors
// the custom AvatarColor/TextColor (now passed through as BadgeClass).
func TestPersonCellCompactCustomColors(t *testing.T) {
	got := renderPersonCell(t, PersonCellProps{
		Name:        "Bob Smith",
		Compact:     true,
		AvatarColor: "bg-secondary",
		TextColor:   "text-secondary-content",
	})
	if !strings.Contains(got, "bg-secondary text-secondary-content") {
		t.Errorf("compact avatar must preserve custom AvatarColor/TextColor:\n%s", got)
	}
}

// TestPersonCellNonCompact pins the non-compact layout is unaffected by the
// bare-return removal: name and subtitle still render.
func TestPersonCellNonCompact(t *testing.T) {
	got := renderPersonCell(t, PersonCellProps{Name: "Alice", Subtitle: "alice@example.com"})
	if !strings.Contains(got, ">Alice</div>") {
		t.Errorf("non-compact PersonCell must render Name:\n%s", got)
	}
	if !strings.Contains(got, "alice@example.com") {
		t.Errorf("non-compact PersonCell must render Subtitle:\n%s", got)
	}
	// Defense-in-depth: the non-compact branch must also stay free of a
	// bare-return literal (the same templ footgun that bit the compact branch).
	if strings.Contains(got, "return") {
		t.Errorf("non-compact PersonCell must not contain a bare-return literal:\n%s", got)
	}
}
