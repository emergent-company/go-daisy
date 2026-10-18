package galleryruntime

import (
	"context"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/shared"
)

// nonceAttrs returns the render-context CSP nonce as an attribute set, or nil
// when there is none.
func nonceAttrs(ctx context.Context) templ.Attributes { return shared.NonceAttrs(ctx) }
