package middlewares

import (
	"HugeSpaceship/internal/config"
	"net/http"
	"strings"
)

// Compressioniddleware is deflate compression for LBP 1/2 (I think)
func Compressioniddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > 1000 && strings.Contains(r.Header.Get("Content-Type"), "text/xml") && r.Header.Get("Accept-Encoding") == "deflate" {

			}
		}
		return http.HandlerFunc(fn)
	}
}
