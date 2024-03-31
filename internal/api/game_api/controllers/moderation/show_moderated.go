package moderation

import (
	"net/http"
)

func ShowModeratedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte("<resources/>"))
	}
}
