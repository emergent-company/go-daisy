package shared

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestNonceAttr(t *testing.T) {
	if got := NonceAttr(""); got != "" {
		t.Errorf("empty nonce must render nothing, got %q", got)
	}
	if got := NonceAttr("abc123"); got != ` nonce="abc123"` {
		t.Errorf("nonce attribute = %q", got)
	}
	if got := NonceAttr(`a"b<c&d`); got != ` nonce="a&#34;b&lt;c&amp;d"` {
		t.Errorf("nonce was not HTML-escaped: %q", got)
	}
}

func TestScriptTagNonce(t *testing.T) {
	if got := ScriptTagNonce("", "var x = 1;"); got != `<script>var x = 1;</script>` {
		t.Errorf("empty nonce must produce a plain tag, got %q", got)
	}
	if got := ScriptTagNonce("n1", "var x = 1;"); got != `<script nonce="n1">var x = 1;</script>` {
		t.Errorf("nonced tag = %q", got)
	}
}

// TestScriptTagReadsContextNonce locks the ctx wrapper to the same escaping and
// omit-when-empty behaviour as the explicit-nonce form, so the three consumers
// (theme runtime, document direction, and the gallery slug map) cannot drift.
func TestScriptTagReadsContextNonce(t *testing.T) {
	plain := ScriptTag(context.Background(), "var x = 1;")
	if plain != ScriptTagNonce("", "var x = 1;") {
		t.Errorf("context without nonce diverged from explicit empty nonce: %q", plain)
	}
	if strings.Contains(plain, "nonce=") {
		t.Errorf("empty nonce leaked an attribute: %q", plain)
	}

	nonced := ScriptTag(templ.WithNonce(context.Background(), `n"1`), "var x = 1;")
	want := ScriptTagNonce(`n"1`, "var x = 1;")
	if nonced != want {
		t.Errorf("context nonce = %q, want %q", nonced, want)
	}
	if !strings.Contains(nonced, `nonce="n&#34;1"`) {
		t.Errorf("context nonce was not escaped: %q", nonced)
	}
}

func TestNonceAttrs(t *testing.T) {
	if got := NonceAttrs(context.Background()); got != nil {
		t.Errorf("no nonce should return nil attributes, got %v", got)
	}
	got := NonceAttrs(templ.WithNonce(context.Background(), "n1"))
	if got["nonce"] != "n1" {
		t.Errorf("NonceAttrs = %v", got)
	}
}
