package moderation

import (
	"io"
	"net/http"
)

// FilterHandler handles all text that can be moderated for the game
// Note: this is not used in LBPVita
func FilterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This is a stub for now
		// TODO: implement text filtering
		io.CopyN(w, r.Body, r.ContentLength)
	}
}
