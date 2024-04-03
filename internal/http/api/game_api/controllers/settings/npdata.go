package settings

import (
	"net/http"
)

func NpDataEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Add NPData support for matching
	}
}
