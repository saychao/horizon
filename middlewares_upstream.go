package horizon

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

const (
	upstreamHeader = "x-upstream"
)

func UpstreamMiddleware(c *web.C, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app := c.Env["app"].(*App)
		w.Header().Set(upstreamHeader, app.config.Hostname)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
