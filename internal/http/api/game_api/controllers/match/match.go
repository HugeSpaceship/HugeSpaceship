package match

import (
	"net/http"
)

func MatchmakingEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}
