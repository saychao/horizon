package horizon

import (
	"net/http"

	gctx "github.com/goji/context"
	"github.com/saychao/horizon/context/requestid"
	"github.com/saychao/horizon/httpx"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

func contextMiddleware(parent context.Context) func(c *web.C, next http.Handler) http.Handler {
	return func(c *web.C, next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := parent
			ctx = requestid.ContextFromC(ctx, c)
			ctx, cancel := httpx.RequestContext(ctx, w, r)

			gctx.Set(c, ctx)
			defer cancel()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
