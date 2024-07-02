//lint:file-ignore U1000 ignore go-swagger template
package docs

import (
	"github.com/emsi-zero/auth_ir/internal/api"
)

// swagger:route GET /verify verify verify-get
// Verifies a sign up.

// swagger:parameters verify-get
type verifyGetParamsWrapper struct {
	// in:query
	api.VerifyParams
}

// swagger:route POST /verify verify verify-post
// Verifies a sign up.

// swagger:parameters verify-post
type verifyPostParamsWrapper struct {
	// in:body
	Body api.VerifyParams
}
