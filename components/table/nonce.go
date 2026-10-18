package table

import (
	"context"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/shared"
)

// nonceAttrs returns the render-context CSP nonce as an attribute set, or nil
// when there is none. Spreading it keeps inline <script> output byte-identical
// for hosts that do not use CSP, instead of emitting nonce="".
func nonceAttrs(ctx context.Context) templ.Attributes { return shared.NonceAttrs(ctx) }
