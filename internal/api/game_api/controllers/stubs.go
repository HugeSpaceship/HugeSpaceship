package controllers

import (
	"net/http"
)

func StubEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}
