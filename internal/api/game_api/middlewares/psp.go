package middlewares

import (
	"net/http"
)

// Header names for psp version information

const PSPExeHeader = "X-exe-v"
const PSPDataHeader = "X-data-v"

func PSPVersionMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// If we're not on PSP, then bail
		if r.Header.Get(PSPExeHeader) == "" {
			return
		}

		// Pass through PSP Data and Exe headers
		// TODO: Make it so you can enforce a version
		w.Header().Set(PSPExeHeader, r.Header.Get(PSPExeHeader))
		w.Header().Set(PSPDataHeader, r.Header.Get(PSPDataHeader))

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
