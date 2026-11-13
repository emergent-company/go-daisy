package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderAccountMenu renders AccountMenu to its HTML string.
func renderAccountMenu(t *testing.T, props AccountMenuProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := AccountMenu(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// signedInProps is the full-featured menu: one active account with a picture,
// one cached account without, and every action wired.
func signedInProps() AccountMenuProps {
	return AccountMenuProps{
		SignedIn: true,
		Active: Account{
			Sub:     "sub-ada",
			Name:    "Ada Lovelace",
			Email:   "ada@example.test",
			Picture: "https://images.example.test/ada.png",
		},
		Others: []Account{
			{Sub: "sub-grace", Name: "Grace Hopper", Email: "grace@example.test"},
		},
		SignInPath:  "/auth/login",
		SwitchPath:  "/auth/switch",
		AddPath:     "/auth/add",
		ProfilePath: "/profile",
		LogoutPath:  "/auth/logout",
	}
}

// TestAccountMenuSignedInRendersEveryRow pins the signed-in markup contract:
// one avatar trigger, the active account, the switch rows and each action, at
// the paths the props name. It is the contract both apps' tests assert against,
// so a rename here cannot silently break one of them.
func TestAccountMenuSignedInRendersEveryRow(t *testing.T) {
	html := renderAccountMenu(t, signedInProps())

	for _, want := range []string{
		`data-testid="account-menu"`,
		`data-testid="account-menu-trigger"`,
		`aria-label="Account"`,
		`aria-haspopup="menu"`,
		`>Signed in as</li>`,
		`data-testid="account-menu-active"`,
		`Ada Lovelace`,
		`ada@example.test`,
		`>Switch account</li>`,
		`data-testid="account-menu-switch"`,
		`action="/auth/switch"`,
		`name="sub" value="sub-grace"`,
		`Grace Hopper`,
		`href="/auth/add"`,
		`Add another account`,
		`href="/profile"`,
		`My Profile`,
		`action="/auth/logout"`,
		`Log out`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q:\n%s", want, html)
		}
	}
	if strings.Contains(html, `data-testid="account-menu-signin"`) {
		t.Errorf("signed-in menu must not render the sign-in link:\n%s", html)
	}
}

// TestAccountMenuAvatarUsesPictureAndInitialsFallback pins the avatar chain:
// the active account's picture becomes an <img>, and an account without one
// falls back to initials rather than an empty circle.
func TestAccountMenuAvatarUsesPictureAndInitialsFallback(t *testing.T) {
	html := renderAccountMenu(t, signedInProps())

	if !strings.Contains(html, `src="https://images.example.test/ada.png"`) {
		t.Errorf("active account picture not rendered as an image:\n%s", html)
	}
	if !strings.Contains(html, ">GH</div>") {
		t.Errorf("cached account without a picture must render initials GH:\n%s", html)
	}
}

// TestAccountMenuOmitsRowsForEmptyPaths pins the "no dead controls" rule: an
// app that has no switch, add, profile or logout route gets a menu without
// them, and no switch section when there is no other account.
func TestAccountMenuOmitsRowsForEmptyPaths(t *testing.T) {
	html := renderAccountMenu(t, AccountMenuProps{
		SignedIn: true,
		Active:   Account{Name: "Ada"},
	})

	for _, unwanted := range []string{
		`/auth/switch`,
		`/auth/add`,
		`/profile`,
		`/auth/logout`,
		`Switch account`,
		`data-testid="account-menu-switch"`,
	} {
		if strings.Contains(html, unwanted) {
			t.Errorf("expected no %q with the path unset:\n%s", unwanted, html)
		}
	}
	if !strings.Contains(html, `data-testid="account-menu-active"`) {
		t.Errorf("active account row must render even with no actions:\n%s", html)
	}
	if strings.Contains(html, `<li aria-hidden="true"></li>`) {
		t.Errorf("no divider may render when nothing follows it:\n%s", html)
	}
}

// TestAccountMenuSkipsAccountsWithoutSub pins the switch form contract: a row
// with no sub would post an empty value, which can only switch to nothing. With
// no switchable account left, the section heading goes too — a heading over an
// empty list reads as a broken menu.
func TestAccountMenuSkipsAccountsWithoutSub(t *testing.T) {
	props := signedInProps()
	props.Others = []Account{{Name: "No Id", Email: "nobody@example.test"}}

	html := renderAccountMenu(t, props)

	if strings.Contains(html, "No Id") {
		t.Errorf("account without a sub must be skipped:\n%s", html)
	}
	if strings.Contains(html, "Switch account") {
		t.Errorf("heading must not render over an empty switch list:\n%s", html)
	}
	if strings.Contains(html, `data-testid="account-menu-switch"`) {
		t.Errorf("no switch form may render:\n%s", html)
	}
}

// TestAccountMenuNoSubKeepsOtherSections pins the split: skipping one unusable
// account must not take the usable ones (or the actions) down with it.
func TestAccountMenuNoSubKeepsOtherSections(t *testing.T) {
	props := signedInProps()
	props.Others = []Account{
		{Name: "No Id", Email: "nobody@example.test"},
		{Sub: "sub-grace", Name: "Grace Hopper"},
	}

	html := renderAccountMenu(t, props)

	if !strings.Contains(html, "Switch account") {
		t.Errorf("heading must render when a usable account remains:\n%s", html)
	}
	if strings.Contains(html, "No Id") {
		t.Errorf("unusable account must be skipped:\n%s", html)
	}
	if !strings.Contains(html, `name="sub" value="sub-grace"`) {
		t.Errorf("usable account must still render:\n%s", html)
	}
}

// TestAccountMenuSignedOutRendersSignInLink pins the signed-out state: the
// catalog's public chrome offers exactly one sign-in control and no menu.
func TestAccountMenuSignedOutRendersSignInLink(t *testing.T) {
	html := renderAccountMenu(t, AccountMenuProps{SignInPath: "/auth/login"})

	if !strings.Contains(html, `href="/auth/login"`) ||
		!strings.Contains(html, `data-testid="account-menu-signin"`) {
		t.Errorf("signed-out menu must render the sign-in link:\n%s", html)
	}
	if strings.Contains(html, "dropdown") {
		t.Errorf("signed-out menu must not render the dropdown:\n%s", html)
	}
}

// TestAccountMenuSignedOutWithoutPathRendersNothing pins the product's
// dev/API-key mode: no identity surface at all, not an empty button.
func TestAccountMenuSignedOutWithoutPathRendersNothing(t *testing.T) {
	html := strings.TrimSpace(renderAccountMenu(t, AccountMenuProps{}))
	if html != "" {
		t.Errorf("expected no markup, got:\n%s", html)
	}
}

// TestAccountMenuCalleeAttrsWin pins the attribute merge: a caller's Attrs land
// on the root dropdown and may override the default testid, because a page that
// hosts two menus needs to tell them apart.
func TestAccountMenuCalleeAttrsWin(t *testing.T) {
	props := signedInProps()
	props.Attrs = map[string]any{"data-testid": "header-account"}

	html := renderAccountMenu(t, props)

	if !strings.Contains(html, `data-testid="header-account"`) {
		t.Errorf("caller attrs must reach the root:\n%s", html)
	}
	if strings.Count(html, `data-testid="header-account"`) != 1 {
		t.Errorf("caller testid must replace, not duplicate, the default:\n%s", html)
	}
}

// TestAccountMenuProfileAttrs pins the htmx passthrough: the product swaps the
// profile page into #main-content, and a shared component that dropped those
// attributes would turn that into a full page load (or a 404 in the container).
func TestAccountMenuProfileAttrs(t *testing.T) {
	props := signedInProps()
	props.ProfileAttrs = map[string]any{
		"hx-get":      "/profile",
		"hx-target":   "#main-content",
		"hx-swap":     "innerHTML",
		"hx-push-url": "true",
	}

	html := renderAccountMenu(t, props)

	for _, want := range []string{
		`hx-get="/profile"`,
		`hx-target="#main-content"`,
		`hx-swap="innerHTML"`,
		`hx-push-url="true"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q:\n%s", want, html)
		}
	}
}

// TestAccountMenuEscapesAccountFields pins escaping: name, email and picture are
// identity data from an IdP, and a component that emitted them raw would turn a
// crafted display name into markup on every page. The picture is also
// scheme-checked, because a plain string attribute does not get templ's URL
// sanitisation.
func TestAccountMenuEscapesAccountFields(t *testing.T) {
	html := renderAccountMenu(t, AccountMenuProps{
		SignedIn: true,
		Active: Account{
			Name:    `<script>alert(1)</script>`,
			Email:   `"x@example.test">`,
			Picture: `javascript:alert(1)`,
		},
	})

	if strings.Contains(html, "<script>alert(1)</script>") {
		t.Errorf("name must be escaped:\n%s", html)
	}
	if strings.Contains(html, `src="javascript:alert(1)"`) {
		t.Errorf("unsafe picture scheme must be neutralised:\n%s", html)
	}
	if !strings.Contains(html, string(templ.FailedSanitizationURL)) {
		t.Errorf("unsafe picture must fall back to the failure URL:\n%s", html)
	}
}

// TestAccountMenuClassIsAppended pins the standard props contract: Class is
// appended to the root's own classes, not swapped for them (the root must stay a
// dropdown for the menu to open at all).
func TestAccountMenuClassIsAppended(t *testing.T) {
	props := signedInProps()
	props.Class = "max-w-xs"

	html := renderAccountMenu(t, props)

	if !strings.Contains(html, "dropdown dropdown-end max-w-xs") {
		t.Errorf("Class must be appended to the root classes:\n%s", html)
	}
}

// TestAccountMenuSingleMenuRoot pins the dev-mode annotation boundary: the root
// is this component's own element, so it carries exactly one menu root and the
// runtime can wire it. Delegating to ui.Dropdown would stamp a second
// data-component and the overlay would attribute the menu to Dropdown.
func TestAccountMenuSingleMenuRoot(t *testing.T) {
	html := renderAccountMenu(t, signedInProps())

	if n := strings.Count(html, "data-gd-menu-root"); n != 1 {
		t.Errorf("want exactly one menu root, got %d:\n%s", n, html)
	}
	if n := strings.Count(html, `class="dropdown dropdown-end`); n != 1 {
		t.Errorf("want exactly one dropdown root, got %d:\n%s", n, html)
	}
}

// TestAccountMenuSignInClassIsAppended pins the same contract for the signed-out
// branch, which is a different root element.
func TestAccountMenuSignInClassIsAppended(t *testing.T) {
	html := renderAccountMenu(t, AccountMenuProps{SignInPath: "/auth/login", Class: "w-full"})

	if !strings.Contains(html, "btn btn-ghost btn-sm w-full") {
		t.Errorf("Class must be appended to the sign-in link classes:\n%s", html)
	}
}

// TestAccountMenuLogoutOnlyRendersOneDivider pins the divider arithmetic: with
// no add/profile row, the logout row is the only thing after the switch section,
// so it must not get a divider of its own on top of the section divider.
func TestAccountMenuLogoutOnlyRendersOneDivider(t *testing.T) {
	html := renderAccountMenu(t, AccountMenuProps{
		SignedIn:   true,
		Active:     Account{Name: "Ada Lovelace"},
		LogoutPath: "/auth/logout",
	})

	if n := strings.Count(html, `<li aria-hidden="true"></li>`); n != 1 {
		t.Errorf("want exactly one divider for a logout-only menu, got %d:\n%s", n, html)
	}
	if !strings.Contains(html, `data-testid="account-menu-logout"`) {
		t.Errorf("logout row must still render:\n%s", html)
	}
}
