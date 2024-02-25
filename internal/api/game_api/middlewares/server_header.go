package middlewares

import (
	"net/http"
)

func ServerHeaderMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "HugeSpaceship")
	}
	return http.HandlerFunc(fn)
}
