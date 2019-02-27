package middleware

import (
	"net/http"

	"github.com/moocss/chi-webserver/src/pkg/version"
)

// Version is a middleware function that appends the Bear version information
// to the HTTP response. This is intended for debugging and troubleshooting.
func Version(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("X-WEBSERVER-VERSION", version.Info.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}