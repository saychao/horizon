package horizon

import (
	"net/http"

	gctx "github.com/goji/context"
	"github.com/saychao/horizon/errors"
	"github.com/saychao/horizon/render/problem"
	"github.com/zenazn/goji/web"
)

// RecoverMiddleware helps the server recover from panics.  It ensures that
// no request can fully bring down the horizon server, and it also logs the
// panics to the logging subsystem.
func RecoverMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.FromC(*c)

		defer func() {
			if rec := recover(); rec != nil {
				err := errors.FromPanic(rec)
				problem.Render(ctx, w, err)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
