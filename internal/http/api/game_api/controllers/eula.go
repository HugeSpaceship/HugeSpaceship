package controllers

import (
	"net/http"
)

const eulaText = `
WHAT'S UP FUCKERS, LBP3 SUCKS
`

func EulaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte(eulaText))
	}
}

func AnnounceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Make this configurable via the config file, or better yet integrate with the DB for a news list
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("")) // If it's an empty string then the client won't see it
	}
}
