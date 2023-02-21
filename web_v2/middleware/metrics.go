package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

type webV2MetricsUpdater interface {
	UpdateWebV2Metrics(duration time.Duration, status int)
}

// WebMetrics - middleware to calculate requests metrics
func WebMetrics(updater webV2MetricsUpdater) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw := middleware.NewWrapResponseWriter(w, 1)
			ts := time.Now()
			next.ServeHTTP(mw.(http.ResponseWriter), r)

			updater.UpdateWebV2Metrics(time.Since(ts), mw.Status())
		})
	}
}
